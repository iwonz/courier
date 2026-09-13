// Package control provides authoritative, PID-independent control of Courier
// data workers discovered through the private delivery registry.
package control

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/worker"
)

// WorkerClient is the private IPC surface required by server control.
type WorkerClient interface {
	Hello(context.Context) (worker.HelloResponse, error)
	List(context.Context) (delivery.Snapshot, error)
	StopDelivery(context.Context, delivery.ID) error
	StopServer(context.Context) error
}

// ClientFactory binds a registry record to its private IPC client.
type ClientFactory func(delivery.Server) (WorkerClient, error)

// ServerView is one deterministic inventory entry. When Live is false, Server
// and Deliveries are explicitly registry observations rather than live state.
type ServerView struct {
	Server     delivery.Server
	Deliveries []delivery.Delivery
	Live       bool
}

// StopRequest selects one UUID or every registered data server.
type StopRequest struct {
	ID  delivery.ID
	All bool
}

// StopResult describes the scope actually selected by a stop operation.
type StopResult struct {
	ID             delivery.ID
	Kind           delivery.TargetKind
	AlreadyStopped bool
	StoppedServers int
}

// Service combines durable discovery with authoritative worker IPC.
type Service struct {
	Store   *delivery.Store
	Connect ClientFactory
}

// DefaultClient creates a client from secret-free registry metadata.
func DefaultClient(server delivery.Server) (WorkerClient, error) {
	client := worker.Client{Endpoint: server.ControlEndpoint, ServerID: server.ID, Compatibility: server.Compatibility}
	if err := client.Validate(); err != nil {
		return nil, err
	}
	return client, nil
}

// List returns every registered server, marking it live only after hello and a
// matching live snapshot. Unreachable entries remain visible for diagnosis.
func (service Service) List(ctx context.Context) ([]ServerView, error) {
	snapshot, err := service.snapshot(ctx)
	if err != nil {
		return nil, err
	}
	servers := append([]delivery.Server(nil), snapshot.Servers...)
	sort.Slice(servers, func(left, right int) bool { return servers[left].ID < servers[right].ID })
	views := make([]ServerView, 0, len(servers))
	for _, recorded := range servers {
		view := ServerView{Server: recorded, Deliveries: deliveriesFor(snapshot, recorded.ID)}
		_, live, liveSnapshot, liveErr := service.authoritative(ctx, recorded, "")
		if liveErr == nil {
			view.Live = true
			view.Server = live
			view.Deliveries = deliveriesFor(liveSnapshot, recorded.ID)
		}
		views = append(views, view)
	}
	return views, nil
}

// Stop resolves and authoritatively stops one target or every data server.
func (service Service) Stop(ctx context.Context, request StopRequest) (StopResult, error) {
	if request.All {
		if request.ID != "" {
			return StopResult{}, fmt.Errorf("%w: --all conflicts with a UUID", delivery.ErrInvalid)
		}
		return service.stopAll(ctx)
	}
	if !request.ID.Valid() {
		return StopResult{}, fmt.Errorf("%w: a canonical target UUID is required", delivery.ErrInvalid)
	}
	snapshot, err := service.snapshot(ctx)
	if err != nil {
		return StopResult{}, err
	}
	for _, tombstone := range snapshot.Tombstones {
		if tombstone.TargetID == request.ID {
			return StopResult{ID: request.ID, Kind: tombstone.Kind, AlreadyStopped: true}, nil
		}
	}
	for _, server := range snapshot.Servers {
		if server.ID == request.ID {
			client, _, _, err := service.authoritative(ctx, server, "")
			if err != nil {
				return StopResult{}, err
			}
			if err := client.StopServer(ctx); err != nil {
				return StopResult{}, err
			}
			return StopResult{ID: request.ID, Kind: delivery.TargetServer, StoppedServers: 1}, nil
		}
	}
	for _, item := range snapshot.Deliveries {
		if item.ID != request.ID {
			continue
		}
		server, _ := serverByID(snapshot, item.ServerID) // registry validation guarantees the owner
		client, _, _, err := service.authoritative(ctx, server, item.ID)
		if err != nil {
			return StopResult{}, err
		}
		if err := client.StopDelivery(ctx, item.ID); err != nil {
			return StopResult{}, err
		}
		return StopResult{ID: request.ID, Kind: delivery.TargetDelivery}, nil
	}
	return StopResult{}, fmt.Errorf("%w: target %s", delivery.ErrNotFound, request.ID)
}

func (service Service) stopAll(ctx context.Context) (StopResult, error) {
	snapshot, err := service.snapshot(ctx)
	if err != nil {
		return StopResult{}, err
	}
	servers := append([]delivery.Server(nil), snapshot.Servers...)
	sort.Slice(servers, func(left, right int) bool { return servers[left].ID < servers[right].ID })
	result := StopResult{Kind: delivery.TargetServer}
	var resultErr error
	for _, server := range servers {
		client, _, _, targetErr := service.authoritative(ctx, server, "")
		if targetErr == nil {
			targetErr = client.StopServer(ctx)
		}
		if targetErr != nil {
			resultErr = errors.Join(resultErr, fmt.Errorf("server %s: %w", server.ID, targetErr))
			continue
		}
		result.StoppedServers++
	}
	return result, resultErr
}

func (service Service) authoritative(ctx context.Context, recorded delivery.Server, target delivery.ID) (WorkerClient, delivery.Server, delivery.Snapshot, error) {
	connect := service.Connect
	if connect == nil {
		connect = DefaultClient
	}
	client, err := connect(recorded)
	if err != nil {
		return nil, delivery.Server{}, delivery.Snapshot{}, err
	}
	if _, err := client.Hello(ctx); err != nil {
		return nil, delivery.Server{}, delivery.Snapshot{}, err
	}
	snapshot, err := client.List(ctx)
	if err != nil {
		return nil, delivery.Server{}, delivery.Snapshot{}, err
	}
	live, found := serverByID(snapshot, recorded.ID)
	if !found || live.ControlEndpoint != recorded.ControlEndpoint || live.Bind != recorded.Bind || live.Compatibility != recorded.Compatibility {
		return nil, delivery.Server{}, delivery.Snapshot{}, errors.New("worker snapshot does not match the registered server")
	}
	if target != "" {
		found = false
		for _, item := range snapshot.Deliveries {
			if item.ID == target && item.ServerID == recorded.ID {
				found = true
				break
			}
		}
		if !found {
			return nil, delivery.Server{}, delivery.Snapshot{}, fmt.Errorf("%w: live delivery %s", delivery.ErrNotFound, target)
		}
	}
	return client, live, snapshot, nil
}

func (service Service) snapshot(ctx context.Context) (delivery.Snapshot, error) {
	if service.Store == nil {
		return delivery.Snapshot{}, fmt.Errorf("%w: control store is required", delivery.ErrInvalid)
	}
	return service.Store.Load(ctx)
}

func deliveriesFor(snapshot delivery.Snapshot, serverID delivery.ID) []delivery.Delivery {
	items := make([]delivery.Delivery, 0)
	for _, item := range snapshot.Deliveries {
		if item.ServerID == serverID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(left, right int) bool { return items[left].ID < items[right].ID })
	return items
}

func serverByID(snapshot delivery.Snapshot, id delivery.ID) (delivery.Server, bool) {
	for _, server := range snapshot.Servers {
		if server.ID == id {
			return server, true
		}
	}
	return delivery.Server{}, false
}

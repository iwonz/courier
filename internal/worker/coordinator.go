package worker

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gofrs/flock"
	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
)

var ErrIncompatible = errors.New("incompatible Courier worker already owns bind")

type lockHandle interface {
	TryLockContext(context.Context, time.Duration) (bool, error)
	Close() error
}

var newLockHandle = func(path string) lockHandle {
	return flock.New(path, flock.SetPermissions(0o600))
}

var (
	makeLockDirectory   = os.MkdirAll
	chmodLockDirectory  = os.Chmod
	makeControlEndpoint = ipc.ControlEndpoint
)

type BindLocker interface {
	Acquire(context.Context, string) (func() error, error)
}

type FileBindLocker struct {
	Directory string
	Retry     time.Duration
}

func (locker FileBindLocker) Acquire(ctx context.Context, bind string) (func() error, error) {
	canonical, err := CanonicalBind(bind)
	if err != nil || strings.TrimSpace(locker.Directory) == "" {
		return nil, fmt.Errorf("%w: invalid bind lock configuration", delivery.ErrInvalid)
	}
	retry := locker.Retry
	if retry <= 0 {
		retry = 10 * time.Millisecond
	}
	digest := sha256.Sum256([]byte(canonical))
	handle := newLockHandle(filepath.Join(locker.Directory, fmt.Sprintf("bind-%x.lock", digest)))
	locked, err := handle.TryLockContext(ctx, retry)
	if err != nil {
		_ = handle.Close()
		return nil, err
	}
	if !locked {
		_ = handle.Close()
		return nil, fmt.Errorf("%w: bind lock was not acquired", delivery.ErrRevisionConflict)
	}
	return handle.Close, nil
}

func CanonicalBind(value string) (string, error) {
	if strings.TrimSpace(value) != value || strings.ContainsAny(value, "\x00\r\n") {
		return "", fmt.Errorf("%w: invalid bind", delivery.ErrInvalid)
	}
	host, port, err := net.SplitHostPort(value)
	if err != nil {
		return "", fmt.Errorf("%w: bind must use host:port", delivery.ErrInvalid)
	}
	number, err := strconv.ParseUint(port, 10, 16)
	if err != nil || number == 0 {
		return "", fmt.Errorf("%w: bind port must be numeric and non-zero", delivery.ErrInvalid)
	}
	if address, parseErr := netip.ParseAddr(host); parseErr == nil {
		host = address.String()
	} else {
		host = strings.ToLower(host)
		if strings.ContainsAny(host, " /\\") {
			return "", fmt.Errorf("%w: invalid bind host", delivery.ErrInvalid)
		}
	}
	return net.JoinHostPort(host, strconv.FormatUint(number, 10)), nil
}

type AcquireRequest struct {
	Bind          string
	Compatibility string
	Route         delivery.Route
	Policy        delivery.Policy
	Foreground    bool
	At            time.Time
}

func (request AcquireRequest) Validate() error {
	if _, err := CanonicalBind(request.Bind); err != nil {
		return err
	}
	if invalidCompatibility(request.Compatibility) || !request.Route.Valid() || request.At.IsZero() {
		return fmt.Errorf("%w: invalid worker acquisition", delivery.ErrInvalid)
	}
	return request.Policy.Validate()
}

type LaunchRequest struct {
	ServerID        delivery.ID
	Bind            string
	ControlEndpoint string
	Compatibility   string
}

type LaunchFunc func(context.Context, LaunchRequest) (Client, error)

type CleanupFunc func(context.Context, delivery.OwnedTemp) error

type Acquired struct {
	ServerID   delivery.ID
	DeliveryID delivery.ID
	Reused     bool
	Lease      *Lease
}

type Coordinator struct {
	Store          *delivery.Store
	StateDirectory string
	Locks          BindLocker
	Launch         LaunchFunc
	Cleanup        CleanupFunc
	Hello          func(context.Context, Client) error
	Register       func(context.Context, Client, delivery.Delivery, bool) (*Lease, error)
	RemoveEndpoint func(string) error
	StaleProbe     func(error) bool
}

func (coordinator *Coordinator) Acquire(ctx context.Context, request AcquireRequest) (result Acquired, resultErr error) {
	if coordinator == nil || coordinator.Store == nil || coordinator.Locks == nil || coordinator.Launch == nil || strings.TrimSpace(coordinator.StateDirectory) == "" || filepath.Clean(coordinator.StateDirectory) != filepath.Clean(coordinator.Store.Directory()) {
		return Acquired{}, fmt.Errorf("%w: incomplete worker coordinator", delivery.ErrInvalid)
	}
	if err := request.Validate(); err != nil {
		return Acquired{}, err
	}
	bind, _ := CanonicalBind(request.Bind)
	release, err := coordinator.Locks.Acquire(ctx, bind)
	if err != nil {
		return Acquired{}, err
	}
	defer func() { resultErr = errors.Join(resultErr, release()) }()

	snapshot, err := coordinator.Store.Load(ctx)
	if err != nil {
		return Acquired{}, err
	}
	servers := append([]delivery.Server(nil), snapshot.Servers...)
	sort.Slice(servers, func(left, right int) bool { return servers[left].ID < servers[right].ID })
	for _, server := range servers {
		serverBind, bindErr := CanonicalBind(server.Bind)
		if bindErr != nil {
			return Acquired{}, bindErr
		}
		if serverBind != bind {
			continue
		}
		client := Client{Endpoint: server.ControlEndpoint, ServerID: server.ID, Compatibility: request.Compatibility}
		probeErr := coordinator.hello(ctx, client)
		if probeErr == nil {
			return coordinator.register(ctx, client, request, true)
		}
		if ipc.IsRemoteError(probeErr, ipc.CodeConflict) {
			return Acquired{}, fmt.Errorf("%w: %s", ErrIncompatible, bind)
		}
		stale := coordinator.StaleProbe
		if stale == nil {
			stale = staleProbeError
		}
		if !stale(probeErr) {
			return Acquired{}, probeErr
		}
		if err := coordinator.reconcileStale(ctx, snapshot, server); err != nil {
			return Acquired{}, err
		}
	}

	serverID := delivery.NewID()
	controlEndpoint, err := makeControlEndpoint(coordinator.StateDirectory, serverID)
	if err != nil {
		return Acquired{}, err
	}
	client, err := coordinator.Launch(ctx, LaunchRequest{ServerID: serverID, Bind: bind, ControlEndpoint: controlEndpoint, Compatibility: request.Compatibility})
	if err != nil {
		return Acquired{}, err
	}
	if err := client.Validate(); err != nil || client.ServerID != serverID || client.Endpoint != controlEndpoint || client.Compatibility != request.Compatibility {
		return Acquired{}, errors.Join(err, fmt.Errorf("%w: launcher returned mismatched worker", ipc.ErrProtocol))
	}
	if err := coordinator.hello(ctx, client); err != nil {
		_ = client.Shutdown(context.Background())
		return Acquired{}, err
	}
	result, err = coordinator.register(ctx, client, request, false)
	if err != nil {
		_ = client.Shutdown(context.Background())
		return Acquired{}, err
	}
	return result, nil
}

func staleProbeError(err error) bool {
	return errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, net.ErrClosed)
}

func (coordinator *Coordinator) hello(ctx context.Context, client Client) error {
	if coordinator.Hello != nil {
		return coordinator.Hello(ctx, client)
	}
	_, err := client.Hello(ctx)
	return err
}

func (coordinator *Coordinator) register(ctx context.Context, client Client, request AcquireRequest, reused bool) (Acquired, error) {
	item := delivery.Delivery{
		ID: delivery.NewID(), ServerID: client.ServerID, Route: request.Route, State: delivery.StateStarting,
		Policy: request.Policy, CreatedAt: request.At.UTC(), UpdatedAt: request.At.UTC(),
	}
	var lease *Lease
	var err error
	if coordinator.Register != nil {
		lease, err = coordinator.Register(ctx, client, item, request.Foreground)
	} else {
		_, lease, err = client.Register(ctx, item, request.Foreground)
	}
	if err != nil {
		return Acquired{}, err
	}
	return Acquired{ServerID: client.ServerID, DeliveryID: item.ID, Reused: reused, Lease: lease}, nil
}

func (coordinator *Coordinator) reconcileStale(ctx context.Context, snapshot delivery.Snapshot, server delivery.Server) error {
	remove := coordinator.RemoveEndpoint
	if remove == nil {
		remove = ipc.RemoveStale
	}
	if err := remove(server.ControlEndpoint); err != nil {
		return err
	}
	owners := map[delivery.ID]bool{server.ID: true}
	items := make([]delivery.Delivery, 0)
	for _, item := range snapshot.Deliveries {
		if item.ServerID == server.ID {
			owners[item.ID] = true
			items = append(items, item)
		}
	}
	for _, temporary := range snapshot.OwnedTemps {
		if !owners[temporary.OwnerID] {
			continue
		}
		if coordinator.Cleanup == nil {
			return fmt.Errorf("%w: stale temporary cleanup is unavailable", delivery.ErrInvalid)
		}
		if err := coordinator.Cleanup(ctx, temporary); err != nil {
			return err
		}
		if _, err := updateStore(ctx, coordinator.Store, func(registry *delivery.Registry) error {
			return registry.RemoveOwnedTemp(temporary.ID)
		}); err != nil {
			return err
		}
	}
	at := time.Now().UTC()
	_, err := updateStore(ctx, coordinator.Store, func(registry *delivery.Registry) error {
		for _, item := range items {
			if err := registry.TombstoneDelivery(item.ID, delivery.Tombstone{ID: delivery.NewID(), TargetID: item.ID, Kind: delivery.TargetDelivery, Reason: delivery.ReasonStale, At: at}); err != nil {
				return err
			}
		}
		return registry.TombstoneServer(server.ID, delivery.Tombstone{ID: delivery.NewID(), TargetID: server.ID, Kind: delivery.TargetServer, Reason: delivery.ReasonStale, At: at})
	})
	return err
}

func EnsurePrivateLockDirectory(directory string) error {
	if strings.TrimSpace(directory) == "" {
		return fmt.Errorf("%w: lock directory is required", delivery.ErrInvalid)
	}
	if err := makeLockDirectory(directory, 0o700); err != nil {
		return err
	}
	return chmodLockDirectory(directory, 0o700)
}

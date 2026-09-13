package control

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/worker"
)

const (
	serverA   = delivery.ID("00000000-0000-4000-8000-000000000101")
	serverB   = delivery.ID("00000000-0000-4000-8000-000000000102")
	deliveryA = delivery.ID("00000000-0000-4000-8000-000000000201")
	deliveryB = delivery.ID("00000000-0000-4000-8000-000000000202")
	stoppedID = delivery.ID("00000000-0000-4000-8000-000000000203")
)

var controlTime = time.Date(2026, 9, 14, 1, 2, 3, 0, time.UTC)

type fakeClient struct {
	helloErr        error
	listErr         error
	snapshot        delivery.Snapshot
	stopDeliveryErr error
	stopServerErr   error
	updatePolicyErr error
	stoppedDelivery delivery.ID
	updatedPolicy   worker.UpdatePolicyRequest
	serverStops     int
}

func (client *fakeClient) Hello(context.Context) (worker.HelloResponse, error) {
	return worker.HelloResponse{}, client.helloErr
}
func (client *fakeClient) List(context.Context) (delivery.Snapshot, error) {
	return client.snapshot, client.listErr
}
func (client *fakeClient) UpdatePolicy(_ context.Context, request worker.UpdatePolicyRequest) error {
	client.updatedPolicy = request
	return client.updatePolicyErr
}
func (client *fakeClient) StopDelivery(_ context.Context, id delivery.ID) error {
	client.stoppedDelivery = id
	return client.stopDeliveryErr
}
func (client *fakeClient) StopServer(context.Context) error {
	client.serverStops++
	return client.stopServerErr
}

func testServer(id delivery.ID, bind string) delivery.Server {
	return delivery.Server{
		ID: id, Bind: bind, ControlEndpoint: filepath.Join(os.TempDir(), string(id)+".sock"), Compatibility: "web-v1/test",
		ProcessID: 42, State: delivery.StateActive, StartedAt: controlTime, UpdatedAt: controlTime,
	}
}

func testDelivery(id, owner delivery.ID) delivery.Delivery {
	return delivery.Delivery{
		ID: id, ServerID: owner, Route: delivery.RoutePathToWeb, Source: "./source", Destination: "web://",
		State: delivery.StateActive, Policy: delivery.DefaultPolicy(), CreatedAt: controlTime, UpdatedAt: controlTime,
	}
}

func seededService(t *testing.T) (Service, delivery.Snapshot) {
	t.Helper()
	store, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	first := testServer(serverA, "127.0.0.1:8101")
	second := testServer(serverB, "127.0.0.1:8102")
	firstDelivery := testDelivery(deliveryA, serverA)
	secondDelivery := testDelivery(deliveryB, serverB)
	stopped := testDelivery(stoppedID, serverA)
	snapshot, err := store.Update(context.Background(), 0, func(registry *delivery.Registry) error {
		for _, server := range []delivery.Server{second, first} {
			if err := registry.RegisterServer(server); err != nil {
				return err
			}
		}
		for _, item := range []delivery.Delivery{secondDelivery, stopped, firstDelivery} {
			if err := registry.RegisterDelivery(item); err != nil {
				return err
			}
		}
		return registry.TombstoneDelivery(stopped.ID, delivery.Tombstone{
			ID: delivery.ID("00000000-0000-4000-8000-000000000301"), TargetID: stopped.ID,
			Kind: delivery.TargetDelivery, Reason: delivery.ReasonStopped, At: controlTime,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return Service{Store: store}, snapshot
}

func TestAuthoritativeInventory(t *testing.T) {
	service, snapshot := seededService(t)
	clients := map[delivery.ID]*fakeClient{
		serverA: {snapshot: snapshot},
		serverB: {helloErr: errors.New("unreachable"), snapshot: snapshot},
	}
	service.Connect = func(server delivery.Server) (WorkerClient, error) { return clients[server.ID], nil }
	views, err := service.List(context.Background())
	if err != nil || len(views) != 2 || views[0].Server.ID != serverA || !views[0].Live || len(views[0].Deliveries) != 1 || views[0].Deliveries[0].ID != deliveryA || views[1].Live || len(views[1].Deliveries) != 1 {
		t.Fatalf("views=%+v err=%v", views, err)
	}
	if _, err := (Service{}).List(context.Background()); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("nil service=%v", err)
	}

	want := errors.New("failure")
	for _, test := range []struct {
		name      string
		connect   ClientFactory
		configure func(*fakeClient, *delivery.Snapshot)
	}{
		{"connect", func(delivery.Server) (WorkerClient, error) { return nil, want }, nil},
		{"list", nil, func(client *fakeClient, _ *delivery.Snapshot) { client.listErr = want }},
		{"mismatch", nil, func(_ *fakeClient, live *delivery.Snapshot) { live.Servers[0].Bind = "127.0.0.1:9999" }},
	} {
		t.Run(test.name, func(t *testing.T) {
			oneStore, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
			if err != nil {
				t.Fatal(err)
			}
			defer oneStore.Close()
			server := testServer(serverA, "127.0.0.1:8101")
			live, err := oneStore.Update(context.Background(), 0, func(registry *delivery.Registry) error { return registry.RegisterServer(server) })
			if err != nil {
				t.Fatal(err)
			}
			client := &fakeClient{snapshot: live}
			if test.configure != nil {
				test.configure(client, &client.snapshot)
			}
			connect := test.connect
			if connect == nil {
				connect = func(delivery.Server) (WorkerClient, error) { return client, nil }
			}
			views, err := (Service{Store: oneStore, Connect: connect}).List(context.Background())
			if err != nil || len(views) != 1 || views[0].Live {
				t.Fatalf("views=%+v err=%v", views, err)
			}
		})
	}
}

func TestScopedStops(t *testing.T) {
	service, snapshot := seededService(t)
	clients := map[delivery.ID]*fakeClient{serverA: {snapshot: snapshot}, serverB: {snapshot: snapshot}}
	service.Connect = func(server delivery.Server) (WorkerClient, error) { return clients[server.ID], nil }

	if _, err := service.Stop(context.Background(), StopRequest{All: true, ID: serverA}); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("all conflict=%v", err)
	}
	if _, err := service.Stop(context.Background(), StopRequest{All: true, Kind: delivery.TargetServer}); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("all kind conflict=%v", err)
	}
	if _, err := service.Stop(context.Background(), StopRequest{Kind: "invalid", ID: serverA}); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid kind=%v", err)
	}
	if _, err := service.Stop(context.Background(), StopRequest{ID: "bad"}); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid UUID=%v", err)
	}
	result, err := service.Stop(context.Background(), StopRequest{ID: stoppedID})
	if err != nil || !result.AlreadyStopped || result.Kind != delivery.TargetDelivery {
		t.Fatalf("tombstone=%+v err=%v", result, err)
	}
	result, err = service.Stop(context.Background(), StopRequest{ID: deliveryA})
	if err != nil || result.Kind != delivery.TargetDelivery || clients[serverA].stoppedDelivery != deliveryA {
		t.Fatalf("delivery=%+v err=%v", result, err)
	}
	result, err = service.Stop(context.Background(), StopRequest{ID: serverB})
	if err != nil || result.Kind != delivery.TargetServer || result.StoppedServers != 1 || clients[serverB].serverStops != 1 {
		t.Fatalf("server=%+v err=%v", result, err)
	}
	if _, err := service.Stop(context.Background(), StopRequest{ID: serverB, Kind: delivery.TargetDelivery}); !errors.Is(err, delivery.ErrNotFound) || clients[serverB].serverStops != 1 {
		t.Fatalf("server as delivery=%v", err)
	}
	if _, err := service.Stop(context.Background(), StopRequest{ID: deliveryA, Kind: delivery.TargetServer}); !errors.Is(err, delivery.ErrNotFound) || clients[serverA].stoppedDelivery != deliveryA {
		t.Fatalf("delivery as server=%v", err)
	}
	if _, err := service.Stop(context.Background(), StopRequest{ID: stoppedID, Kind: delivery.TargetServer}); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("tombstone kind=%v", err)
	}
	unknown := delivery.ID("00000000-0000-4000-8000-000000000999")
	if _, err := service.Stop(context.Background(), StopRequest{ID: unknown}); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("unknown=%v", err)
	}
}

func TestStopFailuresAndStopAll(t *testing.T) {
	service, snapshot := seededService(t)
	want := errors.New("failure")
	clients := map[delivery.ID]*fakeClient{serverA: {snapshot: snapshot}, serverB: {snapshot: snapshot, stopServerErr: want}}
	service.Connect = func(server delivery.Server) (WorkerClient, error) { return clients[server.ID], nil }
	result, err := service.Stop(context.Background(), StopRequest{All: true})
	if !errors.Is(err, want) || result.StoppedServers != 1 || clients[serverA].serverStops != 1 || clients[serverB].serverStops != 1 {
		t.Fatalf("all=%+v err=%v clients=%+v", result, err, clients)
	}

	clients[serverA].stopDeliveryErr = want
	if _, err := service.Stop(context.Background(), StopRequest{ID: deliveryA}); !errors.Is(err, want) {
		t.Fatalf("delivery stop error=%v", err)
	}
	clients[serverA].stopDeliveryErr = nil
	clients[serverA].snapshot.Deliveries = deliveriesFor(clients[serverA].snapshot, serverB)
	if _, err := service.Stop(context.Background(), StopRequest{ID: deliveryA}); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("missing live delivery=%v", err)
	}
	clients[serverA].snapshot = snapshot
	clients[serverA].helloErr = want
	if _, err := service.Stop(context.Background(), StopRequest{ID: serverA}); !errors.Is(err, want) {
		t.Fatalf("hello error=%v", err)
	}
	clients[serverA].helloErr = nil
	clients[serverA].stopServerErr = want
	if _, err := service.Stop(context.Background(), StopRequest{ID: serverA}); !errors.Is(err, want) {
		t.Fatalf("server stop error=%v", err)
	}

	emptyStore, err := delivery.OpenStore(filepath.Join(t.TempDir(), "empty"))
	if err != nil {
		t.Fatal(err)
	}
	defer emptyStore.Close()
	result, err = (Service{Store: emptyStore}).Stop(context.Background(), StopRequest{All: true})
	if err != nil || result.StoppedServers != 0 {
		t.Fatalf("empty all=%+v err=%v", result, err)
	}
	if err := emptyStore.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := (Service{Store: emptyStore}).Stop(context.Background(), StopRequest{All: true}); err == nil {
		t.Fatal("closed store accepted")
	}
	if _, err := (Service{Store: emptyStore}).Stop(context.Background(), StopRequest{ID: serverA}); err == nil {
		t.Fatal("closed store accepted for scoped stop")
	}
}

func TestDefaultClientAndHelpers(t *testing.T) {
	server := testServer(serverA, "127.0.0.1:8101")
	client, err := DefaultClient(server)
	if err != nil || client == nil {
		t.Fatalf("client=%v err=%v", client, err)
	}
	server.Compatibility = ""
	if _, err := DefaultClient(server); err == nil {
		t.Fatal("legacy server unexpectedly produced an authoritative client")
	}
	if _, found := serverByID(delivery.EmptySnapshot(), serverA); found {
		t.Fatal("missing server found")
	}
	service := Service{}
	if _, _, _, err := service.authoritative(context.Background(), server, ""); err == nil {
		t.Fatal("default authoritative client unexpectedly connected")
	}
	snapshot := delivery.EmptySnapshot()
	snapshot.Deliveries = []delivery.Delivery{
		testDelivery(deliveryB, serverA),
		testDelivery(deliveryA, serverA),
	}
	items := deliveriesFor(snapshot, serverA)
	if len(items) != 2 || items[0].ID != deliveryA || items[1].ID != deliveryB {
		t.Fatalf("sorted deliveries=%v", items)
	}
}

func TestOptimisticPolicyUpdate(t *testing.T) {
	service, snapshot := seededService(t)
	client := &fakeClient{snapshot: snapshot}
	service.Connect = func(delivery.Server) (WorkerClient, error) { return client, nil }
	policy := delivery.DefaultPolicy()
	policy.Version = 2
	if err := service.UpdatePolicy(context.Background(), deliveryA, 1, policy); err != nil || client.updatedPolicy.DeliveryID != deliveryA || client.updatedPolicy.ExpectedVersion != 1 || client.updatedPolicy.Policy.Version != 2 {
		t.Fatalf("update=%+v err=%v", client.updatedPolicy, err)
	}

	want := errors.New("update failure")
	client.updatePolicyErr = want
	if err := service.UpdatePolicy(context.Background(), deliveryA, 1, policy); !errors.Is(err, want) {
		t.Fatalf("worker update error=%v", err)
	}
	client.updatePolicyErr = nil
	client.helloErr = want
	if err := service.UpdatePolicy(context.Background(), deliveryA, 1, policy); !errors.Is(err, want) {
		t.Fatalf("authority error=%v", err)
	}

	for _, test := range []struct {
		id       delivery.ID
		expected uint64
		policy   delivery.Policy
	}{
		{"bad", 1, policy},
		{deliveryA, 0, policy},
		{deliveryA, 1, delivery.DefaultPolicy()},
	} {
		if err := service.UpdatePolicy(context.Background(), test.id, test.expected, test.policy); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("invalid update=%v", err)
		}
	}
	invalidPolicy := policy
	invalidPolicy.Auth = "invalid"
	if err := service.UpdatePolicy(context.Background(), deliveryA, 1, invalidPolicy); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid policy=%v", err)
	}
	unknown := delivery.ID("00000000-0000-4000-8000-000000000999")
	if err := service.UpdatePolicy(context.Background(), unknown, 1, policy); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("unknown delivery=%v", err)
	}

	closed, err := delivery.OpenStore(filepath.Join(t.TempDir(), "closed"))
	if err != nil {
		t.Fatal(err)
	}
	if err := closed.Close(); err != nil {
		t.Fatal(err)
	}
	if err := (Service{Store: closed}).UpdatePolicy(context.Background(), deliveryA, 1, policy); err == nil {
		t.Fatal("closed store accepted")
	}
}

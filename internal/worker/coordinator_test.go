package worker

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
)

func TestCanonicalBindAndAcquireValidation(t *testing.T) {
	valid := map[string]string{
		"LOCALHOST:08080":     "localhost:8080",
		"127.0.0.1:8080":      "127.0.0.1:8080",
		"[2001:0db8::1]:8080": "[2001:db8::1]:8080",
		":8080":               ":8080",
	}
	for input, expected := range valid {
		actual, err := CanonicalBind(input)
		if err != nil || actual != expected {
			t.Fatalf("input=%q actual=%q err=%v", input, actual, err)
		}
	}
	for _, input := range []string{" localhost:80", "localhost", "localhost:http", "localhost:0", "bad host:80", "bad\\host:80", "localhost:70000", "localhost:80\n"} {
		if _, err := CanonicalBind(input); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("input=%q err=%v", input, err)
		}
	}
	request := AcquireRequest{Bind: "localhost:8080", Compatibility: "http-v1", Route: delivery.RoutePathToWeb, Policy: delivery.DefaultPolicy(), Foreground: true, At: workerTime}
	if err := request.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*AcquireRequest){
		func(value *AcquireRequest) { value.Bind = "bad" },
		func(value *AcquireRequest) { value.Compatibility = "" },
		func(value *AcquireRequest) { value.Route = "bad" },
		func(value *AcquireRequest) { value.At = time.Time{} },
		func(value *AcquireRequest) { value.Policy.Version = 0 },
	} {
		candidate := request
		mutate(&candidate)
		if err := candidate.Validate(); err == nil {
			t.Fatalf("candidate=%+v accepted", candidate)
		}
	}
}

func TestFileBindLocker(t *testing.T) {
	originalNew := newLockHandle
	t.Cleanup(func() { newLockHandle = originalNew })
	if _, err := (FileBindLocker{}).Acquire(context.Background(), "localhost:8080"); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("directory=%v", err)
	}
	if _, err := (FileBindLocker{Directory: t.TempDir()}).Acquire(context.Background(), "bad"); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("bind=%v", err)
	}
	want := errors.New("failure")
	fake := &fakeLock{err: want}
	newLockHandle = func(string) lockHandle { return fake }
	if _, err := (FileBindLocker{Directory: t.TempDir()}).Acquire(context.Background(), "localhost:8080"); !errors.Is(err, want) || !fake.closed {
		t.Fatalf("lock=%v closed=%v", err, fake.closed)
	}
	fake = &fakeLock{}
	newLockHandle = func(string) lockHandle { return fake }
	if _, err := (FileBindLocker{Directory: t.TempDir()}).Acquire(context.Background(), "localhost:8080"); !errors.Is(err, delivery.ErrRevisionConflict) || !fake.closed {
		t.Fatalf("not locked=%v closed=%v", err, fake.closed)
	}
	fake = &fakeLock{locked: true, closeErr: want}
	newLockHandle = func(string) lockHandle { return fake }
	release, err := (FileBindLocker{Directory: t.TempDir(), Retry: time.Millisecond}).Acquire(context.Background(), "localhost:8080")
	if err != nil || release == nil {
		t.Fatal(err)
	}
	if err := release(); !errors.Is(err, want) {
		t.Fatalf("release=%v", err)
	}

	newLockHandle = originalNew
	directory := t.TempDir()
	locker := FileBindLocker{Directory: directory, Retry: time.Millisecond}
	first, err := locker.Acquire(context.Background(), "localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	canceled, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, err := locker.Acquire(canceled, "localhost:8080"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("contended=%v", err)
	}
	if err := first(); err != nil {
		t.Fatal(err)
	}
	second, err := locker.Acquire(context.Background(), "localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	if err := second(); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(directory)
	if err != nil || len(entries) != 1 {
		t.Fatalf("entries=%d err=%v", len(entries), err)
	}
}

func TestEnsurePrivateLockDirectory(t *testing.T) {
	originalMkdir, originalChmod := makeLockDirectory, chmodLockDirectory
	t.Cleanup(func() { makeLockDirectory, chmodLockDirectory = originalMkdir, originalChmod })
	if err := EnsurePrivateLockDirectory(" "); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("empty=%v", err)
	}
	want := errors.New("failure")
	makeLockDirectory = func(string, fs.FileMode) error { return want }
	if err := EnsurePrivateLockDirectory("locks"); !errors.Is(err, want) {
		t.Fatalf("mkdir=%v", err)
	}
	makeLockDirectory = func(string, fs.FileMode) error { return nil }
	chmodLockDirectory = func(string, fs.FileMode) error { return want }
	if err := EnsurePrivateLockDirectory("locks"); !errors.Is(err, want) {
		t.Fatalf("chmod=%v", err)
	}
	makeLockDirectory, chmodLockDirectory = originalMkdir, originalChmod
	directory := filepath.Join(t.TempDir(), "locks")
	if err := EnsurePrivateLockDirectory(directory); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(directory)
	if err != nil || info.Mode().Perm() != 0o700 {
		t.Fatalf("mode=%v err=%v", info.Mode(), err)
	}
}

func TestCoordinatorReuseAndFailures(t *testing.T) {
	store, server := seededCoordinatorStore(t, false)
	current, err := store.Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	other := server
	other.ID = delivery.ID("00000000-0000-4000-8000-000000000200")
	other.Bind = "127.0.0.1:9090"
	other.ControlEndpoint = "other-control"
	if _, err := store.Update(context.Background(), current.Revision, func(registry *delivery.Registry) error { return registry.RegisterServer(other) }); err != nil {
		t.Fatal(err)
	}
	request := validAcquireRequest(server.Bind)
	locker := &fakeBindLocker{}
	launchCalls := 0
	registered := delivery.Delivery{}
	coordinator := &Coordinator{
		Store: store, StateDirectory: store.Directory(), Locks: locker,
		Launch: func(context.Context, LaunchRequest) (Client, error) {
			launchCalls++
			return Client{}, errors.New("unexpected launch")
		},
		Hello: func(context.Context, Client) error { return nil },
		Register: func(_ context.Context, client Client, item delivery.Delivery, foreground bool) (*Lease, error) {
			registered = item
			if !foreground || client.ServerID != server.ID {
				t.Fatal("registration metadata")
			}
			return &Lease{}, nil
		},
	}
	result, err := coordinator.Acquire(context.Background(), request)
	if err != nil || !result.Reused || result.ServerID != server.ID || result.DeliveryID != registered.ID || result.Lease == nil || launchCalls != 0 || !locker.released {
		t.Fatalf("result=%+v launch=%d release=%v err=%v", result, launchCalls, locker.released, err)
	}

	coordinator.Hello = func(context.Context, Client) error {
		return &ipc.RemoteError{Code: ipc.CodeConflict, Message: "incompatible"}
	}
	if _, err := coordinator.Acquire(context.Background(), request); !errors.Is(err, ErrIncompatible) {
		t.Fatalf("incompatible=%v", err)
	}
	want := errors.New("uncertain")
	coordinator.Hello = func(context.Context, Client) error { return want }
	if _, err := coordinator.Acquire(context.Background(), request); !errors.Is(err, want) {
		t.Fatalf("uncertain=%v", err)
	}

	bad := *coordinator
	bad.Store = nil
	if _, err := bad.Acquire(context.Background(), request); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("incomplete=%v", err)
	}
	bad = *coordinator
	bad.StateDirectory = filepath.Join(store.Directory(), "other")
	if _, err := bad.Acquire(context.Background(), request); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("directory mismatch=%v", err)
	}
	if _, err := coordinator.Acquire(context.Background(), AcquireRequest{}); err == nil {
		t.Fatal("invalid request accepted")
	}
	locker.err = want
	if _, err := coordinator.Acquire(context.Background(), request); !errors.Is(err, want) {
		t.Fatalf("lock=%v", err)
	}
	locker.err = nil
	locker.releaseErr = want
	coordinator.Hello = func(context.Context, Client) error { return nil }
	if _, err := coordinator.Acquire(context.Background(), request); !errors.Is(err, want) {
		t.Fatalf("release=%v", err)
	}
}

func TestCoordinatorInvalidStoredBindAndControlEndpointFailure(t *testing.T) {
	store, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	server := delivery.Server{ID: workerServerID, Bind: "bad host:8080", ControlEndpoint: "control", ProcessID: 1, State: delivery.StateActive, StartedAt: workerTime, UpdatedAt: workerTime}
	if _, err := store.Update(context.Background(), 0, func(registry *delivery.Registry) error { return registry.RegisterServer(server) }); err != nil {
		t.Fatal(err)
	}
	coordinator := &Coordinator{
		Store: store, StateDirectory: store.Directory(), Locks: &fakeBindLocker{},
		Launch: func(context.Context, LaunchRequest) (Client, error) { return Client{}, nil },
	}
	if _, err := coordinator.Acquire(context.Background(), validAcquireRequest("127.0.0.1:8080")); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("stored bind=%v", err)
	}

	originalEndpoint := makeControlEndpoint
	t.Cleanup(func() { makeControlEndpoint = originalEndpoint })
	fresh, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Close()
	want := errors.New("endpoint")
	makeControlEndpoint = func(string, delivery.ID) (string, error) { return "", want }
	coordinator.Store = fresh
	coordinator.StateDirectory = fresh.Directory()
	if _, err := coordinator.Acquire(context.Background(), validAcquireRequest("127.0.0.1:8080")); !errors.Is(err, want) {
		t.Fatalf("control endpoint=%v", err)
	}
}

func TestCoordinatorStaleRecoveryAndLaunch(t *testing.T) {
	store, staleServer := seededCoordinatorStore(t, true)
	request := validAcquireRequest(staleServer.Bind)
	removed := ""
	cleaned := []delivery.ID{}
	launched := LaunchRequest{}
	coordinator := &Coordinator{
		Store: store, StateDirectory: store.Directory(), Locks: &fakeBindLocker{},
		Hello: func(_ context.Context, client Client) error {
			if client.ServerID == staleServer.ID {
				return os.ErrNotExist
			}
			return nil
		},
		RemoveEndpoint: func(endpoint string) error { removed = endpoint; return nil },
		Cleanup: func(_ context.Context, temporary delivery.OwnedTemp) error {
			cleaned = append(cleaned, temporary.ID)
			return nil
		},
		Launch: func(_ context.Context, request LaunchRequest) (Client, error) {
			launched = request
			return Client{Endpoint: request.ControlEndpoint, ServerID: request.ServerID, Compatibility: request.Compatibility}, nil
		},
		Register: func(context.Context, Client, delivery.Delivery, bool) (*Lease, error) { return nil, nil },
	}
	result, err := coordinator.Acquire(context.Background(), request)
	if err != nil || result.Reused || removed != staleServer.ControlEndpoint || len(cleaned) != 2 || launched.ServerID != result.ServerID {
		t.Fatalf("result=%+v removed=%q cleaned=%v launch=%+v err=%v", result, removed, cleaned, launched, err)
	}
	snapshot, err := store.Load(context.Background())
	if err != nil || len(snapshot.Servers) != 0 || len(snapshot.Deliveries) != 0 || len(snapshot.OwnedTemps) != 0 || len(snapshot.Tombstones) != 2 {
		t.Fatalf("snapshot=%+v err=%v", snapshot, err)
	}
}

func TestCoordinatorStaleAndLaunchErrors(t *testing.T) {
	want := errors.New("failure")
	for _, test := range []struct {
		name      string
		configure func(*Coordinator, delivery.Server)
	}{
		{"remove endpoint", func(value *Coordinator, _ delivery.Server) { value.RemoveEndpoint = func(string) error { return want } }},
		{"cleanup unavailable", func(value *Coordinator, _ delivery.Server) { value.Cleanup = nil }},
		{"cleanup failure", func(value *Coordinator, _ delivery.Server) {
			value.Cleanup = func(context.Context, delivery.OwnedTemp) error { return want }
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			store, server := seededCoordinatorStore(t, true)
			coordinator := staleCoordinator(store)
			test.configure(coordinator, server)
			if _, err := coordinator.Acquire(context.Background(), validAcquireRequest(server.Bind)); err == nil {
				t.Fatal("expected stale recovery failure")
			}
		})
	}

	t.Run("cleanup registry update", func(t *testing.T) {
		store, server := seededCoordinatorStore(t, true)
		coordinator := staleCoordinator(store)
		coordinator.Cleanup = func(context.Context, delivery.OwnedTemp) error { return store.Close() }
		if _, err := coordinator.Acquire(context.Background(), validAcquireRequest(server.Bind)); err == nil {
			t.Fatal("closed cleanup store accepted")
		}
	})

	t.Run("concurrent delivery removal", func(t *testing.T) {
		store, server := seededCoordinatorStore(t, true)
		coordinator := staleCoordinator(store)
		calls := 0
		coordinator.Cleanup = func(context.Context, delivery.OwnedTemp) error {
			calls++
			if calls == 2 {
				current, err := store.Load(context.Background())
				if err != nil {
					return err
				}
				_, err = store.Update(context.Background(), current.Revision, func(registry *delivery.Registry) error {
					return registry.TombstoneDelivery(workerDeliveryID, delivery.Tombstone{ID: delivery.NewID(), TargetID: workerDeliveryID, Kind: delivery.TargetDelivery, Reason: delivery.ReasonStale, At: workerTime.Add(time.Minute)})
				})
				return err
			}
			return nil
		}
		if _, err := coordinator.Acquire(context.Background(), validAcquireRequest(server.Bind)); !errors.Is(err, delivery.ErrNotFound) {
			t.Fatalf("concurrent removal=%v", err)
		}
	})

	store, server := seededCoordinatorStore(t, false)
	coordinator := staleCoordinator(store)
	coordinator.StaleProbe = func(error) bool { return false }
	if _, err := coordinator.Acquire(context.Background(), validAcquireRequest(server.Bind)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("custom stale decision=%v", err)
	}
	_ = store.Close()
	if _, err := coordinator.Acquire(context.Background(), validAcquireRequest(server.Bind)); err == nil {
		t.Fatal("closed store accepted")
	}

	freshStore, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
	if err != nil {
		t.Fatal(err)
	}
	defer freshStore.Close()
	request := validAcquireRequest("127.0.0.1:8088")
	base := &Coordinator{Store: freshStore, StateDirectory: freshStore.Directory(), Locks: &fakeBindLocker{}, Hello: func(context.Context, Client) error { return nil }}
	base.Launch = func(context.Context, LaunchRequest) (Client, error) { return Client{}, want }
	if _, err := base.Acquire(context.Background(), request); !errors.Is(err, want) {
		t.Fatalf("launch=%v", err)
	}
	base.Launch = func(_ context.Context, launch LaunchRequest) (Client, error) {
		return Client{Endpoint: "wrong", ServerID: launch.ServerID, Compatibility: launch.Compatibility}, nil
	}
	if _, err := base.Acquire(context.Background(), request); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("launch mismatch=%v", err)
	}

	originalCall := callIPC
	t.Cleanup(func() { callIPC = originalCall })
	shutdown := false
	callIPC = func(_ context.Context, _ string, request ipc.Request) (ipc.Response, error) {
		if request.Operation == ipc.OperationShutdown {
			shutdown = true
			return ipc.Success(request, nil)
		}
		return ipc.Response{}, want
	}
	base.Hello = nil
	base.Launch = func(_ context.Context, launch LaunchRequest) (Client, error) {
		return Client{Endpoint: launch.ControlEndpoint, ServerID: launch.ServerID, Compatibility: launch.Compatibility}, nil
	}
	if _, err := base.Acquire(context.Background(), request); !errors.Is(err, want) || !shutdown {
		t.Fatalf("hello after launch=%v shutdown=%v", err, shutdown)
	}
	shutdown = false
	base.Hello = func(context.Context, Client) error { return nil }
	base.Register = func(context.Context, Client, delivery.Delivery, bool) (*Lease, error) { return nil, want }
	if _, err := base.Acquire(context.Background(), request); !errors.Is(err, want) || !shutdown {
		t.Fatalf("register after launch=%v shutdown=%v", err, shutdown)
	}
}

func TestCoordinatorDefaultRegisterAndEndpointCleanup(t *testing.T) {
	originalOpen := openIPC
	t.Cleanup(func() { openIPC = originalOpen })
	store, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	openIPC = func(context.Context, string) (net.Conn, func() bool, error) {
		connection := newRPCConn(func(request ipc.Request) ipc.Response {
			var payload RegisterRequest
			_ = ipc.DecodePayload(request.Payload, &payload)
			response, _ := ipc.Success(request, RegisterResponse{ServerID: payload.Delivery.ServerID, DeliveryID: payload.Delivery.ID, LeaseID: payload.LeaseID})
			return response
		})
		return connection, func() bool { return true }, nil
	}
	coordinator := &Coordinator{
		Store: store, StateDirectory: store.Directory(), Locks: &fakeBindLocker{}, Hello: func(context.Context, Client) error { return nil },
		Launch: func(_ context.Context, request LaunchRequest) (Client, error) {
			return Client{Endpoint: request.ControlEndpoint, ServerID: request.ServerID, Compatibility: request.Compatibility}, nil
		},
	}
	result, err := coordinator.Acquire(context.Background(), validAcquireRequest("127.0.0.1:8089"))
	if err != nil || result.Lease == nil {
		t.Fatalf("default registration=%+v err=%v", result, err)
	}
	_ = result.Lease.Close()

	staleStore, staleServer := seededCoordinatorStore(t, false)
	current, err := staleStore.Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	other := staleServer
	other.ID = delivery.ID("00000000-0000-4000-8000-000000000299")
	other.Bind = "127.0.0.1:9099"
	other.ControlEndpoint = "other-control"
	if _, err := staleStore.Update(context.Background(), current.Revision, func(registry *delivery.Registry) error {
		if err := registry.RegisterServer(other); err != nil {
			return err
		}
		return registry.AddOwnedTemp(delivery.OwnedTemp{ID: delivery.ID("00000000-0000-4000-8000-000000000298"), OwnerID: other.ID, Location: delivery.TempLocal, Path: "other-temp", CreatedAt: workerTime})
	}); err != nil {
		t.Fatal(err)
	}
	coordinator = staleCoordinator(staleStore)
	coordinator.RemoveEndpoint = nil
	coordinator.Hello = func(_ context.Context, client Client) error {
		if client.ServerID == staleServer.ID {
			return os.ErrNotExist
		}
		return nil
	}
	if _, err := coordinator.Acquire(context.Background(), validAcquireRequest(staleServer.Bind)); err != nil {
		t.Fatalf("default stale endpoint removal=%v", err)
	}
}

func TestCoordinatorDefaultsAndStaleClassification(t *testing.T) {
	store, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	coordinator := DefaultCoordinator(store, store.Directory())
	if coordinator.Store != store || coordinator.Launch == nil || coordinator.Locks == nil || coordinator.Cleanup == nil {
		t.Fatalf("default=%+v", coordinator)
	}
	if !staleProbeError(os.ErrNotExist) || !staleProbeError(syscall.ECONNREFUSED) || !staleProbeError(net.ErrClosed) || staleProbeError(context.DeadlineExceeded) {
		t.Fatal("stale error classification")
	}

	server := delivery.Server{ID: workerServerID, Bind: "127.0.0.1:8080", ControlEndpoint: filepath.Join(store.Directory(), "missing.sock"), ProcessID: 1, State: delivery.StateStarting, StartedAt: workerTime, UpdatedAt: workerTime}
	runtime, err := NewRuntime(context.Background(), store, server, "http-v1")
	if err != nil {
		t.Fatal(err)
	}
	runtime.now = func() time.Time { return workerTime.Add(time.Minute) }
	if err := runtime.StopServer(context.Background(), delivery.ReasonStopped); err != nil {
		t.Fatal(err)
	}
}

func seededCoordinatorStore(t *testing.T, withTemps bool) (*delivery.Store, delivery.Server) {
	t.Helper()
	store, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	server := delivery.Server{ID: workerServerID, Bind: "127.0.0.1:8080", ControlEndpoint: filepath.Join(store.Directory(), "stale.sock"), ProcessID: 99, State: delivery.StateActive, StartedAt: workerTime, UpdatedAt: workerTime}
	item := validWorkerDelivery()
	item.State = delivery.StateActive
	_, err = store.Update(context.Background(), 0, func(registry *delivery.Registry) error {
		if err := registry.RegisterServer(server); err != nil {
			return err
		}
		if withTemps {
			if err := registry.RegisterDelivery(item); err != nil {
				return err
			}
			for _, temporary := range []delivery.OwnedTemp{
				{ID: workerLeaseID, OwnerID: server.ID, Location: delivery.TempLocal, Path: "server-temp", CreatedAt: workerTime},
				{ID: delivery.ID("00000000-0000-4000-8000-000000000204"), OwnerID: item.ID, Location: delivery.TempRemote, Path: "delivery-temp", CreatedAt: workerTime},
			} {
				if err := registry.AddOwnedTemp(temporary); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return store, server
}

func validAcquireRequest(bind string) AcquireRequest {
	return AcquireRequest{Bind: bind, Compatibility: "http-v1", Route: delivery.RoutePathToWeb, Policy: delivery.DefaultPolicy(), Foreground: true, At: workerTime}
}

func staleCoordinator(store *delivery.Store) *Coordinator {
	return &Coordinator{
		Store: store, StateDirectory: store.Directory(), Locks: &fakeBindLocker{},
		Hello:          func(context.Context, Client) error { return os.ErrNotExist },
		RemoveEndpoint: func(string) error { return nil },
		Cleanup:        func(context.Context, delivery.OwnedTemp) error { return nil },
		Launch: func(_ context.Context, request LaunchRequest) (Client, error) {
			return Client{Endpoint: request.ControlEndpoint, ServerID: request.ServerID, Compatibility: request.Compatibility}, nil
		},
		Register: func(context.Context, Client, delivery.Delivery, bool) (*Lease, error) { return nil, nil },
	}
}

type fakeLock struct {
	locked   bool
	err      error
	closed   bool
	closeErr error
}

func (lock *fakeLock) TryLockContext(context.Context, time.Duration) (bool, error) {
	return lock.locked, lock.err
}
func (lock *fakeLock) Close() error { lock.closed = true; return lock.closeErr }

type fakeBindLocker struct {
	err        error
	releaseErr error
	released   bool
}

func (locker *fakeBindLocker) Acquire(context.Context, string) (func() error, error) {
	if locker.err != nil {
		return nil, locker.err
	}
	return func() error { locker.released = true; return locker.releaseErr }, nil
}

package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
)

func newRuntimeForTest(t *testing.T) (*Runtime, *delivery.Store) {
	t.Helper()
	store, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	server := delivery.Server{
		ID: workerServerID, Bind: "127.0.0.1:8080", ControlEndpoint: "control", ProcessID: os.Getpid(),
		State: delivery.StateStarting, StartedAt: workerTime, UpdatedAt: workerTime,
	}
	runtime, err := NewRuntime(context.Background(), store, server, "http-v1")
	if err != nil {
		t.Fatal(err)
	}
	runtime.now = func() time.Time { return workerTime.Add(time.Minute) }
	return runtime, store
}

func TestRuntimeLifecycleDirect(t *testing.T) {
	if _, err := NewRuntime(context.Background(), nil, delivery.Server{}, "bad"); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("nil store=%v", err)
	}
	store, err := delivery.OpenStore(filepath.Join(t.TempDir(), "state"))
	if err != nil {
		t.Fatal(err)
	}
	active := delivery.Server{ID: workerServerID, Bind: "127.0.0.1:8080", ControlEndpoint: "control", ProcessID: 1, State: delivery.StateActive, StartedAt: workerTime, UpdatedAt: workerTime}
	if _, err := NewRuntime(context.Background(), store, active, "http-v1"); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("active=%v", err)
	}
	starting := active
	starting.State = delivery.StateStarting
	invalidStarting := starting
	invalidStarting.ID = "bad"
	if _, err := NewRuntime(context.Background(), store, invalidStarting, "http-v1"); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid starting server=%v", err)
	}
	if _, err := NewRuntime(context.Background(), store, starting, " "); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("compatibility=%v", err)
	}
	_ = store.Close()
	if _, err := NewRuntime(context.Background(), store, starting, "http-v1"); err == nil {
		t.Fatal("closed store accepted")
	}

	runtime, runtimeStore := newRuntimeForTest(t)
	if runtime.ServerID() != workerServerID || runtime.HasDeliveries() {
		t.Fatalf("server=%s deliveries=%v", runtime.ServerID(), runtime.HasDeliveries())
	}
	if _, err := runtime.Snapshot(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.Register(context.Background(), RegisterRequest{}, nil); err == nil {
		t.Fatal("invalid registration accepted")
	}
	wrong := validWorkerDelivery()
	wrong.ServerID = delivery.ID("00000000-0000-4000-8000-000000000299")
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: wrong}, nil); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("wrong server=%v", err)
	}
	wrong = validWorkerDelivery()
	wrong.State = delivery.StateActive
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: wrong}, nil); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("wrong state=%v", err)
	}
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: validWorkerDelivery(), LeaseID: workerLeaseID}, nil); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("missing connection=%v", err)
	}

	first, err := runtime.Register(context.Background(), RegisterRequest{Delivery: validWorkerDelivery()}, nil)
	if err != nil || first.DeliveryID != workerDeliveryID || !runtime.HasDeliveries() {
		t.Fatalf("registration=%+v err=%v", first, err)
	}
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: validWorkerDelivery()}, nil); !errors.Is(err, delivery.ErrDuplicate) {
		t.Fatalf("duplicate=%v", err)
	}
	secondItem := validWorkerDelivery()
	secondItem.ID = delivery.ID("00000000-0000-4000-8000-000000000204")
	connectionA, connectionB := net.Pipe()
	t.Cleanup(func() { _ = connectionA.Close(); _ = connectionB.Close() })
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: secondItem, LeaseID: workerLeaseID}, connectionA); err != nil {
		t.Fatal(err)
	}
	thirdItem := validWorkerDelivery()
	thirdItem.ID = delivery.ID("00000000-0000-4000-8000-000000000205")
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: thirdItem, LeaseID: workerLeaseID}, connectionB); !errors.Is(err, delivery.ErrDuplicate) {
		t.Fatalf("duplicate lease=%v", err)
	}

	if err := runtime.ClaimLease(LeaseRequest{}, connectionA); err == nil {
		t.Fatal("invalid claim accepted")
	}
	if err := runtime.ClaimLease(LeaseRequest{DeliveryID: workerDeliveryID, LeaseID: delivery.ID("00000000-0000-4000-8000-000000000206")}, nil); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("nil claim=%v", err)
	}
	if err := runtime.ClaimLease(LeaseRequest{DeliveryID: thirdItem.ID, LeaseID: delivery.ID("00000000-0000-4000-8000-000000000206")}, connectionB); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("missing delivery=%v", err)
	}
	if err := runtime.ClaimLease(LeaseRequest{DeliveryID: workerDeliveryID, LeaseID: workerLeaseID}, connectionB); !errors.Is(err, delivery.ErrDuplicate) {
		t.Fatalf("duplicate claim=%v", err)
	}
	claimID := delivery.ID("00000000-0000-4000-8000-000000000206")
	if err := runtime.ClaimLease(LeaseRequest{DeliveryID: workerDeliveryID, LeaseID: claimID}, connectionB); err != nil {
		t.Fatal(err)
	}
	if err := runtime.ClaimLease(LeaseRequest{DeliveryID: workerDeliveryID, LeaseID: delivery.ID("00000000-0000-4000-8000-000000000207")}, connectionB); !errors.Is(err, delivery.ErrDuplicate) {
		t.Fatalf("delivery claim=%v", err)
	}
	if err := runtime.ReleaseLease(context.Background(), ReleaseLeaseRequest{}); err == nil {
		t.Fatal("invalid release accepted")
	}
	if err := runtime.ReleaseLease(context.Background(), ReleaseLeaseRequest{DeliveryID: secondItem.ID, LeaseID: claimID}); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("mismatched release=%v", err)
	}
	if err := runtime.ReleaseLease(context.Background(), ReleaseLeaseRequest{DeliveryID: workerDeliveryID, LeaseID: claimID}); err != nil {
		t.Fatal(err)
	}
	if _, exists := runtime.deliveries[workerDeliveryID]; exists {
		t.Fatal("released delivery remains")
	}

	policy := delivery.DefaultPolicy()
	policy.Version = 2
	if err := runtime.UpdatePolicy(context.Background(), UpdatePolicyRequest{}); err == nil {
		t.Fatal("invalid policy update accepted")
	}
	if err := runtime.UpdatePolicy(context.Background(), UpdatePolicyRequest{DeliveryID: workerDeliveryID, ExpectedVersion: 1, Policy: policy}); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("missing policy=%v", err)
	}
	if err := runtime.SetCounters(context.Background(), workerDeliveryID, delivery.CounterSnapshot{}); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("missing counters=%v", err)
	}
	if stopped, err := runtime.StopDelivery(context.Background(), workerDeliveryID, delivery.ReasonStopped); err != nil || stopped {
		t.Fatalf("idempotent stop=%v err=%v", stopped, err)
	}

	if _, _, err := runtime.Subscribe(workerDeliveryID); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("missing subscription=%v", err)
	}
	all, cancel, err := runtime.Subscribe("")
	if err != nil {
		t.Fatal(err)
	}
	if event := <-all; event.Snapshot.SchemaVersion != delivery.SchemaVersion {
		t.Fatalf("event=%+v", event)
	}
	runtime.publishLockedForTest()
	if event := <-all; event.Snapshot.Revision == 0 {
		t.Fatalf("published=%+v", event)
	}
	cancel()
	cancel()

	keepalive, err := runtime.AcquireKeepalive()
	if err != nil {
		t.Fatal(err)
	}
	if err := runtime.ReleaseKeepalive(context.Background(), delivery.ID("00000000-0000-4000-8000-000000000299")); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("missing keepalive=%v", err)
	}
	if err := runtime.ReleaseKeepalive(context.Background(), keepalive); err != nil {
		t.Fatal(err)
	}

	if err := runtimeStore.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: thirdItem}, nil); err == nil {
		t.Fatal("closed store registration succeeded")
	}
}

func TestRuntimeFailureAndKeepalivePaths(t *testing.T) {
	t.Run("registration conflicts with durable state", func(t *testing.T) {
		runtime, store := newRuntimeForTest(t)
		item := validWorkerDelivery()
		current, err := store.Load(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if _, err := store.Update(context.Background(), current.Revision, func(registry *delivery.Registry) error { return registry.RegisterDelivery(item) }); err != nil {
			t.Fatal(err)
		}
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item}, nil); !errors.Is(err, delivery.ErrDuplicate) {
			t.Fatalf("durable duplicate=%v", err)
		}
	})

	t.Run("closed store mutations", func(t *testing.T) {
		runtime, store := newRuntimeForTest(t)
		first := validWorkerDelivery()
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: first}, nil); err != nil {
			t.Fatal(err)
		}
		if err := store.Close(); err != nil {
			t.Fatal(err)
		}
		policy := delivery.DefaultPolicy()
		policy.Version = 2
		if err := runtime.UpdatePolicy(context.Background(), UpdatePolicyRequest{DeliveryID: first.ID, ExpectedVersion: 1, Policy: policy}); err == nil {
			t.Fatal("policy update with closed store succeeded")
		}
		if err := runtime.SetCounters(context.Background(), first.ID, delivery.CounterSnapshot{Read: 1}); err == nil {
			t.Fatal("counter update with closed store succeeded")
		}
		if _, err := runtime.StopDelivery(context.Background(), first.ID, delivery.ReasonStopped); err == nil {
			t.Fatal("delivery stop with closed store succeeded")
		}
		if err := runtime.StopServer(context.Background(), delivery.ReasonStopped); err == nil {
			t.Fatal("server stop with closed store succeeded")
		}
		runtime.publishLockedForTest()
	})

	t.Run("invalid terminal mutations", func(t *testing.T) {
		runtime, _ := newRuntimeForTest(t)
		item := validWorkerDelivery()
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item}, nil); err != nil {
			t.Fatal(err)
		}
		if _, err := runtime.StopDelivery(context.Background(), item.ID, "bad"); err == nil {
			t.Fatal("invalid delivery reason accepted")
		}
		runtime.deliveries[item.ID] = delivery.Delivery{ID: item.ID, State: "bad"}
		if _, err := runtime.StopDelivery(context.Background(), item.ID, delivery.ReasonStopped); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("invalid delivery state=%v", err)
		}
		runtime.deliveries[item.ID] = item
		runtime.server.State = "bad"
		if _, err := runtime.StopDelivery(context.Background(), item.ID, delivery.ReasonStopped); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("invalid server state=%v", err)
		}
	})

	t.Run("stop server with deliveries and leases", func(t *testing.T) {
		runtime, _ := newRuntimeForTest(t)
		first := validWorkerDelivery()
		second := validWorkerDelivery()
		second.ID = delivery.ID("00000000-0000-4000-8000-000000000204")
		firstConnection := &workerConn{}
		secondConnection := &workerConn{}
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: first, LeaseID: workerLeaseID}, firstConnection); err != nil {
			t.Fatal(err)
		}
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: second, LeaseID: delivery.ID("00000000-0000-4000-8000-000000000205")}, secondConnection); err != nil {
			t.Fatal(err)
		}
		if err := runtime.StopServer(context.Background(), "bad"); err == nil {
			t.Fatal("invalid server stop reason accepted")
		}
		if err := runtime.StopServer(context.Background(), delivery.ReasonStopped); err != nil {
			t.Fatal(err)
		}
		if !firstConnection.closed || !secondConnection.closed {
			t.Fatal("lease connections remain open")
		}
		if err := runtime.StopServer(context.Background(), delivery.ReasonStopped); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("stop server transition failures", func(t *testing.T) {
		runtime, _ := newRuntimeForTest(t)
		item := validWorkerDelivery()
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item}, nil); err != nil {
			t.Fatal(err)
		}
		runtime.deliveries[item.ID] = delivery.Delivery{ID: item.ID, State: "bad"}
		if err := runtime.StopServer(context.Background(), delivery.ReasonStopped); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("delivery finish=%v", err)
		}
		runtime.deliveries[item.ID] = item
		runtime.server.State = "bad"
		if err := runtime.StopServer(context.Background(), delivery.ReasonStopped); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("server finish=%v", err)
		}
	})

	t.Run("keepalive defers final server stop", func(t *testing.T) {
		runtime, _ := newRuntimeForTest(t)
		item := validWorkerDelivery()
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item}, nil); err != nil {
			t.Fatal(err)
		}
		keepalive, err := runtime.AcquireKeepalive()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := runtime.StopDelivery(context.Background(), item.ID, delivery.ReasonCompleted); err != nil {
			t.Fatal(err)
		}
		if err := runtime.ReleaseKeepalive(context.Background(), keepalive); err != nil {
			t.Fatal(err)
		}
		select {
		case <-runtime.Done():
		default:
			t.Fatal("last keepalive did not stop server")
		}
	})

	t.Run("pre-stopped run and shutdown checks", func(t *testing.T) {
		runtime, _ := newRuntimeForTest(t)
		item := validWorkerDelivery()
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item}, nil); err != nil {
			t.Fatal(err)
		}
		runtime.mutex.Lock()
		close(runtime.shutdown)
		runtime.mutex.Unlock()
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: delivery.Delivery{
			ID: delivery.ID("00000000-0000-4000-8000-000000000204"), ServerID: workerServerID, Route: delivery.RoutePathToWeb,
			State: delivery.StateStarting, Policy: delivery.DefaultPolicy(), CreatedAt: workerTime, UpdatedAt: workerTime,
		}}, nil); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("shutdown registration=%v", err)
		}
		if err := runtime.ClaimLease(LeaseRequest{DeliveryID: item.ID, LeaseID: workerLeaseID}, &workerConn{}); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("shutdown claim=%v", err)
		}
		listener := &scriptedListener{acceptErr: errors.New("unused")}
		if err := runtime.Run(listener); err != nil || !listener.closed {
			t.Fatalf("pre-stopped run=%v closed=%v", err, listener.closed)
		}
	})
}

func TestRuntimeControlHandlers(t *testing.T) {
	runtime, _ := newRuntimeForTest(t)
	first := validWorkerDelivery()
	second := validWorkerDelivery()
	second.ID = delivery.ID("00000000-0000-4000-8000-000000000204")
	keepalive, _ := runtime.AcquireKeepalive()
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: first, LeaseID: workerLeaseID}, &workerConn{}); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: second}, nil); err != nil {
		t.Fatal(err)
	}

	if response := invokeHandler(t, runtime, ipc.OperationStopDelivery, TargetRequest{ID: first.ID}); !response.OK {
		t.Fatalf("stop response=%+v", response)
	}
	if _, exists := runtime.deliveries[first.ID]; exists {
		t.Fatal("handler did not stop delivery")
	}
	if response := invokeHandler(t, runtime, ipc.OperationStopDelivery, TargetRequest{ID: "bad"}); response.OK || response.Error.Code != ipc.CodeInvalid {
		t.Fatalf("invalid target response=%+v", response)
	}
	if response := invokeHandler(t, runtime, ipc.OperationStopDelivery, map[string]string{"bad": "payload"}); response.OK {
		t.Fatalf("bad payload response=%+v", response)
	}
	if response := invokeHandler(t, runtime, ipc.OperationReleaseLease, nil); response.OK || response.Error.Code != ipc.CodeInvalid {
		t.Fatalf("release response=%+v", response)
	}
	if response := invokeHandler(t, runtime, ipc.OperationSubscribeProgress, ProgressRequest{DeliveryID: delivery.ID("00000000-0000-4000-8000-000000000299")}); response.OK || response.Error.Code != ipc.CodeNotFound {
		t.Fatalf("subscription response=%+v", response)
	}

	serverSide, clientSide := net.Pipe()
	done := make(chan struct{})
	go func() { runtime.handle(serverSide); close(done) }()
	request, _ := ipc.NewRequest(ipc.OperationSubscribeProgress, ProgressRequest{})
	if err := ipc.WriteFrame(clientSide, request); err != nil {
		t.Fatal(err)
	}
	var streamed ipc.Response
	if err := ipc.ReadFrame(clientSide, &streamed); err != nil || !streamed.OK {
		t.Fatalf("streamed=%+v err=%v", streamed, err)
	}
	_ = clientSide.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("closed subscriber was not released")
	}

	if response := invokeHandler(t, runtime, ipc.OperationShutdown, nil); !response.OK {
		t.Fatalf("shutdown response=%+v", response)
	}
	if err := runtime.ReleaseKeepalive(context.Background(), keepalive); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("cleared keepalive=%v", err)
	}

	runtime2, _ := newRuntimeForTest(t)
	runtime2.handle(&workerConn{})
	invalidRequest := ipc.Request{Version: 2, ID: workerLeaseID, Operation: ipc.OperationHello}
	runtime2.handle(requestInputConn(t, invalidRequest, nil))
	if response := invokeHandler(t, runtime2, ipc.OperationStopServer, nil); !response.OK {
		t.Fatalf("stop server response=%+v", response)
	}
}

func TestRuntimeHandlerLeaseAndWriteFailures(t *testing.T) {
	t.Run("register response failure releases lease", func(t *testing.T) {
		runtime, _ := newRuntimeForTest(t)
		item := validWorkerDelivery()
		request, _ := ipc.NewRequest(ipc.OperationRegister, RegisterRequest{Delivery: item, LeaseID: workerLeaseID})
		connection := requestInputConn(t, request, errors.New("write"))
		runtime.handle(connection)
		if runtime.HasDeliveries() {
			t.Fatal("failed response leaked foreground delivery")
		}
	})

	t.Run("lease payload failures and mismatch", func(t *testing.T) {
		runtime, _ := newRuntimeForTest(t)
		item := validWorkerDelivery()
		if _, err := runtime.Register(context.Background(), RegisterRequest{Delivery: item}, nil); err != nil {
			t.Fatal(err)
		}
		if response := invokeHandler(t, runtime, ipc.OperationLease, map[string]string{"bad": "payload"}); response.OK {
			t.Fatalf("bad lease response=%+v", response)
		}
		if response := invokeHandler(t, runtime, ipc.OperationLease, LeaseRequest{DeliveryID: delivery.ID("00000000-0000-4000-8000-000000000299"), LeaseID: workerLeaseID}); response.OK {
			t.Fatalf("missing lease response=%+v", response)
		}
		serverSide, clientSide := net.Pipe()
		if err := runtime.ClaimLease(LeaseRequest{DeliveryID: item.ID, LeaseID: workerLeaseID}, serverSide); err != nil {
			t.Fatal(err)
		}
		done := make(chan struct{})
		go func() { runtime.holdLease(serverSide, item.ID, workerLeaseID); close(done) }()
		request, _ := ipc.NewRequest(ipc.OperationReleaseLease, ReleaseLeaseRequest{DeliveryID: item.ID, LeaseID: delivery.ID("00000000-0000-4000-8000-000000000299")})
		if err := ipc.WriteFrame(clientSide, request); err != nil {
			t.Fatal(err)
		}
		var response ipc.Response
		if err := ipc.ReadFrame(clientSide, &response); err != nil || response.OK {
			t.Fatalf("mismatch response=%+v err=%v", response, err)
		}
		_ = clientSide.Close()
		<-done
		_ = runtime.releaseLostLease(context.Background(), item.ID, workerLeaseID)
		_ = runtime.releaseLostLease(context.Background(), item.ID, delivery.ID("00000000-0000-4000-8000-000000000299"))
	})

	t.Run("subscription write failure", func(t *testing.T) {
		runtime, _ := newRuntimeForTest(t)
		request, _ := ipc.NewRequest(ipc.OperationSubscribeProgress, ProgressRequest{})
		connection := requestInputConn(t, request, errors.New("write"))
		runtime.handle(connection)
	})

	t.Run("respond build and write failures", func(t *testing.T) {
		runtime, _ := newRuntimeForTest(t)
		request, _ := ipc.NewRequest(ipc.OperationHello, nil)
		if runtime.respond(&workerConn{}, request, make(chan int), nil) {
			t.Fatal("unencodable response succeeded")
		}
		if runtime.respond(&failingWriteConn{}, request, nil, nil) {
			t.Fatal("failed write succeeded")
		}
	})
}

func TestRuntimePublishReplacementAndError(t *testing.T) {
	runtime, store := newRuntimeForTest(t)
	channel := make(chan ProgressEvent, 1)
	runtime.publishOneLocked(channel, "")
	runtime.publishOneLocked(channel, "")
	if len(channel) != 1 {
		t.Fatalf("channel length=%d", len(channel))
	}
	runtime.publishOneLocked(make(chan ProgressEvent), "")
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	runtime.publishOneLocked(channel, "")
	long := fmt.Errorf("%s: %w", string(make([]byte, 1100)), delivery.ErrInvalid)
	if len(publicError(long)) != 1024 {
		t.Fatalf("truncated length=%d", len(publicError(long)))
	}
}

func invokeHandler(t *testing.T, runtime *Runtime, operation ipc.Operation, payload any) ipc.Response {
	t.Helper()
	serverSide, clientSide := net.Pipe()
	done := make(chan struct{})
	go func() { runtime.handle(serverSide); close(done) }()
	request, err := ipc.NewRequest(operation, payload)
	if err != nil {
		t.Fatal(err)
	}
	if err := ipc.WriteFrame(clientSide, request); err != nil {
		t.Fatal(err)
	}
	var response ipc.Response
	if err := ipc.ReadFrame(clientSide, &response); err != nil {
		t.Fatal(err)
	}
	_ = clientSide.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handler did not exit")
	}
	return response
}

func requestInputConn(t *testing.T, request ipc.Request, writeErr error) *inputConn {
	t.Helper()
	var input bytes.Buffer
	if err := ipc.WriteFrame(&input, request); err != nil {
		t.Fatal(err)
	}
	return &inputConn{reader: bytes.NewReader(input.Bytes()), writeErr: writeErr}
}

type inputConn struct {
	reader   *bytes.Reader
	writeErr error
	closed   bool
}

func (connection *inputConn) Read(data []byte) (int, error) { return connection.reader.Read(data) }
func (connection *inputConn) Write(data []byte) (int, error) {
	if connection.writeErr != nil {
		return 0, connection.writeErr
	}
	return len(data), nil
}
func (connection *inputConn) Close() error          { connection.closed = true; return nil }
func (*inputConn) LocalAddr() net.Addr              { return workerAddr("local") }
func (*inputConn) RemoteAddr() net.Addr             { return workerAddr("remote") }
func (*inputConn) SetDeadline(time.Time) error      { return nil }
func (*inputConn) SetReadDeadline(time.Time) error  { return nil }
func (*inputConn) SetWriteDeadline(time.Time) error { return nil }

type failingWriteConn struct{ workerConn }

func (*failingWriteConn) Write([]byte) (int, error) { return 0, errors.New("write") }

func (runtime *Runtime) publishLockedForTest() {
	runtime.mutex.Lock()
	defer runtime.mutex.Unlock()
	runtime.publishLocked()
}

func TestRuntimeIPCIntegration(t *testing.T) {
	directory := shortWorkerTempDir(t, "courier-worker-")
	store, err := delivery.OpenStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	endpoint, err := ipc.ControlEndpoint(directory, workerServerID)
	if err != nil {
		t.Fatal(err)
	}
	server := delivery.Server{ID: workerServerID, Bind: "127.0.0.1:8081", ControlEndpoint: endpoint, ProcessID: os.Getpid(), State: delivery.StateStarting, StartedAt: workerTime, UpdatedAt: workerTime}
	runtime, err := NewRuntime(context.Background(), store, server, "http-v1")
	if err != nil {
		t.Fatal(err)
	}
	runtime.now = func() time.Time { return workerTime.Add(time.Minute) }
	listener, err := ipc.Listen(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	runResult := make(chan error, 1)
	go func() { runResult <- runtime.Run(listener) }()
	client := Client{Endpoint: endpoint, ServerID: workerServerID, Compatibility: "http-v1"}
	if hello, err := client.Hello(context.Background()); err != nil || hello.ServerID != workerServerID {
		t.Fatalf("hello=%+v err=%v", hello, err)
	}
	if _, err := (Client{Endpoint: endpoint, ServerID: workerServerID, Compatibility: "other"}).Hello(context.Background()); !ipc.IsRemoteError(err, ipc.CodeConflict) {
		t.Fatalf("compatibility=%v", err)
	}
	if _, err := (Client{Endpoint: endpoint, ServerID: delivery.ID("00000000-0000-4000-8000-000000000299"), Compatibility: "http-v1"}).Hello(context.Background()); !ipc.IsRemoteError(err, ipc.CodeNotFound) {
		t.Fatalf("identity=%v", err)
	}

	first := validWorkerDelivery()
	if _, lease, err := client.Register(context.Background(), first, false); err != nil || lease != nil {
		t.Fatalf("background lease=%v err=%v", lease, err)
	}
	second := validWorkerDelivery()
	second.ID = delivery.ID("00000000-0000-4000-8000-000000000204")
	_, foreground, err := client.Register(context.Background(), second, true)
	if err != nil || foreground == nil || !foreground.ID().Valid() {
		t.Fatalf("foreground=%v err=%v", foreground, err)
	}
	claimed, err := client.ClaimLease(context.Background(), first.ID)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := client.List(context.Background())
	if err != nil || len(snapshot.Deliveries) != 2 {
		t.Fatalf("snapshot=%+v err=%v", snapshot, err)
	}
	subscription, err := client.Subscribe(context.Background(), first.ID)
	if err != nil {
		t.Fatal(err)
	}
	defer subscription.Close()
	if event, err := subscription.Next(context.Background()); err != nil || len(event.Snapshot.Deliveries) != 1 {
		t.Fatalf("initial event=%+v err=%v", event, err)
	}
	if err := runtime.SetCounters(context.Background(), first.ID, delivery.CounterSnapshot{Read: 10, Sent: 8, Confirmed: 8}); err != nil {
		t.Fatal(err)
	}
	if event, err := subscription.Next(context.Background()); err != nil || event.Snapshot.Deliveries[0].Counters.Confirmed != 8 {
		t.Fatalf("counter event=%+v err=%v", event, err)
	}
	policy := delivery.DefaultPolicy()
	policy.Version = 2
	if err := client.UpdatePolicy(context.Background(), UpdatePolicyRequest{DeliveryID: first.ID, ExpectedVersion: 1, Policy: policy}); err != nil {
		t.Fatal(err)
	}
	if err := claimed.Release(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := claimed.Release(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := foreground.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-runResult:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("worker did not stop after lease loss")
	}
}

func TestRuntimeRunAndHelpers(t *testing.T) {
	runtime, _ := newRuntimeForTest(t)
	if err := runtime.Run(nil); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("nil listener=%v", err)
	}
	listener := &scriptedListener{acceptErr: errors.New("accept")}
	if err := runtime.Run(listener); !errors.Is(err, listener.acceptErr) {
		t.Fatalf("accept=%v", err)
	}
	if err := runtime.Run(listener); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("second run=%v", err)
	}
	if _, err := runtime.AcquireKeepalive(); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("shutdown keepalive=%v", err)
	}
	if _, _, err := runtime.Subscribe(""); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("shutdown subscribe=%v", err)
	}
	if err := runtime.ClaimLease(LeaseRequest{DeliveryID: workerDeliveryID, LeaseID: workerLeaseID}, &workerConn{}); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("shutdown claim order=%v", err)
	}

	for _, state := range []delivery.State{delivery.StateStarting, delivery.StateActive, delivery.StateStopping, delivery.StateFailed, delivery.StateStopped} {
		snapshot := delivery.EmptySnapshot()
		server := delivery.Server{ID: workerServerID, Bind: "127.0.0.1:8080", ControlEndpoint: "control", ProcessID: 1, State: delivery.StateStarting, StartedAt: workerTime, UpdatedAt: workerTime}
		item := validWorkerDelivery()
		item.State = state
		if state == delivery.StateStopped {
			item.State = delivery.StateStarting
		}
		snapshot.Servers = []delivery.Server{server}
		snapshot.Deliveries = []delivery.Delivery{item}
		registry, err := delivery.NewRegistry(snapshot)
		if err != nil {
			t.Fatal(err)
		}
		if state == delivery.StateStopped {
			registry, _ = delivery.NewRegistry(delivery.EmptySnapshot())
			if err := finishDelivery(registry, delivery.Delivery{State: delivery.StateStopped}, workerTime); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if err := finishDelivery(registry, item, workerTime.Add(time.Second)); err != nil {
			t.Fatalf("state=%s err=%v", state, err)
		}
	}
	if err := finishDelivery(nil, delivery.Delivery{State: "bad"}, workerTime); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("bad delivery state=%v", err)
	}
	emptyRegistry, _ := delivery.NewRegistry(delivery.EmptySnapshot())
	missingDelivery := validWorkerDelivery()
	missingDelivery.State = delivery.StateActive
	if err := finishDelivery(emptyRegistry, missingDelivery, workerTime); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("missing delivery transition=%v", err)
	}

	for _, state := range []delivery.State{delivery.StateStarting, delivery.StateActive, delivery.StateStopping, delivery.StateFailed} {
		server := delivery.Server{ID: workerServerID, Bind: "127.0.0.1:8080", ControlEndpoint: "control", ProcessID: 1, State: state, StartedAt: workerTime, UpdatedAt: workerTime}
		snapshot := delivery.EmptySnapshot()
		snapshot.Servers = []delivery.Server{server}
		registry, err := delivery.NewRegistry(snapshot)
		if err != nil {
			t.Fatal(err)
		}
		if err := finishServer(registry, server, workerTime.Add(time.Second)); err != nil {
			t.Fatalf("state=%s err=%v", state, err)
		}
	}
	if err := finishServer(nil, delivery.Server{State: delivery.StateStopped}, workerTime); err != nil {
		t.Fatal(err)
	}
	if err := finishServer(nil, delivery.Server{State: "bad"}, workerTime); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("bad server state=%v", err)
	}
	if err := finishServer(emptyRegistry, delivery.Server{ID: workerServerID, State: delivery.StateActive}, workerTime); !errors.Is(err, delivery.ErrNotFound) {
		t.Fatalf("missing server transition=%v", err)
	}

	for _, test := range []struct {
		err  error
		code ipc.ErrorCode
	}{
		{delivery.ErrNotFound, ipc.CodeNotFound}, {delivery.ErrDuplicate, ipc.CodeConflict}, {delivery.ErrRevisionConflict, ipc.CodeConflict},
		{ipc.ErrProtocol, ipc.CodeInvalid}, {delivery.ErrInvalid, ipc.CodeInvalid}, {errors.New("secret"), ipc.CodeInternal},
	} {
		if code := errorCode(test.err); code != test.code {
			t.Fatalf("error=%v code=%s", test.err, code)
		}
	}
	if publicError(nil) != "operation failed" || publicError(errors.New("secret")) != "internal worker error" || len(publicError(errors.New(string(make([]byte, 1100))))) > 1024 {
		t.Fatal("public error sanitization")
	}
	closeConnections([]net.Conn{nil, &workerConn{}})
}

func TestUpdateStoreRetriesConflict(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "state")
	first, err := delivery.OpenStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	second, err := delivery.OpenStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	entered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	done := make(chan error, 1)
	go func() {
		_, updateErr := updateStore(context.Background(), first, func(*delivery.Registry) error {
			once.Do(func() { close(entered); <-release })
			return nil
		})
		done <- updateErr
	}()
	<-entered
	if _, err := second.Update(context.Background(), 0, func(*delivery.Registry) error { return nil }); err != nil {
		close(release)
		t.Fatal(err)
	}
	close(release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	entered = make(chan struct{})
	release = make(chan struct{})
	once = sync.Once{}
	canceled := &delayedCancelContext{cancelAt: 5}
	done = make(chan error, 1)
	go func() {
		_, updateErr := updateStore(canceled, first, func(*delivery.Registry) error {
			once.Do(func() { close(entered); <-release })
			return nil
		})
		done <- updateErr
	}()
	<-entered
	current, err := second.Load(context.Background())
	if err != nil {
		close(release)
		t.Fatal(err)
	}
	if _, err := second.Update(context.Background(), current.Revision, func(*delivery.Registry) error { return nil }); err != nil {
		close(release)
		t.Fatal(err)
	}
	close(release)
	if err := <-done; !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled conflict=%v", err)
	}
}

type delayedCancelContext struct {
	calls    atomic.Int32
	cancelAt int32
}

func (*delayedCancelContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (*delayedCancelContext) Done() <-chan struct{}       { return nil }
func (delayed *delayedCancelContext) Err() error {
	if delayed.calls.Add(1) >= delayed.cancelAt {
		return context.Canceled
	}
	return nil
}
func (*delayedCancelContext) Value(any) any { return nil }

type scriptedListener struct {
	acceptErr error
	closed    bool
}

func (listener *scriptedListener) Accept() (net.Conn, error) { return nil, listener.acceptErr }
func (listener *scriptedListener) Close() error              { listener.closed = true; return nil }
func (*scriptedListener) Addr() net.Addr                     { return workerAddr("listener") }

type workerConn struct{ closed bool }

func (*workerConn) Read([]byte) (int, error)         { return 0, io.EOF }
func (*workerConn) Write(data []byte) (int, error)   { return len(data), nil }
func (connection *workerConn) Close() error          { connection.closed = true; return nil }
func (*workerConn) LocalAddr() net.Addr              { return workerAddr("local") }
func (*workerConn) RemoteAddr() net.Addr             { return workerAddr("remote") }
func (*workerConn) SetDeadline(time.Time) error      { return nil }
func (*workerConn) SetReadDeadline(time.Time) error  { return nil }
func (*workerConn) SetWriteDeadline(time.Time) error { return nil }

type workerAddr string

func (workerAddr) Network() string        { return "worker" }
func (address workerAddr) String() string { return string(address) }

func shortWorkerTempDir(t *testing.T, pattern string) string {
	t.Helper()
	parent := "/tmp"
	if runtime.GOOS == "windows" {
		parent = ""
	}
	directory, err := os.MkdirTemp(parent, pattern)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(directory) })
	return directory
}

package worker

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
)

func TestClientValidationAndHello(t *testing.T) {
	client := Client{Endpoint: "control", ServerID: workerServerID, Compatibility: "http-v1"}
	if err := client.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []Client{{ServerID: workerServerID, Compatibility: "http-v1"}, {Endpoint: "x", ServerID: "bad", Compatibility: "http-v1"}, {Endpoint: "x", ServerID: workerServerID}} {
		if err := invalid.Validate(); !errors.Is(err, ipc.ErrProtocol) {
			t.Fatalf("client=%+v err=%v", invalid, err)
		}
		if _, err := invalid.Hello(context.Background()); !errors.Is(err, ipc.ErrProtocol) {
			t.Fatalf("hello=%v", err)
		}
	}

	originalCall, originalRequest := callIPC, newIPCRequest
	t.Cleanup(func() { callIPC, newIPCRequest = originalCall, originalRequest })
	want := errors.New("failure")
	newIPCRequest = func(ipc.Operation, any) (ipc.Request, error) { return ipc.Request{}, want }
	if _, err := client.Hello(context.Background()); !errors.Is(err, want) {
		t.Fatalf("request=%v", err)
	}
	newIPCRequest = originalRequest
	callIPC = func(context.Context, string, ipc.Request) (ipc.Response, error) { return ipc.Response{}, want }
	if _, err := client.Hello(context.Background()); !errors.Is(err, want) {
		t.Fatalf("call=%v", err)
	}
	callIPC = func(_ context.Context, _ string, request ipc.Request) (ipc.Response, error) {
		return ipc.Success(request, map[string]string{"bad": "payload"})
	}
	if _, err := client.Hello(context.Background()); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("payload=%v", err)
	}
	callIPC = func(_ context.Context, _ string, request ipc.Request) (ipc.Response, error) {
		return ipc.Success(request, HelloResponse{ServerID: workerServerID, Compatibility: "other"})
	}
	if _, err := client.Hello(context.Background()); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("mismatch=%v", err)
	}
	callIPC = func(_ context.Context, _ string, request ipc.Request) (ipc.Response, error) {
		return ipc.Success(request, HelloResponse{ServerID: workerServerID, Compatibility: "http-v1"})
	}
	if result, err := client.Hello(context.Background()); err != nil || result.ServerID != workerServerID {
		t.Fatalf("hello=%+v err=%v", result, err)
	}
}

func TestClientRegistrationAndClaimFailures(t *testing.T) {
	client := Client{Endpoint: "control", ServerID: workerServerID, Compatibility: "http-v1"}
	invalid := client
	invalid.Endpoint = ""
	if _, _, err := invalid.Register(context.Background(), validWorkerDelivery(), false); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("invalid registration=%v", err)
	}
	if _, err := invalid.ClaimLease(context.Background(), workerDeliveryID); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("invalid claim=%v", err)
	}

	originalOpen, originalRequest := openIPC, newIPCRequest
	t.Cleanup(func() { openIPC, newIPCRequest = originalOpen, originalRequest })
	want := errors.New("failure")
	newIPCRequest = func(ipc.Operation, any) (ipc.Request, error) { return ipc.Request{}, want }
	if _, _, err := client.Register(context.Background(), validWorkerDelivery(), false); !errors.Is(err, want) {
		t.Fatalf("registration request=%v", err)
	}
	if _, err := client.ClaimLease(context.Background(), workerDeliveryID); !errors.Is(err, want) {
		t.Fatalf("claim request=%v", err)
	}
	newIPCRequest = originalRequest
	openIPC = func(context.Context, string) (net.Conn, func() bool, error) { return nil, nil, want }
	if _, _, err := client.Register(context.Background(), validWorkerDelivery(), false); !errors.Is(err, want) {
		t.Fatalf("open registration=%v", err)
	}
	if _, err := client.ClaimLease(context.Background(), workerDeliveryID); !errors.Is(err, want) {
		t.Fatalf("open claim=%v", err)
	}

	openWith := func(connection net.Conn) {
		openIPC = func(context.Context, string) (net.Conn, func() bool, error) {
			return connection, func() bool { return true }, nil
		}
	}
	connection := &rpcConn{writeErr: want}
	openWith(connection)
	if _, _, err := client.Register(context.Background(), validWorkerDelivery(), false); !errors.Is(err, want) {
		t.Fatalf("registration exchange=%v", err)
	}
	connection = &rpcConn{writeErr: want}
	openWith(connection)
	if _, err := client.ClaimLease(context.Background(), workerDeliveryID); !errors.Is(err, want) {
		t.Fatalf("claim exchange=%v", err)
	}

	connection = newRPCConn(func(request ipc.Request) ipc.Response {
		response, _ := ipc.Success(request, map[string]string{"bad": "payload"})
		return response
	})
	openWith(connection)
	if _, _, err := client.Register(context.Background(), validWorkerDelivery(), false); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("registration payload=%v", err)
	}
	connection = newRPCConn(func(request ipc.Request) ipc.Response {
		response, _ := ipc.Success(request, RegisterResponse{ServerID: workerServerID, DeliveryID: delivery.ID("00000000-0000-4000-8000-000000000299")})
		return response
	})
	openWith(connection)
	if _, _, err := client.Register(context.Background(), validWorkerDelivery(), false); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("registration mismatch=%v", err)
	}
	connection = newRPCConn(func(request ipc.Request) ipc.Response {
		var payload RegisterRequest
		_ = ipc.DecodePayload(request.Payload, &payload)
		response, _ := ipc.Success(request, RegisterResponse{ServerID: workerServerID, DeliveryID: payload.Delivery.ID, LeaseID: payload.LeaseID})
		return response
	})
	connection.deadlineErrAt = 1
	openWith(connection)
	if _, _, err := client.Register(context.Background(), validWorkerDelivery(), true); !errors.Is(err, errDeadline) {
		t.Fatalf("registration reset=%v", err)
	}

	connection = newRPCConn(func(request ipc.Request) ipc.Response {
		response, _ := ipc.Success(request, nil)
		return response
	})
	connection.deadlineErrAt = 1
	openWith(connection)
	if _, err := client.ClaimLease(context.Background(), workerDeliveryID); !errors.Is(err, errDeadline) {
		t.Fatalf("claim reset=%v", err)
	}
	connection = newRPCConn(func(request ipc.Request) ipc.Response {
		response, _ := ipc.Success(request, nil)
		return response
	})
	openWith(connection)
	lease, err := client.ClaimLease(context.Background(), workerDeliveryID)
	if err != nil || lease.ID() == "" {
		t.Fatalf("lease=%v err=%v", lease, err)
	}
	_ = lease.Close()
}

func TestClientCallsListAndSubscription(t *testing.T) {
	client := Client{Endpoint: "control", ServerID: workerServerID, Compatibility: "http-v1"}
	originalCall, originalOpen := callIPC, openIPC
	t.Cleanup(func() { callIPC, openIPC = originalCall, originalOpen })
	want := errors.New("failure")
	callIPC = func(context.Context, string, ipc.Request) (ipc.Response, error) { return ipc.Response{}, want }
	if _, err := client.List(context.Background()); !errors.Is(err, want) {
		t.Fatalf("list call=%v", err)
	}
	callIPC = func(_ context.Context, _ string, request ipc.Request) (ipc.Response, error) {
		return ipc.Success(request, map[string]string{"bad": "payload"})
	}
	if _, err := client.List(context.Background()); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("list payload=%v", err)
	}
	callIPC = func(_ context.Context, _ string, request ipc.Request) (ipc.Response, error) {
		return ipc.Success(request, ListResponse{Snapshot: delivery.Snapshot{SchemaVersion: 9}})
	}
	if _, err := client.List(context.Background()); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("list invalid=%v", err)
	}
	callIPC = func(_ context.Context, _ string, request ipc.Request) (ipc.Response, error) {
		return ipc.Success(request, ListResponse{Snapshot: delivery.EmptySnapshot()})
	}
	if _, err := client.List(context.Background()); err != nil {
		t.Fatal(err)
	}

	operations := []ipc.Operation{}
	callIPC = func(_ context.Context, _ string, request ipc.Request) (ipc.Response, error) {
		operations = append(operations, request.Operation)
		return ipc.Success(request, nil)
	}
	policy := delivery.DefaultPolicy()
	policy.Version = 2
	if err := client.UpdatePolicy(context.Background(), UpdatePolicyRequest{DeliveryID: workerDeliveryID, ExpectedVersion: 1, Policy: policy}); err != nil {
		t.Fatal(err)
	}
	if err := client.StopDelivery(context.Background(), workerDeliveryID); err != nil {
		t.Fatal(err)
	}
	if err := client.StopServer(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := client.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(operations) != 4 || operations[1] != ipc.OperationStopDelivery || operations[2] != ipc.OperationStopServer || operations[3] != ipc.OperationShutdown {
		t.Fatalf("operations=%v", operations)
	}
	callIPC = func(context.Context, string, ipc.Request) (ipc.Response, error) { return ipc.Response{}, want }
	if err := client.StopServer(context.Background()); !errors.Is(err, want) {
		t.Fatalf("call error=%v", err)
	}

	invalid := client
	invalid.Endpoint = ""
	if _, err := invalid.List(context.Background()); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("invalid list=%v", err)
	}
	if err := invalid.StopServer(context.Background()); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("invalid call=%v", err)
	}
	if _, err := invalid.Subscribe(context.Background(), ""); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("invalid subscribe=%v", err)
	}

	openIPC = func(context.Context, string) (net.Conn, func() bool, error) { return nil, nil, want }
	if _, err := client.Subscribe(context.Background(), ""); !errors.Is(err, want) {
		t.Fatalf("subscribe open=%v", err)
	}
	connection := &rpcConn{writeErr: want}
	openIPC = func(context.Context, string) (net.Conn, func() bool, error) {
		return connection, func() bool { return true }, nil
	}
	if _, err := client.Subscribe(context.Background(), ""); !errors.Is(err, want) {
		t.Fatalf("subscribe write=%v", err)
	}
	connection = &rpcConn{deadlineErrAt: 1}
	openIPC = func(context.Context, string) (net.Conn, func() bool, error) {
		return connection, func() bool { return true }, nil
	}
	if _, err := client.Subscribe(context.Background(), ""); !errors.Is(err, errDeadline) {
		t.Fatalf("subscribe reset=%v", err)
	}
	connection = &rpcConn{}
	openIPC = func(context.Context, string) (net.Conn, func() bool, error) {
		return connection, func() bool { return true }, nil
	}
	subscription, err := client.Subscribe(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if err := subscription.Close(); err != nil {
		t.Fatal(err)
	}
	if err := subscription.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestLeaseSubscriptionAndExchangeFailures(t *testing.T) {
	originalConfigure, originalRequest := configureIPC, newIPCRequest
	t.Cleanup(func() { configureIPC, newIPCRequest = originalConfigure, originalRequest })
	want := errors.New("failure")
	configureIPC = func(context.Context, net.Conn) (func() bool, error) { return nil, want }
	lease := &Lease{connection: &rpcConn{}, deliveryID: workerDeliveryID, leaseID: workerLeaseID}
	if err := lease.Release(context.Background()); !errors.Is(err, want) {
		t.Fatalf("lease configure=%v", err)
	}
	if err := lease.Release(context.Background()); !errors.Is(err, want) {
		t.Fatalf("lease repeat=%v", err)
	}
	lease = &Lease{connection: &rpcConn{closeErr: want}}
	if err := lease.Close(); !errors.Is(err, want) {
		t.Fatalf("lease close=%v", err)
	}

	configureIPC = originalConfigure
	newIPCRequest = func(ipc.Operation, any) (ipc.Request, error) { return ipc.Request{}, want }
	lease = &Lease{connection: &rpcConn{}, deliveryID: workerDeliveryID, leaseID: workerLeaseID}
	if err := lease.Release(context.Background()); !errors.Is(err, want) {
		t.Fatalf("lease request=%v", err)
	}
	newIPCRequest = originalRequest
	connection := &rpcConn{writeErr: want}
	lease = &Lease{connection: connection, deliveryID: workerDeliveryID, leaseID: workerLeaseID}
	if err := lease.Release(context.Background()); !errors.Is(err, want) {
		t.Fatalf("lease exchange=%v", err)
	}
	connection = newRPCConn(func(request ipc.Request) ipc.Response {
		response, _ := ipc.Success(request, nil)
		return response
	})
	lease = &Lease{connection: connection, deliveryID: workerDeliveryID, leaseID: workerLeaseID}
	if err := lease.Release(context.Background()); err != nil {
		t.Fatal(err)
	}

	subscription := &Subscription{connection: &rpcConn{}, requestID: workerLeaseID}
	configureIPC = func(context.Context, net.Conn) (func() bool, error) { return nil, want }
	if _, err := subscription.Next(context.Background()); !errors.Is(err, want) {
		t.Fatalf("subscription configure=%v", err)
	}
	configureIPC = originalConfigure
	if _, err := subscription.Next(context.Background()); !errors.Is(err, io.EOF) {
		t.Fatalf("subscription read=%v", err)
	}

	request := ipc.Request{Version: ipc.ProtocolVersion, ID: workerLeaseID, Operation: ipc.OperationSubscribeProgress}
	responses := []ipc.Response{}
	invalid, _ := ipc.Success(request, ProgressEvent{Snapshot: delivery.EmptySnapshot()})
	invalid.Version = 2
	responses = append(responses, invalid)
	mismatch, _ := ipc.Success(request, ProgressEvent{Snapshot: delivery.EmptySnapshot()})
	mismatch.ID = workerServerID
	responses = append(responses, mismatch)
	remote, _ := ipc.Rejection(request, ipc.CodeConflict, "conflict")
	responses = append(responses, remote)
	badPayload, _ := ipc.Success(request, map[string]string{"bad": "payload"})
	responses = append(responses, badPayload)
	badSnapshot, _ := ipc.Success(request, ProgressEvent{Snapshot: delivery.Snapshot{SchemaVersion: 9}})
	responses = append(responses, badSnapshot)
	valid, _ := ipc.Success(request, ProgressEvent{Snapshot: delivery.EmptySnapshot()})
	responses = append(responses, valid)
	for index, response := range responses {
		connection := responseConn(t, response)
		subscription = &Subscription{connection: connection, requestID: workerLeaseID}
		event, err := subscription.Next(context.Background())
		if index < len(responses)-1 && err == nil {
			t.Fatalf("response %d accepted: %+v", index, event)
		}
		if index == len(responses)-1 && (err != nil || event.Snapshot.SchemaVersion != delivery.SchemaVersion) {
			t.Fatalf("valid event=%+v err=%v", event, err)
		}
	}

	request = ipc.Request{Version: ipc.ProtocolVersion, ID: workerLeaseID, Operation: ipc.OperationHello}
	if _, err := exchange(&rpcConn{writeErr: want}, request); !errors.Is(err, want) {
		t.Fatalf("exchange write=%v", err)
	}
	if _, err := exchange(&rpcConn{}, request); !errors.Is(err, io.EOF) {
		t.Fatalf("exchange read=%v", err)
	}
	connection = newRPCConn(func(request ipc.Request) ipc.Response {
		response, _ := ipc.Success(request, nil)
		response.Version = 2
		return response
	})
	if _, err := exchange(connection, request); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("exchange invalid=%v", err)
	}
	connection = newRPCConn(func(request ipc.Request) ipc.Response {
		response, _ := ipc.Success(request, nil)
		response.ID = workerServerID
		return response
	})
	if _, err := exchange(connection, request); !errors.Is(err, ipc.ErrProtocol) {
		t.Fatalf("exchange mismatch=%v", err)
	}
	connection = newRPCConn(func(request ipc.Request) ipc.Response {
		response, _ := ipc.Rejection(request, ipc.CodeConflict, "conflict")
		return response
	})
	if _, err := exchange(connection, request); !ipc.IsRemoteError(err, ipc.CodeConflict) {
		t.Fatalf("exchange remote=%v", err)
	}
}

var errDeadline = errors.New("deadline")

type rpcConn struct {
	requestBytes  bytes.Buffer
	responseBytes bytes.Buffer
	handler       func(ipc.Request) ipc.Response
	writtenSize   int
	writeErr      error
	readErr       error
	closeErr      error
	deadlineErrAt int
	deadlineCalls int
	closed        bool
}

func newRPCConn(handler func(ipc.Request) ipc.Response) *rpcConn { return &rpcConn{handler: handler} }

func (connection *rpcConn) Read(data []byte) (int, error) {
	if connection.readErr != nil {
		return 0, connection.readErr
	}
	return connection.responseBytes.Read(data)
}

func (connection *rpcConn) Write(data []byte) (int, error) {
	if connection.writeErr != nil {
		return 0, connection.writeErr
	}
	written, _ := connection.requestBytes.Write(data)
	if connection.writtenSize == 0 && connection.requestBytes.Len() >= 4 {
		connection.writtenSize = int(binary.BigEndian.Uint32(connection.requestBytes.Bytes()[:4])) + 4
	}
	if connection.handler != nil && connection.writtenSize != 0 && connection.requestBytes.Len() >= connection.writtenSize {
		var request ipc.Request
		_ = ipc.ReadFrame(bytes.NewReader(connection.requestBytes.Bytes()[:connection.writtenSize]), &request)
		response := connection.handler(request)
		_ = ipc.WriteFrame(&connection.responseBytes, response)
		connection.handler = nil
	}
	return written, nil
}

func (connection *rpcConn) Close() error { connection.closed = true; return connection.closeErr }
func (*rpcConn) LocalAddr() net.Addr     { return workerAddr("local") }
func (*rpcConn) RemoteAddr() net.Addr    { return workerAddr("remote") }
func (connection *rpcConn) SetDeadline(time.Time) error {
	connection.deadlineCalls++
	if connection.deadlineCalls == connection.deadlineErrAt {
		return errDeadline
	}
	return nil
}
func (*rpcConn) SetReadDeadline(time.Time) error  { return nil }
func (*rpcConn) SetWriteDeadline(time.Time) error { return nil }

func responseConn(t *testing.T, response ipc.Response) *rpcConn {
	t.Helper()
	connection := &rpcConn{}
	if err := ipc.WriteFrame(&connection.responseBytes, response); err != nil {
		t.Fatal(err)
	}
	return connection
}

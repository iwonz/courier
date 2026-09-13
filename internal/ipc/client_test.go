package ipc

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
)

func TestCallResults(t *testing.T) {
	originalDial := dialControl
	t.Cleanup(func() { dialControl = originalDial })
	request := validRequest()

	response, _ := Success(request, map[string]bool{"ready": true})
	connection := newScriptedConn(t, response)
	dialControl = func(context.Context, string) (net.Conn, error) { return connection, nil }
	result, err := Call(context.Background(), "control", request)
	if err != nil || !result.OK || !connection.closed || connection.written.Len() == 0 {
		t.Fatalf("result=%+v closed=%v written=%d err=%v", result, connection.closed, connection.written.Len(), err)
	}

	rejection, _ := Rejection(request, CodeConflict, "conflict")
	connection = newScriptedConn(t, rejection)
	dialControl = func(context.Context, string) (net.Conn, error) { return connection, nil }
	if _, err := Call(context.Background(), "control", request); !IsRemoteError(err, CodeConflict) || IsRemoteError(err, CodeInvalid) {
		t.Fatalf("remote error=%v", err)
	}
	if IsRemoteError(errors.New("local"), CodeConflict) {
		t.Fatal("local error classified as remote")
	}
	if text := (&RemoteError{Code: CodeConflict, Message: "conflict"}).Error(); text != "Courier IPC conflict: conflict" {
		t.Fatalf("remote error text=%q", text)
	}

	mismatch := response
	mismatch.ID = delivery.ID("00000000-0000-4000-8000-000000000102")
	connection = newScriptedConn(t, mismatch)
	dialControl = func(context.Context, string) (net.Conn, error) { return connection, nil }
	if _, err := Call(context.Background(), "control", request); !errors.Is(err, ErrProtocol) {
		t.Fatalf("mismatch=%v", err)
	}

	invalid := response
	invalid.Version = 2
	connection = newScriptedConn(t, invalid)
	dialControl = func(context.Context, string) (net.Conn, error) { return connection, nil }
	if _, err := Call(context.Background(), "control", request); !errors.Is(err, ErrProtocol) {
		t.Fatalf("invalid response=%v", err)
	}

	want := errors.New("failure")
	dialControl = func(context.Context, string) (net.Conn, error) { return nil, want }
	if _, err := Call(context.Background(), "control", request); !errors.Is(err, want) {
		t.Fatalf("dial=%v", err)
	}
	badRequest := request
	badRequest.ID = "bad"
	if _, err := Call(context.Background(), "control", badRequest); !errors.Is(err, ErrProtocol) {
		t.Fatalf("request=%v", err)
	}
	connection = &scriptedConn{deadlineErr: want}
	dialControl = func(context.Context, string) (net.Conn, error) { return connection, nil }
	if _, err := Call(context.Background(), "control", request); !errors.Is(err, want) {
		t.Fatalf("deadline=%v", err)
	}
	connection = &scriptedConn{writeErr: want}
	dialControl = func(context.Context, string) (net.Conn, error) { return connection, nil }
	if _, err := Call(context.Background(), "control", request); !errors.Is(err, want) {
		t.Fatalf("write=%v", err)
	}
	connection = &scriptedConn{}
	dialControl = func(context.Context, string) (net.Conn, error) { return connection, nil }
	if _, err := Call(context.Background(), "control", request); !errors.Is(err, io.EOF) {
		t.Fatalf("read=%v", err)
	}
}

func TestConnectionContext(t *testing.T) {
	if _, err := Configure(context.Background(), nil); !errors.Is(err, ErrProtocol) {
		t.Fatalf("nil connection=%v", err)
	}
	connection := &scriptedConn{deadlines: make(chan time.Time, 2)}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := configureConnection(canceled, connection); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled=%v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	stop, err := configureConnection(ctx, connection)
	if err != nil {
		t.Fatal(err)
	}
	initial := <-connection.deadlines
	if time.Until(initial) <= 0 || time.Until(initial) > defaultCallTimeout+time.Second {
		t.Fatalf("default deadline=%v", initial)
	}
	cancel()
	select {
	case interrupted := <-connection.deadlines:
		if time.Until(interrupted) > time.Second {
			t.Fatalf("interrupt deadline=%v", interrupted)
		}
	case <-time.After(time.Second):
		t.Fatal("cancellation did not interrupt connection")
	}
	stop()

	deadline := time.Now().Add(time.Minute)
	withDeadline, cancelDeadline := context.WithDeadline(context.Background(), deadline)
	defer cancelDeadline()
	connection = &scriptedConn{deadlines: make(chan time.Time, 1)}
	stop, err = configureConnection(withDeadline, connection)
	if err != nil {
		t.Fatal(err)
	}
	defer stop()
	if actual := <-connection.deadlines; !actual.Equal(deadline) {
		t.Fatalf("deadline=%v want=%v", actual, deadline)
	}
	connection = &scriptedConn{}
	configuredStop, err := Configure(context.Background(), connection)
	if err != nil {
		t.Fatal(err)
	}
	configuredStop()
}

func newScriptedConn(t *testing.T, response Response) *scriptedConn {
	t.Helper()
	var input bytes.Buffer
	if err := WriteFrame(&input, response); err != nil {
		t.Fatal(err)
	}
	return &scriptedConn{read: bytes.NewReader(input.Bytes())}
}

type scriptedConn struct {
	read        *bytes.Reader
	written     bytes.Buffer
	writeErr    error
	deadlineErr error
	deadlines   chan time.Time
	closed      bool
}

func (connection *scriptedConn) Read(data []byte) (int, error) {
	if connection.read == nil {
		return 0, io.EOF
	}
	return connection.read.Read(data)
}
func (connection *scriptedConn) Write(data []byte) (int, error) {
	if connection.writeErr != nil {
		return 0, connection.writeErr
	}
	return connection.written.Write(data)
}
func (connection *scriptedConn) Close() error { connection.closed = true; return nil }
func (*scriptedConn) LocalAddr() net.Addr     { return testAddr("local") }
func (*scriptedConn) RemoteAddr() net.Addr    { return testAddr("remote") }
func (connection *scriptedConn) SetDeadline(value time.Time) error {
	if connection.deadlineErr != nil {
		return connection.deadlineErr
	}
	if connection.deadlines != nil {
		connection.deadlines <- value
	}
	return nil
}
func (*scriptedConn) SetReadDeadline(time.Time) error  { return nil }
func (*scriptedConn) SetWriteDeadline(time.Time) error { return nil }

type testAddr string

func (testAddr) Network() string        { return "test" }
func (address testAddr) String() string { return string(address) }

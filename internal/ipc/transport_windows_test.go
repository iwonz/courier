//go:build windows

package ipc

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
)

func TestWindowsTransport(t *testing.T) {
	if _, err := ControlEndpoint("", delivery.ID("bad")); !errors.Is(err, ErrProtocol) {
		t.Fatalf("bad ID=%v", err)
	}
	endpoint, err := ControlEndpoint("ignored", delivery.NewID())
	if err != nil || !strings.HasPrefix(endpoint, `\\.\pipe\courier-`) {
		t.Fatalf("endpoint=%q err=%v", endpoint, err)
	}
	if _, err := Listen(""); !errors.Is(err, ErrProtocol) {
		t.Fatalf("empty listen=%v", err)
	}
	if _, err := Dial(context.Background(), ""); !errors.Is(err, ErrProtocol) {
		t.Fatalf("empty dial=%v", err)
	}
	listener, err := Listen(endpoint)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	accepted := make(chan error, 1)
	go func() {
		connection, acceptErr := listener.Accept()
		if acceptErr == nil {
			acceptErr = connection.Close()
		}
		accepted <- acceptErr
	}()
	connection, err := Dial(ctx, endpoint)
	if err != nil {
		t.Fatal(err)
	}
	_ = connection.Close()
	select {
	case err := <-accepted:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatalf("accept did not complete: %v", ctx.Err())
	}
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	if err := RemoveStale(endpoint); err != nil {
		t.Fatal(err)
	}
}

func TestWindowsTransportFailures(t *testing.T) {
	originalSID, originalListen, originalDial := currentUserSID, listenWindows, dialWindows
	t.Cleanup(func() { currentUserSID, listenWindows, dialWindows = originalSID, originalListen, originalDial })
	want := errors.New("failure")
	currentUserSID = func() (string, error) { return "", want }
	if _, err := Listen("pipe"); !errors.Is(err, want) {
		t.Fatalf("SID=%v", err)
	}
	currentUserSID = func() (string, error) { return "S-1-5-18", nil }
	listenWindows = func(endpoint, descriptor string) (net.Listener, error) {
		if endpoint != "pipe" || !strings.Contains(descriptor, "S-1-5-18") {
			t.Fatalf("endpoint=%q descriptor=%q", endpoint, descriptor)
		}
		return nil, want
	}
	if _, err := Listen("pipe"); !errors.Is(err, want) {
		t.Fatalf("listen=%v", err)
	}
	dialWindows = func(context.Context, string) (net.Conn, error) { return nil, want }
	if _, err := Dial(context.Background(), "pipe"); !errors.Is(err, want) {
		t.Fatalf("dial=%v", err)
	}
}

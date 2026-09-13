//go:build !windows

package ipc

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
)

func TestUnixTransport(t *testing.T) {
	if _, err := ControlEndpoint("", requestID); !errors.Is(err, ErrProtocol) {
		t.Fatalf("empty directory=%v", err)
	}
	if _, err := ControlEndpoint(t.TempDir(), delivery.ID("bad")); !errors.Is(err, ErrProtocol) {
		t.Fatalf("bad ID=%v", err)
	}
	directory, err := os.MkdirTemp("/tmp", "courier-ipc-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(directory) })
	endpoint, err := ControlEndpoint(directory, requestID)
	if err != nil || filepath.Dir(endpoint) != directory {
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
	if runtime.GOOS != "windows" {
		info, statErr := os.Stat(endpoint)
		if statErr != nil || info.Mode().Perm() != 0o600 {
			t.Fatalf("mode=%v err=%v", info.Mode(), statErr)
		}
	}
	accepted := make(chan error, 1)
	go func() {
		connection, acceptErr := listener.Accept()
		if acceptErr == nil {
			acceptErr = connection.Close()
		}
		accepted <- acceptErr
	}()
	connection, err := Dial(context.Background(), endpoint)
	if err != nil {
		t.Fatal(err)
	}
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
	if err := <-accepted; err != nil {
		t.Fatal(err)
	}
	if _, err := Listen(endpoint); err == nil {
		t.Fatal("duplicate listener succeeded")
	}
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(endpoint); !os.IsNotExist(err) {
		t.Fatalf("socket remains: %v", err)
	}
	if err := listener.Close(); err == nil {
		t.Fatal("second close unexpectedly succeeded")
	}
}

func TestUnixTransportFailures(t *testing.T) {
	originalListen, originalChmod, originalRemove, originalLstat, originalDial := listenUnix, chmodEndpoint, removeEndpoint, lstatEndpoint, dialUnix
	t.Cleanup(func() {
		listenUnix, chmodEndpoint, removeEndpoint, lstatEndpoint, dialUnix = originalListen, originalChmod, originalRemove, originalLstat, originalDial
	})
	want := errors.New("failure")
	listenUnix = func(string) (net.Listener, error) { return nil, want }
	if _, err := Listen("endpoint"); !errors.Is(err, want) {
		t.Fatalf("listen=%v", err)
	}
	fake := &fakeListener{}
	listenUnix = func(string) (net.Listener, error) { return fake, nil }
	chmodEndpoint = func(string, os.FileMode) error { return want }
	removeEndpoint = func(string) error { return nil }
	if _, err := Listen("endpoint"); !errors.Is(err, want) || !fake.closed {
		t.Fatalf("chmod=%v closed=%v", err, fake.closed)
	}
	dialUnix = func(context.Context, string) (net.Conn, error) { return nil, want }
	if _, err := Dial(context.Background(), "endpoint"); !errors.Is(err, want) {
		t.Fatalf("dial=%v", err)
	}

	removeEndpoint = func(string) error { return want }
	fake = &fakeListener{closeErr: errors.New("close")}
	if err := (&unixListener{Listener: fake, endpoint: "endpoint"}).Close(); !errors.Is(err, want) || !errors.Is(err, fake.closeErr) {
		t.Fatalf("close errors=%v", err)
	}
	removeEndpoint = func(string) error { return os.ErrNotExist }
	fake = &fakeListener{}
	if err := (&unixListener{Listener: fake, endpoint: "endpoint"}).Close(); err != nil {
		t.Fatalf("missing removal=%v", err)
	}

	lstatEndpoint = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	if err := RemoveStale("endpoint"); err != nil {
		t.Fatalf("missing stale=%v", err)
	}
	lstatEndpoint = func(string) (os.FileInfo, error) { return nil, want }
	if err := RemoveStale("endpoint"); !errors.Is(err, want) {
		t.Fatalf("stale lstat=%v", err)
	}
	lstatEndpoint = func(string) (os.FileInfo, error) { return ipcFileInfo{mode: 0o600}, nil }
	if err := RemoveStale("endpoint"); !errors.Is(err, ErrProtocol) {
		t.Fatalf("stale regular=%v", err)
	}
	lstatEndpoint = func(string) (os.FileInfo, error) { return ipcFileInfo{mode: os.ModeSocket | os.ModeSymlink}, nil }
	if err := RemoveStale("endpoint"); !errors.Is(err, ErrProtocol) {
		t.Fatalf("stale symlink=%v", err)
	}
	lstatEndpoint = func(string) (os.FileInfo, error) { return ipcFileInfo{mode: os.ModeSocket}, nil }
	removeEndpoint = func(string) error { return want }
	if err := RemoveStale("endpoint"); !errors.Is(err, want) {
		t.Fatalf("stale remove=%v", err)
	}
}

type fakeListener struct {
	closed   bool
	closeErr error
}

func (*fakeListener) Accept() (net.Conn, error) { return nil, errors.New("unused") }
func (listener *fakeListener) Close() error     { listener.closed = true; return listener.closeErr }
func (*fakeListener) Addr() net.Addr            { return testAddr("listener") }

type ipcFileInfo struct{ mode os.FileMode }

func (ipcFileInfo) Name() string           { return "endpoint" }
func (ipcFileInfo) Size() int64            { return 0 }
func (info ipcFileInfo) Mode() os.FileMode { return info.mode }
func (ipcFileInfo) ModTime() time.Time     { return time.Time{} }
func (info ipcFileInfo) IsDir() bool       { return info.mode.IsDir() }
func (ipcFileInfo) Sys() any               { return nil }

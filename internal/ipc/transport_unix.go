//go:build !windows

package ipc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/iwonz/courier/internal/delivery"
)

var (
	listenUnix     = func(endpoint string) (net.Listener, error) { return net.Listen("unix", endpoint) }
	chmodEndpoint  = os.Chmod
	removeEndpoint = os.Remove
	lstatEndpoint  = os.Lstat
	dialUnix       = func(ctx context.Context, endpoint string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "unix", endpoint)
	}
)

func ControlEndpoint(stateDirectory string, serverID delivery.ID) (string, error) {
	if stateDirectory == "" || !serverID.Valid() {
		return "", fmt.Errorf("%w: state directory and server ID are required", ErrProtocol)
	}
	return filepath.Join(stateDirectory, "control-"+string(serverID)+".sock"), nil
}

func Listen(endpoint string) (net.Listener, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("%w: control endpoint is required", ErrProtocol)
	}
	listener, err := listenUnix(endpoint)
	if err != nil {
		return nil, err
	}
	if err := chmodEndpoint(endpoint, 0o600); err != nil {
		_ = listener.Close()
		_ = removeEndpoint(endpoint)
		return nil, err
	}
	return &unixListener{Listener: listener, endpoint: endpoint}, nil
}

func Dial(ctx context.Context, endpoint string) (net.Conn, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("%w: control endpoint is required", ErrProtocol)
	}
	return dialUnix(ctx, endpoint)
}

func RemoveStale(endpoint string) error {
	info, err := lstatEndpoint(endpoint)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSocket == 0 || info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("%w: stale control endpoint is not a socket", ErrProtocol)
	}
	return removeEndpoint(endpoint)
}

type unixListener struct {
	net.Listener
	endpoint string
}

func (listener *unixListener) Close() error {
	closeErr := listener.Listener.Close()
	removeErr := removeEndpoint(listener.endpoint)
	if os.IsNotExist(removeErr) {
		removeErr = nil
	}
	return errors.Join(closeErr, removeErr)
}

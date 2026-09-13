//go:build windows

package ipc

import (
	"context"
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
	"github.com/iwonz/courier/internal/delivery"
	"golang.org/x/sys/windows"
)

var (
	currentUserSID = func() (string, error) {
		user, err := windows.Token(0).GetTokenUser()
		if err != nil {
			return "", err
		}
		return user.User.Sid.String(), nil
	}
	listenWindows = func(endpoint, descriptor string) (net.Listener, error) {
		return winio.ListenPipe(endpoint, &winio.PipeConfig{SecurityDescriptor: descriptor, InputBufferSize: MaxFrameSize, OutputBufferSize: MaxFrameSize})
	}
	dialWindows = winio.DialPipeContext
)

func ControlEndpoint(_ string, serverID delivery.ID) (string, error) {
	if !serverID.Valid() {
		return "", fmt.Errorf("%w: server ID is required", ErrProtocol)
	}
	return `\\.\pipe\courier-` + string(serverID), nil
}

func Listen(endpoint string) (net.Listener, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("%w: control endpoint is required", ErrProtocol)
	}
	sid, err := currentUserSID()
	if err != nil {
		return nil, err
	}
	descriptor := "D:P(A;;GA;;;" + sid + ")"
	return listenWindows(endpoint, descriptor)
}

func Dial(ctx context.Context, endpoint string) (net.Conn, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("%w: control endpoint is required", ErrProtocol)
	}
	return dialWindows(ctx, endpoint)
}

func RemoveStale(string) error { return nil }

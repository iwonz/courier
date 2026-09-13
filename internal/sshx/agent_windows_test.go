//go:build windows

package sshx

import (
	"context"
	"errors"
	"net"
	"testing"
)

func TestWindowsAgentEndpointAndDial(t *testing.T) {
	if got := AgentEndpoint(""); got != windowsOpenSSHAgentPipe {
		t.Fatalf("default pipe=%q", got)
	}
	if got := AgentEndpoint("custom"); got != "custom" {
		t.Fatalf("custom pipe=%q", got)
	}
	original := dialWindowsAgentPipe
	t.Cleanup(func() { dialWindowsAgentPipe = original })
	server, client := net.Pipe()
	t.Cleanup(func() { _ = server.Close(); _ = client.Close() })
	dialWindowsAgentPipe = func(ctx context.Context, name string) (net.Conn, error) {
		if name != "custom" {
			t.Fatalf("pipe=%q", name)
		}
		if _, ok := ctx.Deadline(); !ok {
			t.Fatal("agent dial has no timeout")
		}
		return client, nil
	}
	connection, err := defaultAgentDial("custom")
	if err != nil || connection != client {
		t.Fatalf("connection=%v err=%v", connection, err)
	}
	dialWindowsAgentPipe = func(context.Context, string) (net.Conn, error) { return nil, errors.New("dial") }
	if _, err := defaultAgentDial("custom"); err == nil {
		t.Fatal("expected pipe dial error")
	}
}

//go:build windows

package sshx

import (
	"context"
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

const windowsOpenSSHAgentPipe = `\\.\pipe\openssh-ssh-agent`

var dialWindowsAgentPipe = winio.DialPipeContext

func defaultAgentEndpoint() string { return windowsOpenSSHAgentPipe }

func defaultAgentDial(name string) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return dialWindowsAgentPipe(ctx, name)
}

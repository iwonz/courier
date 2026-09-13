//go:build !windows

package sshx

import (
	"net"
	"time"
)

func defaultAgentEndpoint() string { return "" }

func defaultAgentDial(name string) (net.Conn, error) {
	return net.DialTimeout("unix", name, 3*time.Second)
}

//go:build !windows

package sshx

import (
	"path/filepath"
	"testing"
)

func TestAgentEndpointAndUnixDial(t *testing.T) {
	if got := AgentEndpoint("explicit"); got != "explicit" {
		t.Fatalf("explicit endpoint=%q", got)
	}
	if got := AgentEndpoint(""); got != "" {
		t.Fatalf("Unix default endpoint=%q", got)
	}
	if _, err := defaultAgentDial(filepath.Join(t.TempDir(), "missing-agent")); err == nil {
		t.Fatal("expected missing Unix agent error")
	}
}

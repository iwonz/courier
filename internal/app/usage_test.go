package app

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/terminalui"
	"github.com/spf13/cobra"
)

func TestUsageDiagnostics(t *testing.T) {
	root := NewRoot(Dependencies{})
	for _, test := range []struct {
		name        string
		args        []string
		cause       error
		want        []string
		forbidden   []string
		interactive bool
	}{
		{
			name: "cobra suggestion", args: []string{"services"},
			cause:     errors.New("unknown command \"services\" for \"courier\"\n\nDid you mean this?\n\tservers"),
			want:      []string{"error: unknown command", "did you mean: courier servers", "help: courier --help"},
			forbidden: []string{"stage:", "read:", "sent:", "confirmed:", "result:"},
		},
		{
			name: "selected command help", args: []string{"ui", "start", "extra"}, cause: errors.New("unknown command \"extra\" for \"courier ui start\""),
			want: []string{"error:", "help: courier ui start --help"},
		},
		{
			name: "nil and controls", cause: nil, interactive: true,
			want: []string{"invalid command usage", "courier --help"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			mode := terminalui.Mode{Interactive: test.interactive, Color: test.interactive, Width: 100}
			renderUsageError(&output, mode, root, test.args, test.cause)
			for _, wanted := range test.want {
				if !strings.Contains(output.String(), wanted) {
					t.Fatalf("missing %q in %q", wanted, output.String())
				}
			}
			for _, forbidden := range test.forbidden {
				if strings.Contains(output.String(), forbidden) {
					t.Fatalf("found %q in %q", forbidden, output.String())
				}
			}
		})
	}
	message, suggestion := usageMessage(errors.New("\n"))
	if message != "invalid command usage" || suggestion != "" {
		t.Fatalf("message=%q suggestion=%q", message, suggestion)
	}
	message, suggestion = usageMessage(errors.New("bad\nDid you mean this?\n\nignored"))
	if message != "bad" || suggestion != "" {
		t.Fatalf("blank suggestion message=%q suggestion=%q", message, suggestion)
	}
	want := errors.New("bad")
	wrapped := newUsageError(want)
	if wrapped.Error() != "bad" || !errors.Is(wrapped, want) {
		t.Fatalf("usage error=%v", wrapped)
	}
	root = &cobra.Command{Use: "courier", SilenceErrors: true, SilenceUsage: true, RunE: func(*cobra.Command, []string) error { return wrapped }}
	var output bytes.Buffer
	if code := Execute(context.Background(), root, nil, io.Discard, &output); code != ExitCLI || !strings.Contains(output.String(), "error: bad") {
		t.Fatalf("usage execute code=%d output=%q", code, output.String())
	}
	if err := newUsageError(errors.New("bad")); err.Error() != "bad" {
		t.Fatalf("usage error=%v", err)
	}
}

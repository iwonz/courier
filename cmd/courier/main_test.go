package main

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/cli"
)

func TestRun(t *testing.T) {
	for _, test := range []struct {
		name       string
		args       []string
		code       int
		wantOutput string
		wantError  string
	}{
		{name: "help", args: nil, code: 0, wantOutput: "courier from <source>"},
		{name: "version", args: []string{"version"}, code: 0, wantOutput: "courier dev"},
		{name: "version args", args: []string{"version", "extra"}, code: 2, wantError: "takes no arguments"},
		{name: "unknown", args: []string{"wat"}, code: 2, wantError: "unknown command"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if got := run(context.Background(), test.args, &stdout, &stderr); got != test.code {
				t.Fatalf("code=%d, want %d", got, test.code)
			}
			if !strings.Contains(stdout.String(), test.wantOutput) || !strings.Contains(stderr.String(), test.wantError) {
				t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
		})
	}
}

func TestVersionCommand(t *testing.T) {
	if versionCommand().Name() != "version" {
		t.Fatal("unexpected version command")
	}
}

func TestRunRegistryFailure(t *testing.T) {
	original := newRegistry
	t.Cleanup(func() { newRegistry = original })
	newRegistry = func(string, ...cli.Command) (*cli.Registry, error) {
		return nil, context.Canceled
	}
	var stderr bytes.Buffer
	if code := run(context.Background(), nil, &bytes.Buffer{}, &stderr); code != 1 || !strings.Contains(stderr.String(), "context canceled") {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
}

func TestMain(t *testing.T) {
	originalExit, originalArgs := exitProcess, os.Args
	t.Cleanup(func() { exitProcess, os.Args = originalExit, originalArgs })
	os.Args = []string{"courier", "version"}
	code := -1
	exitProcess = func(value int) { code = value }
	main()
	if code != 0 {
		t.Fatalf("exit code=%d", code)
	}
}

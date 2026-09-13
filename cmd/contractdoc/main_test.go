package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/iwonz/courier/internal/app"
	"github.com/iwonz/courier/internal/contract"
	"github.com/spf13/cobra"
)

func toolContract() contract.Contract {
	return contract.Contract{
		Commands: []contract.Command{
			{Name: "from", Path: "from", Status: "shipped", Flags: []string{"archive"}},
			{Name: "help", Path: "help", Status: "system", System: true},
			{Name: "update", Path: "update", Status: "shipped", System: true},
			{Name: "version", Path: "version", Status: "shipped", System: true},
		},
		Flags: []contract.Flag{{Name: "archive"}},
	}
}

func restoreGlobals(t *testing.T) {
	t.Helper()
	originalExit, originalRead, originalWrite, originalLoad, originalRoot := exitProcess, readFile, writeFile, load, root
	t.Cleanup(func() {
		exitProcess, readFile, writeFile, load, root = originalExit, originalRead, originalWrite, originalLoad, originalRoot
	})
}

func TestRunModes(t *testing.T) {
	restoreGlobals(t)
	_ = root()
	value := toolContract()
	load = func(string) (contract.Contract, error) { return value, nil }
	root = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
	reference := value.Reference()
	readFile = func(string) ([]byte, error) { return reference, nil }
	var stdout, stderr bytes.Buffer
	if code := run(nil, &stdout, &stderr); code != 0 || stdout.String() == "" || stderr.Len() != 0 {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}

	written := false
	writeFile = func(string, []byte, os.FileMode) error { written = true; return nil }
	stdout.Reset()
	if code := run([]string{"--write"}, &stdout, &stderr); code != 0 || !written {
		t.Fatalf("code=%d written=%v", code, written)
	}

	if code := run([]string{"bad"}, io.Discard, io.Discard); code != 2 {
		t.Fatalf("code=%d", code)
	}
	if code := run([]string{"--check", "extra"}, io.Discard, io.Discard); code != 2 {
		t.Fatalf("code=%d", code)
	}
}

func TestRunFailures(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		args  []string
	}{
		{"load", func() {
			load = func(string) (contract.Contract, error) { return contract.Contract{}, errors.New("load") }
		}, nil},
		{"parity", func() {
			load = func(string) (contract.Contract, error) { return toolContract(), nil }
			root = func() *cobra.Command { return &cobra.Command{Use: "courier"} }
		}, nil},
		{"read", func() {
			value := toolContract()
			load = func(string) (contract.Contract, error) { return value, nil }
			root = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
			readFile = func(string) ([]byte, error) { return nil, errors.New("read") }
		}, nil},
		{"stale", func() {
			value := toolContract()
			load = func(string) (contract.Contract, error) { return value, nil }
			root = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
			readFile = func(string) ([]byte, error) { return []byte("stale"), nil }
		}, nil},
		{"write", func() {
			value := toolContract()
			load = func(string) (contract.Contract, error) { return value, nil }
			root = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
			writeFile = func(string, []byte, os.FileMode) error { return errors.New("write") }
		}, []string{"--write"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restoreGlobals(t)
			test.setup()
			var stderr bytes.Buffer
			if code := run(test.args, io.Discard, &stderr); code != 1 || stderr.Len() == 0 {
				t.Fatalf("code=%d stderr=%q", code, stderr.String())
			}
		})
	}
}

func TestMain(t *testing.T) {
	restoreGlobals(t)
	value := toolContract()
	load = func(string) (contract.Contract, error) { return value, nil }
	root = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
	readFile = func(string) ([]byte, error) { return value.Reference(), nil }
	exited := -1
	exitProcess = func(code int) { exited = code }
	originalArgs := os.Args
	os.Args = []string{"contractdoc", "--check"}
	t.Cleanup(func() { os.Args = originalArgs })
	main()
	if exited != 0 {
		t.Fatalf("exit=%d", exited)
	}
}

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
	value, err := contract.Load("../../docs/cli-contract.yaml")
	if err != nil {
		panic(err)
	}
	return value
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
	readFile = func(name string) ([]byte, error) {
		if name == landingPath {
			return value.LandingJSON(), nil
		}
		return reference, nil
	}
	var stdout, stderr bytes.Buffer
	if code := run(nil, &stdout, &stderr); code != 0 || stdout.String() == "" || stderr.Len() != 0 {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}

	written := 0
	writeFile = func(string, []byte, os.FileMode) error { written++; return nil }
	stdout.Reset()
	if code := run([]string{"--write"}, &stdout, &stderr); code != 0 || written != 2 {
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
		{"landing read", func() {
			value := toolContract()
			load = func(string) (contract.Contract, error) { return value, nil }
			root = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
			readFile = func(name string) ([]byte, error) {
				if name == landingPath {
					return nil, errors.New("read")
				}
				return value.Reference(), nil
			}
		}, nil},
		{"landing stale", func() {
			value := toolContract()
			load = func(string) (contract.Contract, error) { return value, nil }
			root = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
			readFile = func(name string) ([]byte, error) {
				if name == landingPath {
					return []byte("stale"), nil
				}
				return value.Reference(), nil
			}
		}, nil},
		{"write", func() {
			value := toolContract()
			load = func(string) (contract.Contract, error) { return value, nil }
			root = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
			writeFile = func(string, []byte, os.FileMode) error { return errors.New("write") }
		}, []string{"--write"}},
		{"landing write", func() {
			value := toolContract()
			load = func(string) (contract.Contract, error) { return value, nil }
			root = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
			writeFile = func(name string, _ []byte, _ os.FileMode) error {
				if name == landingPath {
					return errors.New("write")
				}
				return nil
			}
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
	readFile = func(name string) ([]byte, error) {
		if name == landingPath {
			return value.LandingJSON(), nil
		}
		return value.Reference(), nil
	}
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

func TestRunPlannerFailure(t *testing.T) {
	restoreGlobals(t)
	value := toolContract()
	value.EndpointKinds = append(value.EndpointKinds, contract.Endpoint{Name: "extra", Status: "planned", Syntax: "extra://"})
	load = func(string) (contract.Contract, error) { return value, nil }
	root = func() *cobra.Command { return app.NewRoot(app.Dependencies{}) }
	var stderr bytes.Buffer
	if code := run(nil, io.Discard, &stderr); code != 1 || !bytes.Contains(stderr.Bytes(), []byte("check planner parity")) {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
}

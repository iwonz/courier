package app

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/admin"
)

func TestUICommands(t *testing.T) {
	want := errors.New("UI failure")
	var started admin.StartRequest
	dependencies := Dependencies{
		StartUI: func(_ context.Context, request admin.StartRequest) (admin.StartResult, error) {
			started = request
			state := admin.State{Bind: request.Bind}
			if request.Ready != nil {
				if err := request.Ready(state); err != nil {
					return admin.StartResult{}, err
				}
			}
			return admin.StartResult{State: state}, nil
		},
		StopUI: func(context.Context) (admin.StopResult, error) { return admin.StopResult{}, nil },
	}
	for _, test := range []struct {
		name       string
		args       []string
		code       int
		wantOutput string
		wantError  string
	}{
		{"group", []string{"ui"}, ExitOK, "Control the local", ""},
		{"start", []string{"ui", "start"}, ExitOK, "http://127.0.0.1:9090/", ""},
		{"background", []string{"ui", "start", "--background", "--listen", "[::1]:9091"}, ExitOK, "http://[::1]:9091/", ""},
		{"start args", []string{"ui", "start", "extra"}, ExitCLI, "", "unknown command"},
		{"invalid bind", []string{"ui", "start", "--listen", "0.0.0.0:9090"}, ExitCLI, "", "loopback"},
		{"repeated listen", []string{"ui", "start", "--listen", "127.0.0.1:1", "--listen", "127.0.0.1:2"}, ExitCLI, "", "value may only be set once"},
		{"repeated background", []string{"ui", "start", "--background", "--background"}, ExitCLI, "", "value may only be set once"},
		{"stop", []string{"ui", "stop"}, ExitOK, "Stopped Courier", ""},
		{"stop args", []string{"ui", "stop", "extra"}, ExitCLI, "", "unknown command"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Execute(context.Background(), NewRoot(dependencies), test.args, &stdout, &stderr)
			if code != test.code || !strings.Contains(stdout.String(), test.wantOutput) || !strings.Contains(stderr.String(), test.wantError) {
				t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
		})
	}
	if !started.Background || started.Bind != "[::1]:9091" {
		t.Fatalf("started=%+v", started)
	}

	missing := dependencies
	missing.StartUI = nil
	if code := Execute(context.Background(), NewRoot(missing), []string{"ui", "start"}, io.Discard, io.Discard); code != ExitControl {
		t.Fatalf("missing start=%d", code)
	}
	missing = dependencies
	missing.StopUI = nil
	if code := Execute(context.Background(), NewRoot(missing), []string{"ui", "stop"}, io.Discard, io.Discard); code != ExitControl {
		t.Fatalf("missing stop=%d", code)
	}
	failing := dependencies
	failing.StartUI = func(context.Context, admin.StartRequest) (admin.StartResult, error) { return admin.StartResult{}, want }
	if code := Execute(context.Background(), NewRoot(failing), []string{"ui", "start"}, io.Discard, io.Discard); code != ExitControl {
		t.Fatalf("failed start=%d", code)
	}
	failing.StopUI = func(context.Context) (admin.StopResult, error) { return admin.StopResult{}, want }
	if code := Execute(context.Background(), NewRoot(failing), []string{"ui", "stop"}, io.Discard, io.Discard); code != ExitControl {
		t.Fatalf("failed stop=%d", code)
	}
}

func TestUICommandOutputAndAlreadyStopped(t *testing.T) {
	dependencies := Dependencies{
		StartUI: func(_ context.Context, request admin.StartRequest) (admin.StartResult, error) {
			return admin.StartResult{}, request.Ready(admin.State{Bind: request.Bind})
		},
		StopUI: func(context.Context) (admin.StopResult, error) { return admin.StopResult{AlreadyStopped: true}, nil },
	}
	var output bytes.Buffer
	if code := Execute(context.Background(), NewRoot(dependencies), []string{"ui", "stop"}, &output, io.Discard); code != ExitOK || !strings.Contains(output.String(), "already stopped") {
		t.Fatalf("code/output=%d %q", code, output.String())
	}
	want := errors.New("write")
	if code := Execute(context.Background(), NewRoot(dependencies), []string{"ui", "start"}, failingWriter{err: want}, io.Discard); code != ExitControl {
		t.Fatalf("start writer=%d", code)
	}
	if code := Execute(context.Background(), NewRoot(dependencies), []string{"ui", "stop"}, failingWriter{err: want}, io.Discard); code != ExitControl {
		t.Fatalf("stop writer=%d", code)
	}
}

func TestDefaultUIFunctions(t *testing.T) {
	original := defaultStateDirectory
	t.Cleanup(func() { defaultStateDirectory = original })
	want := errors.New("directory")
	defaultStateDirectory = func() (string, error) { return "", want }
	if _, err := adminManager(); !errors.Is(err, want) {
		t.Fatalf("manager=%v", err)
	}
	if _, err := startUI(context.Background(), admin.StartRequest{Bind: "bad"}); !errors.Is(err, want) {
		t.Fatalf("start=%v", err)
	}
	if _, err := stopUI(context.Background()); !errors.Is(err, want) {
		t.Fatalf("stop=%v", err)
	}
	directory := filepath.Join(t.TempDir(), "state")
	defaultStateDirectory = func() (string, error) { return directory, nil }
	manager, err := adminManager()
	if err != nil || manager.StateDirectory != directory || !strings.HasPrefix(manager.Compatibility, "admin-v1/") {
		t.Fatalf("manager=%+v err=%v", manager, err)
	}
	if _, err := startUI(context.Background(), admin.StartRequest{Bind: "bad"}); err == nil {
		t.Fatal("invalid default start accepted")
	}
	result, err := stopUI(context.Background())
	if err != nil || !result.AlreadyStopped {
		t.Fatalf("default stop=%+v err=%v", result, err)
	}
}

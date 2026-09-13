package app

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/archive"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/report"
	"github.com/iwonz/courier/internal/selection"
	"github.com/iwonz/courier/internal/sshx"
	"github.com/iwonz/courier/internal/transfer"
	"github.com/iwonz/courier/internal/update"
	"github.com/spf13/cobra"
)

func TestCommandsAndExitCodes(t *testing.T) {
	base := Dependencies{
		Build: BuildIdentity{Version: "v1.2.3", Commit: "abc", Date: "today"},
		Update: func(context.Context) (update.Result, error) {
			return update.Result{Current: true, From: "v1.2.3"}, nil
		},
		Hosted: func(_ context.Context, _ operation.Plan, output io.Writer) error {
			_, err := io.WriteString(output, "delivery: http://127.0.0.1:8080/d/token/\n")
			return err
		},
	}
	for _, test := range []struct {
		name       string
		args       []string
		code       int
		wantOutput string
		wantError  string
	}{
		{name: "root", args: nil, code: ExitOK, wantOutput: "Safely transfer"},
		{name: "help", args: []string{"--help"}, code: ExitOK, wantOutput: "Safely transfer"},
		{name: "version", args: []string{"version"}, code: ExitOK, wantOutput: "courier v1.2.3 (commit abc, built today)"},
		{name: "version args", args: []string{"version", "extra"}, code: ExitCLI, wantError: "stage: preflight"},
		{name: "unknown", args: []string{"unknown"}, code: ExitCLI, wantError: "unknown command"},
		{name: "invalid transfer grammar", args: []string{"from", "a", "into", "b"}, code: ExitCLI, wantError: "expected: courier from"},
		{name: "unsupported transfer route", args: []string{"from", "https://example.test/file", "to", "out"}, code: ExitCLI, wantError: "unsupported operation route"},
		{name: "browser transfer route", args: []string{"from", "web://", "to", "out", "--background"}, code: ExitOK, wantOutput: "delivery: http://"},
		{name: "webhook transfer route", args: []string{"from", "webhook://", "to", "out"}, code: ExitOK, wantOutput: "delivery: http://"},
		{name: "duplicate archive", args: []string{"from", "a", "to", "b", "--archive", "--archive"}, code: ExitCLI, wantError: "value may only be set once"},
		{name: "current update", args: []string{"update"}, code: ExitOK, wantOutput: "courier v1.2.3 is current"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Execute(context.Background(), NewRoot(base), test.args, &stdout, &stderr)
			if code != test.code || !strings.Contains(stdout.String(), test.wantOutput) || !strings.Contains(stderr.String(), test.wantError) {
				t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
			}
		})
	}

	t.Run("updated with notes", func(t *testing.T) {
		dependencies := base
		dependencies.Update = func(context.Context) (update.Result, error) {
			return update.Result{From: "v1", To: "v2", Notes: "changes"}, nil
		}
		var output bytes.Buffer
		if code := Execute(context.Background(), NewRoot(dependencies), []string{"update"}, &output, io.Discard); code != ExitOK || !strings.Contains(output.String(), "updated courier from v1 to v2\nchanges") {
			t.Fatalf("code=%d output=%q", code, output.String())
		}
	})

	for _, test := range []struct {
		name string
		run  func(context.Context) (update.Result, error)
	}{
		{name: "unavailable"},
		{name: "failed", run: func(context.Context) (update.Result, error) { return update.Result{}, errors.New("network") }},
	} {
		t.Run("update "+test.name, func(t *testing.T) {
			dependencies := base
			dependencies.Update = test.run
			var stderr bytes.Buffer
			if code := Execute(context.Background(), NewRoot(dependencies), []string{"update"}, io.Discard, &stderr); code != ExitUpdate || !strings.Contains(stderr.String(), "stage: update") {
				t.Fatalf("code=%d stderr=%q", code, stderr.String())
			}
		})
	}
}

func TestUnknownRuntimeRoute(t *testing.T) {
	err := runRoute(context.Background(), Dependencies{}, operation.Plan{Route: operation.Route(255)}, selection.All(), io.Discard, io.Discard)
	var commandErr *commandError
	if !errors.As(err, &commandErr) || commandErr.code != ExitCLI {
		t.Fatalf("unknown route=%v", err)
	}
}

func TestExecuteInterrupted(t *testing.T) {
	for _, err := range []error{
		context.Canceled,
		&commandError{code: ExitTransfer, stage: string(progress.StageTransfer), read: 3, sent: 2, confirmed: 1, cause: context.Canceled},
	} {
		root := &cobra.Command{Use: "courier", RunE: func(*cobra.Command, []string) error { return err }}
		var stderr bytes.Buffer
		if code := Execute(context.Background(), root, nil, io.Discard, &stderr); code != ExitInterrupted || !strings.Contains(stderr.String(), "result: failed") {
			t.Fatalf("code=%d stderr=%q", code, stderr.String())
		}
	}
}

func TestNoOpAndDestinationCollision(t *testing.T) {
	t.Run("plain identity is a no-op", func(t *testing.T) {
		called := false
		dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path}, nil
		})
		dependencies.Transfer = func(context.Context, transfer.Request) (transfer.Result, error) {
			called = true
			return transfer.Result{}, nil
		}
		var stdout bytes.Buffer
		code := Execute(context.Background(), NewRoot(dependencies), []string{"from", "same", "to", "same"}, &stdout, io.Discard)
		if code != ExitOK || called || !strings.Contains(stdout.String(), "transferred: 0 bytes") {
			t.Fatalf("code=%d called=%v stdout=%q", code, called, stdout.String())
		}
	})

	t.Run("resolved SSH aliases are a no-op", func(t *testing.T) {
		called := false
		dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			value.Host = "canonical.example"
			value.User = "same-user"
			return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path}, nil
		})
		dependencies.Transfer = func(context.Context, transfer.Request) (transfer.Result, error) {
			called = true
			return transfer.Result{}, nil
		}
		if code := Execute(context.Background(), NewRoot(dependencies), []string{"from", "first:/same", "to", "second:/same"}, io.Discard, io.Discard); code != ExitOK || called {
			t.Fatalf("code=%d called=%v", code, called)
		}
	})

	t.Run("container resolution is a no-op", func(t *testing.T) {
		root := t.TempDir()
		source := filepath.Join(root, "source")
		if err := os.WriteFile(source, []byte("keep"), 0o600); err != nil {
			t.Fatal(err)
		}
		called := false
		dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path}, nil
		})
		dependencies.Transfer = func(context.Context, transfer.Request) (transfer.Result, error) {
			called = true
			return transfer.Result{}, nil
		}
		if code := Execute(context.Background(), NewRoot(dependencies), []string{"from", source, "to", root + string(filepath.Separator)}, io.Discard, io.Discard); code != ExitOK || called {
			t.Fatalf("code=%d called=%v", code, called)
		}
	})

	t.Run("archive identity is unsafe", func(t *testing.T) {
		archiveCalled := false
		dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path}, nil
		})
		dependencies.Archive = func(context.Context, fsx.Backend, string, string, string, selection.Selector, progress.Sink) (*archive.Artifact, error) {
			archiveCalled = true
			return nil, nil
		}
		var stderr bytes.Buffer
		code := Execute(context.Background(), NewRoot(dependencies), []string{"from", "same", "to", "same", "--archive"}, io.Discard, &stderr)
		if code != ExitTransfer || archiveCalled || !strings.Contains(stderr.String(), "transformed output collides") {
			t.Fatalf("code=%d archiveCalled=%v stderr=%q", code, archiveCalled, stderr.String())
		}
	})

	t.Run("container child collision preserves entries", func(t *testing.T) {
		root := t.TempDir()
		source := filepath.Join(root, "source")
		container := filepath.Join(root, "container")
		if err := os.WriteFile(source, []byte("new"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(container, 0o700); err != nil {
			t.Fatal(err)
		}
		collision := filepath.Join(container, "source")
		unrelated := filepath.Join(container, "unrelated")
		if err := os.WriteFile(collision, []byte("keep"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(unrelated, []byte("also keep"), 0o600); err != nil {
			t.Fatal(err)
		}
		called := false
		dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path}, nil
		})
		dependencies.Transfer = func(context.Context, transfer.Request) (transfer.Result, error) {
			called = true
			return transfer.Result{}, nil
		}
		var stderr bytes.Buffer
		code := Execute(context.Background(), NewRoot(dependencies), []string{"from", source, "to", container + string(filepath.Separator)}, io.Discard, &stderr)
		collisionData, collisionErr := os.ReadFile(collision)
		unrelatedData, unrelatedErr := os.ReadFile(unrelated)
		if code != ExitTransfer || called || collisionErr != nil || unrelatedErr != nil || string(collisionData) != "keep" || string(unrelatedData) != "also keep" || !strings.Contains(stderr.String(), "destination already exists") {
			t.Fatalf("code=%d called=%v collision=%q,%v unrelated=%q,%v stderr=%q", code, called, collisionData, collisionErr, unrelatedData, unrelatedErr, stderr.String())
		}
	})
}

func TestSelectionCommandPreflight(t *testing.T) {
	t.Run("ordered compilation failure precedes endpoint opening", func(t *testing.T) {
		opened := false
		var got []operation.SelectionRule
		dependencies := transferDependencies(t, func(context.Context, endpoint.Endpoint) (*Resource, error) {
			opened = true
			return nil, errors.New("must not open")
		})
		dependencies.Select = func(rules []operation.SelectionRule) (selection.Selector, error) {
			got = append([]operation.SelectionRule(nil), rules...)
			return nil, errors.New("selection compile")
		}
		args := []string{"from", "in", "to", "out", "--exclude", "*.tmp", "--exclude-from", "rules", "--exclude-regex", "^private/", "--exclude", "!keep.tmp"}
		var stderr bytes.Buffer
		code := Execute(context.Background(), NewRoot(dependencies), args, io.Discard, &stderr)
		wantKinds := []operation.SelectionKind{operation.SelectionGitignore, operation.SelectionFile, operation.SelectionRegex, operation.SelectionGitignore}
		if code != ExitCLI || opened || len(got) != len(wantKinds) || !strings.Contains(stderr.String(), "selection compile") {
			t.Fatalf("code=%d opened=%v rules=%v stderr=%q", code, opened, got, stderr.String())
		}
		for index, kind := range wantKinds {
			if got[index].Kind != kind || got[index].Position != index {
				t.Fatalf("rules=%v", got)
			}
		}
	})

	t.Run("missing compiler", func(t *testing.T) {
		dependencies := transferDependencies(t, nil)
		var stderr bytes.Buffer
		code := Execute(context.Background(), NewRoot(dependencies), []string{"from", "in", "to", "out", "--exclude", "*.tmp"}, io.Discard, &stderr)
		if code != ExitCLI || !strings.Contains(stderr.String(), "selection dependencies are incomplete") {
			t.Fatalf("code=%d stderr=%q", code, stderr.String())
		}
	})

	t.Run("compiled selector reaches transfer", func(t *testing.T) {
		root := t.TempDir()
		source := filepath.Join(root, "source")
		destination := filepath.Join(root, "destination")
		if err := os.WriteFile(source, []byte("source"), 0o600); err != nil {
			t.Fatal(err)
		}
		dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path}, nil
		})
		dependencies.Select = func([]operation.SelectionRule) (selection.Selector, error) {
			return selection.All(), nil
		}
		dependencies.Transfer = func(_ context.Context, request transfer.Request) (transfer.Result, error) {
			if request.Selector == nil || !request.Selector.Include("keep", false) {
				t.Fatal("compiled selector was not propagated")
			}
			return transfer.Result{Bytes: 6, Destination: request.Destination}, nil
		}
		if code := Execute(context.Background(), NewRoot(dependencies), []string{"from", source, "to", destination, "--exclude", "*.tmp"}, io.Discard, io.Discard); code != ExitOK {
			t.Fatalf("code=%d", code)
		}
	})
}

func TestAllTransferDirections(t *testing.T) {
	tests := []struct {
		name              string
		sourceRemote      bool
		destinationRemote bool
	}{
		{"local to local", false, false},
		{"local to remote", false, true},
		{"remote to local", true, false},
		{"remote to remote", true, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			sourcePath := filepath.Join(root, "source.txt")
			destinationPath := filepath.Join(root, "destination.txt")
			if err := os.WriteFile(sourcePath, []byte("payload"), 0o640); err != nil {
				t.Fatal(err)
			}
			sourceText := sourcePath
			if test.sourceRemote {
				sourceText = "source-host:/source.txt"
			}
			destinationText := destinationPath
			if test.destinationRemote {
				destinationText = "destination-host:/destination.txt"
			}
			var closes atomic.Int32
			dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
				path := value.Path
				if value.Remote && value.Host == "source-host" {
					path = sourcePath
				}
				if value.Remote && value.Host == "destination-host" {
					path = destinationPath
				}
				return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: path, Close: func() error { closes.Add(1); return nil }}, nil
			})
			var stdout, stderr bytes.Buffer
			code := Execute(context.Background(), NewRoot(dependencies), []string{"from", sourceText, "to", destinationText}, &stdout, &stderr)
			data, readErr := os.ReadFile(destinationPath)
			if code != ExitOK || readErr != nil || string(data) != "payload" || closes.Load() != 2 {
				t.Fatalf("code=%d data=%q readErr=%v closes=%d stdout=%q stderr=%q", code, data, readErr, closes.Load(), stdout.String(), stderr.String())
			}
			if !strings.Contains(stdout.String(), "destination: "+destinationText) || !strings.Contains(stderr.String(), "stage=transfer") {
				t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
		})
	}
}

func TestArchiveDestinationAndResolvedAliasSafety(t *testing.T) {
	t.Run("archive into remote directory", func(t *testing.T) {
		root := t.TempDir()
		source := filepath.Join(root, "photos")
		destination := filepath.Join(root, "uploads")
		if err := os.WriteFile(source, []byte("photo"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(destination, 0o700); err != nil {
			t.Fatal(err)
		}
		dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			path := value.Path
			if value.Remote {
				path = destination
			}
			return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: path, Close: func() error { return nil }}, nil
		})
		var stdout bytes.Buffer
		code := Execute(context.Background(), NewRoot(dependencies), []string{"from", source, "to", "server:/uploads/", "--archive"}, &stdout, io.Discard)
		archivePath := filepath.Join(destination, "photos.tar.gz")
		if code != ExitOK {
			t.Fatalf("code=%d output=%q", code, stdout.String())
		}
		if err := archive.Verify(archivePath); err != nil {
			t.Fatalf("archive verification: %v", err)
		}
		if !strings.Contains(stdout.String(), "destination: server:/uploads/photos.tar.gz") {
			t.Fatalf("output=%q", stdout.String())
		}
	})

	t.Run("aliases resolve before safety check", func(t *testing.T) {
		backend := lstatBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return fakeInfo{name: "source", directory: true}, nil }}
		dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			value.Host = "canonical.example"
			value.User = "user"
			return &Resource{Endpoint: value, Backend: backend, Path: value.Path}, nil
		})
		called := false
		dependencies.Transfer = func(context.Context, transfer.Request) (transfer.Result, error) {
			called = true
			return transfer.Result{}, nil
		}
		var stderr bytes.Buffer
		code := Execute(context.Background(), NewRoot(dependencies), []string{"from", "first:/tree", "to", "second:/tree/child"}, io.Discard, &stderr)
		if code != ExitTransfer || called || !strings.Contains(stderr.String(), "destination is inside source") {
			t.Fatalf("code=%d called=%v stderr=%q", code, called, stderr.String())
		}
	})
}

func TestExtractCommand(t *testing.T) {
	root := t.TempDir()
	archivePath := filepath.Join(root, "source.tar.gz")
	sourceDirectory := filepath.Join(root, "source")
	if err := os.Mkdir(sourceDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDirectory, "file.txt"), []byte("payload"), 0o600); err != nil {
		t.Fatal(err)
	}
	artifact, err := archive.Create(context.Background(), fsx.Local{}, sourceDirectory, "source", root, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(artifact.Path, archivePath); err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(root, "extracted")
	dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
		return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path}, nil
	})
	dependencies.Transfer = func(context.Context, transfer.Request) (transfer.Result, error) {
		t.Fatal("copy engine must not run for extraction")
		return transfer.Result{}, nil
	}
	var stdout, stderr bytes.Buffer
	code := Execute(context.Background(), NewRoot(dependencies), []string{"from", archivePath, "to", destination, "--extract", "--max-extracted-size", "1GiB"}, &stdout, &stderr)
	data, readErr := os.ReadFile(filepath.Join(destination, "source", "file.txt"))
	if code != ExitOK || readErr != nil || string(data) != "payload" || !strings.Contains(stdout.String(), "transferred: 7 bytes") || !strings.Contains(stderr.String(), "stage=extract") {
		t.Fatalf("code=%d data=%q readErr=%v stdout=%q stderr=%q", code, data, readErr, stdout.String(), stderr.String())
	}

	for _, test := range []struct {
		name      string
		args      []string
		mutate    func(*Dependencies)
		wantCode  int
		wantStage string
		wantText  string
	}{
		{name: "identity", args: []string{"from", archivePath, "to", archivePath, "--extract"}, wantCode: ExitTransfer, wantStage: "preflight", wantText: "transformed output collides"},
		{name: "archive conflict", args: []string{"from", archivePath, "to", destination + "-a", "--archive", "--extract"}, wantCode: ExitCLI, wantStage: "preflight", wantText: "conflicts"},
		{name: "limit requires extract", args: []string{"from", archivePath, "to", destination + "-b", "--max-extracted-size", "1GiB"}, wantCode: ExitCLI, wantStage: "preflight", wantText: "requires --extract"},
		{name: "duplicate extract", args: []string{"from", archivePath, "to", destination + "-c", "--extract", "--extract"}, wantCode: ExitCLI, wantStage: "preflight", wantText: "value may only be set once"},
		{name: "duplicate limit", args: []string{"from", archivePath, "to", destination + "-d", "--extract", "--max-extracted-size", "1GiB", "--max-extracted-size", "2GiB"}, wantCode: ExitCLI, wantStage: "preflight", wantText: "value may only be set once"},
		{name: "incomplete", args: []string{"from", archivePath, "to", destination + "-e", "--extract"}, mutate: func(value *Dependencies) { value.Extract = nil }, wantCode: ExitTransfer, wantStage: "preflight", wantText: "incomplete"},
		{name: "typed failure", args: []string{"from", archivePath, "to", destination + "-f", "--extract"}, mutate: func(value *Dependencies) {
			value.Extract = func(context.Context, archive.ExtractionRequest) (archive.ExtractionResult, error) {
				cause := errors.New("late collision")
				return archive.ExtractionResult{Bytes: 3}, &archive.ExtractionError{Stage: progress.StageCommit, Confirmed: 3, Cause: cause}
			}
		}, wantCode: ExitTransfer, wantStage: "commit", wantText: "confirmed: 3 bytes"},
		{name: "generic failure", args: []string{"from", archivePath, "to", destination + "-g", "--extract"}, mutate: func(value *Dependencies) {
			value.Extract = func(context.Context, archive.ExtractionRequest) (archive.ExtractionResult, error) {
				return archive.ExtractionResult{Bytes: 2}, errors.New("extract failure")
			}
		}, wantCode: ExitTransfer, wantStage: "extract", wantText: "confirmed: 2 bytes"},
	} {
		t.Run(test.name, func(t *testing.T) {
			configured := dependencies
			if test.mutate != nil {
				test.mutate(&configured)
			}
			var stderr bytes.Buffer
			code := Execute(context.Background(), NewRoot(configured), test.args, io.Discard, &stderr)
			if code != test.wantCode || !strings.Contains(stderr.String(), "stage: "+test.wantStage) || !strings.Contains(stderr.String(), test.wantText) {
				t.Fatalf("code=%d stderr=%q", code, stderr.String())
			}
		})
	}
}

func TestTransferFailurePathsAndCleanup(t *testing.T) {
	validOpen := func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
		return &Resource{Endpoint: value, Backend: lstatBackend{Backend: fsx.Local{}, lstat: func(name string) (fs.FileInfo, error) {
			if name == "out" {
				return nil, fs.ErrNotExist
			}
			return fakeInfo{name: "file", size: 4}, nil
		}}, Path: value.Path}, nil
	}
	for _, test := range []struct {
		name      string
		args      []string
		mutate    func(*Dependencies)
		wantCode  int
		wantStage string
		wantText  string
	}{
		{name: "invalid source", args: []string{"from", "", "to", "out"}, wantCode: ExitCLI, wantStage: "preflight", wantText: "empty value"},
		{name: "invalid destination", args: []string{"from", "in", "to", "bad\x00path"}, wantCode: ExitCLI, wantStage: "preflight", wantText: "control character"},
		{name: "incomplete", args: []string{"from", "in", "to", "out"}, mutate: func(dependencies *Dependencies) { dependencies.Transfer = nil }, wantCode: ExitTransfer, wantStage: "preflight", wantText: "incomplete"},
		{name: "local open", args: []string{"from", "in", "to", "out"}, mutate: func(dependencies *Dependencies) {
			dependencies.Open = func(context.Context, endpoint.Endpoint) (*Resource, error) { return nil, errors.New("local open") }
		}, wantCode: ExitTransfer, wantStage: "preflight", wantText: "local open"},
		{name: "remote open", args: []string{"from", "host:/in", "to", "out"}, mutate: func(dependencies *Dependencies) {
			dependencies.Open = func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
				if value.Remote {
					return nil, errors.New("connect")
				}
				return validOpen(context.Background(), value)
			}
		}, wantCode: ExitConnection, wantStage: "connection", wantText: "connect"},
		{name: "source stat", args: []string{"from", "in", "to", "out"}, mutate: func(dependencies *Dependencies) {
			dependencies.Open = func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
				backend := lstatBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return nil, errors.New("source stat") }}
				if value.Path == "out" {
					backend.lstat = func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
				}
				return &Resource{Endpoint: value, Backend: backend, Path: value.Path}, nil
			}
		}, wantCode: ExitTransfer, wantStage: "preflight", wantText: "source stat"},
		{name: "destination stat", args: []string{"from", "in", "to", "out"}, mutate: func(dependencies *Dependencies) {
			dependencies.Open = func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
				backend := lstatBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return fakeInfo{name: "file"}, nil }}
				if value.Path == "out" {
					backend.lstat = func(string) (fs.FileInfo, error) { return nil, errors.New("destination stat") }
				}
				return &Resource{Endpoint: value, Backend: backend, Path: value.Path}, nil
			}
		}, wantCode: ExitTransfer, wantStage: "preflight", wantText: "destination stat"},
		{name: "resolved destination stat", args: []string{"from", "in", "to", "out/"}, mutate: func(dependencies *Dependencies) {
			dependencies.Open = func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
				backend := lstatBackend{Backend: fsx.Local{}, lstat: func(name string) (fs.FileInfo, error) {
					switch name {
					case "in":
						return fakeInfo{name: "in"}, nil
					case "out/":
						return nil, fs.ErrNotExist
					default:
						return nil, errors.New("resolved destination stat")
					}
				}}
				return &Resource{Endpoint: value, Backend: backend, Path: value.Path}, nil
			}
		}, wantCode: ExitTransfer, wantStage: "preflight", wantText: "resolved destination stat"},
		{name: "nameless source", args: []string{"from", "/", "to", "out/"}, wantCode: ExitTransfer, wantStage: "preflight", wantText: "no transferable name"},
		{name: "typed transfer", args: []string{"from", "in", "to", "out"}, mutate: func(dependencies *Dependencies) {
			dependencies.Transfer = func(context.Context, transfer.Request) (transfer.Result, error) {
				return transfer.Result{}, &transfer.Error{Stage: progress.StageCommit, Confirmed: 3, Cause: errors.New("commit")}
			}
		}, wantCode: ExitTransfer, wantStage: "commit", wantText: "confirmed: 3 bytes"},
		{name: "generic transfer", args: []string{"from", "in", "to", "out"}, mutate: func(dependencies *Dependencies) {
			dependencies.Transfer = func(context.Context, transfer.Request) (transfer.Result, error) {
				return transfer.Result{}, errors.New("copy")
			}
		}, wantCode: ExitTransfer, wantStage: "transfer", wantText: "copy"},
	} {
		t.Run(test.name, func(t *testing.T) {
			dependencies := transferDependencies(t, validOpen)
			if test.mutate != nil {
				test.mutate(&dependencies)
			}
			var stderr bytes.Buffer
			code := Execute(context.Background(), NewRoot(dependencies), test.args, io.Discard, &stderr)
			if code != test.wantCode || !strings.Contains(stderr.String(), "stage: "+test.wantStage) || !strings.Contains(stderr.String(), test.wantText) {
				t.Fatalf("code=%d stderr=%q", code, stderr.String())
			}
		})
	}

	t.Run("connected resource closes after peer failure", func(t *testing.T) {
		var closed atomic.Bool
		ready := make(chan struct{})
		dependencies := transferDependencies(t, func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			if value.Remote {
				<-ready
				return nil, errors.New("remote failure")
			}
			close(ready)
			return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path, Close: func() error { closed.Store(true); return nil }}, nil
		})
		if code := Execute(context.Background(), NewRoot(dependencies), []string{"from", "local", "to", "host:/remote"}, io.Discard, io.Discard); code != ExitConnection || !closed.Load() {
			t.Fatalf("code=%d closed=%v", code, closed.Load())
		}
	})

	t.Run("resource cleanup blocks success report", func(t *testing.T) {
		dependencies := transferDependencies(t, validOpen)
		dependencies.Transfer = func(context.Context, transfer.Request) (transfer.Result, error) {
			return transfer.Result{Bytes: 4}, nil
		}
		dependencies.Open = func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
			resource, _ := validOpen(context.Background(), value)
			resource.Close = func() error { return errors.New("close") }
			return resource, nil
		}
		var stdout, stderr bytes.Buffer
		code := Execute(context.Background(), NewRoot(dependencies), []string{"from", "in", "to", "out"}, &stdout, &stderr)
		if code != ExitTransfer || strings.Contains(stdout.String(), "result: success") || !strings.Contains(stderr.String(), "stage: cleanup") {
			t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
	})
}

func TestArchiveFailurePaths(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	if err := os.WriteFile(source, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	open := func(_ context.Context, value endpoint.Endpoint) (*Resource, error) {
		return &Resource{Endpoint: value, Backend: fsx.Local{}, Path: value.Path}, nil
	}
	args := []string{"from", source, "to", destination, "--archive"}

	t.Run("archive creation", func(t *testing.T) {
		dependencies := transferDependencies(t, open)
		dependencies.Archive = func(context.Context, fsx.Backend, string, string, string, selection.Selector, progress.Sink) (*archive.Artifact, error) {
			return nil, errors.New("archive")
		}
		var stderr bytes.Buffer
		if code := Execute(context.Background(), NewRoot(dependencies), args, io.Discard, &stderr); code != ExitTransfer || !strings.Contains(stderr.String(), "stage: archive") {
			t.Fatalf("code=%d stderr=%q", code, stderr.String())
		}
	})

	t.Run("artifact open", func(t *testing.T) {
		dependencies := transferDependencies(t, open)
		dependencies.OpenArtifact = func(string) (fsx.Backend, string, func() error, error) {
			return nil, "", nil, errors.New("artifact open")
		}
		var stderr bytes.Buffer
		if code := Execute(context.Background(), NewRoot(dependencies), args, io.Discard, &stderr); code != ExitTransfer || !strings.Contains(stderr.String(), "artifact open") {
			t.Fatalf("code=%d stderr=%q", code, stderr.String())
		}
	})

	t.Run("artifact and backend cleanup", func(t *testing.T) {
		artifactPath := filepath.Join(root, "cleanup.tar.gz")
		if err := os.WriteFile(artifactPath, []byte("artifact"), 0o600); err != nil {
			t.Fatal(err)
		}
		dependencies := transferDependencies(t, open)
		dependencies.Archive = func(context.Context, fsx.Backend, string, string, string, selection.Selector, progress.Sink) (*archive.Artifact, error) {
			return &archive.Artifact{Path: artifactPath, Name: "source.tar.gz"}, nil
		}
		dependencies.OpenArtifact = func(string) (fsx.Backend, string, func() error, error) {
			return fsx.Local{}, artifactPath, func() error { return errors.New("backend close") }, nil
		}
		dependencies.Transfer = func(context.Context, transfer.Request) (transfer.Result, error) {
			if err := os.Remove(artifactPath); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(artifactPath, 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(artifactPath, "child"), []byte("x"), 0o600); err != nil {
				t.Fatal(err)
			}
			return transfer.Result{Bytes: 8}, nil
		}
		var stdout, stderr bytes.Buffer
		code := Execute(context.Background(), NewRoot(dependencies), args, &stdout, &stderr)
		if code != ExitTransfer || strings.Contains(stdout.String(), "result: success") || !strings.Contains(stderr.String(), "backend close") {
			t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
	})
}

func TestDefaultDependenciesAndTerminalPrompt(t *testing.T) {
	originalHome, originalLoad, originalEmpty := userHomeDirectory, loadSSHConfig, emptySSHConfig
	originalUser, originalEnvironment := currentUser, environmentValue
	originalRoot, originalOpenSSH, originalDetect := openRootedPath, openSSHConnection, detectSSHPlatform
	originalTerminal, originalRead := terminalAttached, readTerminalSecret
	t.Cleanup(func() {
		userHomeDirectory, loadSSHConfig, emptySSHConfig = originalHome, originalLoad, originalEmpty
		currentUser, environmentValue = originalUser, originalEnvironment
		openRootedPath, openSSHConnection, detectSSHPlatform = originalRoot, originalOpenSSH, originalDetect
		terminalAttached, readTerminalSecret = originalTerminal, originalRead
	})

	userHomeDirectory = func() (string, error) { return t.TempDir(), nil }
	loadSSHConfig = func(string, string) (*sshx.Config, error) { return nil, os.ErrNotExist }
	currentUser = func() (*user.User, error) { return &user.User{Username: "local-user"}, nil }
	environmentValue = func(name string) string { return "value-for-" + name }
	dependencies, err := DefaultDependencies(nil, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	selected, err := dependencies.Select([]operation.SelectionRule{{Kind: operation.SelectionRegex, Value: "secret"}})
	if err != nil || selected.Include("secret.txt", false) {
		t.Fatalf("default selector=%v err=%v", selected, err)
	}
	localPath := filepath.Join(t.TempDir(), "file")
	local, err := dependencies.Open(context.Background(), endpoint.Endpoint{Path: localPath})
	if err != nil || local.Path != "file" {
		t.Fatalf("local=%+v err=%v", local, err)
	}
	_ = local.Close()
	backend, relative, closeBackend, err := dependencies.OpenArtifact(localPath)
	if err != nil || backend == nil || relative != "file" || closeBackend == nil {
		t.Fatalf("backend=%v relative=%q close=%v err=%v", backend, relative, closeBackend != nil, err)
	}
	_ = closeBackend()

	if _, err := originalOpenSSH(context.Background(), sshx.Factory{}, "host", "user"); err == nil {
		t.Fatal("expected default SSH wrapper failure")
	}
	openSSHConnection = func(_ context.Context, _ sshx.Factory, host, username string) (*sshx.Connection, error) {
		if host == "fail" {
			return nil, errors.New("dial")
		}
		return &sshx.Connection{Target: sshx.Target{Host: "canonical", User: username}}, nil
	}
	detectSSHPlatform = func(context.Context, sshx.Runner) (sshx.Platform, sshx.Archiver, error) {
		return sshx.Platform{OS: "linux", Arch: "amd64"}, sshx.Archiver{Name: "builtin", BuiltIn: true}, nil
	}
	remote, err := dependencies.Open(context.Background(), endpoint.Endpoint{Remote: true, Host: "alias", User: "user", Path: "/file"})
	if err != nil || remote.Endpoint.Host != "canonical" || remote.Endpoint.User != "user" {
		t.Fatalf("remote=%+v err=%v", remote, err)
	}
	_ = remote.Close()
	if _, err := dependencies.Open(context.Background(), endpoint.Endpoint{Remote: true, Host: "fail", Path: "/file"}); err == nil {
		t.Fatal("expected remote open failure")
	}
	detectSSHPlatform = func(context.Context, sshx.Runner) (sshx.Platform, sshx.Archiver, error) {
		return sshx.Platform{}, sshx.Archiver{}, errors.New("probe")
	}
	if _, err := dependencies.Open(context.Background(), endpoint.Endpoint{Remote: true, Host: "alias", Path: "/file"}); err == nil {
		t.Fatal("expected probe failure")
	}

	openRootedPath = func(string) (*fsx.RootedLocal, string, error) { return nil, "", errors.New("root") }
	if _, err := dependencies.Open(context.Background(), endpoint.Endpoint{Path: "local"}); err == nil {
		t.Fatal("expected local root failure")
	}
	if _, _, _, err := dependencies.OpenArtifact("artifact"); err == nil {
		t.Fatal("expected artifact root failure")
	}

	if _, err := terminalPrompt(nil, io.Discard)("Password: "); err == nil {
		t.Fatal("expected non-terminal prompt failure")
	}
	input, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	terminalAttached = func(int) bool { return true }
	readTerminalSecret = func(int) ([]byte, error) { return []byte("secret"), nil }
	var promptOutput bytes.Buffer
	secret, err := terminalPrompt(input, &promptOutput)("Password: ")
	if err != nil || string(secret) != "secret" || promptOutput.String() != "Password: \n" || !writerIsTerminal(input) || writerIsTerminal(&bytes.Buffer{}) {
		t.Fatalf("secret=%q err=%v output=%q", secret, err, promptOutput.String())
	}
	readTerminalSecret = func(int) ([]byte, error) { return nil, errors.New("read") }
	if _, err := terminalPrompt(input, io.Discard)("Password: "); err == nil {
		t.Fatal("expected terminal read failure")
	}
}

func TestTerminalConfirmation(t *testing.T) {
	originalTerminal := terminalAttached
	t.Cleanup(func() { terminalAttached = originalTerminal })
	terminalAttached = func(int) bool { return true }
	if accepted, err := terminalConfirmation(nil, io.Discard)(context.Background(), "question"); err == nil || accepted {
		t.Fatalf("missing terminal accepted=%v err=%v", accepted, err)
	}
	for _, test := range []struct {
		answer string
		want   bool
	}{{"y\n", true}, {" YES \n", true}, {"no\n", false}, {"y", true}} {
		name := filepath.Join(t.TempDir(), "answer")
		if err := os.WriteFile(name, []byte(test.answer), 0o600); err != nil {
			t.Fatal(err)
		}
		input, err := os.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		var output bytes.Buffer
		accepted, err := terminalConfirmation(input, &output)(context.Background(), "Deploy?")
		_ = input.Close()
		if err != nil || accepted != test.want || !strings.Contains(output.String(), "[y/N]") {
			t.Fatalf("answer=%q accepted=%v output=%q err=%v", test.answer, accepted, output.String(), err)
		}
	}
	name := filepath.Join(t.TempDir(), "closed")
	input, err := os.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	_ = input.Close()
	if _, err := terminalConfirmation(input, io.Discard)(context.Background(), "question"); err == nil {
		t.Fatal("expected closed input error")
	}
	input, err = os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := terminalConfirmation(input, io.Discard)(ctx, "question"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation: %v", err)
	}
}

func TestDefaultDependencyFailuresAndErrorHelpers(t *testing.T) {
	originalHome, originalLoad := userHomeDirectory, loadSSHConfig
	t.Cleanup(func() { userHomeDirectory, loadSSHConfig = originalHome, originalLoad })
	userHomeDirectory = func() (string, error) { return "", errors.New("home") }
	if _, err := DefaultDependencies(nil, io.Discard); err == nil {
		t.Fatal("expected home failure")
	}
	userHomeDirectory = func() (string, error) { return "/home", nil }
	loadSSHConfig = func(string, string) (*sshx.Config, error) { return nil, errors.New("config") }
	if _, err := DefaultDependencies(nil, io.Discard); err == nil {
		t.Fatal("expected config failure")
	}

	cause := errors.New("cause")
	commandErr := &commandError{cause: cause}
	openErr := &resourceOpenError{role: "source", cause: cause}
	if commandErr.Error() != "cause" || !errors.Is(commandErr, cause) || openErr.Error() != "open source: cause" || !errors.Is(openErr, cause) {
		t.Fatal("error wrappers are inconsistent")
	}
	if err := closeResources(nil, &Resource{}, &Resource{Close: func() error { return cause }}); !errors.Is(err, cause) {
		t.Fatalf("closeResources error=%v", err)
	}
	if got := transferCommandError(progress.StageCleanup, cause, 7); got.code != ExitTransfer || got.confirmed != 7 {
		t.Fatalf("transfer error=%+v", got)
	}
}

func transferDependencies(t *testing.T, open func(context.Context, endpoint.Endpoint) (*Resource, error)) Dependencies {
	t.Helper()
	registry := archive.DefaultRegistry()
	return Dependencies{
		Open: open,
		OpenArtifact: func(name string) (fsx.Backend, string, func() error, error) {
			return fsx.Local{}, name, func() error { return nil }, nil
		},
		Transfer: (transfer.Engine{Token: func() (string, error) { return "test", nil }}).Run,
		Archive:  archive.CreateSelected,
		Extract:  registry.Extract,
		Reporter: report.New,
		Terminal: func(io.Writer) bool { return false },
		TempDir:  t.TempDir(),
	}
}

type lstatBackend struct {
	fsx.Backend
	lstat func(string) (fs.FileInfo, error)
}

func (b lstatBackend) Lstat(name string) (fs.FileInfo, error) { return b.lstat(name) }

type fakeInfo struct {
	name      string
	size      int64
	directory bool
	mode      fs.FileMode
}

func (f fakeInfo) Name() string { return f.name }
func (f fakeInfo) Size() int64  { return f.size }
func (f fakeInfo) Mode() fs.FileMode {
	if f.mode != 0 {
		return f.mode
	}
	if f.directory {
		return fs.ModeDir | 0o700
	}
	return 0o600
}
func (f fakeInfo) ModTime() time.Time { return time.Unix(1, 0) }
func (f fakeInfo) IsDir() bool        { return f.directory }
func (f fakeInfo) Sys() any           { return nil }

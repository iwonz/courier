package app_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/app"
	"github.com/iwonz/courier/internal/archive"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/progress"
	"github.com/iwonz/courier/internal/report"
	"github.com/iwonz/courier/internal/selection"
	"github.com/iwonz/courier/internal/transfer"
)

func TestPathRouteAcceptance(t *testing.T) {
	for _, test := range []struct {
		name              string
		sourceRemote      bool
		destinationRemote bool
	}{
		{name: "local-to-local"},
		{name: "local-to-ssh", destinationRemote: true},
		{name: "ssh-to-local", sourceRemote: true},
		{name: "ssh-to-ssh", sourceRemote: true, destinationRemote: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			source := filepath.Join(root, "source tree")
			destination := filepath.Join(root, "destination tree")
			nested := filepath.Join(source, "nested directory")
			if err := os.MkdirAll(nested, 0o750); err != nil {
				t.Fatal(err)
			}
			name := "courier [#] 世界.txt"
			if err := os.WriteFile(filepath.Join(nested, name), []byte("accepted payload"), 0o640); err != nil {
				t.Fatal(err)
			}

			remotes := map[string]string{
				"source.test:/payload":     source,
				"destination.test:/result": destination,
			}
			dependencies := acceptanceDependencies(t, remotes)
			sourceArgument := source
			if test.sourceRemote {
				sourceArgument = "source.test:/payload"
			}
			destinationArgument := destination
			if test.destinationRemote {
				destinationArgument = "destination.test:/result"
			}

			var stdout, stderr bytes.Buffer
			code := app.Execute(context.Background(), app.NewRoot(dependencies), []string{"from", sourceArgument, "to", destinationArgument}, &stdout, &stderr)
			data, err := os.ReadFile(filepath.Join(destination, "nested directory", name))
			if code != app.ExitOK || err != nil || string(data) != "accepted payload" {
				t.Fatalf("code=%d data=%q err=%v stdout=%q stderr=%q", code, data, err, stdout.String(), stderr.String())
			}
			if !strings.Contains(stdout.String(), "result: success") || !strings.Contains(stderr.String(), "stage=complete") {
				t.Fatalf("stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
			assertNoCourierTemporaryPaths(t, root)
		})
	}
}

func TestArchiveExtractionAcceptance(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "release payload")
	if err := os.MkdirAll(filepath.Join(source, "nested"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "nested", "файл.txt"), []byte("archive payload"), 0o640); err != nil {
		t.Fatal(err)
	}
	archivePath := filepath.Join(root, "release payload.tar.gz")
	dependencies := acceptanceDependencies(t, nil)

	if code := app.Execute(context.Background(), app.NewRoot(dependencies), []string{"from", source, "to", archivePath, "--archive"}, io.Discard, io.Discard); code != app.ExitOK {
		t.Fatalf("archive exit code=%d", code)
	}
	if err := archive.Verify(archivePath); err != nil {
		t.Fatalf("verify archive: %v", err)
	}

	extractionRoot := filepath.Join(root, "extracted")
	if code := app.Execute(context.Background(), app.NewRoot(dependencies), []string{"from", archivePath, "to", extractionRoot, "--extract"}, io.Discard, io.Discard); code != app.ExitOK {
		t.Fatalf("extract exit code=%d", code)
	}
	data, err := os.ReadFile(filepath.Join(extractionRoot, "release payload", "nested", "файл.txt"))
	if err != nil || string(data) != "archive payload" {
		t.Fatalf("round-trip data=%q err=%v", data, err)
	}
	assertNoCourierTemporaryPaths(t, root)
}

func TestSafeFailureAcceptance(t *testing.T) {
	t.Run("identity is a successful no-op", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "same.txt")
		if err := os.WriteFile(path, []byte("unchanged"), 0o600); err != nil {
			t.Fatal(err)
		}
		var stdout bytes.Buffer
		code := app.Execute(context.Background(), app.NewRoot(acceptanceDependencies(t, nil)), []string{"from", path, "to", path}, &stdout, io.Discard)
		data, err := os.ReadFile(path)
		if code != app.ExitOK || err != nil || string(data) != "unchanged" || !strings.Contains(stdout.String(), "transferred: 0 bytes") {
			t.Fatalf("code=%d data=%q err=%v stdout=%q", code, data, err, stdout.String())
		}
	})

	t.Run("collision preserves final data", func(t *testing.T) {
		root := t.TempDir()
		source := filepath.Join(root, "source.txt")
		destination := filepath.Join(root, "destination.txt")
		if err := os.WriteFile(source, []byte("new"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(destination, []byte("keep"), 0o600); err != nil {
			t.Fatal(err)
		}
		var stderr bytes.Buffer
		code := app.Execute(context.Background(), app.NewRoot(acceptanceDependencies(t, nil)), []string{"from", source, "to", destination}, io.Discard, &stderr)
		data, err := os.ReadFile(destination)
		if code != app.ExitTransfer || err != nil || string(data) != "keep" || !strings.Contains(stderr.String(), "destination already exists") {
			t.Fatalf("code=%d data=%q err=%v stderr=%q", code, data, err, stderr.String())
		}
		assertNoCourierTemporaryPaths(t, root)
	})

	t.Run("interruption reports confirmed data and cleans staging", func(t *testing.T) {
		root := t.TempDir()
		source := filepath.Join(root, "large source.bin")
		destination := filepath.Join(root, "final.bin")
		if err := os.WriteFile(source, bytes.Repeat([]byte("x"), 256*1024), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "unrelated.txt"), []byte("keep"), 0o600); err != nil {
			t.Fatal(err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		dependencies := acceptanceDependencies(t, nil)
		dependencies.Transfer = func(runContext context.Context, request transfer.Request) (transfer.Result, error) {
			original := request.Progress
			request.Progress = func(event progress.Event) {
				if original != nil {
					original(event)
				}
				if event.Confirmed > 0 {
					cancel()
				}
			}
			return (transfer.Engine{Token: func() (string, error) { return "interrupted", nil }, BufferSize: 1024}).Run(runContext, request)
		}
		var stderr bytes.Buffer
		code := app.Execute(ctx, app.NewRoot(dependencies), []string{"from", source, "to", destination}, io.Discard, &stderr)
		_, destinationErr := os.Lstat(destination)
		unrelated, unrelatedErr := os.ReadFile(filepath.Join(root, "unrelated.txt"))
		if code != app.ExitInterrupted || !os.IsNotExist(destinationErr) || unrelatedErr != nil || string(unrelated) != "keep" || !strings.Contains(stderr.String(), "confirmed:") || strings.Contains(stderr.String(), "confirmed: 0 bytes") {
			t.Fatalf("code=%d destinationErr=%v unrelated=%q,%v stderr=%q", code, destinationErr, unrelated, unrelatedErr, stderr.String())
		}
		assertNoCourierTemporaryPaths(t, root)
	})
}

func acceptanceDependencies(t *testing.T, remotes map[string]string) app.Dependencies {
	t.Helper()
	temporaryDirectory := t.TempDir()
	registry := archive.DefaultRegistry()
	return app.Dependencies{
		Open: func(_ context.Context, value endpoint.Endpoint) (*app.Resource, error) {
			path := value.Path
			if value.Remote {
				var ok bool
				path, ok = remotes[value.Host+":"+value.Path]
				if !ok {
					return nil, fmt.Errorf("unmapped acceptance endpoint %s:%s", value.Host, value.Path)
				}
			}
			return &app.Resource{Endpoint: value, Backend: fsx.Local{}, Path: path, Close: func() error { return nil }}, nil
		},
		OpenArtifact: func(path string) (fsx.Backend, string, func() error, error) {
			return fsx.Local{}, path, func() error { return nil }, nil
		},
		Transfer: (transfer.Engine{Token: func() (string, error) { return "acceptance", nil }, BufferSize: 4 * 1024}).Run,
		Archive:  archive.CreateSelected,
		Extract:  registry.Extract,
		Select: func(rules []operation.SelectionRule) (selection.Selector, error) {
			return selection.Compile(rules, nil)
		},
		Reporter: report.New,
		Terminal: func(io.Writer) bool { return false },
		TempDir:  temporaryDirectory,
	}
}

func assertNoCourierTemporaryPaths(t *testing.T, root string) {
	t.Helper()
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(entry.Name(), ".courier-partial-") || strings.Contains(entry.Name(), ".courier-backup-") {
			t.Errorf("Courier temporary path remains: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

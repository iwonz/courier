package app

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/fsx"
)

func localPathResource(path string) *Resource {
	return &Resource{Endpoint: endpoint.Endpoint{Raw: path, Path: path}, Backend: fsx.Local{}, Path: path}
}

func TestPathPreflightAuthorizationAndCreation(t *testing.T) {
	root := t.TempDir()
	sourcePath := filepath.Join(root, "source")
	destinationPath := filepath.Join(root, "destination")
	source := localPathResource(sourcePath)
	destination := localPathResource(destinationPath)
	var prompts []string
	confirm := func(_ context.Context, question string) (bool, error) {
		if _, err := os.Stat(sourcePath); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("source changed before authorization completed: %v", err)
		}
		if _, err := os.Stat(destinationPath); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("destination changed before authorization completed: %v", err)
		}
		prompts = append(prompts, question)
		return true, nil
	}
	sourceInfo, destinationInfo, err := preflightTransferPaths(context.Background(), confirm, false, source, destination, false, false)
	if err != nil || !sourceInfo.IsDir() || !destinationInfo.IsDir() {
		t.Fatalf("source=%v destination=%v err=%v", sourceInfo, destinationInfo, err)
	}
	wantPrompts := []string{
		"Source directory " + sourcePath + " does not exist. Create it?",
		"Destination directory " + destinationPath + " does not exist. Create it?",
	}
	if !reflect.DeepEqual(prompts, wantPrompts) {
		t.Fatalf("prompts=%q", prompts)
	}
	for _, path := range []string{sourcePath, destinationPath} {
		info, statErr := os.Stat(path)
		if statErr != nil || !info.IsDir() || info.Mode().Perm() != 0o700 {
			t.Fatalf("path=%q info=%v err=%v", path, info, statErr)
		}
	}
}

func TestPathPreflightDeclineAndNoninteractive(t *testing.T) {
	for _, test := range []struct {
		name    string
		confirm DirectoryConfirm
		want    string
	}{
		{name: "noninteractive", want: "--force-source-creation"},
		{name: "unavailable terminal", confirm: func(context.Context, string) (bool, error) { return false, errDirectoryConfirmationUnavailable }, want: "--force-source-creation"},
		{name: "declined", confirm: func(context.Context, string) (bool, error) { return false, nil }, want: "was not created"},
		{name: "confirmation error", confirm: func(context.Context, string) (bool, error) { return false, io.ErrUnexpectedEOF }, want: io.ErrUnexpectedEOF.Error()},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "missing")
			_, err := preflightSinglePath(context.Background(), test.confirm, false, localPathResource(path), "Destination", pathDirectory)
			if err == nil || !strings.Contains(err.Error(), test.want) || !isPathPreflightError(err) && test.name != "confirmation error" {
				t.Fatalf("error=%v", err)
			}
			if _, statErr := os.Stat(path); !errors.Is(statErr, fs.ErrNotExist) {
				t.Fatalf("path was created: %v", statErr)
			}
		})
	}

	t.Run("all confirmations precede mutation", func(t *testing.T) {
		root := t.TempDir()
		sourcePath := filepath.Join(root, "source")
		destinationPath := filepath.Join(root, "destination")
		calls := 0
		_, _, err := preflightTransferPaths(context.Background(), func(context.Context, string) (bool, error) {
			calls++
			return calls == 1, nil
		}, false, localPathResource(sourcePath), localPathResource(destinationPath), false, false)
		if err == nil || calls != 2 {
			t.Fatalf("calls=%d err=%v", calls, err)
		}
		for _, path := range []string{sourcePath, destinationPath} {
			if _, statErr := os.Stat(path); !errors.Is(statErr, fs.ErrNotExist) {
				t.Fatalf("path %q changed: %v", path, statErr)
			}
		}
	})
}

func TestPathPreflightRolesAndRaces(t *testing.T) {
	t.Run("exact file destination remains absent", func(t *testing.T) {
		root := t.TempDir()
		sourcePath := filepath.Join(root, "report.pdf")
		destinationPath := filepath.Join(root, "renamed.pdf")
		if err := os.WriteFile(sourcePath, []byte("report"), 0o600); err != nil {
			t.Fatal(err)
		}
		called := false
		sourceInfo, destinationInfo, err := preflightTransferPaths(context.Background(), func(context.Context, string) (bool, error) {
			called = true
			return true, nil
		}, false, localPathResource(sourcePath), localPathResource(destinationPath), false, false)
		if err != nil || !sourceInfo.Mode().IsRegular() || destinationInfo != nil || called {
			t.Fatalf("source=%v destination=%v called=%v err=%v", sourceInfo, destinationInfo, called, err)
		}
		if _, statErr := os.Stat(destinationPath); !errors.Is(statErr, fs.ErrNotExist) {
			t.Fatalf("exact destination changed: %v", statErr)
		}
	})

	t.Run("forced directory creation", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "nested", "directory")
		info, err := preflightSinglePath(context.Background(), nil, true, localPathResource(path), "Source", pathFlexibleSource)
		if err != nil || !info.IsDir() {
			t.Fatalf("info=%v err=%v", info, err)
		}
	})

	t.Run("mock SSH directory creation", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "remote", "directory")
		mode := fs.FileMode(0)
		backend := sequenceBackend{Backend: fsx.Local{}, mkdir: func(name string, requested fs.FileMode) error {
			mode = requested
			return os.MkdirAll(name, requested)
		}}
		resource := &Resource{
			Endpoint: endpoint.Endpoint{Raw: "courier@host:/remote/directory", Remote: true, Host: "host", User: "courier", Path: "/remote/directory"},
			Backend:  backend,
			Path:     path,
		}
		info, err := preflightSinglePath(context.Background(), nil, true, resource, "Destination", pathDirectory)
		if err != nil || !info.IsDir() || mode.Perm() != 0o700 {
			t.Fatalf("info=%v mode=%#o err=%v", info, mode.Perm(), err)
		}
	})

	t.Run("concurrent file appearance", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "race")
		_, err := preflightSinglePath(context.Background(), func(context.Context, string) (bool, error) {
			return true, os.WriteFile(path, []byte("file"), 0o600)
		}, false, localPathResource(path), "Destination", pathDirectory)
		if err == nil || !isPathPreflightError(err) || !strings.Contains(err.Error(), "not a directory") {
			t.Fatalf("error=%v", err)
		}
	})

	t.Run("concurrent directory appearance", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "race-directory")
		backend := sequenceBackend{Backend: fsx.Local{}}
		backend.mkdir = func(name string, mode fs.FileMode) error {
			if err := os.MkdirAll(name, mode); err != nil {
				return err
			}
			return fs.ErrExist
		}
		resource := localPathResource(path)
		resource.Backend = backend
		info, err := preflightSinglePath(context.Background(), nil, true, resource, "Destination", pathDirectory)
		if err != nil || !info.IsDir() {
			t.Fatalf("info=%v err=%v", info, err)
		}
	})

	t.Run("revalidation failures", func(t *testing.T) {
		for _, test := range []struct {
			name  string
			lstat func(int) (fs.FileInfo, error)
			want  string
		}{
			{name: "stat", lstat: func(call int) (fs.FileInfo, error) {
				if call == 1 {
					return nil, fs.ErrNotExist
				}
				return nil, io.ErrClosedPipe
			}, want: "revalidate source directory"},
			{name: "type", lstat: func(call int) (fs.FileInfo, error) {
				if call == 1 {
					return nil, fs.ErrNotExist
				}
				return fakeInfo{name: "file"}, nil
			}, want: "not a directory"},
		} {
			t.Run(test.name, func(t *testing.T) {
				calls := 0
				backend := sequenceBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) {
					calls++
					return test.lstat(calls)
				}, mkdir: func(string, fs.FileMode) error { return nil }}
				resource := &Resource{Endpoint: endpoint.Endpoint{Raw: "mock", Path: "mock"}, Backend: backend, Path: "mock"}
				if _, err := preflightSinglePath(context.Background(), nil, true, resource, "Source", pathFlexibleSource); err == nil || !strings.Contains(err.Error(), test.want) {
					t.Fatalf("error=%v", err)
				}
			})
		}
	})

	t.Run("wrong existing types", func(t *testing.T) {
		root := t.TempDir()
		file := filepath.Join(root, "file")
		if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := preflightSinglePath(context.Background(), nil, true, localPathResource(file), "Destination", pathDirectory); err == nil || !strings.Contains(err.Error(), "not a directory") {
			t.Fatalf("directory role error=%v", err)
		}
		if _, err := preflightSinglePath(context.Background(), nil, true, localPathResource(root), "Source", pathRegularFile); err == nil || !strings.Contains(err.Error(), "not a regular file") {
			t.Fatalf("file role error=%v", err)
		}
		missing := filepath.Join(root, "archive.tar.gz")
		if _, err := preflightSinglePath(context.Background(), nil, true, localPathResource(missing), "Source", pathRequiredFile); err == nil || !strings.Contains(err.Error(), "does not exist") {
			t.Fatalf("required file error=%v", err)
		}
	})
}

func TestPathPreflightFailureKeepsEarlierCreation(t *testing.T) {
	root := t.TempDir()
	sourcePath := filepath.Join(root, "source")
	destinationPath := filepath.Join(root, "destination")
	destination := localPathResource(destinationPath)
	destination.Backend = mkdirErrorBackend{Backend: destination.Backend, err: io.ErrClosedPipe}
	_, _, err := preflightTransferPaths(context.Background(), nil, true, localPathResource(sourcePath), destination, false, false)
	if err == nil || !strings.Contains(err.Error(), "create destination directory") {
		t.Fatalf("error=%v", err)
	}
	if info, statErr := os.Stat(sourcePath); statErr != nil || !info.IsDir() {
		t.Fatalf("source was rolled back: info=%v err=%v", info, statErr)
	}
}

type mkdirErrorBackend struct {
	fsx.Backend
	err error
}

func (backend mkdirErrorBackend) MkdirAll(string, fs.FileMode) error { return backend.err }

type sequenceBackend struct {
	fsx.Backend
	lstat func(string) (fs.FileInfo, error)
	mkdir func(string, fs.FileMode) error
}

func (backend sequenceBackend) Lstat(name string) (fs.FileInfo, error) {
	if backend.lstat != nil {
		return backend.lstat(name)
	}
	return backend.Backend.Lstat(name)
}

func (backend sequenceBackend) MkdirAll(name string, mode fs.FileMode) error {
	return backend.mkdir(name, mode)
}

func TestPathPreflightContextAndInspectionFailures(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing")
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := preflightSinglePath(canceled, nil, true, localPathResource(path), "Source", pathFlexibleSource); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation=%v", err)
	}
	resource := localPathResource(path)
	resource.Backend = lstatBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return nil, io.ErrClosedPipe }}
	if _, err := preflightSinglePath(context.Background(), nil, true, resource, "Source", pathFlexibleSource); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("inspection=%v", err)
	}
	if err := inspectPath(nil); err == nil {
		t.Fatal("nil resource accepted")
	}
}

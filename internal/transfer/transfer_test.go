package transfer

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/progress"
)

func TestLocalDirectorySynchronization(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	if err := os.MkdirAll(filepath.Join(source, "nested"), 0o750); err != nil {
		t.Fatal(err)
	}
	stamp := time.Unix(1_700_000_000, 0)
	file := filepath.Join(source, "nested", "file.txt")
	if err := os.WriteFile(file, []byte("hello courier"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(file, stamp, stamp); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("nested/file.txt", filepath.Join(source, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(destination, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(destination, "stale"), []byte("remove"), 0o600); err != nil {
		t.Fatal(err)
	}
	current := time.Unix(100, 0)
	var events []progress.Event
	engine := Engine{Token: func() (string, error) { return "fixed", nil }, Now: func() time.Time { current = current.Add(time.Second); return current }, BufferSize: 3}
	result, err := engine.Run(context.Background(), Request{SourceFS: fsx.Local{}, SourcePath: source, DestinationFS: fsx.Local{}, Destination: destination, Progress: func(event progress.Event) { events = append(events, event) }})
	if err != nil {
		t.Fatal(err)
	}
	if result.Bytes != int64(len("hello courier")) || result.Destination != destination || result.Elapsed <= 0 {
		t.Fatalf("result=%+v", result)
	}
	data, err := os.ReadFile(filepath.Join(destination, "nested", "file.txt"))
	if err != nil || string(data) != "hello courier" {
		t.Fatalf("data=%q err=%v", data, err)
	}
	info, err := os.Stat(filepath.Join(destination, "nested", "file.txt"))
	if err != nil || info.Mode().Perm() != 0o640 || !info.ModTime().Equal(stamp) {
		t.Fatalf("info=%v err=%v", info, err)
	}
	if target, err := os.Readlink(filepath.Join(destination, "link")); err != nil || target != "nested/file.txt" {
		t.Fatalf("link=%q err=%v", target, err)
	}
	if _, err := os.Stat(filepath.Join(destination, "stale")); !os.IsNotExist(err) {
		t.Fatalf("stale destination remains: %v", err)
	}
	if _, err := os.Stat(source); err != nil {
		t.Fatalf("source changed: %v", err)
	}
	if events[0].Stage != progress.StagePreflight || events[len(events)-1].Stage != progress.StageComplete {
		t.Fatalf("events=%v", events)
	}
	assertNoTemporaryPaths(t, root)
}

func TestLocalFileAndCancellation(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	destination := filepath.Join(root, "destination")
	if err := os.WriteFile(source, []byte(strings.Repeat("x", 128*1024)), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(destination, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	engine := Engine{Token: func() (string, error) { return "cancel", nil }, BufferSize: 2}
	_, err := engine.Run(ctx, Request{SourceFS: fsx.Local{}, SourcePath: source, DestinationFS: fsx.Local{}, Destination: destination, Progress: func(event progress.Event) {
		if event.Current > 0 {
			cancel()
		}
	}})
	var transferErr *Error
	if !errors.As(err, &transferErr) || !errors.Is(err, context.Canceled) || transferErr.Stage != progress.StageTransfer || transferErr.Confirmed == 0 {
		t.Fatalf("error=%v", err)
	}
	data, readErr := os.ReadFile(destination)
	if readErr != nil || string(data) != "old" {
		t.Fatalf("destination=%q err=%v", data, readErr)
	}
	assertNoTemporaryPaths(t, root)

	result, err := (Engine{Token: func() (string, error) { return "retry", nil }}).Run(context.Background(), Request{SourceFS: fsx.Local{}, SourcePath: source, DestinationFS: fsx.Local{}, Destination: destination})
	if err != nil || result.Bytes != 128*1024 {
		t.Fatalf("retry result=%+v err=%v", result, err)
	}
}

func TestPreflightFailures(t *testing.T) {
	engine := Engine{}
	if _, err := engine.Run(context.Background(), Request{}); err == nil {
		t.Fatal("expected invalid request")
	}
	backend := fsx.Local{}
	if _, err := engine.Run(context.Background(), Request{SourceFS: backend, SourcePath: filepath.Join(t.TempDir(), "missing"), DestinationFS: backend, Destination: "out"}); err == nil {
		t.Fatal("expected missing source")
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := engine.Run(canceled, Request{SourceFS: backend, SourcePath: "source", DestinationFS: backend, Destination: "out"}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
	root := t.TempDir()
	source := filepath.Join(root, "file")
	if err := os.WriteFile(source, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := (Engine{Token: func() (string, error) { return "", errors.New("token") }}).Run(context.Background(), Request{SourceFS: backend, SourcePath: source, DestinationFS: backend, Destination: filepath.Join(root, "out")}); err == nil {
		t.Fatal("expected token error")
	}
	if _, err := (Engine{}).Run(context.Background(), Request{SourceFS: backend, SourcePath: source, DestinationFS: backend, Destination: filepath.Join(root, "default-token")}); err != nil {
		t.Fatalf("default token transfer: %v", err)
	}
}

func TestTransferErrorFormatting(t *testing.T) {
	cause := errors.New("boom")
	err := &Error{Stage: progress.StageCommit, Confirmed: 7, Cause: cause}
	if !strings.Contains(err.Error(), "commit failed after 7") || !errors.Is(err, cause) {
		t.Fatalf("error=%v", err)
	}
}

func TestScanSpecialObject(t *testing.T) {
	backend := stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return fakeInfo{mode: fs.ModeNamedPipe}, nil }}
	if _, err := scan(context.Background(), backend, "pipe"); err == nil {
		t.Fatal("expected unsupported object")
	}
}

func TestScanAndCopyNodeFailures(t *testing.T) {
	tracker := progress.New(0, nil, nil)
	directory := fakeInfo{mode: fs.ModeDir | 0o755}
	link := fakeInfo{mode: fs.ModeSymlink}
	readFailure := stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return directory, nil }, readDir: func(string) ([]fs.DirEntry, error) { return nil, errors.New("read-dir") }}
	if _, err := scan(context.Background(), readFailure, "dir"); err == nil {
		t.Fatal("expected scan read-dir error")
	}
	entry := fakeEntry{name: "child"}
	recursiveFailure := stubBackend{Backend: fsx.Local{}, lstat: func(name string) (fs.FileInfo, error) {
		if name == "dir" {
			return directory, nil
		}
		return nil, errors.New("child")
	}, readDir: func(string) ([]fs.DirEntry, error) { return []fs.DirEntry{entry}, nil }, join: filepath.Join}
	if _, err := scan(context.Background(), recursiveFailure, "dir"); err == nil {
		t.Fatal("expected recursive scan error")
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := copyNode(canceled, recursiveFailure, "dir", fsx.Local{}, "out", nil, tracker); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	lstatFailure := stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return nil, errors.New("lstat") }}
	if err := copyNode(context.Background(), lstatFailure, "x", fsx.Local{}, "out", nil, tracker); err == nil {
		t.Fatal("expected lstat error")
	}
	readlinkFailure := stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return link, nil }, readlink: func(string) (string, error) { return "", errors.New("readlink") }}
	if err := copyNode(context.Background(), readlinkFailure, "x", fsx.Local{}, "out", nil, tracker); err == nil {
		t.Fatal("expected readlink error")
	}
	special := stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return fakeInfo{mode: fs.ModeNamedPipe}, nil }}
	if err := copyNode(context.Background(), special, "x", fsx.Local{}, "out", nil, tracker); err == nil {
		t.Fatal("expected special object error")
	}
	directorySource := stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return directory, nil }, readDir: func(string) ([]fs.DirEntry, error) { return nil, nil }}
	for _, destination := range []stubBackend{
		{Backend: fsx.Local{}, mkdirAll: func(string, fs.FileMode) error { return errors.New("mkdir") }},
		{Backend: fsx.Local{}, chmod: func(string, fs.FileMode) error { return errors.New("chmod") }},
		{Backend: fsx.Local{}, chtimes: func(string, time.Time, time.Time) error { return errors.New("chtimes") }},
	} {
		if err := copyNode(context.Background(), directorySource, "x", destination, filepath.Join(t.TempDir(), "out"), nil, tracker); err == nil {
			t.Fatal("expected directory destination error")
		}
	}
	copyReadFailure := stubBackend{Backend: directorySource, readDir: func(string) ([]fs.DirEntry, error) { return nil, errors.New("read-dir") }}
	if err := copyNode(context.Background(), copyReadFailure, "x", fsx.Local{}, t.TempDir(), nil, tracker); err == nil {
		t.Fatal("expected copy read-dir error")
	}
	childFailure := stubBackend{Backend: directorySource, readDir: func(string) ([]fs.DirEntry, error) { return []fs.DirEntry{entry}, nil }, join: filepath.Join, lstat: func(name string) (fs.FileInfo, error) {
		if name == "x" {
			return directory, nil
		}
		return nil, errors.New("child")
	}}
	if err := copyNode(context.Background(), childFailure, "x", fsx.Local{}, t.TempDir(), nil, tracker); err == nil {
		t.Fatal("expected child copy error")
	}
}

func TestCopyFileFailures(t *testing.T) {
	tracker := progress.New(3, nil, nil)
	info := fakeInfo{mode: 0o640, size: 3}
	openFailure := stubBackend{Backend: fsx.Local{}, open: func(string) (io.ReadCloser, error) { return nil, errors.New("open") }}
	if err := copyFile(context.Background(), openFailure, "x", fsx.Local{}, "out", info, make([]byte, 2), tracker); err == nil {
		t.Fatal("expected open error")
	}
	source := stubBackend{Backend: fsx.Local{}, open: func(string) (io.ReadCloser, error) { return io.NopCloser(bytes.NewBufferString("abc")), nil }}
	for _, writer := range []fsx.Writable{
		&fakeWriter{writeErr: errors.New("write")},
		&fakeWriter{syncErr: errors.New("sync")},
		&fakeWriter{closeErr: errors.New("close")},
	} {
		destination := stubBackend{Backend: fsx.Local{}, create: func(string, fs.FileMode) (fsx.Writable, error) { return writer, nil }}
		if err := copyFile(context.Background(), source, "x", destination, "out", info, make([]byte, 2), tracker); err == nil {
			t.Fatal("expected writer error")
		}
	}
	createFailure := stubBackend{Backend: fsx.Local{}, create: func(string, fs.FileMode) (fsx.Writable, error) { return nil, errors.New("create") }}
	if err := copyFile(context.Background(), source, "x", createFailure, "out", info, make([]byte, 2), tracker); err == nil {
		t.Fatal("expected create error")
	}
	for _, destination := range []stubBackend{
		{Backend: fsx.Local{}, create: func(string, fs.FileMode) (fsx.Writable, error) { return &fakeWriter{}, nil }, chmod: func(string, fs.FileMode) error { return errors.New("chmod") }},
		{Backend: fsx.Local{}, create: func(string, fs.FileMode) (fsx.Writable, error) { return &fakeWriter{}, nil }, chtimes: func(string, time.Time, time.Time) error { return errors.New("chtimes") }},
	} {
		if err := copyFile(context.Background(), source, "x", destination, "out", info, make([]byte, 2), tracker); err == nil {
			t.Fatal("expected metadata error")
		}
	}
}

func TestRunDestinationFailures(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.WriteFile(source, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, destination := range []stubBackend{
		{Backend: fsx.Local{}, removeAll: func(string) error { return errors.New("remove") }},
		{Backend: fsx.Local{}, mkdirAll: func(string, fs.FileMode) error { return errors.New("mkdir") }},
		{Backend: fsx.Local{}, lstat: func(name string) (fs.FileInfo, error) {
			if strings.Contains(name, ".courier-partial-") {
				return os.Stat(name)
			}
			return nil, errors.New("destination stat")
		}},
	} {
		_, err := (Engine{Token: func() (string, error) { return "fault", nil }}).Run(context.Background(), Request{SourceFS: fsx.Local{}, SourcePath: source, DestinationFS: destination, Destination: filepath.Join(root, "out")})
		if err == nil {
			t.Fatal("expected destination error")
		}
	}
}

func TestCommitFailures(t *testing.T) {
	exists := fakeInfo{mode: 0o600}
	for _, test := range []struct {
		name    string
		backend stubBackend
	}{
		{"destination lstat", stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return nil, errors.New("lstat") }}},
		{"remove backup", stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return exists, nil }, removeAll: func(string) error { return errors.New("remove") }}},
		{"move old", stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return exists, nil }, rename: func(string, string) error { return errors.New("rename") }}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := commit(test.backend, "stage", "destination", "token"); err == nil {
				t.Fatal("expected commit error")
			}
		})
	}
	for _, restoreFails := range []bool{false, true} {
		calls := 0
		backend := stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return exists, nil }, rename: func(string, string) error {
			calls++
			if calls == 2 || (calls == 3 && restoreFails) {
				return errors.New("rename")
			}
			return nil
		}}
		if err := commit(backend, "stage", "destination", "token"); err == nil {
			t.Fatal("expected stage/restore error")
		}
	}
	calls := 0
	removeCommitted := stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return exists, nil }, removeAll: func(string) error {
		calls++
		if calls == 2 {
			return errors.New("remove committed")
		}
		return nil
	}, rename: func(string, string) error { return nil }}
	if err := commit(removeCommitted, "stage", "destination", "token"); err == nil {
		t.Fatal("expected committed backup error")
	}
	absent := stubBackend{Backend: fsx.Local{}, lstat: func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }, rename: func(string, string) error { return errors.New("rename absent") }}
	if err := commit(absent, "stage", "destination", "token"); err == nil {
		t.Fatal("expected absent rename error")
	}
}

func TestRandomToken(t *testing.T) {
	token, err := randomToken()
	if err != nil || len(token) != 16 {
		t.Fatalf("token=%q err=%v", token, err)
	}
	original := randomRead
	t.Cleanup(func() { randomRead = original })
	randomRead = func([]byte) (int, error) { return 0, errors.New("random") }
	if _, err := randomToken(); err == nil {
		t.Fatal("expected random error")
	}
}

func assertNoTemporaryPaths(t *testing.T, root string) {
	t.Helper()
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.Contains(entry.Name(), ".courier-partial-") || strings.Contains(entry.Name(), ".courier-backup-") {
			t.Errorf("temporary path remains: %s", entry.Name())
		}
	}
}

type fakeInfo struct {
	mode fs.FileMode
	size int64
}

func (f fakeInfo) Name() string       { return "fake" }
func (f fakeInfo) Size() int64        { return f.size }
func (f fakeInfo) Mode() fs.FileMode  { return f.mode }
func (f fakeInfo) ModTime() time.Time { return time.Time{} }
func (f fakeInfo) IsDir() bool        { return f.mode.IsDir() }
func (f fakeInfo) Sys() any           { return nil }

type stubBackend struct {
	fsx.Backend
	lstat     func(string) (fs.FileInfo, error)
	readDir   func(string) ([]fs.DirEntry, error)
	open      func(string) (io.ReadCloser, error)
	create    func(string, fs.FileMode) (fsx.Writable, error)
	mkdirAll  func(string, fs.FileMode) error
	removeAll func(string) error
	rename    func(string, string) error
	chmod     func(string, fs.FileMode) error
	chtimes   func(string, time.Time, time.Time) error
	readlink  func(string) (string, error)
	join      func(...string) string
}

func (s stubBackend) Lstat(name string) (fs.FileInfo, error) {
	if s.lstat != nil {
		return s.lstat(name)
	}
	return s.Backend.Lstat(name)
}

func (s stubBackend) ReadDir(name string) ([]fs.DirEntry, error) {
	if s.readDir != nil {
		return s.readDir(name)
	}
	return s.Backend.ReadDir(name)
}
func (s stubBackend) Open(name string) (io.ReadCloser, error) {
	if s.open != nil {
		return s.open(name)
	}
	return s.Backend.Open(name)
}
func (s stubBackend) Create(name string, mode fs.FileMode) (fsx.Writable, error) {
	if s.create != nil {
		return s.create(name, mode)
	}
	return s.Backend.Create(name, mode)
}
func (s stubBackend) MkdirAll(name string, mode fs.FileMode) error {
	if s.mkdirAll != nil {
		return s.mkdirAll(name, mode)
	}
	return s.Backend.MkdirAll(name, mode)
}
func (s stubBackend) RemoveAll(name string) error {
	if s.removeAll != nil {
		return s.removeAll(name)
	}
	return s.Backend.RemoveAll(name)
}
func (s stubBackend) Rename(oldPath, newPath string) error {
	if s.rename != nil {
		return s.rename(oldPath, newPath)
	}
	return s.Backend.Rename(oldPath, newPath)
}
func (s stubBackend) Chmod(name string, mode fs.FileMode) error {
	if s.chmod != nil {
		return s.chmod(name, mode)
	}
	return s.Backend.Chmod(name, mode)
}
func (s stubBackend) Chtimes(name string, atime, mtime time.Time) error {
	if s.chtimes != nil {
		return s.chtimes(name, atime, mtime)
	}
	return s.Backend.Chtimes(name, atime, mtime)
}
func (s stubBackend) Readlink(name string) (string, error) {
	if s.readlink != nil {
		return s.readlink(name)
	}
	return s.Backend.Readlink(name)
}
func (s stubBackend) Join(parts ...string) string {
	if s.join != nil {
		return s.join(parts...)
	}
	return s.Backend.Join(parts...)
}

type fakeEntry struct{ name string }

func (f fakeEntry) Name() string             { return f.name }
func (fakeEntry) IsDir() bool                { return false }
func (fakeEntry) Type() fs.FileMode          { return 0 }
func (fakeEntry) Info() (fs.FileInfo, error) { return fakeInfo{}, nil }

type fakeWriter struct {
	bytes.Buffer
	writeErr error
	syncErr  error
	closeErr error
}

func (w *fakeWriter) Write(data []byte) (int, error) {
	if w.writeErr != nil {
		return 0, w.writeErr
	}
	return w.Buffer.Write(data)
}
func (w *fakeWriter) Sync() error  { return w.syncErr }
func (w *fakeWriter) Close() error { return w.closeErr }

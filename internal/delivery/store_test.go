package delivery

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestStoreRegistryAndHistory(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "state")
	store, err := OpenStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	if store.Directory() != directory {
		t.Fatalf("directory=%q", store.Directory())
	}
	if runtime.GOOS != "windows" {
		info, err := os.Stat(directory)
		if err != nil || info.Mode().Perm() != 0o700 {
			t.Fatalf("directory mode=%v err=%v", info.Mode(), err)
		}
	}
	if err := os.WriteFile(filepath.Join(directory, ".registry-interrupted.tmp"), []byte("partial"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "unrelated"), []byte("ignored"), 0o600); err != nil {
		t.Fatal(err)
	}
	empty, err := store.Load(context.Background())
	if err != nil || empty.SchemaVersion != SchemaVersion || empty.Revision != 0 || len(empty.Servers)+len(empty.Deliveries)+len(empty.Tombstones)+len(empty.OwnedTemps) != 0 {
		t.Fatalf("empty=%+v err=%v", empty, err)
	}

	next, err := store.Update(context.Background(), 0, func(registry *Registry) error {
		if err := registry.RegisterServer(validServer()); err != nil {
			return err
		}
		return registry.RegisterDelivery(validDelivery())
	})
	if err != nil || next.Revision != 1 || len(next.Servers) != 1 || len(next.Deliveries) != 1 {
		t.Fatalf("next=%+v err=%v", next, err)
	}
	loaded, err := store.Load(context.Background())
	if err != nil || !reflectSnapshots(next, loaded) {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
	registryPath := filepath.Join(directory, registryFilename(1))
	if runtime.GOOS != "windows" {
		info, err := os.Stat(registryPath)
		if err != nil || info.Mode().Perm() != 0o600 {
			t.Fatalf("registry mode=%v err=%v", info.Mode(), err)
		}
	}
	if _, err := store.Update(context.Background(), 0, func(*Registry) error { return nil }); !errors.Is(err, ErrRevisionConflict) {
		t.Fatalf("revision=%v", err)
	}
	if _, err := store.Update(context.Background(), 1, nil); !errors.Is(err, ErrInvalid) {
		t.Fatalf("nil mutation=%v", err)
	}
	mutateErr := errors.New("mutate")
	if _, err := store.Update(context.Background(), 1, func(*Registry) error { return mutateErr }); !errors.Is(err, mutateErr) {
		t.Fatalf("mutation=%v", err)
	}
	if _, err := store.Update(context.Background(), 1, func(registry *Registry) error { registry.snapshot.SchemaVersion = 2; return nil }); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid mutation=%v", err)
	}

	events := []HistoryEvent{
		{ID: fourthID, TargetID: deliveryID, Kind: HistoryActivated, At: testTime.Add(time.Second), Counters: CounterSnapshot{Read: 2}},
		{ID: thirdID, TargetID: deliveryID, Kind: HistoryRegistered, At: testTime, Counters: CounterSnapshot{Read: 1}},
		{ID: ID("00000000-0000-4000-8000-000000000005"), TargetID: deliveryID, Kind: HistoryStopped, At: testTime.Add(time.Second), Counters: CounterSnapshot{Read: 3}},
	}
	for _, event := range events {
		if err := store.AppendHistory(context.Background(), event); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.AppendHistory(context.Background(), events[0]); !errors.Is(err, ErrRevisionConflict) {
		t.Fatalf("duplicate history=%v", err)
	}
	history, err := store.History(context.Background())
	if err != nil || len(history) != 3 || history[0].ID != thirdID || history[1].ID != fourthID || history[2].ID != events[2].ID {
		t.Fatalf("history=%+v err=%v", history, err)
	}
	if err := store.AppendHistory(context.Background(), HistoryEvent{}); !errors.Is(err, ErrInvalid) {
		t.Fatalf("invalid history=%v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := store.Load(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("load cancellation=%v", err)
	}
	if err := store.AppendHistory(ctx, events[0]); !errors.Is(err, context.Canceled) {
		t.Fatalf("append cancellation=%v", err)
	}
	if _, err := store.History(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("history cancellation=%v", err)
	}
	if _, err := store.Update(ctx, 1, func(*Registry) error { return nil }); !errors.Is(err, context.Canceled) {
		t.Fatalf("update cancellation=%v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Load(context.Background()); !errors.Is(err, ErrInvalid) {
		t.Fatalf("closed load=%v", err)
	}
}

func TestConcurrentStoreRevisionCommit(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "state")
	first, err := OpenStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	second, err := OpenStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	ready := make(chan struct{}, 2)
	release := make(chan struct{})
	mutate := func(registry *Registry) error {
		ready <- struct{}{}
		<-release
		return registry.RegisterServer(validServer())
	}
	results := make(chan error, 2)
	go func() { _, err := first.Update(context.Background(), 0, mutate); results <- err }()
	go func() { _, err := second.Update(context.Background(), 0, mutate); results <- err }()
	<-ready
	<-ready
	close(release)
	errorsSeen := []error{<-results, <-results}
	successes, conflicts := 0, 0
	for _, err := range errorsSeen {
		if err == nil {
			successes++
		} else if errors.Is(err, ErrRevisionConflict) {
			conflicts++
		} else {
			t.Fatal(err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("results=%v", errorsSeen)
	}
}

func TestStoreCorruptionFailsClosed(t *testing.T) {
	write := func(t *testing.T, name string, data []byte, mode fs.FileMode) *Store {
		t.Helper()
		directory := filepath.Join(t.TempDir(), "state")
		store, err := OpenStore(directory)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = store.Close() })
		if err := os.WriteFile(filepath.Join(directory, name), data, mode); err != nil {
			t.Fatal(err)
		}
		return store
	}
	valid := EmptySnapshot()
	valid.Revision = 1
	validData, err := encodeJSON(valid)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name, file string
		data       []byte
		mode       fs.FileMode
	}{
		{"overflow filename", "registry-99999999999999999999.json", validData, 0o600},
		{"invalid json", registryFilename(1), []byte("{"), 0o600},
		{"unknown field", registryFilename(1), []byte(`{"schemaVersion":1,"revision":1,"servers":[],"deliveries":[],"tombstones":[],"ownedTemps":[],"secret":"x"}`), 0o600},
		{"trailing data", registryFilename(1), append(append([]byte(nil), validData...), []byte("{}")...), 0o600},
		{"revision mismatch", registryFilename(2), validData, 0o600},
		{"invalid snapshot", registryFilename(1), []byte(`{"schemaVersion":2,"revision":1,"servers":[],"deliveries":[],"tombstones":[],"ownedTemps":[]}`), 0o600},
		{"insecure file", registryFilename(1), validData, 0o644},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := write(t, test.file, test.data, test.mode)
			if _, err := store.Load(context.Background()); err == nil && (runtime.GOOS != "windows" || test.name != "insecure file") {
				t.Fatal("corruption accepted")
			}
		})
	}

	t.Run("symlink", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("symlink privileges are not guaranteed")
		}
		directory := filepath.Join(t.TempDir(), "state")
		store, err := OpenStore(directory)
		if err != nil {
			t.Fatal(err)
		}
		defer store.Close()
		target := filepath.Join(directory, "target")
		if err := os.WriteFile(target, validData, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("target", filepath.Join(directory, registryFilename(1))); err != nil {
			t.Fatal(err)
		}
		if _, err := store.Load(context.Background()); !errors.Is(err, ErrInvalid) {
			t.Fatalf("symlink=%v", err)
		}
	})

	t.Run("history metadata and content", func(t *testing.T) {
		directory := filepath.Join(t.TempDir(), "state")
		store, err := OpenStore(directory)
		if err != nil {
			t.Fatal(err)
		}
		defer store.Close()
		event := HistoryEvent{ID: thirdID, TargetID: deliveryID, Kind: HistoryStopped, At: testTime}
		data, _ := encodeJSON(event)
		wrongName := HistoryEvent{ID: fourthID, At: testTime}
		if err := os.WriteFile(filepath.Join(directory, historyFilename(wrongName)), data, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := store.History(context.Background()); !errors.Is(err, ErrInvalid) {
			t.Fatalf("metadata=%v", err)
		}
	})

	t.Run("history invalid json and record", func(t *testing.T) {
		for _, data := range [][]byte{[]byte("bad"), []byte(`{"id":"bad"}`)} {
			directory := filepath.Join(t.TempDir(), "state")
			store, err := OpenStore(directory)
			if err != nil {
				t.Fatal(err)
			}
			name := historyFilename(HistoryEvent{ID: thirdID, At: testTime})
			if err := os.WriteFile(filepath.Join(directory, name), data, 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := store.History(context.Background()); err == nil {
				t.Fatal("invalid history accepted")
			}
			_ = store.Close()
		}
	})

	t.Run("history valid metadata with invalid event", func(t *testing.T) {
		directory := filepath.Join(t.TempDir(), "state")
		store, err := OpenStore(directory)
		if err != nil {
			t.Fatal(err)
		}
		defer store.Close()
		event := HistoryEvent{ID: thirdID, TargetID: deliveryID, Kind: "invalid", At: testTime}
		data, _ := encodeJSON(event)
		if err := os.WriteFile(filepath.Join(directory, historyFilename(event)), data, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := store.History(context.Background()); !errors.Is(err, ErrInvalid) {
			t.Fatalf("event=%v", err)
		}
	})

	t.Run("maximum revision", func(t *testing.T) {
		directory := filepath.Join(t.TempDir(), "state")
		store, err := OpenStore(directory)
		if err != nil {
			t.Fatal(err)
		}
		defer store.Close()
		snapshot := EmptySnapshot()
		snapshot.Revision = math.MaxUint64
		data, _ := encodeJSON(snapshot)
		if err := os.WriteFile(filepath.Join(directory, registryFilename(math.MaxUint64)), data, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := store.Update(context.Background(), math.MaxUint64, func(*Registry) error { return nil }); !errors.Is(err, ErrInvalid) {
			t.Fatalf("overflow=%v", err)
		}
	})
}

func TestOpenStoreFailuresAndHelpers(t *testing.T) {
	if _, err := OpenStore(" "); !errors.Is(err, ErrInvalid) {
		t.Fatalf("empty=%v", err)
	}
	originalMkdir, originalLstat, originalChmod, originalOpen, originalConfig := makeStateDirectory, lstatStateDirectory, chmodStateDirectory, openStateRoot, userConfigDirectory
	t.Cleanup(func() {
		makeStateDirectory, lstatStateDirectory, chmodStateDirectory, openStateRoot, userConfigDirectory = originalMkdir, originalLstat, originalChmod, originalOpen, originalConfig
	})
	want := errors.New("failure")
	for _, test := range []struct {
		name      string
		configure func()
		want      error
	}{
		{"lstat", func() { lstatStateDirectory = func(string) (fs.FileInfo, error) { return nil, want } }, want},
		{"mkdir", func() {
			lstatStateDirectory = func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
			makeStateDirectory = func(string, fs.FileMode) error { return want }
		}, want},
		{"post-create lstat", func() {
			calls := 0
			lstatStateDirectory = func(string) (fs.FileInfo, error) {
				calls++
				if calls == 1 {
					return nil, fs.ErrNotExist
				}
				return nil, want
			}
			makeStateDirectory = func(string, fs.FileMode) error { return nil }
		}, want},
		{"post-create replacement", func() {
			calls := 0
			lstatStateDirectory = func(string) (fs.FileInfo, error) {
				calls++
				if calls == 1 {
					return nil, fs.ErrNotExist
				}
				return fakeFileInfo{mode: 0o600}, nil
			}
			makeStateDirectory = func(string, fs.FileMode) error { return nil }
		}, ErrInvalid},
		{"chmod", func() {
			lstatStateDirectory = newDirectoryLstat()
			makeStateDirectory = func(string, fs.FileMode) error { return nil }
			chmodStateDirectory = func(string, fs.FileMode) error { return want }
		}, want},
		{"open", func() {
			lstatStateDirectory = newDirectoryLstat()
			makeStateDirectory = func(string, fs.FileMode) error { return nil }
			chmodStateDirectory = func(string, fs.FileMode) error { return nil }
			openStateRoot = func(string) (stateRoot, error) { return nil, want }
		}, want},
	} {
		t.Run(test.name, func(t *testing.T) {
			makeStateDirectory, lstatStateDirectory, chmodStateDirectory, openStateRoot = originalMkdir, originalLstat, originalChmod, originalOpen
			test.configure()
			if _, err := OpenStore("state"); !errors.Is(err, test.want) {
				t.Fatalf("error=%v", err)
			}
		})
	}
	makeStateDirectory, lstatStateDirectory, chmodStateDirectory, openStateRoot = originalMkdir, originalLstat, originalChmod, originalOpen
	if _, err := openStateRoot(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("missing root opened")
	}
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenStore(file); !errors.Is(err, ErrInvalid) {
		t.Fatalf("file=%v", err)
	}

	if registryFilename(12) != "registry-00000000000000000012.json" {
		t.Fatal("registry filename")
	}
	event := HistoryEvent{ID: thirdID, At: testTime}
	if !strings.Contains(historyFilename(event), string(thirdID)) {
		t.Fatal("history filename")
	}
	data, err := encodeJSON(map[string]string{"html": "<safe>"})
	if err != nil || !bytes.Contains(data, []byte("<safe>")) {
		t.Fatalf("encoded=%q err=%v", data, err)
	}
	if _, err := encodeJSON(math.Inf(1)); err == nil {
		t.Fatal("unsupported JSON accepted")
	}
	var decoded map[string]string
	if err := decodeStrict(data, &decoded); err != nil || decoded["html"] != "<safe>" {
		t.Fatalf("decoded=%v err=%v", decoded, err)
	}
	userConfigDirectory = func() (string, error) { return "/config", nil }
	if directory, err := DefaultStateDirectory(); err != nil || directory != filepath.Join("/config", "courier", "state") {
		t.Fatalf("default=%q err=%v", directory, err)
	}
	userConfigDirectory = func() (string, error) { return "", want }
	if _, err := DefaultStateDirectory(); !errors.Is(err, want) {
		t.Fatalf("config error=%v", err)
	}
}

func TestStoreInjectedReadAndEncodeFailures(t *testing.T) {
	want := errors.New("failure")
	root := &fakeStateRoot{readDirErr: want}
	store := &Store{root: root}
	if _, err := store.Load(context.Background()); !errors.Is(err, want) {
		t.Fatalf("load read dir=%v", err)
	}
	if _, err := store.History(context.Background()); !errors.Is(err, want) {
		t.Fatalf("history read dir=%v", err)
	}
	root.readDirErr = nil
	root.lstatErr = want
	if _, err := store.readPrivate("state"); !errors.Is(err, want) {
		t.Fatalf("lstat=%v", err)
	}
	root.entries = []fs.DirEntry{fakeDirEntry{name: historyFilename(HistoryEvent{ID: thirdID, At: testTime})}}
	if _, err := store.History(context.Background()); !errors.Is(err, want) {
		t.Fatalf("history private=%v", err)
	}

	originalRegistry, originalHistory := encodeRegistryState, encodeHistoryState
	t.Cleanup(func() { encodeRegistryState, encodeHistoryState = originalRegistry, originalHistory })
	directory := filepath.Join(t.TempDir(), "state")
	realStore, err := OpenStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer realStore.Close()
	encodeRegistryState = func(any) ([]byte, error) { return nil, want }
	if _, err := realStore.Update(context.Background(), 0, func(*Registry) error { return nil }); !errors.Is(err, want) {
		t.Fatalf("registry encode=%v", err)
	}
	encodeHistoryState = func(any) ([]byte, error) { return nil, want }
	event := HistoryEvent{ID: thirdID, TargetID: deliveryID, Kind: HistoryStopped, At: testTime}
	if err := realStore.AppendHistory(context.Background(), event); !errors.Is(err, want) {
		t.Fatalf("history encode=%v", err)
	}
}

func TestWriteImmutableFailureCleanup(t *testing.T) {
	want := errors.New("failure")
	newStore := func(root stateRoot) *Store { return &Store{directory: "state", root: root} }
	for _, test := range []struct {
		name         string
		setup        func(context.CancelFunc) (*fakeStateRoot, *fakeStateFile)
		canceled     bool
		wantConflict bool
	}{
		{"create", func(context.CancelFunc) (*fakeStateRoot, *fakeStateFile) { return &fakeStateRoot{createErr: want}, nil }, false, false},
		{"write", func(context.CancelFunc) (*fakeStateRoot, *fakeStateFile) {
			file := &fakeStateFile{writeErr: want}
			return &fakeStateRoot{file: file}, file
		}, false, false},
		{"sync file", func(context.CancelFunc) (*fakeStateRoot, *fakeStateFile) {
			file := &fakeStateFile{syncErr: want}
			return &fakeStateRoot{file: file}, file
		}, false, false},
		{"close file", func(context.CancelFunc) (*fakeStateRoot, *fakeStateFile) {
			file := &fakeStateFile{closeErr: want}
			return &fakeStateRoot{file: file}, file
		}, false, false},
		{"cancel after close", func(cancel context.CancelFunc) (*fakeStateRoot, *fakeStateFile) {
			file := &fakeStateFile{onClose: cancel}
			return &fakeStateRoot{file: file}, file
		}, false, false},
		{"commit conflict", func(context.CancelFunc) (*fakeStateRoot, *fakeStateFile) {
			file := &fakeStateFile{}
			return &fakeStateRoot{file: file, commitErr: fs.ErrExist}, file
		}, false, true},
		{"commit", func(context.CancelFunc) (*fakeStateRoot, *fakeStateFile) {
			file := &fakeStateFile{}
			return &fakeStateRoot{file: file, commitErr: want}, file
		}, false, false},
		{"remove", func(context.CancelFunc) (*fakeStateRoot, *fakeStateFile) {
			file := &fakeStateFile{}
			return &fakeStateRoot{file: file, removeErr: want}, file
		}, false, false},
		{"sync root", func(context.CancelFunc) (*fakeStateRoot, *fakeStateFile) {
			file := &fakeStateFile{}
			return &fakeStateRoot{file: file, syncErr: want}, file
		}, false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			root, _ := test.setup(cancel)
			err := newStore(root).writeImmutable(ctx, "final", ".temp-", []byte("data"))
			if test.wantConflict {
				if !errors.Is(err, ErrRevisionConflict) {
					t.Fatalf("error=%v", err)
				}
			} else if err == nil {
				t.Fatal("expected error")
			}
		})
	}

	removeErr := errors.New("remove")
	file := &fakeStateFile{writeErr: want, closeErr: errors.New("close")}
	root := &fakeStateRoot{file: file, removeErr: removeErr}
	if err := newStore(root).writeImmutable(context.Background(), "final", ".temp-", nil); !errors.Is(err, want) || !errors.Is(err, removeErr) {
		t.Fatalf("joined=%v", err)
	}

	closedRoot, err := os.OpenRoot(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	_ = closedRoot.Close()
	if err := (&osStateRoot{root: closedRoot}).Sync(); err == nil {
		t.Fatal("closed root sync succeeded")
	}
}

func reflectSnapshots(left, right Snapshot) bool {
	leftData, _ := encodeJSON(left)
	rightData, _ := encodeJSON(right)
	return bytes.Equal(leftData, rightData)
}

type fakeStateFile struct {
	bytes.Buffer
	writeErr error
	syncErr  error
	closeErr error
	onClose  func()
}

func (file *fakeStateFile) Write(data []byte) (int, error) {
	if file.writeErr != nil {
		return 0, file.writeErr
	}
	return file.Buffer.Write(data)
}
func (file *fakeStateFile) Sync() error { return file.syncErr }
func (file *fakeStateFile) Close() error {
	if file.onClose != nil {
		file.onClose()
		file.onClose = nil
	}
	return file.closeErr
}

type fakeStateRoot struct {
	file       *fakeStateFile
	createErr  error
	commitErr  error
	removeErr  error
	syncErr    error
	closeErr   error
	entries    []fs.DirEntry
	readDirErr error
	lstatInfo  fs.FileInfo
	lstatErr   error
	readData   []byte
	readErr    error
}

type fakeDirEntry struct{ name string }

type fakeFileInfo struct{ mode fs.FileMode }

func (fakeFileInfo) Name() string           { return "state" }
func (fakeFileInfo) Size() int64            { return 0 }
func (info fakeFileInfo) Mode() fs.FileMode { return info.mode }
func (fakeFileInfo) ModTime() time.Time     { return testTime }
func (info fakeFileInfo) IsDir() bool       { return info.mode.IsDir() }
func (fakeFileInfo) Sys() any               { return nil }

func newDirectoryLstat() func(string) (fs.FileInfo, error) {
	calls := 0
	return func(string) (fs.FileInfo, error) {
		calls++
		if calls == 1 {
			return nil, fs.ErrNotExist
		}
		return fakeFileInfo{mode: fs.ModeDir | 0o700}, nil
	}
}

func (entry fakeDirEntry) Name() string         { return entry.name }
func (fakeDirEntry) IsDir() bool                { return false }
func (fakeDirEntry) Type() fs.FileMode          { return 0 }
func (fakeDirEntry) Info() (fs.FileInfo, error) { return nil, errors.New("unused") }

func (root *fakeStateRoot) ReadDir() ([]fs.DirEntry, error)   { return root.entries, root.readDirErr }
func (root *fakeStateRoot) Lstat(string) (fs.FileInfo, error) { return root.lstatInfo, root.lstatErr }
func (root *fakeStateRoot) ReadFile(string) ([]byte, error)   { return root.readData, root.readErr }
func (root *fakeStateRoot) CreateExclusive(string, fs.FileMode) (stateFile, error) {
	if root.createErr != nil {
		return nil, root.createErr
	}
	return root.file, nil
}
func (root *fakeStateRoot) CommitAbsent(string, string) error { return root.commitErr }
func (root *fakeStateRoot) Remove(string) error               { return root.removeErr }
func (root *fakeStateRoot) Sync() error                       { return root.syncErr }
func (root *fakeStateRoot) Close() error                      { return root.closeErr }

var _ io.Writer = (*fakeStateFile)(nil)
var _ sync.Locker = (*sync.Mutex)(nil)

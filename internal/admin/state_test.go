package admin

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
)

type fakeAdminInfo struct{ mode fs.FileMode }

func (info fakeAdminInfo) Name() string       { return "admin" }
func (info fakeAdminInfo) Size() int64        { return 0 }
func (info fakeAdminInfo) Mode() fs.FileMode  { return info.mode }
func (info fakeAdminInfo) ModTime() time.Time { return adminTestTime }
func (info fakeAdminInfo) IsDir() bool        { return info.mode.IsDir() }
func (info fakeAdminInfo) Sys() any           { return nil }

type fakeAdminFile struct {
	bytes.Buffer
	name      string
	chmodErr  error
	writeErr  error
	syncErr   error
	closeErr  error
	closeCall int
}

type fakeLockFile struct {
	locked    bool
	lockErr   error
	unlockErr error
	closeErr  error
}

func (file *fakeLockFile) TryLock() (bool, error) { return file.locked, file.lockErr }
func (file *fakeLockFile) Unlock() error          { return file.unlockErr }
func (file *fakeLockFile) Close() error           { return file.closeErr }

func (file *fakeAdminFile) Name() string            { return file.name }
func (file *fakeAdminFile) Chmod(fs.FileMode) error { return file.chmodErr }
func (file *fakeAdminFile) Write(data []byte) (int, error) {
	if file.writeErr != nil {
		return 0, file.writeErr
	}
	return file.Buffer.Write(data)
}
func (file *fakeAdminFile) Sync() error { return file.syncErr }
func (file *fakeAdminFile) Close() error {
	file.closeCall++
	return file.closeErr
}

func validAdminState(directory string) State {
	id := delivery.ID("00000000-0000-4000-8000-000000000701")
	return State{
		SchemaVersion: adminStateSchema, ID: id, Bind: "127.0.0.1:9090",
		ControlEndpoint: filepath.Join(directory, "control-"+string(id)+".sock"), Compatibility: "admin-v1/test",
		ProcessID: 42, StartedAt: adminTestTime,
	}
}

func TestAdminStateAndBindValidation(t *testing.T) {
	directory := t.TempDir()
	state := validAdminState(directory)
	if err := state.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*State){
		func(value *State) { value.SchemaVersion = 2 },
		func(value *State) { value.ID = "bad" },
		func(value *State) { value.ProcessID = 0 },
		func(value *State) { value.StartedAt = time.Time{} },
		func(value *State) { value.ControlEndpoint = " " },
		func(value *State) { value.ControlEndpoint = "bad\nvalue" },
		func(value *State) { value.Compatibility = " " },
		func(value *State) { value.Compatibility = "bad\x00value" },
		func(value *State) { value.Bind = "127.0.0.1:09090" },
	} {
		candidate := state
		mutate(&candidate)
		if err := candidate.Validate(); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("state=%+v err=%v", candidate, err)
		}
	}
	for _, value := range []string{"127.0.0.1:9090", "[::1]:9090", "127.1.2.3:443"} {
		if _, err := CanonicalBind(value); err != nil {
			t.Fatalf("bind=%q err=%v", value, err)
		}
	}
	for _, value := range []string{"localhost:9090", "0.0.0.0:9090", "127.0.0.1:0", "127.0.0.1:http", "missing"} {
		if _, err := CanonicalBind(value); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("bind=%q err=%v", value, err)
		}
	}
	if !strings.HasSuffix(StatePath(directory), "admin.json") {
		t.Fatalf("state path=%q", StatePath(directory))
	}
}

func TestAdminStateRoundTripAndLock(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "state")
	state := validAdminState(directory)
	if _, err := LoadState(directory); !errors.Is(err, ErrNotRunning) {
		t.Fatalf("missing=%v", err)
	}
	if err := writeState(directory, state); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadState(directory)
	if err != nil || loaded != state {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
	info, err := os.Stat(StatePath(directory))
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("mode=%v err=%v", info.Mode(), err)
	}
	if err := removeState(directory, delivery.NewID()); !errors.Is(err, delivery.ErrRevisionConflict) {
		t.Fatalf("owner mismatch=%v", err)
	}
	if err := removeState(directory, state.ID); err != nil {
		t.Fatal(err)
	}
	if err := removeState(directory, state.ID); err != nil {
		t.Fatalf("idempotent remove=%v", err)
	}

	lock, err := acquireSingleton(directory)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := acquireSingleton(directory); !errors.Is(err, ErrAlreadyRunning) {
		t.Fatalf("second lock=%v", err)
	}
	if err := lock.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestLoadStateFailures(t *testing.T) {
	directory := t.TempDir()
	path := StatePath(directory)
	state := validAdminState(directory)
	data, err := jsonBytes(state)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		data []byte
		mode fs.FileMode
	}{
		{"directory", data, fs.ModeDir | 0o700},
		{"symlink", data, fs.ModeSymlink | 0o777},
		{"malformed", []byte("{"), 0o600},
		{"unknown", []byte(`{"unknown":true}`), 0o600},
		{"trailing", append(data, []byte(`{}`)...), 0o600},
		{"invalid", []byte(`{}`), 0o600},
	}
	if runtime.GOOS != "windows" {
		tests = append(tests, struct {
			name string
			data []byte
			mode fs.FileMode
		}{"public", data, 0o644})
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			originalLstat, originalRead := stateLstat, stateReadFile
			defer func() { stateLstat, stateReadFile = originalLstat, originalRead }()
			stateLstat = func(string) (fs.FileInfo, error) { return fakeAdminInfo{mode: test.mode}, nil }
			stateReadFile = func(string) ([]byte, error) { return test.data, nil }
			if _, err := LoadState(directory); err == nil {
				t.Fatal("invalid state accepted")
			}
		})
	}
	want := errors.New("failure")
	originalLstat, originalRead := stateLstat, stateReadFile
	t.Cleanup(func() { stateLstat, stateReadFile = originalLstat, originalRead })
	stateLstat = func(string) (fs.FileInfo, error) { return nil, want }
	if _, err := LoadState(directory); !errors.Is(err, want) {
		t.Fatalf("lstat=%v", err)
	}
	stateLstat = func(string) (fs.FileInfo, error) { return fakeAdminInfo{mode: 0o600}, nil }
	stateReadFile = func(string) ([]byte, error) { return nil, want }
	if _, err := LoadState(path); !errors.Is(err, want) {
		t.Fatalf("read=%v", err)
	}
}

func TestStateWriteAndDirectoryFailures(t *testing.T) {
	directory := "state"
	state := validAdminState(directory)
	want := errors.New("failure")
	originalMkdir, originalLstat, originalRead := stateMkdirAll, stateLstat, stateReadFile
	originalCreate, originalRename, originalRemove, originalChmod, originalLock := stateCreateFile, stateRename, stateRemove, stateChmod, newSingletonFile
	t.Cleanup(func() {
		stateMkdirAll, stateLstat, stateReadFile = originalMkdir, originalLstat, originalRead
		stateCreateFile, stateRename, stateRemove, stateChmod = originalCreate, originalRename, originalRemove, originalChmod
		newSingletonFile = originalLock
	})
	invalid := state
	invalid.ID = "bad"
	if err := writeState(directory, invalid); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid state=%v", err)
	}

	if err := ensurePrivateDirectory(" "); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("empty directory=%v", err)
	}
	stateMkdirAll = func(string, fs.FileMode) error { return want }
	if err := ensurePrivateDirectory(directory); !errors.Is(err, want) {
		t.Fatalf("mkdir=%v", err)
	}
	if err := writeState(directory, state); !errors.Is(err, want) {
		t.Fatalf("write directory=%v", err)
	}
	stateMkdirAll = func(string, fs.FileMode) error { return nil }
	stateLstat = func(string) (fs.FileInfo, error) { return nil, want }
	if err := ensurePrivateDirectory(directory); !errors.Is(err, want) {
		t.Fatalf("lstat=%v", err)
	}
	stateLstat = func(string) (fs.FileInfo, error) { return fakeAdminInfo{mode: 0o600}, nil }
	if err := ensurePrivateDirectory(directory); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("file directory=%v", err)
	}
	stateLstat = func(string) (fs.FileInfo, error) { return fakeAdminInfo{mode: fs.ModeDir | 0o700}, nil }
	stateChmod = func(string, fs.FileMode) error { return want }
	if err := ensurePrivateDirectory(directory); !errors.Is(err, want) {
		t.Fatalf("chmod=%v", err)
	}
	stateChmod = func(string, fs.FileMode) error { return nil }
	stateCreateFile = func(string, string) (stateFile, error) { return nil, want }
	if err := writeState(directory, state); !errors.Is(err, want) {
		t.Fatalf("create=%v", err)
	}

	for _, test := range []struct {
		name string
		file *fakeAdminFile
	}{
		{"chmod", &fakeAdminFile{name: "temporary", chmodErr: want}},
		{"write", &fakeAdminFile{name: "temporary", writeErr: want}},
		{"sync", &fakeAdminFile{name: "temporary", syncErr: want}},
		{"close", &fakeAdminFile{name: "temporary", closeErr: want}},
	} {
		t.Run(test.name, func(t *testing.T) {
			stateCreateFile = func(string, string) (stateFile, error) { return test.file, nil }
			stateRemove = func(string) error { return nil }
			if err := writeState(directory, state); !errors.Is(err, want) {
				t.Fatalf("error=%v", err)
			}
		})
	}
	file := &fakeAdminFile{name: "temporary"}
	stateCreateFile = func(string, string) (stateFile, error) { return file, nil }
	stateRename = func(string, string) error { return want }
	if err := writeState(directory, state); !errors.Is(err, want) {
		t.Fatalf("rename=%v", err)
	}

	stateLstat = func(string) (fs.FileInfo, error) { return nil, want }
	if err := removeState(directory, state.ID); !errors.Is(err, want) {
		t.Fatalf("remove load=%v", err)
	}
	stateLstat = func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }
	stateMkdirAll = func(string, fs.FileMode) error { return want }
	if _, err := acquireSingleton(directory); !errors.Is(err, want) {
		t.Fatalf("lock directory=%v", err)
	}
	stateMkdirAll = func(string, fs.FileMode) error { return nil }
	stateLstat = func(string) (fs.FileInfo, error) { return fakeAdminInfo{mode: fs.ModeDir | 0o700}, nil }
	stateChmod = func(string, fs.FileMode) error { return nil }
	newSingletonFile = func(string) lockFile { return &fakeLockFile{lockErr: want} }
	if _, err := acquireSingleton(directory); !errors.Is(err, want) {
		t.Fatalf("lock error=%v", err)
	}
	newSingletonFile = func(string) lockFile { return &fakeLockFile{locked: true, unlockErr: want} }
	lock, err := acquireSingleton(directory)
	if err != nil || !errors.Is(lock.Close(), want) {
		t.Fatalf("unlock error=%v close=%v", err, lock)
	}
	data, err := jsonBytes(state)
	if err != nil {
		t.Fatal(err)
	}
	stateLstat = func(string) (fs.FileInfo, error) { return fakeAdminInfo{mode: 0o600}, nil }
	stateReadFile = func(string) ([]byte, error) { return data, nil }
	stateRemove = func(string) error { return want }
	if err := removeState(directory, state.ID); !errors.Is(err, want) {
		t.Fatalf("remove error=%v", err)
	}
}

func jsonBytes(value any) ([]byte, error) {
	var output bytes.Buffer
	err := json.NewEncoder(&output).Encode(value)
	return output.Bytes(), err
}

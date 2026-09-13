package admin

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
	"github.com/iwonz/courier/internal/worker"
)

func freeAdminBind(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	bind := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	return bind
}

func shortAdminDirectory(t *testing.T) string {
	t.Helper()
	base := os.TempDir()
	if runtime.GOOS == "darwin" {
		base = "/private/tmp"
	}
	directory, err := os.MkdirTemp(base, "courier-admin-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(directory) })
	return filepath.Join(directory, "state")
}

func validAdminConfig(t *testing.T, directory, bind string) processConfig {
	t.Helper()
	id := delivery.NewID()
	endpoint, err := ipc.ControlEndpoint(directory, id)
	if err != nil {
		t.Fatal(err)
	}
	return processConfig{StateDirectory: directory, State: State{
		SchemaVersion: adminStateSchema, ID: id, Bind: bind, ControlEndpoint: endpoint,
		Compatibility: "admin-v1/test", ProcessID: 1, StartedAt: adminTestTime,
	}}
}

func TestAdminForegroundLifecycleAndSingleton(t *testing.T) {
	directory := shortAdminDirectory(t)
	manager := Manager{StateDirectory: directory, Compatibility: "admin-v1/test", ReadyTimeout: time.Second}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ready := make(chan State, 1)
	done := make(chan error, 1)
	bind := freeAdminBind(t)
	go func() {
		_, err := manager.Start(ctx, StartRequest{Bind: bind, Ready: func(state State) error { ready <- state; return nil }})
		done <- err
	}()
	var state State
	select {
	case state = <-ready:
	case err := <-done:
		t.Fatalf("start before ready: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("start did not become ready")
	}
	if state.Bind != bind || state.ProcessID != os.Getpid() || URL(state) != "http://"+bind+"/" {
		t.Fatalf("state=%+v url=%s", state, URL(state))
	}
	loaded, err := LoadState(directory)
	if err != nil || loaded.ID != state.ID {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
	existingReady := false
	result, err := manager.Start(context.Background(), StartRequest{Bind: bind, Ready: func(current State) error { existingReady = current.ID == state.ID; return nil }})
	if err != nil || !result.AlreadyRunning || !existingReady {
		t.Fatalf("existing=%+v ready=%v err=%v", result, existingReady, err)
	}
	if _, err := manager.Start(context.Background(), StartRequest{Bind: freeAdminBind(t)}); !errors.Is(err, ErrAlreadyRunning) {
		t.Fatalf("different bind=%v", err)
	}
	want := errors.New("ready failure")
	if _, err := manager.Start(context.Background(), StartRequest{Bind: bind, Ready: func(State) error { return want }}); !errors.Is(err, want) {
		t.Fatalf("existing ready=%v", err)
	}
	stopped, err := manager.Stop(context.Background())
	if err != nil || stopped.AlreadyStopped {
		t.Fatalf("stop=%+v err=%v", stopped, err)
	}
	if err := <-done; err != nil {
		t.Fatalf("foreground=%v", err)
	}
	stopped, err = manager.Stop(context.Background())
	if err != nil || !stopped.AlreadyStopped {
		t.Fatalf("repeated stop=%+v err=%v", stopped, err)
	}
}

func TestAdminRunInternal(t *testing.T) {
	if handled, err := RunInternal(context.Background(), nil); handled || err != nil {
		t.Fatalf("unhandled=%v err=%v", handled, err)
	}
	if handled, err := RunInternal(context.Background(), []string{"other"}); handled || err != nil {
		t.Fatalf("other=%v err=%v", handled, err)
	}
	directory := shortAdminDirectory(t)
	config := validAdminConfig(t, directory, freeAdminBind(t))
	path, err := writeProcessConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv(configEnvironment, path)
	done := make(chan error, 1)
	go func() {
		handled, runErr := RunInternal(context.Background(), []string{internalArgument})
		if !handled && runErr == nil {
			runErr = errors.New("internal mode was not handled")
		}
		done <- runErr
	}()
	manager := Manager{StateDirectory: directory, Compatibility: config.State.Compatibility, ReadyTimeout: time.Second}
	deadline := time.Now().Add(time.Second)
	for {
		if _, found, _ := manager.probe(context.Background()); found {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("internal admin did not start")
		}
		time.Sleep(time.Millisecond)
	}
	if _, err := manager.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("launch file remains: %v", err)
	}

	t.Setenv(configEnvironment, "")
	if handled, err := RunInternal(context.Background(), []string{internalArgument}); !handled || err == nil {
		t.Fatalf("invalid internal handled=%v err=%v", handled, err)
	}
}

func TestAdminConfigRoundTripAndValidation(t *testing.T) {
	directory := t.TempDir()
	config := validAdminConfig(t, directory, "127.0.0.1:9090")
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*processConfig){
		func(value *processConfig) { value.StateDirectory = " " },
		func(value *processConfig) { value.State.Bind = "bad" },
		func(value *processConfig) { value.State.ControlEndpoint = "other" },
	} {
		candidate := config
		mutate(&candidate)
		if err := candidate.Validate(); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("config=%+v err=%v", candidate, err)
		}
	}
	path, err := writeProcessConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)
	loaded, err := readProcessConfig(path)
	if err != nil || loaded.State.ID != config.State.ID {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
	if _, err := readProcessConfig(" "); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("empty path=%v", err)
	}
}

func TestAdminManagerValidationAndInjectedLifecycle(t *testing.T) {
	want := errors.New("failure")
	manager := Manager{StateDirectory: t.TempDir(), Compatibility: "admin-v1/test", Executable: os.Args[0], ReadyTimeout: 100 * time.Millisecond}
	for _, test := range []struct {
		manager Manager
		request StartRequest
	}{
		{manager, StartRequest{Bind: "bad"}},
		{Manager{}, StartRequest{Bind: "127.0.0.1:9090"}},
		{Manager{StateDirectory: manager.StateDirectory, Compatibility: "bad\nvalue"}, StartRequest{Bind: "127.0.0.1:9090"}},
	} {
		if _, err := test.manager.Start(context.Background(), test.request); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("manager=%+v request=%+v err=%v", test.manager, test.request, err)
		}
	}

	originalStart, originalWait, originalExecutable, originalEndpoint := adminStartCommand, adminWaitReady, adminExecutable, adminControlEndpoint
	t.Cleanup(func() {
		adminStartCommand, adminWaitReady, adminExecutable, adminControlEndpoint = originalStart, originalWait, originalExecutable, originalEndpoint
	})
	adminControlEndpoint = func(string, delivery.ID) (string, error) { return "", want }
	if _, err := manager.Start(context.Background(), StartRequest{Bind: "127.0.0.1:9090"}); !errors.Is(err, want) {
		t.Fatalf("control endpoint=%v", err)
	}
	adminControlEndpoint = originalEndpoint
	adminStartCommand = func(string, ...string) *exec.Cmd {
		return exec.Command(os.Args[0], "-test.run=TestAdminChildProcess", "--", "return")
	}
	adminWaitReady = func(_ context.Context, _ string, id delivery.ID) (State, error) {
		state := validAdminState(manager.StateDirectory)
		state.ID = id
		state.Bind = "127.0.0.1:9090"
		return state, nil
	}
	result, err := manager.Start(context.Background(), StartRequest{Bind: "127.0.0.1:9090", Background: true})
	if err != nil || result.State.ID == "" {
		t.Fatalf("background=%+v err=%v", result, err)
	}
	if _, err := manager.Start(context.Background(), StartRequest{Bind: "127.0.0.1:9090", Background: true, Ready: func(State) error { return want }}); !errors.Is(err, want) {
		t.Fatalf("background ready=%v", err)
	}
	manager.ReadyTimeout = 0
	if _, err := manager.Start(context.Background(), StartRequest{Bind: "127.0.0.1:9090", Background: true}); err != nil {
		t.Fatalf("default background timeout=%v", err)
	}

	adminExecutable = func() (string, error) { return "", want }
	manager.Executable = ""
	if _, err := manager.Start(context.Background(), StartRequest{Bind: "127.0.0.1:9090", Background: true}); !errors.Is(err, want) {
		t.Fatalf("executable=%v", err)
	}
	adminStartCommand = originalStart
	manager.Executable = "missing-courier-executable"
	if _, err := manager.Start(context.Background(), StartRequest{Bind: "127.0.0.1:9090", Background: true}); err == nil {
		t.Fatal("missing executable started")
	}
}

func TestAdminChildProcess(t *testing.T) {
	for index, argument := range os.Args {
		if argument == "--" && index+1 < len(os.Args) && os.Args[index+1] == "sleep" {
			time.Sleep(10 * time.Second)
		}
	}
}

func TestAdminLaunchTimeoutAndWaitForReady(t *testing.T) {
	directory := t.TempDir()
	manager := Manager{StateDirectory: directory, Compatibility: "admin-v1/test", Executable: os.Args[0], ReadyTimeout: 20 * time.Millisecond}
	originalStart, originalWait, originalPoll, originalHello := adminStartCommand, adminWaitReady, adminPollInterval, adminHello
	t.Cleanup(func() {
		adminStartCommand, adminWaitReady, adminPollInterval, adminHello = originalStart, originalWait, originalPoll, originalHello
	})
	adminStartCommand = func(string, ...string) *exec.Cmd {
		return exec.Command(os.Args[0], "-test.run=TestAdminChildProcess", "--", "sleep")
	}
	adminWaitReady = waitForReady
	adminPollInterval = time.Millisecond
	if _, err := manager.Start(context.Background(), StartRequest{Bind: "127.0.0.1:9090", Background: true}); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("timeout=%v", err)
	}

	state := validAdminState(directory)
	if err := writeState(directory, state); err != nil {
		t.Fatal(err)
	}
	adminHello = func(context.Context, State) error { return nil }
	ready, err := waitForReady(context.Background(), directory, state.ID)
	if err != nil || ready.ID != state.ID {
		t.Fatalf("ready=%+v err=%v", ready, err)
	}
	if _, err := waitForReady(context.Background(), directory, delivery.NewID()); !errors.Is(err, ErrAlreadyRunning) {
		t.Fatalf("competing owner=%v", err)
	}
}

func TestAdminProcessConfigFailures(t *testing.T) {
	directory := t.TempDir()
	config := validAdminConfig(t, directory, "127.0.0.1:9090")
	want := errors.New("failure")
	originalCreate, originalLstat, originalRead, originalRemove := adminCreateFile, adminLstatFile, adminReadFile, adminRemoveFile
	originalMkdir := stateMkdirAll
	t.Cleanup(func() {
		adminCreateFile, adminLstatFile, adminReadFile, adminRemoveFile = originalCreate, originalLstat, originalRead, originalRemove
		stateMkdirAll = originalMkdir
	})
	if _, err := writeProcessConfig(processConfig{}); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid config=%v", err)
	}
	stateMkdirAll = func(string, fs.FileMode) error { return want }
	if _, err := writeProcessConfig(config); !errors.Is(err, want) {
		t.Fatalf("directory=%v", err)
	}
	stateMkdirAll = originalMkdir

	adminCreateFile = func(string, string) (processFile, error) { return nil, want }
	if _, err := writeProcessConfig(config); !errors.Is(err, want) {
		t.Fatalf("create=%v", err)
	}
	for _, test := range []struct {
		name string
		file *fakeAdminFile
	}{
		{"chmod", &fakeAdminFile{name: "config", chmodErr: want}},
		{"encode", &fakeAdminFile{name: "config", writeErr: want}},
		{"sync", &fakeAdminFile{name: "config", syncErr: want}},
		{"close", &fakeAdminFile{name: "config", closeErr: want}},
	} {
		t.Run(test.name, func(t *testing.T) {
			adminCreateFile = func(string, string) (processFile, error) { return test.file, nil }
			adminRemoveFile = func(string) error { return nil }
			if _, err := writeProcessConfig(config); !errors.Is(err, want) {
				t.Fatalf("error=%v", err)
			}
		})
	}

	adminLstatFile = func(string) (fs.FileInfo, error) { return nil, want }
	if _, err := readProcessConfig("config"); !errors.Is(err, want) {
		t.Fatalf("lstat=%v", err)
	}
	for _, mode := range []fs.FileMode{fs.ModeDir | 0o700, fs.ModeSymlink | 0o777} {
		adminLstatFile = func(string) (fs.FileInfo, error) { return fakeAdminInfo{mode: mode}, nil }
		if _, err := readProcessConfig("config"); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("mode=%v err=%v", mode, err)
		}
	}
	if runtime.GOOS != "windows" {
		adminLstatFile = func(string) (fs.FileInfo, error) { return fakeAdminInfo{mode: 0o644}, nil }
		if _, err := readProcessConfig("config"); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("public config=%v", err)
		}
	}
	adminLstatFile = func(string) (fs.FileInfo, error) { return fakeAdminInfo{mode: 0o600}, nil }
	adminReadFile = func(string) ([]byte, error) { return nil, want }
	if _, err := readProcessConfig("config"); !errors.Is(err, want) {
		t.Fatalf("read=%v", err)
	}
	for _, data := range [][]byte{[]byte("{"), []byte(`{"unknown":true}`), []byte(`{} {}`), []byte(`{}`)} {
		adminReadFile = func(string) ([]byte, error) { return data, nil }
		if _, err := readProcessConfig("config"); err == nil {
			t.Fatalf("data=%q accepted", data)
		}
	}
}

func TestAdminStopRegistryReadAndInvalidLaunch(t *testing.T) {
	directory := t.TempDir()
	state := validAdminState(directory)
	if err := writeState(directory, state); err != nil {
		t.Fatal(err)
	}
	want := errors.New("registry read")
	originalLstat, originalHello, originalShutdown := stateLstat, adminHello, adminShutdown
	t.Cleanup(func() { stateLstat, adminHello, adminShutdown = originalLstat, originalHello, originalShutdown })
	adminHello = func(context.Context, State) error { return nil }
	adminShutdown = func(context.Context, State) error { return nil }
	calls := 0
	stateLstat = func(path string) (fs.FileInfo, error) {
		calls++
		if calls > 1 {
			return nil, want
		}
		return originalLstat(path)
	}
	if _, err := (Manager{StateDirectory: directory, ReadyTimeout: time.Second}).Stop(context.Background()); !errors.Is(err, want) {
		t.Fatalf("read error=%v", err)
	}
	if _, err := (Manager{}).launch(context.Background(), processConfig{}); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid launch=%v", err)
	}
}

func TestNormalizeAndClosedListener(t *testing.T) {
	if normalizeServerError(nil) != nil || normalizeServerError(net.ErrClosed) != nil || normalizeServerError(errors.New("failure")) == nil {
		t.Fatal("server error normalization mismatch")
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	if err := closeListener(listener); err != nil || closeListener(listener) != nil {
		t.Fatalf("close listener=%v", err)
	}
	if !strings.HasPrefix(worker.ProcessEnvironment("TEST_ADMIN", "value")[0], "TEST_ADMIN=value") {
		t.Fatal("shared process environment missing")
	}
}

func TestAdminStopAndProbeFailures(t *testing.T) {
	want := errors.New("failure")
	originalHello, originalShutdown, originalReady, originalPoll := adminHello, adminShutdown, adminReadyTimeout, adminPollInterval
	t.Cleanup(func() {
		adminHello, adminShutdown, adminReadyTimeout, adminPollInterval = originalHello, originalShutdown, originalReady, originalPoll
	})
	adminPollInterval = time.Millisecond

	t.Run("load", func(t *testing.T) {
		directory := t.TempDir()
		if err := os.WriteFile(StatePath(directory), []byte(`{}`), 0o600); err != nil {
			t.Fatal(err)
		}
		manager := Manager{StateDirectory: directory}
		if _, err := manager.Stop(context.Background()); err == nil {
			t.Fatal("invalid state accepted")
		}
		if _, _, err := manager.probe(context.Background()); err == nil {
			t.Fatal("invalid probe state accepted")
		}
		if _, err := (Manager{StateDirectory: directory, Compatibility: "admin-v1/test"}).Start(context.Background(), StartRequest{Bind: "127.0.0.1:9090"}); err == nil {
			t.Fatal("invalid start state accepted")
		}
	})

	for _, test := range []struct {
		name       string
		helloErr   error
		stopErr    error
		mutate     bool
		timeout    bool
		useDefault bool
	}{
		{name: "hello", helloErr: want},
		{name: "shutdown", stopErr: want},
		{name: "owner change", mutate: true},
		{name: "timeout", timeout: true},
		{name: "default timeout", timeout: true, useDefault: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			state := validAdminState(directory)
			if err := writeState(directory, state); err != nil {
				t.Fatal(err)
			}
			adminHello = func(context.Context, State) error { return test.helloErr }
			adminShutdown = func(context.Context, State) error {
				if test.mutate {
					replacement := state
					replacement.ID = delivery.ID("00000000-0000-4000-8000-000000000799")
					if err := writeState(directory, replacement); err != nil {
						t.Fatal(err)
					}
				}
				return test.stopErr
			}
			manager := Manager{StateDirectory: directory, ReadyTimeout: 10 * time.Millisecond}
			if test.useDefault {
				manager.ReadyTimeout = 0
				adminReadyTimeout = 10 * time.Millisecond
			}
			_, err := manager.Stop(context.Background())
			switch {
			case test.helloErr != nil || test.stopErr != nil:
				if !errors.Is(err, want) {
					t.Fatalf("error=%v", err)
				}
			case test.mutate:
				if !errors.Is(err, delivery.ErrRevisionConflict) {
					t.Fatalf("owner error=%v", err)
				}
			case test.timeout:
				if !errors.Is(err, context.DeadlineExceeded) {
					t.Fatalf("timeout=%v", err)
				}
			}
		})
	}

	directory := t.TempDir()
	state := validAdminState(directory)
	if err := writeState(directory, state); err != nil {
		t.Fatal(err)
	}
	adminHello = func(context.Context, State) error { return want }
	if current, found, err := (Manager{StateDirectory: directory}).probe(context.Background()); err != nil || found || current.ID != state.ID {
		t.Fatalf("stale probe=%+v found=%v err=%v", current, found, err)
	}
}

type fakeAdminListener struct {
	accept net.Conn
	err    error
	closed error
}

func (listener *fakeAdminListener) Accept() (net.Conn, error) { return listener.accept, listener.err }
func (listener *fakeAdminListener) Close() error              { return listener.closed }
func (*fakeAdminListener) Addr() net.Addr                     { return &net.TCPAddr{} }

func TestRunProcessFailures(t *testing.T) {
	want := errors.New("failure")
	originalAcquire, originalHello, originalRemoveStale := adminAcquire, adminHello, adminRemoveStale
	originalRemoveState, originalControl, originalTCP := adminRemoveState, adminListenControl, adminListenTCP
	originalStore, originalAPI, originalWrite := adminOpenStore, adminNewAPI, adminWriteState
	originalHTTP, originalServeControl, originalShutdownHTTP := adminServeHTTP, adminServeControl, adminShutdownHTTP
	t.Cleanup(func() {
		adminAcquire, adminHello, adminRemoveStale = originalAcquire, originalHello, originalRemoveStale
		adminRemoveState, adminListenControl, adminListenTCP = originalRemoveState, originalControl, originalTCP
		adminOpenStore, adminNewAPI, adminWriteState = originalStore, originalAPI, originalWrite
		adminServeHTTP, adminServeControl, adminShutdownHTTP = originalHTTP, originalServeControl, originalShutdownHTTP
	})

	if err := runProcess(context.Background(), processConfig{}, nil); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid=%v", err)
	}
	directory := t.TempDir()
	config := validAdminConfig(t, directory, "127.0.0.1:9090")
	adminAcquire = func(string) (*singletonLock, error) { return nil, want }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, want) {
		t.Fatalf("lock=%v", err)
	}
	adminAcquire = func(string) (*singletonLock, error) { return &singletonLock{file: &fakeLockFile{locked: true}}, nil }

	previous := validAdminState(directory)
	if err := writeState(directory, previous); err != nil {
		t.Fatal(err)
	}
	adminHello = func(context.Context, State) error { return nil }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, ErrAlreadyRunning) {
		t.Fatalf("live previous=%v", err)
	}
	adminHello = func(context.Context, State) error { return want }
	adminRemoveStale = func(string) error { return want }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, want) {
		t.Fatalf("remove stale=%v", err)
	}
	adminRemoveStale = func(string) error { return nil }
	adminRemoveState = func(string, delivery.ID) error { return want }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, want) {
		t.Fatalf("remove state=%v", err)
	}
	adminRemoveState = originalRemoveState
	if err := os.WriteFile(StatePath(directory), []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := runProcess(context.Background(), config, nil); err == nil {
		t.Fatal("invalid previous state accepted")
	}
	if err := os.Remove(StatePath(directory)); err != nil {
		t.Fatal(err)
	}

	adminListenControl = func(string) (net.Listener, error) { return nil, want }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, want) {
		t.Fatalf("control listen=%v", err)
	}
	controlListener := &fakeAdminListener{}
	adminListenControl = func(string) (net.Listener, error) { return controlListener, nil }
	adminListenTCP = func(string) (net.Listener, error) { return nil, want }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, want) {
		t.Fatalf("TCP listen=%v", err)
	}
	dataListener := &fakeAdminListener{}
	adminListenTCP = func(string) (net.Listener, error) { return dataListener, nil }
	adminOpenStore = func(string) (*delivery.Store, error) { return nil, want }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, want) {
		t.Fatalf("store=%v", err)
	}
	store, err := delivery.OpenStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	adminOpenStore = func(string) (*delivery.Store, error) { return store, nil }
	adminNewAPI = func(APIConfig) (*API, error) { return nil, want }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, want) {
		t.Fatalf("API=%v", err)
	}
	adminNewAPI = originalAPI
	adminWriteState = func(string, State) error { return want }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, want) {
		t.Fatalf("state write=%v", err)
	}
	adminWriteState = originalWrite
	if err := runProcess(context.Background(), config, func(State) error { return want }); !errors.Is(err, want) {
		t.Fatalf("ready=%v", err)
	}

	adminServeHTTP = func(*http.Server, net.Listener) error { return want }
	adminServeControl = func(ctx context.Context, _ net.Listener, _ State, _ context.CancelFunc) error {
		<-ctx.Done()
		return nil
	}
	adminShutdownHTTP = func(*http.Server, context.Context) error { return nil }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, want) {
		t.Fatalf("HTTP serve=%v", err)
	}
	controlBranchHTTP := make(chan struct{})
	adminServeHTTP = func(_ *http.Server, _ net.Listener) error { <-controlBranchHTTP; return nil }
	adminServeControl = func(context.Context, net.Listener, State, context.CancelFunc) error { return want }
	adminShutdownHTTP = func(*http.Server, context.Context) error { close(controlBranchHTTP); return nil }
	if err := runProcess(context.Background(), config, nil); !errors.Is(err, want) {
		t.Fatalf("control serve=%v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cancelBranchHTTP := make(chan struct{})
	adminServeHTTP = func(_ *http.Server, _ net.Listener) error { <-cancelBranchHTTP; return nil }
	adminServeControl = func(ctx context.Context, _ net.Listener, _ State, _ context.CancelFunc) error {
		<-ctx.Done()
		return nil
	}
	adminShutdownHTTP = func(*http.Server, context.Context) error { close(cancelBranchHTTP); return want }
	if err := runProcess(ctx, config, nil); !errors.Is(err, want) {
		t.Fatalf("HTTP shutdown=%v", err)
	}
}

func TestControlProtocolFailures(t *testing.T) {
	state := validAdminState(t.TempDir())
	want := errors.New("accept")
	if err := serveControl(context.Background(), &fakeAdminListener{err: want}, state, func() {}); !errors.Is(err, want) {
		t.Fatalf("accept=%v", err)
	}
	if err := serveControl(context.Background(), &fakeAdminListener{err: net.ErrClosed}, state, func() {}); err != nil {
		t.Fatalf("closed=%v", err)
	}

	request := func(operation ipc.Operation, payload any) ipc.Request {
		value, err := ipc.NewRequest(operation, payload)
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	for _, value := range []ipc.Request{
		request(ipc.OperationHello, map[string]string{"bad": "payload"}),
		request(ipc.OperationHello, worker.HelloRequest{ServerID: delivery.NewID(), Compatibility: state.Compatibility}),
		request(ipc.OperationShutdown, map[string]bool{"unexpected": true}),
		request(ipc.OperationList, nil),
	} {
		client, server := net.Pipe()
		done := make(chan struct{})
		go func() { handleControl(server, state, func() {}); close(done) }()
		if err := ipc.WriteFrame(client, value); err != nil {
			t.Fatal(err)
		}
		var response ipc.Response
		if err := ipc.ReadFrame(client, &response); err != nil || response.OK {
			t.Fatalf("response=%+v err=%v", response, err)
		}
		_ = client.Close()
		<-done
	}

	shutdownCalled := make(chan struct{})
	client, server := net.Pipe()
	go handleControl(server, state, func() { close(shutdownCalled) })
	value := request(ipc.OperationShutdown, nil)
	if err := ipc.WriteFrame(client, value); err != nil {
		t.Fatal(err)
	}
	var response ipc.Response
	if err := ipc.ReadFrame(client, &response); err != nil || !response.OK {
		t.Fatalf("shutdown=%+v err=%v", response, err)
	}
	_ = client.Close()
	<-shutdownCalled

	client, server = net.Pipe()
	done := make(chan struct{})
	go func() { handleControl(server, state, func() {}); close(done) }()
	_, _ = client.Write([]byte{0})
	_ = client.Close()
	<-done

	listener := &fakeAdminListener{closed: want}
	if err := closeListener(listener); !errors.Is(err, want) {
		t.Fatalf("close error=%v", err)
	}
}

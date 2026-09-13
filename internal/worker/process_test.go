package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
	"github.com/iwonz/courier/internal/webdelivery"
)

func validProcessConfig(t *testing.T, directory, bind string) processConfig {
	t.Helper()
	endpoint, err := ipc.ControlEndpoint(directory, workerServerID)
	if err != nil {
		t.Fatal(err)
	}
	return processConfig{StateDirectory: directory, Launch: LaunchRequest{ServerID: workerServerID, Bind: bind, ControlEndpoint: endpoint, Compatibility: "http-v1"}}
}

func TestProcessConfigValidation(t *testing.T) {
	directory := t.TempDir()
	config := validProcessConfig(t, directory, "127.0.0.1:8080")
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*processConfig){
		func(value *processConfig) { value.StateDirectory = "" },
		func(value *processConfig) { value.Launch.ServerID = "bad" },
		func(value *processConfig) { value.Launch.Compatibility = "" },
		func(value *processConfig) { value.Launch.Bind = "bad" },
		func(value *processConfig) { value.Launch.Bind = "LOCALHOST:8080" },
		func(value *processConfig) { value.Launch.ControlEndpoint = "wrong" },
	} {
		candidate := config
		mutate(&candidate)
		if err := candidate.Validate(); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("config=%+v err=%v", candidate, err)
		}
	}
}

func TestOwnedTemporaryCleanup(t *testing.T) {
	originalDirectory, originalOpen := ownedTempDirectory, openTemporaryRoot
	t.Cleanup(func() { ownedTempDirectory, openTemporaryRoot = originalDirectory, originalOpen })
	directory := t.TempDir()
	ownedTempDirectory = func() string { return directory }
	name := ".courier-123456.tar.gz"
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, []byte("archive"), 0o600); err != nil {
		t.Fatal(err)
	}
	temporary := delivery.OwnedTemp{Location: delivery.TempLocal, Path: path}
	if err := cleanupOwnedTemp(context.Background(), temporary); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("temporary remains: %v", err)
	}
	if err := cleanupOwnedTemp(context.Background(), temporary); err != nil {
		t.Fatalf("missing temporary is not idempotent: %v", err)
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := cleanupOwnedTemp(canceled, temporary); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled cleanup=%v", err)
	}
	remote := temporary
	remote.Location = delivery.TempRemote
	if err := cleanupOwnedTemp(context.Background(), remote); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("remote cleanup=%v", err)
	}
	outside := temporary
	outside.Path = filepath.Join(t.TempDir(), name)
	if err := cleanupOwnedTemp(context.Background(), outside); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("outside cleanup=%v", err)
	}
	for _, invalid := range []string{"archive.tar.gz", ".courier-.tar.gz", ".courier-token.tar.gz"} {
		if validOwnedArchiveName(invalid) {
			t.Fatalf("invalid owned archive name accepted: %q", invalid)
		}
	}
	if !validOwnedArchiveName(name) {
		t.Fatalf("valid owned archive name rejected: %q", name)
	}

	want := errors.New("failure")
	openTemporaryRoot = func(string) (temporaryRoot, error) { return nil, want }
	if err := cleanupOwnedTemp(context.Background(), temporary); !errors.Is(err, want) {
		t.Fatalf("open error=%v", err)
	}
	for _, test := range []struct {
		name string
		root *fakeTemporaryRoot
	}{
		{"lstat", &fakeTemporaryRoot{lstatErr: want}},
		{"non-regular", &fakeTemporaryRoot{info: fakeProcessInfo{mode: fs.ModeDir}}},
		{"remove", &fakeTemporaryRoot{info: fakeProcessInfo{mode: 0o600}, removeErr: want}},
		{"close", &fakeTemporaryRoot{lstatErr: fs.ErrNotExist, closeErr: want}},
	} {
		t.Run(test.name, func(t *testing.T) {
			openTemporaryRoot = func(string) (temporaryRoot, error) { return test.root, nil }
			if err := cleanupOwnedTemp(context.Background(), temporary); err == nil {
				t.Fatal("cleanup failure ignored")
			}
		})
	}
}

func TestProcessConfigFiles(t *testing.T) {
	directory := t.TempDir()
	config := validProcessConfig(t, directory, "127.0.0.1:8080")
	path, err := writeProcessConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(path) })
	info, err := os.Stat(path)
	if err != nil || runtime.GOOS != "windows" && info.Mode().Perm() != 0o600 {
		t.Fatalf("mode=%v err=%v", info.Mode(), err)
	}
	loaded, err := readProcessConfig(path)
	if err != nil || loaded.Launch != config.Launch {
		t.Fatalf("loaded=%+v err=%v", loaded, err)
	}
	if _, err := readProcessConfig(""); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("empty=%v", err)
	}
	if _, err := readProcessConfig(filepath.Join(directory, "missing")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing=%v", err)
	}

	write := func(name string, data []byte, mode fs.FileMode) string {
		file := filepath.Join(directory, name)
		if err := os.WriteFile(file, data, mode); err != nil {
			t.Fatal(err)
		}
		return file
	}
	validData, _ := os.ReadFile(path)
	for _, test := range []struct {
		name string
		data []byte
		mode fs.FileMode
	}{
		{"invalid", []byte("{"), 0o600},
		{"unknown", append(bytes.TrimSpace(validData), []byte{}...), 0o600},
		{"trailing", append(validData, []byte("{}")...), 0o600},
		{"bad config", []byte(`{"stateDirectory":"","launch":{}}`), 0o600},
		{"insecure", validData, 0o644},
	} {
		t.Run(test.name, func(t *testing.T) {
			data := test.data
			if test.name == "unknown" {
				data = []byte(`{"stateDirectory":"x","launch":{},"unknown":true}`)
			}
			file := write(test.name, data, test.mode)
			_, err := readProcessConfig(file)
			if err == nil && (runtime.GOOS != "windows" || test.name != "insecure") {
				t.Fatal("invalid configuration accepted")
			}
		})
	}
	directoryPath := filepath.Join(directory, "directory")
	if err := os.Mkdir(directoryPath, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := readProcessConfig(directoryPath); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("directory=%v", err)
	}
	if runtime.GOOS != "windows" {
		symlink := filepath.Join(directory, "link")
		if err := os.Symlink(path, symlink); err != nil {
			t.Fatal(err)
		}
		if _, err := readProcessConfig(symlink); !errors.Is(err, delivery.ErrInvalid) {
			t.Fatalf("symlink=%v", err)
		}
	}
}

func TestProcessConfigInjectedFailures(t *testing.T) {
	originalCreate, originalLstat, originalRead, originalRemove := createProcessFile, lstatProcessFile, readProcessFile, removeProcessFile
	t.Cleanup(func() {
		createProcessFile, lstatProcessFile, readProcessFile, removeProcessFile = originalCreate, originalLstat, originalRead, originalRemove
	})
	config := validProcessConfig(t, t.TempDir(), "127.0.0.1:8080")
	want := errors.New("failure")
	createProcessFile = func(string, string) (processConfigFile, error) { return nil, want }
	if _, err := writeProcessConfig(config); !errors.Is(err, want) {
		t.Fatalf("create=%v", err)
	}
	for _, test := range []struct {
		name string
		file *fakeProcessFile
	}{
		{"chmod", &fakeProcessFile{name: "config", chmodErr: want}},
		{"encode", &fakeProcessFile{name: "config", writeErr: want}},
		{"sync", &fakeProcessFile{name: "config", syncErr: want}},
		{"close", &fakeProcessFile{name: "config", closeErr: want}},
	} {
		t.Run(test.name, func(t *testing.T) {
			removed := ""
			createProcessFile = func(string, string) (processConfigFile, error) { return test.file, nil }
			removeProcessFile = func(path string) error { removed = path; return nil }
			if _, err := writeProcessConfig(config); !errors.Is(err, want) || removed != "config" {
				t.Fatalf("error=%v removed=%q", err, removed)
			}
		})
	}
	createProcessFile, removeProcessFile = originalCreate, originalRemove
	lstatProcessFile = func(string) (fs.FileInfo, error) { return fakeProcessInfo{mode: 0o600}, nil }
	readProcessFile = func(string) ([]byte, error) { return nil, want }
	if _, err := readProcessConfig("config"); !errors.Is(err, want) {
		t.Fatalf("read=%v", err)
	}
}

func TestProcessLauncherFailures(t *testing.T) {
	originalExecutable, originalCreate, originalTimeout, originalPoll := executablePath, createProcessFile, readyTimeoutDefault, readyPollInterval
	t.Cleanup(func() {
		executablePath, createProcessFile, readyTimeoutDefault, readyPollInterval = originalExecutable, originalCreate, originalTimeout, originalPoll
	})
	directory := t.TempDir()
	config := validProcessConfig(t, directory, "127.0.0.1:8080")
	launcher := ProcessLauncher{StateDirectory: directory}
	invalid := config.Launch
	invalid.Bind = "bad"
	if _, err := launcher.Launch(context.Background(), invalid); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid=%v", err)
	}
	want := errors.New("failure")
	executablePath = func() (string, error) { return "", want }
	if _, err := launcher.Launch(context.Background(), config.Launch); !errors.Is(err, want) {
		t.Fatalf("executable=%v", err)
	}
	executablePath = originalExecutable
	createProcessFile = func(string, string) (processConfigFile, error) { return nil, want }
	if _, err := launcher.Launch(context.Background(), config.Launch); !errors.Is(err, want) {
		t.Fatalf("config write=%v", err)
	}
	createProcessFile = originalCreate
	launcher.Executable = filepath.Join(directory, "missing-executable")
	if _, err := launcher.Launch(context.Background(), config.Launch); err == nil {
		t.Fatal("missing executable started")
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	readyTimeoutDefault = 30 * time.Millisecond
	readyPollInterval = 5 * time.Millisecond
	launcher = ProcessLauncher{StateDirectory: directory, Executable: executable, Arguments: []string{"-test.run=TestNoopWorkerProcess", "--"}}
	if _, err := launcher.Launch(context.Background(), config.Launch); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("timeout=%v", err)
	}
}

func TestNoopWorkerProcess(t *testing.T) {}

func TestProcessLauncherEndToEnd(t *testing.T) {
	if hasArgument(os.Args, internalWorkerArgument) {
		handled, err := RunInternal(context.Background(), []string{internalWorkerArgument})
		if !handled || err != nil {
			t.Fatalf("child handled=%v err=%v", handled, err)
		}
		return
	}
	directory := shortWorkerTempDir(t, "courier-process-")
	store, err := delivery.OpenStore(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	probe, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	bind := probe.Addr().String()
	_ = probe.Close()
	config := validProcessConfig(t, directory, bind)
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	launcher := ProcessLauncher{StateDirectory: directory, Executable: executable, Arguments: []string{"-test.run=TestProcessLauncherEndToEnd", "--"}, ReadyTimeout: 5 * time.Second}
	client, err := launcher.Launch(context.Background(), config.Launch)
	if err != nil {
		t.Fatal(err)
	}
	item := validWorkerDelivery()
	item.ServerID = config.Launch.ServerID
	item.CreatedAt = time.Now().UTC()
	item.UpdatedAt = item.CreatedAt
	source := t.TempDir()
	if err := os.WriteFile(filepath.Join(source, "background.txt"), []byte("background"), 0o600); err != nil {
		t.Fatal(err)
	}
	backgroundToken, err := webdelivery.NewToken(strings.NewReader(strings.Repeat("b", webdelivery.ResourceTokenBytes)))
	if err != nil {
		t.Fatal(err)
	}
	backgroundDefinition, _ := json.Marshal(webdelivery.Definition{
		Version: webdelivery.DefinitionVersion, Token: backgroundToken, Source: source, Destination: "web://",
	})
	if _, lease, err := client.RegisterDefinition(context.Background(), item, false, backgroundDefinition); err != nil || lease != nil {
		t.Fatal(err)
	}
	httpClient := &http.Client{Timeout: 2 * time.Second}
	fetch := func(token, name string) (int, string) {
		response, requestErr := httpClient.Get("http://" + bind + "/d/" + token + "/api/v1/download?path=" + name)
		if requestErr != nil {
			t.Fatal(requestErr)
		}
		defer response.Body.Close()
		body, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			t.Fatal(readErr)
		}
		return response.StatusCode, string(body)
	}
	if status, body := fetch(backgroundToken, "background.txt"); status != http.StatusOK || body != "background" {
		t.Fatalf("background delivery=%d %q", status, body)
	}
	foreground := item
	foreground.ID = delivery.ID("00000000-0000-4000-8000-000000000098")
	foreground.CreatedAt = time.Now().UTC()
	foreground.UpdatedAt = foreground.CreatedAt
	foregroundToken, err := webdelivery.NewToken(strings.NewReader(strings.Repeat("f", webdelivery.ResourceTokenBytes)))
	if err != nil {
		t.Fatal(err)
	}
	foregroundDefinition, _ := json.Marshal(webdelivery.Definition{
		Version: webdelivery.DefinitionVersion, Token: foregroundToken, Source: source, Destination: "web://",
	})
	if _, lease, err := client.RegisterDefinition(context.Background(), foreground, true, foregroundDefinition); err != nil || lease == nil {
		t.Fatalf("foreground registration lease=%v err=%v", lease, err)
	} else if err := lease.Release(context.Background()); err != nil {
		t.Fatal(err)
	}
	if status, _ := fetch(foregroundToken, "background.txt"); status != http.StatusNotFound {
		t.Fatalf("released foreground delivery=%d", status)
	}
	if status, body := fetch(backgroundToken, "background.txt"); status != http.StatusOK || body != "background" {
		t.Fatalf("background delivery after foreground release=%d %q", status, body)
	}
	if err := client.StopServer(context.Background()); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Lstat(config.Launch.ControlEndpoint); os.IsNotExist(err) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("worker process did not remove its control endpoint")
}

func TestRunInternalAndProcessFailures(t *testing.T) {
	if handled, err := RunInternal(context.Background(), nil); handled || err != nil {
		t.Fatalf("nil arguments handled=%v err=%v", handled, err)
	}
	if handled, err := RunInternal(context.Background(), []string{"other"}); handled || err != nil {
		t.Fatalf("other argument handled=%v err=%v", handled, err)
	}
	originalEnvironment, originalRemove, originalControl, originalData, originalStartup := readEnvironment, removeProcessFile, listenControl, listenData, startupTimeout
	t.Cleanup(func() {
		readEnvironment, removeProcessFile, listenControl, listenData, startupTimeout = originalEnvironment, originalRemove, originalControl, originalData, originalStartup
	})
	readEnvironment = func(string) string { return "" }
	if handled, err := RunInternal(context.Background(), []string{internalWorkerArgument}); !handled || err == nil {
		t.Fatalf("missing config handled=%v err=%v", handled, err)
	}
	directory := t.TempDir()
	config := validProcessConfig(t, directory, "127.0.0.1:8080")
	path, err := writeProcessConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	readEnvironment = func(string) string { return path }
	want := errors.New("remove")
	removeProcessFile = func(string) error { return want }
	if handled, err := RunInternal(context.Background(), []string{internalWorkerArgument}); !handled || !errors.Is(err, want) {
		t.Fatalf("remove handled=%v err=%v", handled, err)
	}
	_ = os.Remove(path)
	removeProcessFile = originalRemove
	path, err = writeProcessConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	readEnvironment = func(string) string { return path }
	control := newBlockingListener()
	data := newBlockingListener()
	listenControl = func(string) (net.Listener, error) { return control, nil }
	listenData = func(string) (net.Listener, error) { return data, nil }
	startupTimeout = 10 * time.Millisecond
	if handled, err := RunInternal(context.Background(), []string{internalWorkerArgument}); !handled || err != nil {
		t.Fatalf("valid internal worker handled=%v err=%v", handled, err)
	}
}

func TestRunProcessInjectedPaths(t *testing.T) {
	originalControl, originalData, originalStore, originalStartup := listenControl, listenData, openWorkerStore, startupTimeout
	t.Cleanup(func() {
		listenControl, listenData, openWorkerStore, startupTimeout = originalControl, originalData, originalStore, originalStartup
	})
	directory := t.TempDir()
	config := validProcessConfig(t, directory, "127.0.0.1:8080")
	invalid := config
	invalid.Launch.Bind = "bad"
	if err := runProcess(context.Background(), invalid); !errors.Is(err, delivery.ErrInvalid) {
		t.Fatalf("invalid=%v", err)
	}
	want := errors.New("failure")
	listenControl = func(string) (net.Listener, error) { return nil, want }
	if err := runProcess(context.Background(), config); !errors.Is(err, want) {
		t.Fatalf("control=%v", err)
	}
	control := newBlockingListener()
	listenControl = func(string) (net.Listener, error) { return control, nil }
	listenData = func(string) (net.Listener, error) { return nil, want }
	if err := runProcess(context.Background(), config); !errors.Is(err, want) || !control.closed {
		t.Fatalf("data=%v closed=%v", err, control.closed)
	}
	control = newBlockingListener()
	data := newBlockingListener()
	listenControl = func(string) (net.Listener, error) { return control, nil }
	listenData = func(string) (net.Listener, error) { return data, nil }
	openWorkerStore = func(string) (*delivery.Store, error) { return nil, want }
	if err := runProcess(context.Background(), config); !errors.Is(err, want) || !control.closed || !data.closed {
		t.Fatalf("store=%v control=%v data=%v", err, control.closed, data.closed)
	}

	openWorkerStore = originalStore
	control = newBlockingListener()
	data = newBlockingListener()
	listenControl = func(string) (net.Listener, error) { return control, nil }
	listenData = func(string) (net.Listener, error) { return data, nil }
	startupTimeout = 10 * time.Millisecond
	if err := runProcess(context.Background(), config); err != nil {
		t.Fatalf("startup timeout=%v", err)
	}
	if !control.closed || !data.closed {
		t.Fatalf("startup cleanup control=%v data=%v", control.closed, data.closed)
	}

	closedStore, err := delivery.OpenStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	_ = closedStore.Close()
	control = newBlockingListener()
	data = newBlockingListener()
	listenControl = func(string) (net.Listener, error) { return control, nil }
	listenData = func(string) (net.Listener, error) { return data, nil }
	openWorkerStore = func(string) (*delivery.Store, error) { return closedStore, nil }
	if err := runProcess(context.Background(), config); err == nil {
		t.Fatal("closed runtime store accepted")
	}
	openWorkerStore = originalStore

	directory2 := t.TempDir()
	config2 := validProcessConfig(t, directory2, "127.0.0.1:8081")
	control = newBlockingListener()
	data = newBlockingListener()
	listenControl = func(string) (net.Listener, error) { return control, nil }
	listenData = func(string) (net.Listener, error) { return data, nil }
	startupTimeout = time.Hour
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runProcess(ctx, config2) }()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		store, openErr := delivery.OpenStore(directory2)
		if openErr == nil {
			snapshot, loadErr := store.Load(context.Background())
			_ = store.Close()
			if loadErr == nil && len(snapshot.Servers) == 1 {
				break
			}
		}
		time.Sleep(time.Millisecond)
	}
	cancel()
	if err := <-done; err != nil {
		t.Fatalf("cancellation=%v", err)
	}
}

func TestWorkerEnvironmentAndDetach(t *testing.T) {
	t.Setenv("NPM_TOKEN", "secret")
	t.Setenv("HOME", "/home/test")
	environment := workerEnvironment("config")
	joined := strings.Join(environment, "\n")
	if !strings.Contains(joined, workerConfigEnvironment+"=config") || !strings.Contains(joined, "HOME=/home/test") || strings.Contains(joined, "NPM_TOKEN") {
		t.Fatalf("environment=%v", environment)
	}
	command := exec.Command("ignored")
	configureDetached(command)
	if command.SysProcAttr == nil {
		t.Fatal("detached process attributes missing")
	}
	listener, err := listenData("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
}

func hasArgument(arguments []string, expected string) bool {
	for _, argument := range arguments {
		if argument == expected {
			return true
		}
	}
	return false
}

type fakeProcessFile struct {
	bytes.Buffer
	name     string
	chmodErr error
	writeErr error
	syncErr  error
	closeErr error
}

func (file *fakeProcessFile) Name() string            { return file.name }
func (file *fakeProcessFile) Chmod(fs.FileMode) error { return file.chmodErr }
func (file *fakeProcessFile) Write(data []byte) (int, error) {
	if file.writeErr != nil {
		return 0, file.writeErr
	}
	return file.Buffer.Write(data)
}
func (file *fakeProcessFile) Sync() error  { return file.syncErr }
func (file *fakeProcessFile) Close() error { return file.closeErr }

type fakeProcessInfo struct{ mode fs.FileMode }

func (fakeProcessInfo) Name() string           { return "config" }
func (fakeProcessInfo) Size() int64            { return 0 }
func (info fakeProcessInfo) Mode() fs.FileMode { return info.mode }
func (fakeProcessInfo) ModTime() time.Time     { return workerTime }
func (info fakeProcessInfo) IsDir() bool       { return info.mode.IsDir() }
func (fakeProcessInfo) Sys() any               { return nil }

type fakeTemporaryRoot struct {
	info      fs.FileInfo
	lstatErr  error
	removeErr error
	closeErr  error
}

func (root *fakeTemporaryRoot) Lstat(string) (fs.FileInfo, error) { return root.info, root.lstatErr }
func (root *fakeTemporaryRoot) Remove(string) error               { return root.removeErr }
func (root *fakeTemporaryRoot) Close() error                      { return root.closeErr }

type blockingListener struct {
	closedChannel chan struct{}
	closeOnce     sync.Once
	closed        bool
	closeErr      error
}

func newBlockingListener() *blockingListener {
	return &blockingListener{closedChannel: make(chan struct{})}
}

func (listener *blockingListener) Accept() (net.Conn, error) {
	<-listener.closedChannel
	return nil, net.ErrClosed
}
func (listener *blockingListener) Close() error {
	listener.closeOnce.Do(func() { listener.closed = true; close(listener.closedChannel) })
	return listener.closeErr
}
func (*blockingListener) Addr() net.Addr { return workerAddr("listener") }

var _ io.Reader = (*bytes.Buffer)(nil)

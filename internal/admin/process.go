package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/iwonz/courier/internal/control"
	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
	"github.com/iwonz/courier/internal/worker"
)

const (
	internalArgument    = "_admin"
	configEnvironment   = "COURIER_ADMIN_LAUNCH"
	defaultStartTimeout = 5 * time.Second
)

var (
	adminExecutable      = os.Executable
	adminStartCommand    = func(name string, arguments ...string) *exec.Cmd { return exec.Command(name, arguments...) }
	adminListenTCP       = func(bind string) (net.Listener, error) { return net.Listen("tcp", bind) }
	adminListenControl   = ipc.Listen
	adminOpenStore       = delivery.OpenStore
	adminReadEnvironment = os.Getenv
	adminCreateFile      = func(directory, pattern string) (processFile, error) { return os.CreateTemp(directory, pattern) }
	adminLstatFile       = os.Lstat
	adminReadFile        = os.ReadFile
	adminRemoveFile      = os.Remove
	adminNow             = time.Now
	adminProcessID       = os.Getpid
	adminReadyTimeout    = defaultStartTimeout
	adminPollInterval    = 20 * time.Millisecond
	adminHello           = helloState
	adminShutdown        = shutdownState
	adminWaitReady       = waitForReady
	adminControlEndpoint = ipc.ControlEndpoint
	adminAcquire         = acquireSingleton
	adminRemoveStale     = ipc.RemoveStale
	adminWriteState      = writeState
	adminRemoveState     = removeState
	adminNewAPI          = NewAPI
	adminServeHTTP       = func(server *http.Server, listener net.Listener) error { return server.Serve(listener) }
	adminServeControl    = serveControl
	adminShutdownHTTP    = func(server *http.Server, ctx context.Context) error { return server.Shutdown(ctx) }
)

type processFile interface {
	io.Writer
	Name() string
	Chmod(fs.FileMode) error
	Sync() error
	Close() error
}

type StartRequest struct {
	Bind       string
	Background bool
	Ready      func(State) error
}

type StartResult struct {
	State          State
	AlreadyRunning bool
}

type StopResult struct{ AlreadyStopped bool }

type Manager struct {
	StateDirectory string
	Compatibility  string
	Executable     string
	Arguments      []string
	ReadyTimeout   time.Duration
}

type processConfig struct {
	StateDirectory string `json:"stateDirectory"`
	State          State  `json:"state"`
}

func (config processConfig) Validate() error {
	if strings.TrimSpace(config.StateDirectory) == "" {
		return fmt.Errorf("%w: administration state directory is required", delivery.ErrInvalid)
	}
	if err := config.State.Validate(); err != nil {
		return err
	}
	expected, err := adminControlEndpoint(config.StateDirectory, config.State.ID)
	if err != nil || expected != config.State.ControlEndpoint {
		return fmt.Errorf("%w: administration control endpoint mismatch", delivery.ErrInvalid)
	}
	return nil
}

func (manager Manager) Start(ctx context.Context, request StartRequest) (StartResult, error) {
	bind, err := CanonicalBind(request.Bind)
	if err != nil {
		return StartResult{}, err
	}
	if strings.TrimSpace(manager.StateDirectory) == "" || strings.TrimSpace(manager.Compatibility) == "" || strings.ContainsAny(manager.Compatibility, "\x00\r\n") {
		return StartResult{}, fmt.Errorf("%w: administration manager is invalid", delivery.ErrInvalid)
	}
	if existing, found, probeErr := manager.probe(ctx); probeErr != nil {
		return StartResult{}, probeErr
	} else if found {
		if existing.Bind != bind {
			return StartResult{}, fmt.Errorf("%w on %s", ErrAlreadyRunning, existing.Bind)
		}
		if request.Ready != nil {
			if err := request.Ready(existing); err != nil {
				return StartResult{}, err
			}
		}
		return StartResult{State: existing, AlreadyRunning: true}, nil
	}
	id := delivery.NewID()
	controlEndpoint, err := adminControlEndpoint(manager.StateDirectory, id)
	if err != nil {
		return StartResult{}, err
	}
	config := processConfig{StateDirectory: manager.StateDirectory, State: State{
		SchemaVersion: adminStateSchema, ID: id, Bind: bind, ControlEndpoint: controlEndpoint,
		Compatibility: manager.Compatibility, ProcessID: 1, StartedAt: adminNow().UTC(),
	}}
	if request.Background {
		state, err := manager.launch(ctx, config)
		if err != nil {
			return StartResult{}, err
		}
		if request.Ready != nil {
			if err := request.Ready(state); err != nil {
				return StartResult{}, err
			}
		}
		return StartResult{State: state}, nil
	}
	err = runProcess(ctx, config, request.Ready)
	return StartResult{State: config.State}, err
}

func (manager Manager) Stop(ctx context.Context) (StopResult, error) {
	state, err := LoadState(manager.StateDirectory)
	if errors.Is(err, ErrNotRunning) {
		return StopResult{AlreadyStopped: true}, nil
	}
	if err != nil {
		return StopResult{}, err
	}
	if err := adminHello(ctx, state); err != nil {
		return StopResult{}, err
	}
	if err := adminShutdown(ctx, state); err != nil {
		return StopResult{}, err
	}
	timeout := manager.ReadyTimeout
	if timeout <= 0 {
		timeout = adminReadyTimeout
	}
	wait, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for wait.Err() == nil {
		current, loadErr := LoadState(manager.StateDirectory)
		if errors.Is(loadErr, ErrNotRunning) {
			return StopResult{}, nil
		}
		if loadErr != nil {
			return StopResult{}, loadErr
		}
		if current.ID != state.ID {
			return StopResult{}, fmt.Errorf("%w: administration owner changed during stop", delivery.ErrRevisionConflict)
		}
		select {
		case <-wait.Done():
		case <-time.After(adminPollInterval):
		}
	}
	return StopResult{}, wait.Err()
}

func (manager Manager) probe(ctx context.Context) (State, bool, error) {
	state, err := LoadState(manager.StateDirectory)
	if errors.Is(err, ErrNotRunning) {
		return State{}, false, nil
	}
	if err != nil {
		return State{}, false, err
	}
	if err := adminHello(ctx, state); err != nil {
		return state, false, nil
	}
	return state, true, nil
}

func (manager Manager) launch(ctx context.Context, config processConfig) (State, error) {
	path, err := writeProcessConfig(config)
	if err != nil {
		return State{}, err
	}
	removeConfig := true
	defer func() {
		if removeConfig {
			_ = adminRemoveFile(path)
		}
	}()
	executable := manager.Executable
	if executable == "" {
		executable, err = adminExecutable()
		if err != nil {
			return State{}, err
		}
	}
	arguments := append(append([]string(nil), manager.Arguments...), internalArgument)
	command := adminStartCommand(executable, arguments...)
	command.Env = worker.ProcessEnvironment(configEnvironment, path)
	worker.ConfigureDetached(command)
	if err := command.Start(); err != nil {
		return State{}, err
	}
	timeout := manager.ReadyTimeout
	if timeout <= 0 {
		timeout = adminReadyTimeout
	}
	ready, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	state, readyErr := adminWaitReady(ready, manager.StateDirectory, config.State.ID)
	if readyErr == nil {
		removeConfig = false
		go func() { _ = command.Wait() }()
		return state, nil
	}
	_ = command.Process.Kill()
	waitErr := command.Wait()
	return State{}, errors.Join(readyErr, waitErr)
}

func waitForReady(ctx context.Context, directory string, id delivery.ID) (State, error) {
	for ctx.Err() == nil {
		state, err := LoadState(directory)
		if err == nil {
			if state.ID == id && adminHello(ctx, state) == nil {
				return state, nil
			}
			if state.ID != id && adminHello(ctx, state) == nil {
				return State{}, fmt.Errorf("%w on %s", ErrAlreadyRunning, state.Bind)
			}
		}
		select {
		case <-ctx.Done():
		case <-time.After(adminPollInterval):
		}
	}
	return State{}, ctx.Err()
}

func helloState(ctx context.Context, state State) error {
	client := worker.Client{Endpoint: state.ControlEndpoint, ServerID: state.ID, Compatibility: state.Compatibility}
	_, err := client.Hello(ctx)
	return err
}

func shutdownState(ctx context.Context, state State) error {
	client := worker.Client{Endpoint: state.ControlEndpoint, ServerID: state.ID, Compatibility: state.Compatibility}
	return client.Shutdown(ctx)
}

func writeProcessConfig(config processConfig) (string, error) {
	if err := config.Validate(); err != nil {
		return "", err
	}
	if err := ensurePrivateDirectory(config.StateDirectory); err != nil {
		return "", err
	}
	file, err := adminCreateFile(config.StateDirectory, ".admin-launch-*.json")
	if err != nil {
		return "", err
	}
	path := file.Name()
	success := false
	defer func() {
		_ = file.Close()
		if !success {
			_ = adminRemoveFile(path)
		}
	}()
	if err := file.Chmod(0o600); err != nil {
		return "", err
	}
	if err := json.NewEncoder(file).Encode(config); err != nil {
		return "", err
	}
	if err := file.Sync(); err != nil {
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	success = true
	return path, nil
}

func readProcessConfig(path string) (processConfig, error) {
	if strings.TrimSpace(path) == "" {
		return processConfig{}, fmt.Errorf("%w: administration launch path is required", delivery.ErrInvalid)
	}
	info, err := adminLstatFile(path)
	if err != nil {
		return processConfig{}, err
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 || runtime.GOOS != "windows" && info.Mode().Perm()&0o077 != 0 {
		return processConfig{}, fmt.Errorf("%w: administration launch file is not private", delivery.ErrInvalid)
	}
	data, err := adminReadFile(path)
	if err != nil {
		return processConfig{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var config processConfig
	if err := decoder.Decode(&config); err != nil {
		return processConfig{}, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return processConfig{}, fmt.Errorf("%w: trailing administration configuration", delivery.ErrInvalid)
	}
	if err := config.Validate(); err != nil {
		return processConfig{}, err
	}
	return config, nil
}

func RunInternal(ctx context.Context, arguments []string) (bool, error) {
	if len(arguments) != 1 || arguments[0] != internalArgument {
		return false, nil
	}
	path := adminReadEnvironment(configEnvironment)
	config, err := readProcessConfig(path)
	if path != "" {
		err = errors.Join(err, adminRemoveFile(path))
	}
	if err != nil {
		return true, err
	}
	return true, runProcess(ctx, config, nil)
}

func runProcess(ctx context.Context, config processConfig, ready func(State) error) (resultErr error) {
	if err := config.Validate(); err != nil {
		return err
	}
	lock, err := adminAcquire(config.StateDirectory)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, lock.Close()) }()
	if previous, loadErr := LoadState(config.StateDirectory); loadErr == nil {
		if helloErr := adminHello(ctx, previous); helloErr == nil {
			return ErrAlreadyRunning
		}
		if err := adminRemoveStale(previous.ControlEndpoint); err != nil {
			return err
		}
		if err := adminRemoveState(config.StateDirectory, previous.ID); err != nil {
			return err
		}
	} else if !errors.Is(loadErr, ErrNotRunning) {
		return loadErr
	}
	controlListener, err := adminListenControl(config.State.ControlEndpoint)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, closeListener(controlListener)) }()
	dataListener, err := adminListenTCP(config.State.Bind)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, closeListener(dataListener)) }()
	store, err := adminOpenStore(config.StateDirectory)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, store.Close()) }()
	api, err := adminNewAPI(APIConfig{Host: config.State.Bind, Control: control.Service{Store: store}})
	if err != nil {
		return err
	}
	state := config.State
	state.ProcessID = adminProcessID()
	state.StartedAt = adminNow().UTC()
	if err := adminWriteState(config.StateDirectory, state); err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, adminRemoveState(config.StateDirectory, state.ID)) }()
	if ready != nil {
		if err := ready(state); err != nil {
			return err
		}
	}
	runContext, cancel := context.WithCancel(ctx)
	defer cancel()
	server := &http.Server{
		Handler: api, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		IdleTimeout: 30 * time.Second, MaxHeaderBytes: 32 << 10,
	}
	httpResult := make(chan error, 1)
	controlResult := make(chan error, 1)
	serveHTTP, serveAdminControl := adminServeHTTP, adminServeControl
	go func() { httpResult <- serveHTTP(server, dataListener) }()
	go func() { controlResult <- serveAdminControl(runContext, controlListener, state, cancel) }()
	select {
	case <-runContext.Done():
	case err := <-httpResult:
		resultErr = normalizeServerError(err)
		cancel()
	case err := <-controlResult:
		resultErr = normalizeServerError(err)
		cancel()
	}
	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	resultErr = errors.Join(resultErr, normalizeServerError(adminShutdownHTTP(server, shutdownContext)), closeListener(controlListener))
	return resultErr
}

func serveControl(ctx context.Context, listener net.Listener, state State, shutdown context.CancelFunc) error {
	for {
		connection, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		handleControl(connection, state, shutdown)
	}
}

func handleControl(connection net.Conn, state State, shutdown context.CancelFunc) {
	defer connection.Close()
	_ = connection.SetDeadline(time.Now().Add(5 * time.Second))
	var request ipc.Request
	if err := ipc.ReadFrame(connection, &request); err != nil || request.Validate() != nil {
		return
	}
	var response ipc.Response
	var err error
	switch request.Operation {
	case ipc.OperationHello:
		var hello worker.HelloRequest
		err = ipc.DecodePayload(request.Payload, &hello)
		if err == nil {
			err = hello.Validate()
		}
		if err == nil && (hello.ServerID != state.ID || hello.Compatibility != state.Compatibility) {
			err = delivery.ErrNotFound
		}
		if err == nil {
			response, err = ipc.Success(request, worker.HelloResponse{ServerID: state.ID, Compatibility: state.Compatibility})
		}
	case ipc.OperationShutdown:
		if len(request.Payload) != 0 {
			err = delivery.ErrInvalid
		} else {
			response, err = ipc.Success(request, nil)
		}
	default:
		err = errors.New("unsupported administration operation")
	}
	if err != nil {
		response, err = ipc.Rejection(request, ipc.CodeUnsupported, "request rejected")
	}
	if err == nil && ipc.WriteFrame(connection, response) == nil && request.Operation == ipc.OperationShutdown && response.OK {
		shutdown()
	}
}

func closeListener(listener net.Listener) error {
	if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		return err
	}
	return nil
}

func normalizeServerError(err error) error {
	if errors.Is(err, http.ErrServerClosed) || errors.Is(err, net.ErrClosed) {
		return nil
	}
	return err
}

func URL(state State) string { return "http://" + state.Bind + "/" }

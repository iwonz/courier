package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/ipc"
)

const (
	internalWorkerArgument  = "_worker"
	workerConfigEnvironment = "COURIER_WORKER_LAUNCH"
	defaultStartupTimeout   = 5 * time.Second
)

type processConfigFile interface {
	io.Writer
	Name() string
	Chmod(fs.FileMode) error
	Sync() error
	Close() error
}

var (
	executablePath      = os.Executable
	startCommand        = func(name string, arguments ...string) *exec.Cmd { return exec.Command(name, arguments...) }
	listenData          = func(bind string) (net.Listener, error) { return net.Listen("tcp", bind) }
	newProcessHost      = newWorkerWebHost
	openWorkerStore     = delivery.OpenStore
	listenControl       = ipc.Listen
	readEnvironment     = os.Getenv
	createProcessFile   = func(directory, pattern string) (processConfigFile, error) { return os.CreateTemp(directory, pattern) }
	lstatProcessFile    = os.Lstat
	readProcessFile     = os.ReadFile
	removeProcessFile   = os.Remove
	startupTimeout      = defaultStartupTimeout
	readyTimeoutDefault = defaultStartupTimeout
	readyPollInterval   = 20 * time.Millisecond
)

type processConfig struct {
	StateDirectory string        `json:"stateDirectory"`
	Launch         LaunchRequest `json:"launch"`
}

func (config processConfig) Validate() error {
	if strings.TrimSpace(config.StateDirectory) == "" || !config.Launch.ServerID.Valid() || invalidCompatibility(config.Launch.Compatibility) {
		return fmt.Errorf("%w: invalid worker process configuration", delivery.ErrInvalid)
	}
	bind, err := CanonicalBind(config.Launch.Bind)
	if err != nil || bind != config.Launch.Bind {
		return fmt.Errorf("%w: worker bind is not canonical", delivery.ErrInvalid)
	}
	expected, err := ipc.ControlEndpoint(config.StateDirectory, config.Launch.ServerID)
	if err != nil || expected != config.Launch.ControlEndpoint {
		return fmt.Errorf("%w: worker control endpoint mismatch", delivery.ErrInvalid)
	}
	return nil
}

type ProcessLauncher struct {
	StateDirectory string
	Executable     string
	Arguments      []string
	ReadyTimeout   time.Duration
}

func (launcher ProcessLauncher) Launch(ctx context.Context, request LaunchRequest) (Client, error) {
	config := processConfig{StateDirectory: launcher.StateDirectory, Launch: request}
	if err := config.Validate(); err != nil {
		return Client{}, err
	}
	executable := launcher.Executable
	if executable == "" {
		var err error
		executable, err = executablePath()
		if err != nil {
			return Client{}, err
		}
	}
	configPath, err := writeProcessConfig(config)
	if err != nil {
		return Client{}, err
	}
	removeConfig := true
	defer func() {
		if removeConfig {
			_ = removeProcessFile(configPath)
		}
	}()
	arguments := append(append([]string(nil), launcher.Arguments...), internalWorkerArgument)
	command := startCommand(executable, arguments...)
	command.Env = workerEnvironment(configPath)
	configureDetached(command)
	if err := command.Start(); err != nil {
		return Client{}, err
	}
	client := Client{Endpoint: request.ControlEndpoint, ServerID: request.ServerID, Compatibility: request.Compatibility}
	timeout := launcher.ReadyTimeout
	if timeout <= 0 {
		timeout = readyTimeoutDefault
	}
	readyContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	var probeErr error
	for readyContext.Err() == nil {
		_, probeErr = client.Hello(readyContext)
		if probeErr == nil {
			removeConfig = false
			go func() { _ = command.Wait() }()
			return client, nil
		}
		select {
		case <-readyContext.Done():
		case <-time.After(readyPollInterval):
		}
	}
	_ = command.Process.Kill()
	waitErr := command.Wait()
	return Client{}, errors.Join(readyContext.Err(), probeErr, waitErr)
}

func writeProcessConfig(config processConfig) (string, error) {
	file, err := createProcessFile(config.StateDirectory, ".worker-launch-*.json")
	if err != nil {
		return "", err
	}
	path := file.Name()
	success := false
	defer func() {
		_ = file.Close()
		if !success {
			_ = removeProcessFile(path)
		}
	}()
	if err := file.Chmod(0o600); err != nil {
		return "", err
	}
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
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
		return processConfig{}, fmt.Errorf("%w: worker launch path is required", delivery.ErrInvalid)
	}
	info, err := lstatProcessFile(path)
	if err != nil {
		return processConfig{}, err
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 || runtime.GOOS != "windows" && info.Mode().Perm()&0o077 != 0 {
		return processConfig{}, fmt.Errorf("%w: worker launch file is not private", delivery.ErrInvalid)
	}
	data, err := readProcessFile(path)
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
		return processConfig{}, fmt.Errorf("%w: trailing worker configuration", delivery.ErrInvalid)
	}
	if err := config.Validate(); err != nil {
		return processConfig{}, err
	}
	return config, nil
}

func RunInternal(ctx context.Context, arguments []string) (bool, error) {
	if len(arguments) != 1 || arguments[0] != internalWorkerArgument {
		return false, nil
	}
	path := readEnvironment(workerConfigEnvironment)
	config, err := readProcessConfig(path)
	if path != "" {
		err = errors.Join(err, removeProcessFile(path))
	}
	if err != nil {
		return true, err
	}
	return true, runProcess(ctx, config)
}

func runProcess(ctx context.Context, config processConfig) (resultErr error) {
	if err := config.Validate(); err != nil {
		return err
	}
	control, err := listenControl(config.Launch.ControlEndpoint)
	if err != nil {
		return err
	}
	controlOwned := true
	defer func() {
		if controlOwned {
			resultErr = errors.Join(resultErr, control.Close())
		}
	}()
	data, err := listenData(config.Launch.Bind)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, data.Close()) }()
	store, err := openWorkerStore(config.StateDirectory)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, store.Close()) }()
	now := time.Now().UTC()
	server := delivery.Server{
		ID: config.Launch.ServerID, Bind: config.Launch.Bind, ControlEndpoint: config.Launch.ControlEndpoint,
		ProcessID: os.Getpid(), State: delivery.StateStarting, StartedAt: now, UpdatedAt: now,
	}
	runtime, err := NewRuntime(ctx, store, server, config.Launch.Compatibility)
	if err != nil {
		return err
	}
	host, err := newProcessHost(func(id delivery.ID) {
		go func() { _, _ = runtime.StopDelivery(context.Background(), id, delivery.ReasonStopped) }()
	})
	if err != nil {
		return err
	}
	if host == nil {
		return fmt.Errorf("%w: process delivery host is required", delivery.ErrInvalid)
	}
	_ = runtime.AttachHost(host)
	serveResult := make(chan error, 1)
	go func() { serveResult <- host.Serve(data) }()
	defer func() {
		resultErr = errors.Join(resultErr, host.Close(context.Background()), <-serveResult)
	}()
	controlOwned = false
	startup := time.AfterFunc(startupTimeout, func() {
		if !runtime.HasDeliveries() {
			_ = runtime.StopServer(context.Background(), delivery.ReasonStale)
		}
	})
	defer startup.Stop()
	go func() {
		select {
		case <-ctx.Done():
			_ = runtime.StopServer(context.Background(), delivery.ReasonStopped)
		case <-runtime.Done():
		}
	}()
	return runtime.Run(control)
}

func workerEnvironment(configPath string) []string {
	allowed := map[string]bool{
		"HOME": true, "USERPROFILE": true, "TMPDIR": true, "TMP": true, "TEMP": true,
		"SYSTEMROOT": true, "WINDIR": true, "LANG": true, "LC_ALL": true, "SSH_AUTH_SOCK": true,
	}
	result := []string{workerConfigEnvironment + "=" + configPath}
	for _, entry := range os.Environ() {
		name, _, found := strings.Cut(entry, "=")
		if found && allowed[strings.ToUpper(name)] {
			result = append(result, entry)
		}
	}
	return result
}

func DefaultCoordinator(store *delivery.Store, stateDirectory string) *Coordinator {
	launcher := ProcessLauncher{StateDirectory: stateDirectory}
	return &Coordinator{
		Store: store, StateDirectory: stateDirectory, Locks: FileBindLocker{Directory: stateDirectory},
		Launch: launcher.Launch, Cleanup: cleanupOwnedTemp,
	}
}

type temporaryRoot interface {
	Lstat(string) (fs.FileInfo, error)
	Remove(string) error
	Close() error
}

var (
	ownedTempDirectory = os.TempDir
	openTemporaryRoot  = func(path string) (temporaryRoot, error) { return os.OpenRoot(path) }
)

func cleanupOwnedTemp(ctx context.Context, temporary delivery.OwnedTemp) (resultErr error) {
	if err := ctx.Err(); err != nil {
		return err
	}
	if temporary.Location != delivery.TempLocal {
		return fmt.Errorf("%w: remote temporary cleanup requires its owning transport", delivery.ErrInvalid)
	}
	directory := filepath.Clean(ownedTempDirectory())
	path := filepath.Clean(temporary.Path)
	name := filepath.Base(path)
	if filepath.Dir(path) != directory || !validOwnedArchiveName(name) {
		return fmt.Errorf("%w: temporary archive is outside the private cache", delivery.ErrInvalid)
	}
	root, err := openTemporaryRoot(directory)
	if err != nil {
		return err
	}
	defer func() { resultErr = errors.Join(resultErr, root.Close()) }()
	info, err := root.Lstat(name)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%w: temporary archive is not a regular file", delivery.ErrInvalid)
	}
	return root.Remove(name)
}

func validOwnedArchiveName(name string) bool {
	const prefix, suffix = ".courier-", ".tar.gz"
	if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, suffix) {
		return false
	}
	random := strings.TrimSuffix(strings.TrimPrefix(name, prefix), suffix)
	if random == "" {
		return false
	}
	for _, character := range random {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

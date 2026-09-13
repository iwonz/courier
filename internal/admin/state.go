package admin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/netip"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/iwonz/courier/internal/delivery"
)

const adminStateSchema = 1

var (
	ErrAlreadyRunning = errors.New("Courier administration UI is already running")
	ErrNotRunning     = errors.New("Courier administration UI is not running")
	stateMkdirAll     = os.MkdirAll
	stateLstat        = os.Lstat
	stateReadFile     = os.ReadFile
	stateCreateFile   = func(directory, pattern string) (stateFile, error) { return os.CreateTemp(directory, pattern) }
	stateRename       = os.Rename
	stateRemove       = os.Remove
	stateChmod        = os.Chmod
	newSingletonFile  = func(path string) lockFile { return flock.New(path, flock.SetPermissions(0o600)) }
)

type stateFile interface {
	io.Writer
	Name() string
	Chmod(fs.FileMode) error
	Sync() error
	Close() error
}

type State struct {
	SchemaVersion   int         `json:"schemaVersion"`
	ID              delivery.ID `json:"id"`
	Bind            string      `json:"bind"`
	ControlEndpoint string      `json:"controlEndpoint"`
	Compatibility   string      `json:"compatibility"`
	ProcessID       int         `json:"processId"`
	StartedAt       time.Time   `json:"startedAt"`
}

func (state State) Validate() error {
	if state.SchemaVersion != adminStateSchema || !state.ID.Valid() || state.ProcessID <= 0 || state.StartedAt.IsZero() || strings.TrimSpace(state.ControlEndpoint) == "" || strings.ContainsAny(state.ControlEndpoint, "\x00\r\n") || strings.TrimSpace(state.Compatibility) == "" || strings.ContainsAny(state.Compatibility, "\x00\r\n") {
		return fmt.Errorf("%w: invalid administration state", delivery.ErrInvalid)
	}
	bind, err := CanonicalBind(state.Bind)
	if err != nil || bind != state.Bind {
		return fmt.Errorf("%w: administration bind is not canonical", delivery.ErrInvalid)
	}
	return nil
}

func CanonicalBind(value string) (string, error) {
	host, port, err := net.SplitHostPort(value)
	if err != nil {
		return "", fmt.Errorf("%w: administration listen address must use IP:port", delivery.ErrInvalid)
	}
	address, addressErr := netip.ParseAddr(host)
	number, portErr := strconv.ParseUint(port, 10, 16)
	if addressErr != nil || !address.IsLoopback() || portErr != nil || number == 0 {
		return "", fmt.Errorf("%w: administration listen address must be loopback with a non-zero port", delivery.ErrInvalid)
	}
	return net.JoinHostPort(address.String(), strconv.FormatUint(number, 10)), nil
}

func StatePath(directory string) string { return filepath.Join(directory, "admin.json") }

func LoadState(directory string) (State, error) {
	path := StatePath(directory)
	info, err := stateLstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return State{}, ErrNotRunning
	}
	if err != nil {
		return State{}, err
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 || runtime.GOOS != "windows" && info.Mode().Perm()&0o077 != 0 {
		return State{}, fmt.Errorf("%w: administration state file is not private", delivery.ErrInvalid)
	}
	data, err := stateReadFile(path)
	if err != nil {
		return State{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var state State
	if err := decoder.Decode(&state); err != nil {
		return State{}, fmt.Errorf("%w: decode administration state: %v", delivery.ErrInvalid, err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return State{}, fmt.Errorf("%w: trailing administration state", delivery.ErrInvalid)
	}
	if err := state.Validate(); err != nil {
		return State{}, err
	}
	return state, nil
}

func writeState(directory string, state State) (resultErr error) {
	if err := state.Validate(); err != nil {
		return err
	}
	if err := ensurePrivateDirectory(directory); err != nil {
		return err
	}
	file, err := stateCreateFile(directory, ".admin-state-*.json")
	if err != nil {
		return err
	}
	temporary := file.Name()
	committed := false
	closed := false
	defer func() {
		if !closed {
			resultErr = errors.Join(resultErr, file.Close())
		}
		if !committed {
			resultErr = errors.Join(resultErr, stateRemove(temporary))
		}
	}()
	if err := file.Chmod(0o600); err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(state); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	err = file.Close()
	closed = true
	if err != nil {
		return err
	}
	if err := stateRename(temporary, StatePath(directory)); err != nil {
		return err
	}
	committed = true
	return nil
}

func removeState(directory string, owner delivery.ID) error {
	state, err := LoadState(directory)
	if errors.Is(err, ErrNotRunning) {
		return nil
	}
	if err != nil {
		return err
	}
	if state.ID != owner {
		return fmt.Errorf("%w: administration state owner changed", delivery.ErrRevisionConflict)
	}
	return stateRemove(StatePath(directory))
}

func ensurePrivateDirectory(directory string) error {
	if strings.TrimSpace(directory) == "" {
		return fmt.Errorf("%w: administration state directory is required", delivery.ErrInvalid)
	}
	if err := stateMkdirAll(directory, 0o700); err != nil {
		return err
	}
	info, err := stateLstat(directory)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&fs.ModeSymlink != 0 {
		return fmt.Errorf("%w: administration state path is not a directory", delivery.ErrInvalid)
	}
	return stateChmod(directory, 0o700)
}

type lockFile interface {
	TryLock() (bool, error)
	Unlock() error
	Close() error
}

type singletonLock struct{ file lockFile }

func acquireSingleton(directory string) (*singletonLock, error) {
	if err := ensurePrivateDirectory(directory); err != nil {
		return nil, err
	}
	file := newSingletonFile(filepath.Join(directory, "admin.lock"))
	locked, err := file.TryLock()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	if !locked {
		_ = file.Close()
		return nil, ErrAlreadyRunning
	}
	return &singletonLock{file: file}, nil
}

func (lock *singletonLock) Close() error {
	return errors.Join(lock.file.Unlock(), lock.file.Close())
}

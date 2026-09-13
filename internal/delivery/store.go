package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	registryPrefix = "registry-"
	historyPrefix  = "history-"
)

var (
	registryNamePattern = regexp.MustCompile(`^registry-([0-9]{20})[.]json$`)
	historyNamePattern  = regexp.MustCompile(`^history-([0-9]{20})-([0-9a-f-]{36})[.]json$`)
	makeStateDirectory  = os.MkdirAll
	lstatStateDirectory = os.Lstat
	chmodStateDirectory = os.Chmod
	encodeRegistryState = encodeJSON
	encodeHistoryState  = encodeJSON
	userConfigDirectory = os.UserConfigDir
	openStateRoot       = func(name string) (stateRoot, error) {
		root, err := os.OpenRoot(name)
		if err != nil {
			return nil, err
		}
		return &osStateRoot{root: root}, nil
	}
)

type stateFile interface {
	io.WriteCloser
	Sync() error
}

type stateRoot interface {
	ReadDir() ([]fs.DirEntry, error)
	Lstat(string) (fs.FileInfo, error)
	ReadFile(string) ([]byte, error)
	CreateExclusive(string, fs.FileMode) (stateFile, error)
	CommitAbsent(string, string) error
	Remove(string) error
	Sync() error
	Close() error
}

type osStateRoot struct{ root *os.Root }

func (root *osStateRoot) ReadDir() ([]fs.DirEntry, error)        { return fs.ReadDir(root.root.FS(), ".") }
func (root *osStateRoot) Lstat(name string) (fs.FileInfo, error) { return root.root.Lstat(name) }
func (root *osStateRoot) ReadFile(name string) ([]byte, error)   { return root.root.ReadFile(name) }
func (root *osStateRoot) CreateExclusive(name string, mode fs.FileMode) (stateFile, error) {
	return root.root.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, mode.Perm())
}
func (root *osStateRoot) CommitAbsent(temporary, final string) error {
	return root.root.Link(temporary, final)
}
func (root *osStateRoot) Remove(name string) error { return root.root.Remove(name) }
func (root *osStateRoot) Sync() error {
	directory, err := root.root.Open(".")
	if err != nil {
		return err
	}
	return errors.Join(directory.Sync(), directory.Close())
}
func (root *osStateRoot) Close() error { return root.root.Close() }

type Store struct {
	directory string
	root      stateRoot
	mutex     sync.Mutex
}

func OpenStore(directory string) (*Store, error) {
	if strings.TrimSpace(directory) == "" {
		return nil, fmt.Errorf("%w: state directory is required", ErrInvalid)
	}
	if info, err := lstatStateDirectory(directory); err == nil {
		if !info.IsDir() || info.Mode()&fs.ModeSymlink != 0 {
			return nil, fmt.Errorf("%w: state path is not a real directory", ErrInvalid)
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	if err := makeStateDirectory(directory, 0o700); err != nil {
		return nil, err
	}
	info, err := lstatStateDirectory(directory)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() || info.Mode()&fs.ModeSymlink != 0 {
		return nil, fmt.Errorf("%w: state path changed during creation", ErrInvalid)
	}
	if err := chmodStateDirectory(directory, 0o700); err != nil {
		return nil, err
	}
	root, err := openStateRoot(directory)
	if err != nil {
		return nil, err
	}
	return &Store{directory: directory, root: root}, nil
}

func (store *Store) Directory() string { return store.directory }

func (store *Store) Close() error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if store.root == nil {
		return nil
	}
	err := store.root.Close()
	store.root = nil
	return err
}

func (store *Store) Load(ctx context.Context) (Snapshot, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return store.load(ctx)
}

func (store *Store) load(ctx context.Context) (Snapshot, error) {
	if err := store.ready(ctx); err != nil {
		return Snapshot{}, err
	}
	entries, err := store.root.ReadDir()
	if err != nil {
		return Snapshot{}, err
	}
	var selected string
	var revision uint64
	for _, entry := range entries {
		matches := registryNamePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}
		value, err := strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			return Snapshot{}, fmt.Errorf("%w: invalid registry filename", ErrInvalid)
		}
		if selected == "" || value > revision {
			selected, revision = entry.Name(), value
		}
	}
	if selected == "" {
		return EmptySnapshot(), nil
	}
	data, err := store.readPrivate(selected)
	if err != nil {
		return Snapshot{}, err
	}
	var snapshot Snapshot
	if err := decodeStrict(data, &snapshot); err != nil {
		return Snapshot{}, err
	}
	if snapshot.Revision != revision {
		return Snapshot{}, fmt.Errorf("%w: registry filename revision mismatch", ErrInvalid)
	}
	if err := snapshot.Validate(); err != nil {
		return Snapshot{}, err
	}
	return cloneSnapshot(snapshot), nil
}

func (store *Store) Update(ctx context.Context, expected uint64, mutate func(*Registry) error) (Snapshot, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if mutate == nil {
		return Snapshot{}, fmt.Errorf("%w: registry mutation is required", ErrInvalid)
	}
	current, err := store.load(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	if current.Revision != expected {
		return Snapshot{}, fmt.Errorf("%w: expected %d, current %d", ErrRevisionConflict, expected, current.Revision)
	}
	if expected == math.MaxUint64 {
		return Snapshot{}, fmt.Errorf("%w: registry revision overflow", ErrInvalid)
	}
	registry := &Registry{snapshot: cloneSnapshot(current)}
	if err := mutate(registry); err != nil {
		return Snapshot{}, err
	}
	next := registry.Snapshot()
	next.Revision = expected + 1
	if err := next.Validate(); err != nil {
		return Snapshot{}, err
	}
	data, err := encodeRegistryState(next)
	if err != nil {
		return Snapshot{}, err
	}
	if err := store.writeImmutable(ctx, registryFilename(next.Revision), ".registry-", data); err != nil {
		return Snapshot{}, err
	}
	return cloneSnapshot(next), nil
}

func (store *Store) AppendHistory(ctx context.Context, event HistoryEvent) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if err := event.Validate(); err != nil {
		return err
	}
	data, err := encodeHistoryState(event)
	if err != nil {
		return err
	}
	return store.writeImmutable(ctx, historyFilename(event), ".history-", data)
}

func (store *Store) History(ctx context.Context) ([]HistoryEvent, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if err := store.ready(ctx); err != nil {
		return nil, err
	}
	entries, err := store.root.ReadDir()
	if err != nil {
		return nil, err
	}
	result := make([]HistoryEvent, 0)
	for _, entry := range entries {
		matches := historyNamePattern.FindStringSubmatch(entry.Name())
		if matches == nil {
			continue
		}
		data, err := store.readPrivate(entry.Name())
		if err != nil {
			return nil, err
		}
		var event HistoryEvent
		if err := decodeStrict(data, &event); err != nil {
			return nil, err
		}
		if event.ID != ID(matches[2]) || event.At.UnixNano() < 0 || fmt.Sprintf("%020d", event.At.UnixNano()) != matches[1] {
			return nil, fmt.Errorf("%w: history filename metadata mismatch", ErrInvalid)
		}
		if err := event.Validate(); err != nil {
			return nil, err
		}
		result = append(result, event)
	}
	sort.Slice(result, func(left, right int) bool {
		if result[left].At.Equal(result[right].At) {
			return result[left].ID < result[right].ID
		}
		return result[left].At.Before(result[right].At)
	})
	return result, nil
}

func (store *Store) writeImmutable(ctx context.Context, final, temporaryPrefix string, data []byte) (resultErr error) {
	if err := store.ready(ctx); err != nil {
		return err
	}
	temporary := temporaryPrefix + string(NewID()) + ".tmp"
	file, err := store.root.CreateExclusive(temporary, 0o600)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			resultErr = errors.Join(resultErr, file.Close())
		}
		if removeErr := store.root.Remove(temporary); removeErr != nil && !errors.Is(removeErr, fs.ErrNotExist) {
			resultErr = errors.Join(resultErr, removeErr)
		}
	}()
	if _, err := file.Write(data); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	closed = true
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := store.root.CommitAbsent(temporary, final); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("%w: %s", ErrRevisionConflict, final)
		}
		return err
	}
	if err := store.root.Remove(temporary); err != nil {
		return err
	}
	if err := store.root.Sync(); err != nil {
		return err
	}
	return nil
}

func (store *Store) readPrivate(name string) ([]byte, error) {
	info, err := store.root.Lstat(name)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() || runtime.GOOS != "windows" && info.Mode().Perm()&0o077 != 0 {
		return nil, fmt.Errorf("%w: state object %s is not a private regular file", ErrInvalid, name)
	}
	return store.root.ReadFile(name)
}

func (store *Store) ready(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if store.root == nil {
		return fmt.Errorf("%w: store is closed", ErrInvalid)
	}
	return nil
}

func registryFilename(revision uint64) string {
	return fmt.Sprintf("%s%020d.json", registryPrefix, revision)
}

func historyFilename(event HistoryEvent) string {
	return fmt.Sprintf("%s%020d-%s.json", historyPrefix, event.At.UnixNano(), event.ID)
}

func encodeJSON(value any) ([]byte, error) {
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func decodeStrict(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("%w: decode state: %v", ErrInvalid, err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: trailing state data", ErrInvalid)
	}
	return nil
}

func DefaultStateDirectory() (string, error) {
	root, err := userConfigDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "courier", "state"), nil
}

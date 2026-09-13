// Package fsx defines the filesystem boundary shared by transfer transports.
package fsx

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// ErrDestinationExists identifies a no-replace commit collision.
var ErrDestinationExists = errors.New("destination already exists")

// Writable is a destination file that can flush confirmed data.
type Writable interface {
	io.WriteCloser
	Sync() error
}

// Backend is the minimal filesystem contract required by Courier.
type Backend interface {
	Lstat(string) (fs.FileInfo, error)
	ReadDir(string) ([]fs.DirEntry, error)
	Open(string) (io.ReadCloser, error)
	Create(string, fs.FileMode) (Writable, error)
	MkdirAll(string, fs.FileMode) error
	RemoveAll(string) error
	Rename(string, string) error
	Chmod(string, fs.FileMode) error
	Chtimes(string, time.Time, time.Time) error
	Symlink(string, string) error
	Readlink(string) (string, error)
	Join(...string) string
	Dir(string) string
	CommitAbsent(string, string) error
}

// Local implements Backend with the Go standard library.
type Local struct{}

func (Local) Lstat(name string) (fs.FileInfo, error)       { return os.Lstat(name) }
func (Local) ReadDir(name string) ([]fs.DirEntry, error)   { return os.ReadDir(name) }
func (Local) Open(name string) (io.ReadCloser, error)      { return os.Open(name) }
func (Local) MkdirAll(name string, mode fs.FileMode) error { return os.MkdirAll(name, mode) }
func (Local) RemoveAll(name string) error                  { return os.RemoveAll(name) }
func (Local) Rename(oldPath, newPath string) error         { return os.Rename(oldPath, newPath) }
func (Local) Link(oldPath, newPath string) error           { return os.Link(oldPath, newPath) }
func (Local) Mkdir(name string, mode fs.FileMode) error    { return os.Mkdir(name, mode) }
func (Local) Remove(name string) error                     { return os.Remove(name) }
func (Local) Chmod(name string, mode fs.FileMode) error    { return os.Chmod(name, mode) }
func (Local) Chtimes(name string, atime, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}
func (Local) Symlink(target, name string) error    { return os.Symlink(target, name) }
func (Local) Readlink(name string) (string, error) { return os.Readlink(name) }
func (Local) Join(parts ...string) string          { return filepath.Join(parts...) }
func (Local) Dir(name string) string               { return filepath.Dir(name) }
func (Local) Create(name string, mode fs.FileMode) (Writable, error) {
	return os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode.Perm())
}

// CommitAbsent commits a staged local object without replacing destination.
func (backend Local) CommitAbsent(stagePath, destination string) error {
	return commitAbsent(backend, stagePath, destination)
}

type localCommitBackend interface {
	Lstat(string) (fs.FileInfo, error)
	ReadDir(string) ([]fs.DirEntry, error)
	Readlink(string) (string, error)
	Link(string, string) error
	Symlink(string, string) error
	Mkdir(string, fs.FileMode) error
	Rename(string, string) error
	Remove(string) error
	RemoveAll(string) error
	Chmod(string, fs.FileMode) error
	Chtimes(string, time.Time, time.Time) error
	Join(...string) string
}

func commitAbsent(backend localCommitBackend, stagePath, destination string) (resultErr error) {
	if _, err := backend.Lstat(destination); err == nil {
		return fmt.Errorf("%w: %q", ErrDestinationExists, destination)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	info, err := backend.Lstat(stagePath)
	if err != nil {
		return err
	}
	if info.Mode().IsRegular() {
		return commitLink(backend, stagePath, destination, func() error { return backend.Link(stagePath, destination) })
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		target, err := backend.Readlink(stagePath)
		if err != nil {
			return err
		}
		return commitLink(backend, stagePath, destination, func() error { return backend.Symlink(target, destination) })
	}
	if !info.IsDir() {
		return fmt.Errorf("unsupported staged object %q", stagePath)
	}
	if err := backend.Mkdir(destination, 0o700); err != nil {
		return classifyCommitError(destination, err)
	}
	defer func() {
		if resultErr != nil {
			resultErr = errors.Join(resultErr, backend.RemoveAll(destination))
		}
	}()
	entries, err := backend.ReadDir(stagePath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := backend.Rename(backend.Join(stagePath, entry.Name()), backend.Join(destination, entry.Name())); err != nil {
			return err
		}
	}
	if err := backend.Chmod(destination, info.Mode().Perm()); err != nil {
		return err
	}
	if err := backend.Chtimes(destination, info.ModTime(), info.ModTime()); err != nil {
		return err
	}
	return backend.Remove(stagePath)
}

func commitLink(backend localCommitBackend, stagePath, destination string, create func() error) error {
	if err := create(); err != nil {
		return classifyCommitError(destination, err)
	}
	if err := backend.Remove(stagePath); err != nil {
		return errors.Join(err, backend.RemoveAll(destination))
	}
	return nil
}

func classifyCommitError(destination string, err error) error {
	if errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("%w: %q", ErrDestinationExists, destination)
	}
	return err
}

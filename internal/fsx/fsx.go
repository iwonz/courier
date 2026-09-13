// Package fsx defines the filesystem boundary shared by transfer transports.
package fsx

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

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
}

// Local implements Backend with the Go standard library.
type Local struct{}

func (Local) Lstat(name string) (fs.FileInfo, error)       { return os.Lstat(name) }
func (Local) ReadDir(name string) ([]fs.DirEntry, error)   { return os.ReadDir(name) }
func (Local) Open(name string) (io.ReadCloser, error)      { return os.Open(name) }
func (Local) MkdirAll(name string, mode fs.FileMode) error { return os.MkdirAll(name, mode) }
func (Local) RemoveAll(name string) error                  { return os.RemoveAll(name) }
func (Local) Rename(oldPath, newPath string) error         { return os.Rename(oldPath, newPath) }
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

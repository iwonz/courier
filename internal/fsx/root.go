package fsx

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// RootedLocal is a local backend constrained by os.Root.
type RootedLocal struct{ root *os.Root }

var (
	absoluteRootPath = filepath.Abs
	rootPathDir      = filepath.Dir
	statRootPath     = os.Stat
	openRootPath     = os.OpenRoot
	relativeRootPath = filepath.Rel
)

// OpenRootedForPath opens the nearest existing parent and returns a relative
// endpoint name that cannot escape that root through traversal or symlinks.
func OpenRootedForPath(name string) (*RootedLocal, string, error) {
	absolute, err := absoluteRootPath(name)
	if err != nil {
		return nil, "", err
	}
	parent := rootPathDir(absolute)
	for {
		info, statErr := statRootPath(parent)
		if statErr == nil && info.IsDir() {
			break
		}
		if statErr != nil && !os.IsNotExist(statErr) {
			return nil, "", statErr
		}
		next := rootPathDir(parent)
		if next == parent {
			return nil, "", &os.PathError{Op: "openroot", Path: name, Err: os.ErrNotExist}
		}
		parent = next
	}
	root, err := openRootPath(parent)
	if err != nil {
		return nil, "", err
	}
	relative, err := relativeRootPath(parent, absolute)
	if err != nil {
		_ = root.Close()
		return nil, "", err
	}
	return &RootedLocal{root: root}, relative, nil
}

// Close releases the root handle.
func (r *RootedLocal) Close() error { return r.root.Close() }

func (r *RootedLocal) Lstat(name string) (fs.FileInfo, error) { return r.root.Lstat(name) }
func (r *RootedLocal) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(r.root.FS(), name)
}
func (r *RootedLocal) Open(name string) (io.ReadCloser, error) { return r.root.Open(name) }
func (r *RootedLocal) Create(name string, mode fs.FileMode) (Writable, error) {
	return r.root.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode.Perm())
}
func (r *RootedLocal) MkdirAll(name string, mode fs.FileMode) error {
	return r.root.MkdirAll(name, mode.Perm())
}
func (r *RootedLocal) RemoveAll(name string) error { return r.root.RemoveAll(name) }
func (r *RootedLocal) Rename(oldPath, newPath string) error {
	return r.root.Rename(oldPath, newPath)
}
func (r *RootedLocal) Chmod(name string, mode fs.FileMode) error {
	return r.root.Chmod(name, mode.Perm())
}
func (r *RootedLocal) Chtimes(name string, atime, mtime time.Time) error {
	return r.root.Chtimes(name, atime, mtime)
}
func (r *RootedLocal) Symlink(target, name string) error { return r.root.Symlink(target, name) }
func (r *RootedLocal) Readlink(name string) (string, error) {
	return r.root.Readlink(name)
}
func (*RootedLocal) Join(parts ...string) string { return filepath.Join(parts...) }
func (*RootedLocal) Dir(name string) string      { return filepath.Dir(name) }

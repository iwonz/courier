package sshx

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"time"

	"github.com/iwonz/courier/internal/fsx"
	"github.com/pkg/sftp"
)

type sftpOperations interface {
	Lstat(string) (fs.FileInfo, error)
	ReadDir(string) ([]fs.FileInfo, error)
	Open(string) (io.ReadCloser, error)
	OpenFile(string, int) (fsx.Writable, error)
	MkdirAll(string) error
	RemoveAll(string) error
	PosixRename(string, string) error
	Rename(string, string) error
	Chmod(string, fs.FileMode) error
	Chtimes(string, time.Time, time.Time) error
	Symlink(string, string) error
	ReadLink(string) (string, error)
}

type realSFTP struct{ client *sftp.Client }

func (r realSFTP) Lstat(name string) (fs.FileInfo, error)     { return r.client.Lstat(name) }
func (r realSFTP) ReadDir(name string) ([]fs.FileInfo, error) { return r.client.ReadDir(name) }
func (r realSFTP) Open(name string) (io.ReadCloser, error)    { return r.client.Open(name) }
func (r realSFTP) OpenFile(name string, flags int) (fsx.Writable, error) {
	file, err := r.client.OpenFile(name, flags)
	if err != nil {
		return nil, err
	}
	return &sftpWritable{File: file}, nil
}
func (r realSFTP) MkdirAll(name string) error  { return r.client.MkdirAll(name) }
func (r realSFTP) RemoveAll(name string) error { return r.client.RemoveAll(name) }
func (r realSFTP) PosixRename(oldPath, newPath string) error {
	return r.client.PosixRename(oldPath, newPath)
}
func (r realSFTP) Rename(oldPath, newPath string) error      { return r.client.Rename(oldPath, newPath) }
func (r realSFTP) Chmod(name string, mode fs.FileMode) error { return r.client.Chmod(name, mode) }
func (r realSFTP) Chtimes(name string, atime, mtime time.Time) error {
	return r.client.Chtimes(name, atime, mtime)
}
func (r realSFTP) Symlink(target, name string) error    { return r.client.Symlink(target, name) }
func (r realSFTP) ReadLink(name string) (string, error) { return r.client.ReadLink(name) }

// SFTPBackend adapts an authenticated SFTP client to fsx.Backend.
type SFTPBackend struct{ operations sftpOperations }

// NewSFTPBackend constructs a POSIX-path remote filesystem.
func NewSFTPBackend(client *sftp.Client) *SFTPBackend {
	return &SFTPBackend{operations: realSFTP{client: client}}
}

func (s *SFTPBackend) Lstat(name string) (fs.FileInfo, error) { return s.operations.Lstat(name) }
func (s *SFTPBackend) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := s.operations.ReadDir(name)
	if err != nil {
		return nil, err
	}
	result := make([]fs.DirEntry, len(entries))
	for index := range entries {
		result[index] = fileInfoEntry{entries[index]}
	}
	return result, nil
}
func (s *SFTPBackend) Open(name string) (io.ReadCloser, error) { return s.operations.Open(name) }
func (s *SFTPBackend) Create(name string, mode fs.FileMode) (fsx.Writable, error) {
	file, err := s.operations.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY)
	if err != nil {
		return nil, err
	}
	if err := s.operations.Chmod(name, mode.Perm()); err != nil {
		_ = file.Close()
		return nil, err
	}
	return file, nil
}
func (s *SFTPBackend) MkdirAll(name string, mode fs.FileMode) error {
	if err := s.operations.MkdirAll(name); err != nil {
		return err
	}
	return s.operations.Chmod(name, mode.Perm())
}
func (s *SFTPBackend) RemoveAll(name string) error { return s.operations.RemoveAll(name) }
func (s *SFTPBackend) Rename(oldPath, newPath string) error {
	if err := s.operations.PosixRename(oldPath, newPath); err == nil {
		return nil
	}
	return s.operations.Rename(oldPath, newPath)
}

// CommitAbsent uses the standard SFTP v3 rename, whose contract requires the
// destination not to exist, instead of the overwrite-capable POSIX extension.
func (s *SFTPBackend) CommitAbsent(stagePath, destination string) error {
	if _, err := s.operations.Lstat(destination); err == nil {
		return fsx.ErrDestinationExists
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err := s.operations.Rename(stagePath, destination); err != nil {
		if _, collisionErr := s.operations.Lstat(destination); collisionErr == nil {
			return errors.Join(fsx.ErrDestinationExists, err)
		}
		return err
	}
	return nil
}
func (s *SFTPBackend) Chmod(name string, mode fs.FileMode) error {
	return s.operations.Chmod(name, mode.Perm())
}
func (s *SFTPBackend) Chtimes(name string, atime, mtime time.Time) error {
	return s.operations.Chtimes(name, atime, mtime)
}
func (s *SFTPBackend) Symlink(target, name string) error    { return s.operations.Symlink(target, name) }
func (s *SFTPBackend) Readlink(name string) (string, error) { return s.operations.ReadLink(name) }
func (*SFTPBackend) Join(parts ...string) string            { return path.Join(parts...) }
func (*SFTPBackend) Dir(name string) string                 { return path.Dir(name) }

type fileInfoEntry struct{ fs.FileInfo }

func (e fileInfoEntry) Type() fs.FileMode          { return e.Mode().Type() }
func (e fileInfoEntry) Info() (fs.FileInfo, error) { return e.FileInfo, nil }

type sftpWritable struct{ *sftp.File }

func (w *sftpWritable) Sync() error {
	err := w.File.Sync()
	var status *sftp.StatusError
	if errors.As(err, &status) && status.FxCode() == sftp.ErrSSHFxOpUnsupported {
		return nil
	}
	return err
}

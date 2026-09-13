package sshx

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/fsx"
)

func TestSFTPBackend(t *testing.T) {
	info := sshInfo{name: "file", mode: 0o640}
	writer := &sshWriter{}
	operations := &fakeSFTPOperations{info: info, entries: []fs.FileInfo{info}, reader: io.NopCloser(bytes.NewBufferString("data")), writer: writer, link: "target"}
	backend := &SFTPBackend{operations: operations}
	if got, err := backend.Lstat("/file"); err != nil || got.Name() != "file" {
		t.Fatalf("lstat=%v,%v", got, err)
	}
	entries, err := backend.ReadDir("/")
	if err != nil || len(entries) != 1 || entries[0].Name() != "file" || entries[0].Type() != 0 || entries[0].IsDir() || mustInfo(t, entries[0]).Name() != "file" {
		t.Fatalf("entries=%v err=%v", entries, err)
	}
	reader, err := backend.Open("/file")
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(reader)
	_ = reader.Close()
	if string(data) != "data" {
		t.Fatalf("data=%q", data)
	}
	created, err := backend.Create("/new", 0o640)
	if err != nil || created != writer || operations.chmodMode != 0o640 {
		t.Fatalf("created=%v chmod=%o err=%v", created, operations.chmodMode, err)
	}
	stamp := time.Unix(1, 0)
	for _, err := range []error{
		backend.MkdirAll("/dir", 0o700), backend.RemoveAll("/old"), backend.Rename("/old", "/new"),
		backend.Chmod("/new", 0o600), backend.Chtimes("/new", stamp, stamp), backend.Symlink("target", "/link"),
	} {
		if err != nil {
			t.Fatal(err)
		}
	}
	if link, err := backend.Readlink("/link"); err != nil || link != "target" {
		t.Fatalf("link=%q err=%v", link, err)
	}
	if backend.Join("/a", "b") != "/a/b" || backend.Dir("/a/b") != "/a" {
		t.Fatal("POSIX path operations failed")
	}
	if operations.renameCalls != 0 {
		t.Fatal("fallback rename should not run after POSIX rename success")
	}
}

func TestSFTPBackendFailures(t *testing.T) {
	operations := &fakeSFTPOperations{err: errors.New("operation")}
	backend := &SFTPBackend{operations: operations}
	if _, err := backend.ReadDir("/"); err == nil {
		t.Fatal("expected read-dir error")
	}
	if _, err := backend.Create("/new", 0o600); err == nil {
		t.Fatal("expected open error")
	}
	if err := backend.MkdirAll("/dir", 0o700); err == nil {
		t.Fatal("expected mkdir error")
	}
	operations.err = nil
	writer := &sshWriter{}
	operations.writer = writer
	operations.chmodErr = errors.New("chmod")
	if _, err := backend.Create("/new", 0o600); err == nil || !writer.closed {
		t.Fatal("expected chmod error and close")
	}
	operations.chmodErr = nil
	operations.posixErr = errors.New("unsupported")
	if err := backend.Rename("/old", "/new"); err != nil || operations.renameCalls != 1 {
		t.Fatalf("fallback rename failed: %v", err)
	}
}

func TestSFTPCommitAbsent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		operations := &fakeSFTPOperations{lstat: func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }}
		backend := &SFTPBackend{operations: operations}
		if err := backend.CommitAbsent("/stage", "/final"); err != nil || operations.renameCalls != 1 {
			t.Fatalf("error=%v renameCalls=%d", err, operations.renameCalls)
		}
	})
	t.Run("preflight collision", func(t *testing.T) {
		operations := &fakeSFTPOperations{info: sshInfo{name: "final"}}
		if err := (&SFTPBackend{operations: operations}).CommitAbsent("/stage", "/final"); !errors.Is(err, fsx.ErrDestinationExists) || operations.renameCalls != 0 {
			t.Fatalf("error=%v renameCalls=%d", err, operations.renameCalls)
		}
	})
	t.Run("preflight failure", func(t *testing.T) {
		operations := &fakeSFTPOperations{err: errors.New("lstat")}
		if err := (&SFTPBackend{operations: operations}).CommitAbsent("/stage", "/final"); err == nil {
			t.Fatal("expected lstat error")
		}
	})
	t.Run("late collision", func(t *testing.T) {
		calls := 0
		operations := &fakeSFTPOperations{
			renameErr: errors.New("rename"),
			lstat: func(string) (fs.FileInfo, error) {
				calls++
				if calls == 1 {
					return nil, fs.ErrNotExist
				}
				return sshInfo{name: "final"}, nil
			},
		}
		if err := (&SFTPBackend{operations: operations}).CommitAbsent("/stage", "/final"); !errors.Is(err, fsx.ErrDestinationExists) {
			t.Fatalf("error=%v", err)
		}
	})
	t.Run("rename failure", func(t *testing.T) {
		operations := &fakeSFTPOperations{renameErr: errors.New("rename"), lstat: func(string) (fs.FileInfo, error) { return nil, fs.ErrNotExist }}
		if err := (&SFTPBackend{operations: operations}).CommitAbsent("/stage", "/final"); err == nil || errors.Is(err, fsx.ErrDestinationExists) {
			t.Fatalf("error=%v", err)
		}
	})
}

func TestNewSFTPBackend(t *testing.T) {
	backend := NewSFTPBackend(nil)
	if backend == nil {
		t.Fatal("expected backend")
	}
}

func mustInfo(t *testing.T, entry fs.DirEntry) fs.FileInfo {
	t.Helper()
	info, err := entry.Info()
	if err != nil {
		t.Fatal(err)
	}
	return info
}

type sshInfo struct {
	name string
	mode fs.FileMode
}

func (s sshInfo) Name() string      { return s.name }
func (sshInfo) Size() int64         { return 0 }
func (s sshInfo) Mode() fs.FileMode { return s.mode }
func (sshInfo) ModTime() time.Time  { return time.Time{} }
func (s sshInfo) IsDir() bool       { return s.mode.IsDir() }
func (sshInfo) Sys() any            { return nil }

type sshWriter struct {
	bytes.Buffer
	closed bool
}

func (s *sshWriter) Sync() error  { return nil }
func (s *sshWriter) Close() error { s.closed = true; return nil }

type fakeSFTPOperations struct {
	info        fs.FileInfo
	entries     []fs.FileInfo
	reader      io.ReadCloser
	writer      fsx.Writable
	link        string
	err         error
	chmodErr    error
	posixErr    error
	chmodMode   fs.FileMode
	renameCalls int
	renameErr   error
	lstat       func(string) (fs.FileInfo, error)
}

func (f *fakeSFTPOperations) Lstat(name string) (fs.FileInfo, error) {
	if f.lstat != nil {
		return f.lstat(name)
	}
	return f.info, f.err
}
func (f *fakeSFTPOperations) ReadDir(string) ([]fs.FileInfo, error)      { return f.entries, f.err }
func (f *fakeSFTPOperations) Open(string) (io.ReadCloser, error)         { return f.reader, f.err }
func (f *fakeSFTPOperations) OpenFile(string, int) (fsx.Writable, error) { return f.writer, f.err }
func (f *fakeSFTPOperations) MkdirAll(string) error                      { return f.err }
func (f *fakeSFTPOperations) RemoveAll(string) error                     { return f.err }
func (f *fakeSFTPOperations) PosixRename(string, string) error           { return f.posixErr }
func (f *fakeSFTPOperations) Rename(string, string) error {
	f.renameCalls++
	if f.renameErr != nil {
		return f.renameErr
	}
	return f.err
}
func (f *fakeSFTPOperations) Chmod(_ string, mode fs.FileMode) error {
	f.chmodMode = mode
	return f.chmodErr
}
func (f *fakeSFTPOperations) Chtimes(string, time.Time, time.Time) error { return f.err }
func (f *fakeSFTPOperations) Symlink(string, string) error               { return f.err }
func (f *fakeSFTPOperations) ReadLink(string) (string, error)            { return f.link, f.err }

var _ sftpOperations = (*fakeSFTPOperations)(nil)

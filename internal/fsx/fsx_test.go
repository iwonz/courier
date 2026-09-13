package fsx

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestLocalBackend(t *testing.T) {
	backend := Local{}
	root := t.TempDir()
	directory := backend.Join(root, "dir")
	if backend.Dir(directory) != root {
		t.Fatal("unexpected parent")
	}
	if err := backend.MkdirAll(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	file := backend.Join(directory, "file")
	writer, err := backend.Create(file, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("content")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Sync(); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := backend.Chmod(file, 0o640); err != nil {
		t.Fatal(err)
	}
	stamp := time.Unix(1_700_000_000, 0)
	if err := backend.Chtimes(file, stamp, stamp); err != nil {
		t.Fatal(err)
	}
	info, err := backend.Lstat(file)
	if err != nil || runtime.GOOS != "windows" && info.Mode().Perm() != 0o640 || !info.ModTime().Equal(stamp) {
		t.Fatalf("info=%v err=%v", info, err)
	}
	reader, err := backend.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	data, err := io.ReadAll(reader)
	if err != nil || string(data) != "content" {
		t.Fatalf("data=%q err=%v", data, err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	entries, err := backend.ReadDir(directory)
	if err != nil || len(entries) != 1 || entries[0].Name() != "file" {
		t.Fatalf("entries=%v err=%v", entries, err)
	}
	link := backend.Join(root, "link")
	if err := backend.Symlink(file, link); err != nil {
		t.Fatal(err)
	}
	if target, err := backend.Readlink(link); err != nil || target != file {
		t.Fatalf("target=%q err=%v", target, err)
	}
	renamed := filepath.Join(directory, "renamed")
	if err := backend.Rename(file, renamed); err != nil {
		t.Fatal(err)
	}
	if err := backend.RemoveAll(directory); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(directory); !os.IsNotExist(err) {
		t.Fatalf("directory remains: %v", err)
	}
}

func TestLocalCommitAbsent(t *testing.T) {
	backend := Local{}
	root := t.TempDir()

	fileStage := filepath.Join(root, "file-stage")
	fileFinal := filepath.Join(root, "file-final")
	if err := os.WriteFile(fileStage, []byte("file"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := backend.CommitAbsent(fileStage, fileFinal); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(fileFinal); err != nil || string(data) != "file" {
		t.Fatalf("file=%q err=%v", data, err)
	}

	linkStage := filepath.Join(root, "link-stage")
	linkFinal := filepath.Join(root, "link-final")
	if err := os.Symlink("file-final", linkStage); err != nil {
		t.Fatal(err)
	}
	if err := backend.CommitAbsent(linkStage, linkFinal); err != nil {
		t.Fatal(err)
	}
	if target, err := os.Readlink(linkFinal); err != nil || target != "file-final" {
		t.Fatalf("link=%q err=%v", target, err)
	}

	directoryStage := filepath.Join(root, "directory-stage")
	directoryFinal := filepath.Join(root, "directory-final")
	if err := os.Mkdir(directoryStage, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directoryStage, "child"), []byte("child"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := backend.CommitAbsent(directoryStage, directoryFinal); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(directoryFinal, "child")); err != nil || string(data) != "child" {
		t.Fatalf("directory child=%q err=%v", data, err)
	}

	collisionStage := filepath.Join(root, "collision-stage")
	if err := os.WriteFile(collisionStage, []byte("new"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := backend.CommitAbsent(collisionStage, fileFinal); !errors.Is(err, ErrDestinationExists) {
		t.Fatalf("collision error=%v", err)
	}
}

func TestCommitAbsentFailures(t *testing.T) {
	regular := commitInfo{mode: 0o600}
	directory := commitInfo{mode: fs.ModeDir | 0o700}
	symlink := commitInfo{mode: fs.ModeSymlink}
	tests := []struct {
		name    string
		backend faultCommitBackend
	}{
		{"destination stat", faultCommitBackend{destinationErr: errors.New("destination stat")}},
		{"stage stat", faultCommitBackend{destinationErr: fs.ErrNotExist, stageErr: errors.New("stage stat")}},
		{"late collision", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: regular, linkErr: fs.ErrExist}},
		{"link", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: regular, linkErr: errors.New("link")}},
		{"linked stage remove", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: regular, removeErr: errors.New("remove"), removeAllErr: errors.New("rollback")}},
		{"readlink", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: symlink, readlinkErr: errors.New("readlink")}},
		{"unsupported", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: commitInfo{mode: fs.ModeNamedPipe}}},
		{"mkdir", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: directory, mkdirErr: errors.New("mkdir")}},
		{"read directory", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: directory, readDirErr: errors.New("read directory")}},
		{"move child", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: directory, entries: []fs.DirEntry{commitEntry{"child"}}, renameErr: errors.New("rename")}},
		{"chmod", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: directory, chmodErr: errors.New("chmod")}},
		{"chtimes", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: directory, chtimesErr: errors.New("chtimes")}},
		{"remove directory stage", faultCommitBackend{destinationErr: fs.ErrNotExist, stageInfo: directory, removeErr: errors.New("remove")}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := commitAbsent(&test.backend, "stage", "destination"); err == nil {
				t.Fatal("expected commit failure")
			}
		})
	}
}

type commitInfo struct{ mode fs.FileMode }

func (commitInfo) Name() string        { return "stage" }
func (commitInfo) Size() int64         { return 0 }
func (i commitInfo) Mode() fs.FileMode { return i.mode }
func (commitInfo) ModTime() time.Time  { return time.Unix(1, 0) }
func (i commitInfo) IsDir() bool       { return i.mode.IsDir() }
func (commitInfo) Sys() any            { return nil }

type commitEntry struct{ name string }

func (e commitEntry) Name() string             { return e.name }
func (commitEntry) IsDir() bool                { return false }
func (commitEntry) Type() fs.FileMode          { return 0 }
func (commitEntry) Info() (fs.FileInfo, error) { return commitInfo{}, nil }

type faultCommitBackend struct {
	destinationInfo fs.FileInfo
	destinationErr  error
	stageInfo       fs.FileInfo
	stageErr        error
	linkErr         error
	readlinkErr     error
	symlinkErr      error
	mkdirErr        error
	readDirErr      error
	entries         []fs.DirEntry
	renameErr       error
	removeErr       error
	removeAllErr    error
	chmodErr        error
	chtimesErr      error
}

func (b *faultCommitBackend) Lstat(name string) (fs.FileInfo, error) {
	if name == "destination" {
		return b.destinationInfo, b.destinationErr
	}
	return b.stageInfo, b.stageErr
}
func (b *faultCommitBackend) ReadDir(string) ([]fs.DirEntry, error) { return b.entries, b.readDirErr }
func (b *faultCommitBackend) Readlink(string) (string, error)       { return "target", b.readlinkErr }
func (b *faultCommitBackend) Link(string, string) error             { return b.linkErr }
func (b *faultCommitBackend) Symlink(string, string) error          { return b.symlinkErr }
func (b *faultCommitBackend) Mkdir(string, fs.FileMode) error       { return b.mkdirErr }
func (b *faultCommitBackend) Rename(string, string) error           { return b.renameErr }
func (b *faultCommitBackend) Remove(string) error                   { return b.removeErr }
func (b *faultCommitBackend) RemoveAll(string) error                { return b.removeAllErr }
func (b *faultCommitBackend) Chmod(string, fs.FileMode) error       { return b.chmodErr }
func (b *faultCommitBackend) Chtimes(string, time.Time, time.Time) error {
	return b.chtimesErr
}
func (b *faultCommitBackend) Join(parts ...string) string { return filepath.Join(parts...) }

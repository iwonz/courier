package fsx

import (
	"io"
	"os"
	"path/filepath"
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
	if err != nil || info.Mode().Perm() != 0o640 || !info.ModTime().Equal(stamp) {
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

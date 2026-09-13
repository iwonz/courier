package fsx

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestRootedLocal(t *testing.T) {
	root := t.TempDir()
	backend, relative, err := OpenRootedForPath(filepath.Join(root, "missing", "target"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = backend.Close() })
	if relative != filepath.Join("missing", "target") || backend.Dir(relative) != "missing" || backend.Join("missing", "child") != filepath.Join("missing", "child") {
		t.Fatalf("relative=%q", relative)
	}
	if err := backend.MkdirAll(backend.Dir(relative), 0o700); err != nil {
		t.Fatal(err)
	}
	writer, err := backend.Create(relative, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = writer.Write([]byte("content"))
	_ = writer.Sync()
	_ = writer.Close()
	stamp := time.Unix(1_700_000_000, 0)
	if err := backend.Chmod(relative, 0o640); err != nil {
		t.Fatal(err)
	}
	if err := backend.Chtimes(relative, stamp, stamp); err != nil {
		t.Fatal(err)
	}
	info, err := backend.Lstat(relative)
	if err != nil || runtime.GOOS != "windows" && info.Mode().Perm() != 0o640 {
		t.Fatalf("info=%v err=%v", info, err)
	}
	entries, err := backend.ReadDir("missing")
	if err != nil || len(entries) != 1 {
		t.Fatalf("entries=%v err=%v", entries, err)
	}
	reader, err := backend.Open(relative)
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(reader)
	_ = reader.Close()
	if string(data) != "content" {
		t.Fatalf("data=%q", data)
	}
	link := filepath.Join("missing", "link")
	if err := backend.Symlink("target", link); err != nil {
		t.Fatal(err)
	}
	if target, err := backend.Readlink(link); err != nil || target != "target" {
		t.Fatalf("target=%q err=%v", target, err)
	}
	renamed := filepath.Join("missing", "renamed")
	if err := backend.Rename(relative, renamed); err != nil {
		t.Fatal(err)
	}
	committed := filepath.Join("missing", "committed")
	if err := backend.CommitAbsent(renamed, committed); err != nil {
		t.Fatal(err)
	}
	stageDirectory := filepath.Join("missing", "stage-directory")
	if err := backend.Mkdir(stageDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	child, err := backend.Create(filepath.Join(stageDirectory, "child"), 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := child.Close(); err != nil {
		t.Fatal(err)
	}
	if err := backend.CommitAbsent(stageDirectory, filepath.Join("missing", "committed-directory")); err != nil {
		t.Fatal(err)
	}
	if err := backend.RemoveAll("missing"); err != nil {
		t.Fatal(err)
	}
}

func TestRootedLocalRejectsEscape(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret"), []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
		t.Fatal(err)
	}
	backend, _, err := OpenRootedForPath(filepath.Join(root, "target"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = backend.Close() })
	if _, err := backend.Open(filepath.Join("escape", "secret")); err == nil {
		t.Fatal("os.Root allowed a symlink escape")
	}
}

func TestOpenRootedFailures(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	backend, relative, err := OpenRootedForPath(filepath.Join(file, "child"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = backend.Close() })
	if _, err := backend.Open(relative); err == nil {
		t.Fatal("expected non-directory path error")
	}
}

func TestOpenRootedInjectedFailures(t *testing.T) {
	originalAbsolute, originalDir, originalStat := absoluteRootPath, rootPathDir, statRootPath
	originalOpen, originalRelative := openRootPath, relativeRootPath
	t.Cleanup(func() {
		absoluteRootPath, rootPathDir, statRootPath = originalAbsolute, originalDir, originalStat
		openRootPath, relativeRootPath = originalOpen, originalRelative
	})

	absoluteRootPath = func(string) (string, error) { return "", os.ErrInvalid }
	if _, _, err := OpenRootedForPath("target"); err == nil {
		t.Fatal("expected absolute path failure")
	}
	absoluteRootPath = func(string) (string, error) { return "/target", nil }
	rootPathDir = func(string) string { return "/" }
	statRootPath = func(string) (os.FileInfo, error) { return nil, os.ErrPermission }
	if _, _, err := OpenRootedForPath("target"); err == nil {
		t.Fatal("expected stat failure")
	}
	statRootPath = func(string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	if _, _, err := OpenRootedForPath("target"); err == nil {
		t.Fatal("expected missing root failure")
	}

	rootPathDir = originalDir
	statRootPath = originalStat
	root := t.TempDir()
	absoluteRootPath = func(string) (string, error) { return filepath.Join(root, "target"), nil }
	openRootPath = func(string) (*os.Root, error) { return nil, os.ErrPermission }
	if _, _, err := OpenRootedForPath("target"); err == nil {
		t.Fatal("expected open root failure")
	}
	openRootPath = originalOpen
	relativeRootPath = func(string, string) (string, error) { return "", os.ErrInvalid }
	if _, _, err := OpenRootedForPath("target"); err == nil {
		t.Fatal("expected relative path failure")
	}
}

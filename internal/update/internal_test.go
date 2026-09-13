package update

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRunInternalHandoffAndCleanup(t *testing.T) {
	isolateInternalHooks(t)
	root := t.TempDir()
	target := filepath.Join(root, "courier.exe")
	staged := filepath.Join(root, stagedPrefix+"one.exe")
	if err := os.WriteFile(target, []byte("old"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(staged, []byte("new"), 0o700); err != nil {
		t.Fatal(err)
	}
	internalExecutable = func() (string, error) { return staged, nil }
	waits := 0
	waitForUpdatePID = func(pid int) error {
		if pid != 123 {
			t.Fatalf("wait PID=%d", pid)
		}
		waits++
		return nil
	}
	currentProcessID = func() int { return 456 }
	var launchedExecutable string
	var launchedArguments []string
	launchUpdateProcess = func(executable string, arguments ...string) error {
		launchedExecutable, launchedArguments = executable, arguments
		return nil
	}
	handled, err := RunInternal([]string{handoffMode, target, "123"})
	if !handled || err != nil || waits != 1 || launchedExecutable != target {
		t.Fatalf("handled=%v waits=%d executable=%q args=%v err=%v", handled, waits, launchedExecutable, launchedArguments, err)
	}
	if len(launchedArguments) != 3 || launchedArguments[0] != cleanupMode || launchedArguments[1] != staged || launchedArguments[2] != "456" {
		t.Fatalf("cleanup arguments=%v", launchedArguments)
	}
	if data, err := os.ReadFile(target); err != nil || string(data) != "new" {
		t.Fatalf("target=%q err=%v", data, err)
	}
	internalExecutable = func() (string, error) { return target, nil }
	waitForUpdatePID = func(pid int) error {
		if pid != 456 {
			t.Fatalf("cleanup PID=%d", pid)
		}
		return nil
	}
	handled, err = RunInternal([]string{cleanupMode, staged, "456"})
	if !handled || err != nil {
		t.Fatalf("cleanup handled=%v err=%v", handled, err)
	}
	if _, err := os.Stat(staged); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("staged helper remains: %v", err)
	}
}

func TestRunInternalValidation(t *testing.T) {
	if handled, err := RunInternal(nil); handled || err != nil {
		t.Fatalf("normal arguments handled=%v err=%v", handled, err)
	}
	if handled, err := RunInternal([]string{"version"}); handled || err != nil {
		t.Fatalf("public arguments handled=%v err=%v", handled, err)
	}
	for _, arguments := range [][]string{{handoffMode}, {cleanupMode, "path", "bad"}, {handoffMode, "path", "0"}} {
		if handled, err := RunInternal(arguments); !handled || err == nil {
			t.Fatalf("arguments=%v handled=%v err=%v", arguments, handled, err)
		}
	}
}

func TestLaunchUpdateHandoffAndPathValidation(t *testing.T) {
	defaultLaunch := launchUpdateProcess
	if err := defaultLaunch(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("expected process start error")
	}
	currentTest, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if err := defaultLaunch(currentTest, "-test.run=^$"); err != nil {
		t.Fatalf("release child process: %v", err)
	}
	isolateInternalHooks(t)
	root := t.TempDir()
	staged := filepath.Join(root, stagedPrefix+"x")
	target := filepath.Join(root, "courier")
	called := false
	launchUpdateProcess = func(executable string, arguments ...string) error {
		called = executable == staged && len(arguments) == 3 && arguments[0] == handoffMode && arguments[1] == target && arguments[2] == "7"
		return nil
	}
	if err := launchUpdateHandoff(staged, target, 7); err != nil || !called {
		t.Fatalf("called=%v err=%v", called, err)
	}
	if err := launchUpdateHandoff(target, target, 7); err == nil {
		t.Fatal("expected unsafe handoff paths")
	}
	launchUpdateProcess = func(string, ...string) error { return errors.New("launch") }
	if err := launchUpdateHandoff(staged, target, 7); err == nil {
		t.Fatal("expected launch error")
	}
	for _, paths := range [][2]string{{target, target}, {staged, filepath.Join(t.TempDir(), "courier")}, {filepath.Join(root, "other"), target}} {
		if err := validateUpdatePaths(paths[0], paths[1]); err == nil {
			t.Fatalf("expected unsafe paths %v", paths)
		}
	}
	calls := 0
	absoluteUpdatePath = func(name string) (string, error) {
		calls++
		if calls == 1 {
			return "", errors.New("first absolute")
		}
		return name, nil
	}
	if err := validateUpdatePaths(staged, target); err == nil {
		t.Fatal("expected first absolute error")
	}
	calls = 0
	absoluteUpdatePath = func(name string) (string, error) {
		calls++
		if calls == 2 {
			return "", errors.New("second absolute")
		}
		return name, nil
	}
	if err := validateUpdatePaths(staged, target); err == nil {
		t.Fatal("expected second absolute error")
	}
}

func TestRunHandoffFailures(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "courier")
	staged := filepath.Join(root, stagedPrefix+"source")
	if err := os.WriteFile(staged, []byte("new"), 0o700); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name  string
		setup func(*testing.T)
	}{
		{"executable", func(*testing.T) { internalExecutable = func() (string, error) { return "", errors.New("executable") } }},
		{"unsafe path", func(*testing.T) {
			internalExecutable = func() (string, error) { return filepath.Join(root, "unsafe"), nil }
		}},
		{"wait", func(*testing.T) { waitForUpdatePID = func(int) error { return errors.New("wait") } }},
		{"create", func(*testing.T) {
			createStagedFile = func(string, string) (stagedFile, error) { return nil, errors.New("create") }
		}},
		{"chmod", func(*testing.T) {
			createStagedFile = func(string, string) (stagedFile, error) {
				return &fakeStagedFile{name: filepath.Join(root, ".courier-update-commit-x"), chmodErr: errors.New("chmod")}, nil
			}
		}},
		{"source open", func(*testing.T) {
			openUpdateSource = func(string) (io.ReadCloser, error) { return nil, errors.New("open") }
		}},
		{"source copy", func(*testing.T) {
			openUpdateSource = func(string) (io.ReadCloser, error) { return &failingReadCloser{readErr: errors.New("copy")}, nil }
		}},
		{"source close", func(*testing.T) {
			openUpdateSource = func(string) (io.ReadCloser, error) {
				return &failingReadCloser{Reader: bytes.NewBufferString("new"), closeErr: errors.New("source close")}, nil
			}
		}},
		{"sync", func(*testing.T) {
			createStagedFile = func(string, string) (stagedFile, error) {
				return &fakeStagedFile{name: filepath.Join(root, ".courier-update-commit-x"), syncErr: errors.New("sync")}, nil
			}
		}},
		{"commit close", func(*testing.T) {
			createStagedFile = func(string, string) (stagedFile, error) {
				return &fakeStagedFile{name: filepath.Join(root, ".courier-update-commit-x"), closeErr: errors.New("close")}, nil
			}
		}},
		{"replace", func(*testing.T) { replaceStagedUpdate = func(string, string) error { return errors.New("replace") } }},
		{"cleanup launch", func(*testing.T) {
			launchUpdateProcess = func(string, ...string) error { return errors.New("launch") }
			deferUpdateRemoval = func(string) error { return nil }
		}},
		{"cleanup launch and defer", func(*testing.T) {
			launchUpdateProcess = func(string, ...string) error { return errors.New("launch") }
			deferUpdateRemoval = func(string) error { return errors.New("defer") }
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			isolateUpdateHooks(t)
			isolateInternalHooks(t)
			internalExecutable = func() (string, error) { return staged, nil }
			waitForUpdatePID = func(int) error { return nil }
			launchUpdateProcess = func(string, ...string) error { return nil }
			test.setup(t)
			if err := runHandoff(target, 1); err == nil {
				t.Fatal("expected handoff error")
			}
		})
	}
}

func TestRunCleanupFailures(t *testing.T) {
	isolateInternalHooks(t)
	root := t.TempDir()
	target := filepath.Join(root, "courier")
	staged := filepath.Join(root, stagedPrefix+"x")
	internalExecutable = func() (string, error) { return "", errors.New("executable") }
	if err := runCleanup(staged, 1); err == nil {
		t.Fatal("expected executable error")
	}
	internalExecutable = func() (string, error) { return filepath.Join(root, "other", "courier"), nil }
	if err := runCleanup(staged, 1); err == nil {
		t.Fatal("expected path error")
	}
	internalExecutable = func() (string, error) { return target, nil }
	waitForUpdatePID = func(int) error { return errors.New("wait") }
	if err := runCleanup(staged, 1); err == nil {
		t.Fatal("expected wait error")
	}
	waitForUpdatePID = func(int) error { return nil }
	removeStagedUpdate = func(string) error { return errors.New("remove") }
	if err := runCleanup(staged, 1); err == nil {
		t.Fatal("expected removal error")
	}
}

func isolateInternalHooks(t *testing.T) {
	t.Helper()
	originalPID, originalExecutable := currentProcessID, internalExecutable
	originalWait, originalLaunch := waitForUpdatePID, launchUpdateProcess
	originalRemove, originalAbsolute := removeStagedUpdate, absoluteUpdatePath
	originalReplace, originalDefer := replaceStagedUpdate, deferUpdateRemoval
	t.Cleanup(func() {
		currentProcessID, internalExecutable = originalPID, originalExecutable
		waitForUpdatePID, launchUpdateProcess = originalWait, originalLaunch
		removeStagedUpdate, absoluteUpdatePath = originalRemove, originalAbsolute
		replaceStagedUpdate, deferUpdateRemoval = originalReplace, originalDefer
	})
}

func TestWindowsHandoffFailureFromUpdater(t *testing.T) {
	isolateUpdateHooks(t)
	updater, _, _ := validUpdater(t, "windows", "amd64", "courier.exe")
	updater.Handoff = func(string, string, int) error { return errors.New("handoff") }
	if _, err := updater.Run(t.Context()); err == nil {
		t.Fatal("expected handoff error")
	}
	if strconv.Itoa(currentProcessID()) == "" {
		t.Fatal("PID should render")
	}
}

func TestWindowsDefaultHandoffFailureFromUpdater(t *testing.T) {
	isolateUpdateHooks(t)
	updater, _, _ := validUpdater(t, "windows", "amd64", "courier.exe")
	if _, err := updater.Run(t.Context()); err == nil {
		t.Fatal("expected non-executable staged handoff failure")
	}
}

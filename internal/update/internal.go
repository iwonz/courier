package update

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	handoffMode  = "_update-handoff"
	cleanupMode  = "_update-cleanup"
	stagedPrefix = ".courier-update-partial-"
)

var (
	currentProcessID    = os.Getpid
	internalExecutable  = os.Executable
	waitForUpdatePID    = waitForProcess
	launchUpdateProcess = func(executable string, arguments ...string) error {
		command := exec.Command(executable, arguments...)
		if err := command.Start(); err != nil {
			return err
		}
		return command.Process.Release()
	}
	removeStagedUpdate  = os.Remove
	absoluteUpdatePath  = filepath.Abs
	replaceStagedUpdate = replaceUpdateFile
	deferUpdateRemoval  = deferStagedRemoval
)

// RunInternal handles private process modes used by Windows self-update. The
// boolean is false for normal public CLI arguments.
func RunInternal(arguments []string) (bool, error) {
	if len(arguments) == 0 || arguments[0] != handoffMode && arguments[0] != cleanupMode {
		return false, nil
	}
	if len(arguments) != 3 {
		return true, errors.New("invalid internal update arguments")
	}
	pid, err := strconv.Atoi(arguments[2])
	if err != nil || pid <= 0 {
		return true, errors.New("invalid internal update PID")
	}
	if arguments[0] == handoffMode {
		return true, runHandoff(arguments[1], pid)
	}
	return true, runCleanup(arguments[1], pid)
}

func launchUpdateHandoff(staged, target string, pid int) error {
	if err := validateUpdatePaths(staged, target); err != nil {
		return err
	}
	return launchUpdateProcess(staged, handoffMode, target, strconv.Itoa(pid))
}

func runHandoff(target string, pid int) error {
	staged, err := internalExecutable()
	if err != nil {
		return err
	}
	if err := validateUpdatePaths(staged, target); err != nil {
		return err
	}
	if err := waitForUpdatePID(pid); err != nil {
		return err
	}
	commit, err := createStagedFile(filepath.Dir(target), ".courier-update-commit-*")
	if err != nil {
		return err
	}
	commitPath := commit.Name()
	committed := false
	defer func() {
		_ = commit.Close()
		if !committed {
			_ = os.Remove(commitPath)
		}
	}()
	if err := commit.Chmod(0o700); err != nil {
		return err
	}
	source, err := openUpdateSource(staged)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(commit, source)
	closeSourceErr := source.Close()
	if copyErr != nil || closeSourceErr != nil {
		return errors.Join(copyErr, closeSourceErr)
	}
	if err := commit.Sync(); err != nil {
		return err
	}
	if err := commit.Close(); err != nil {
		return err
	}
	if err := replaceStagedUpdate(commitPath, target); err != nil {
		return fmt.Errorf("replace executable after handoff: %w", err)
	}
	committed = true
	if err := launchUpdateProcess(target, cleanupMode, staged, strconv.Itoa(currentProcessID())); err != nil {
		return errors.Join(fmt.Errorf("start update cleanup: %w", err), deferUpdateRemoval(staged))
	}
	return nil
}

func runCleanup(staged string, pid int) error {
	target, err := internalExecutable()
	if err != nil {
		return err
	}
	if err := validateUpdatePaths(staged, target); err != nil {
		return err
	}
	if err := waitForUpdatePID(pid); err != nil {
		return err
	}
	return removeStagedUpdate(staged)
}

func validateUpdatePaths(staged, target string) error {
	stagedAbsolute, err := absoluteUpdatePath(staged)
	if err != nil {
		return err
	}
	targetAbsolute, err := absoluteUpdatePath(target)
	if err != nil {
		return err
	}
	if filepath.Clean(stagedAbsolute) == filepath.Clean(targetAbsolute) || filepath.Dir(stagedAbsolute) != filepath.Dir(targetAbsolute) {
		return errors.New("update staging and target paths are unsafe")
	}
	if !strings.HasPrefix(filepath.Base(stagedAbsolute), stagedPrefix) {
		return errors.New("update staging filename is unsafe")
	}
	return nil
}

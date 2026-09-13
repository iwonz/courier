//go:build windows

package update

import (
	"errors"
	"testing"

	"golang.org/x/sys/windows"
)

func TestWaitForProcessWindows(t *testing.T) {
	originalOpen, originalClose := windowsOpenProcess, windowsCloseHandle
	originalWait := windowsWaitProcess
	t.Cleanup(func() {
		windowsOpenProcess, windowsCloseHandle, windowsWaitProcess = originalOpen, originalClose, originalWait
	})
	windowsOpenProcess = func(uint32, bool, uint32) (windows.Handle, error) { return 0, windows.ERROR_INVALID_PARAMETER }
	if err := waitForProcess(1); err != nil {
		t.Fatal(err)
	}
	windowsOpenProcess = func(uint32, bool, uint32) (windows.Handle, error) { return 0, windows.ERROR_ACCESS_DENIED }
	if err := waitForProcess(1); err == nil {
		t.Fatal("expected process open error")
	}
	closed := 0
	windowsOpenProcess = func(access uint32, inherit bool, pid uint32) (windows.Handle, error) {
		if access != windows.SYNCHRONIZE || inherit || pid != 7 {
			t.Fatalf("access=%d inherit=%v pid=%d", access, inherit, pid)
		}
		return windows.Handle(1), nil
	}
	windowsCloseHandle = func(windows.Handle) error { closed++; return nil }
	windowsWaitProcess = func(windows.Handle, uint32) (uint32, error) { return 0, errors.New("wait") }
	if err := waitForProcess(7); err == nil || closed != 1 {
		t.Fatalf("closed=%d err=%v", closed, err)
	}
	windowsWaitProcess = func(windows.Handle, uint32) (uint32, error) { return uint32(windows.WAIT_TIMEOUT), nil }
	if err := waitForProcess(7); err == nil {
		t.Fatal("expected process timeout")
	}
	windowsWaitProcess = func(handle windows.Handle, milliseconds uint32) (uint32, error) {
		if handle != windows.Handle(1) || milliseconds != 120_000 {
			t.Fatalf("handle=%v milliseconds=%d", handle, milliseconds)
		}
		return uint32(windows.WAIT_OBJECT_0), nil
	}
	if err := waitForProcess(7); err != nil {
		t.Fatal(err)
	}
}

func TestWindowsFileReplacement(t *testing.T) {
	originalMove := windowsMoveFile
	t.Cleanup(func() { windowsMoveFile = originalMove })
	if err := replaceUpdateFile("bad\x00source", "target"); err == nil {
		t.Fatal("expected invalid source path")
	}
	if err := replaceUpdateFile("source", "bad\x00target"); err == nil {
		t.Fatal("expected invalid target path")
	}
	if err := deferStagedRemoval("bad\x00path"); err == nil {
		t.Fatal("expected invalid deferred path")
	}
	calls := 0
	windowsMoveFile = func(from, to *uint16, flags uint32) error {
		calls++
		if from == nil {
			t.Fatal("missing source path")
		}
		if calls == 1 {
			if to == nil || flags != windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH {
				t.Fatalf("replacement flags=%d", flags)
			}
			return errors.New("move")
		}
		if to != nil || flags != windows.MOVEFILE_DELAY_UNTIL_REBOOT {
			t.Fatalf("deferred flags=%d target=%v", flags, to)
		}
		return nil
	}
	if err := replaceUpdateFile("source", "target"); err == nil {
		t.Fatal("expected move error")
	}
	if err := deferStagedRemoval("source"); err != nil || calls != 2 {
		t.Fatalf("calls=%d err=%v", calls, err)
	}
}

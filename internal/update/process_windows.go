//go:build windows

package update

import (
	"errors"

	"golang.org/x/sys/windows"
)

var (
	windowsOpenProcess = windows.OpenProcess
	windowsCloseHandle = windows.CloseHandle
	windowsWaitProcess = windows.WaitForSingleObject
	windowsMoveFile    = windows.MoveFileEx
)

func waitForProcess(pid int) error {
	handle, err := windowsOpenProcess(windows.SYNCHRONIZE, false, uint32(pid))
	if err != nil {
		if errors.Is(err, windows.ERROR_INVALID_PARAMETER) {
			return nil
		}
		return err
	}
	defer windowsCloseHandle(handle)
	status, err := windowsWaitProcess(handle, 120_000)
	if err != nil {
		return err
	}
	if status == uint32(windows.WAIT_TIMEOUT) {
		return errors.New("timed out waiting for update process")
	}
	return nil
}

func replaceUpdateFile(source, target string) error {
	from, err := windows.UTF16PtrFromString(source)
	if err != nil {
		return err
	}
	to, err := windows.UTF16PtrFromString(target)
	if err != nil {
		return err
	}
	return windowsMoveFile(from, to, windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH)
}

func deferStagedRemoval(name string) error {
	path, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}
	return windowsMoveFile(path, nil, windows.MOVEFILE_DELAY_UNTIL_REBOOT)
}

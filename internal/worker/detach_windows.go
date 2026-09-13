//go:build windows

package worker

import (
	"os/exec"
	"syscall"
)

func configureDetached(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008, HideWindow: true}
}

// ConfigureDetached applies Courier's platform-specific background process
// attributes.
func ConfigureDetached(command *exec.Cmd) { configureDetached(command) }

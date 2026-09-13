//go:build !windows

package update

import (
	"errors"
	"os"
	"syscall"
	"time"
)

var (
	unixProcessSignal = syscall.Kill
	unixProcessNow    = time.Now
	unixProcessSleep  = time.Sleep
)

func waitForProcess(pid int) error {
	deadline := unixProcessNow().Add(2 * time.Minute)
	for unixProcessNow().Before(deadline) {
		err := unixProcessSignal(pid, 0)
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		if err != nil && !errors.Is(err, syscall.EPERM) {
			return err
		}
		unixProcessSleep(50 * time.Millisecond)
	}
	return errors.New("timed out waiting for update process")
}

func replaceUpdateFile(source, target string) error { return os.Rename(source, target) }

func deferStagedRemoval(name string) error { return os.Remove(name) }

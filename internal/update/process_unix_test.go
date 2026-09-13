//go:build !windows

package update

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestWaitForProcessUnixBranches(t *testing.T) {
	originalSignal, originalNow, originalSleep := unixProcessSignal, unixProcessNow, unixProcessSleep
	t.Cleanup(func() {
		unixProcessSignal, unixProcessNow, unixProcessSleep = originalSignal, originalNow, originalSleep
	})
	base := time.Unix(1, 0)
	unixProcessNow = func() time.Time { return base }
	unixProcessSignal = func(int, syscall.Signal) error { return syscall.ESRCH }
	if err := waitForProcess(999); err != nil {
		t.Fatal(err)
	}
	unixProcessSignal = func(int, syscall.Signal) error { return syscall.EINVAL }
	if err := waitForProcess(999); !errors.Is(err, syscall.EINVAL) {
		t.Fatalf("unexpected signal error: %v", err)
	}
	nowCalls := 0
	unixProcessNow = func() time.Time {
		nowCalls++
		if nowCalls == 1 {
			return base
		}
		if nowCalls == 2 {
			return base.Add(time.Minute)
		}
		return base.Add(3 * time.Minute)
	}
	sleeps := 0
	unixProcessSleep = func(time.Duration) { sleeps++ }
	unixProcessSignal = func(int, syscall.Signal) error { return syscall.EPERM }
	if err := waitForProcess(999); err == nil || sleeps != 1 {
		t.Fatalf("expected timeout after permission retry: sleeps=%d err=%v", sleeps, err)
	}
	name := filepath.Join(t.TempDir(), "staged")
	if err := os.WriteFile(name, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := deferStagedRemoval(name); err != nil {
		t.Fatal(err)
	}
	if err := launchUpdateProcess("/usr/bin/true"); err != nil {
		t.Fatalf("release child process: %v", err)
	}
}

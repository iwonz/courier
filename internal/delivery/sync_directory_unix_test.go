//go:build !windows

package delivery

import (
	"os"
	"testing"
)

func assertClosedRootSync(t *testing.T, root *os.Root) {
	t.Helper()
	if err := (&osStateRoot{root: root}).Sync(); err == nil {
		t.Fatal("closed root sync succeeded")
	}
}

package buildinfo

import "testing"

func TestDefaultsAreNonEmpty(t *testing.T) {
	if Version == "" || Commit == "" || Date == "" {
		t.Fatal("build identity defaults must be non-empty")
	}
}

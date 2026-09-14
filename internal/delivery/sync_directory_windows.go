//go:build windows

package delivery

import "os"

// Windows does not expose a portable directory fsync through os.File. The
// immutable state file itself is flushed before its no-replace commit.
func syncStateDirectory(_ *os.Root) error { return nil }

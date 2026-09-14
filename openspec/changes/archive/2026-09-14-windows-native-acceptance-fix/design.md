# Design: Native Windows runtime boundaries

## Context

The cross-platform interfaces are intentionally shared, but Windows differs at the OS boundary: token handle zero is invalid for token queries, `fs.FS` always requires slash-separated names even when `filepath` produces backslashes, tar link names are POSIX paths, and chmod does not expose Unix permission bits.

## Decisions

### Named-pipe ownership

Use `windows.GetCurrentProcessToken`, a supported pseudo token that needs no close, then obtain the user SID for the existing owner-only SDDL. The ACL remains `D:P(A;;GA;;;SID)` and no broader principal is added.

### Rooted filesystem path dialect

Only the adapter call from Courier's native rooted backend into `fs.ReadDir` converts with `filepath.ToSlash`. Other `os.Root` operations continue to receive native paths, and the root confinement boundary remains unchanged.

### Portable archive links

Convert a source symlink target with `filepath.ToSlash` before constructing its tar header. On Windows this makes a safe relative target portable; drive-absolute, root-absolute, parent traversal, and cross-top-level targets are still rejected by the shared archive inspector. On POSIX, backslashes remain unchanged and therefore remain invalid archive syntax.

### Permission assertions

Runtime code continues requesting private modes. Native tests follow the existing policy of asserting exact mode bits only outside Windows, where Go and the filesystem expose those bits.

## Verification

Run local exact coverage, strict OpenSpec, full release verification, Windows test cross-compilation, and a fresh native Windows CI suite including admin lifecycle, named pipes, worker subprocesses, archive creation, and rooted directory commit.

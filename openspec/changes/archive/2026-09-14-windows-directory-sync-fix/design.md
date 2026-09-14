# Design: Platform-correct state-directory synchronization

## Context

Courier's immutable registry protocol writes and flushes a private temporary file, creates the final revision with no replacement, removes the temporary link, and synchronizes the containing directory. The final step is valid on POSIX filesystems but Go's portable Windows file API cannot sync a directory handle.

## Decisions

### Platform boundary

`osStateRoot.Sync` delegates to a build-tagged implementation. The non-Windows implementation retains the existing open-directory, sync, and close behavior. The Windows implementation is an explicit no-op because the state file was already flushed before commit and there is no portable directory flush to perform. This keeps unsupported platform behavior out of the registry algorithm and avoids runtime conditionals in covered first-party code.

### Coverage and error behavior

Platform-specific tests assert that a closed rooted directory fails on systems with directory sync and succeeds through the Windows compatibility implementation. The full Windows suite therefore executes and covers its compiled path, while macOS/Linux retain exact coverage of their implementation.

### Fail-fast coordination tests

The optimistic-update concurrency test closes its release channel before calling `Fatal` on setup mutations. This does not change the success path; it ensures an unexpected persistence error cannot leave the first store mutex held until Go's global test timeout.

## Verification

Run formatting, exact local coverage, strict OpenSpec, the full release-candidate gate, cross-build all Windows artifacts, and require a fresh native Windows CI job to complete before delivery.

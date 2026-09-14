# Change: Make registry directory durability Windows-compatible

## Why

The native Windows acceptance suite showed that `os.File.Sync` on a directory handle returns `Access is denied`. Courier had already flushed every immutable state file, but treated the unsupported post-commit directory flush as a failed registry update. A concurrency test then waited for its blocked mutation during cleanup and hid the original failure behind the global test timeout.

## What Changes

- Isolate state-directory synchronization behind platform-specific implementations.
- Keep real directory `fsync` behavior on supported systems and treat it as unsupported on Windows after the immutable file itself has been flushed.
- Preserve a closed-root error assertion on supported systems and exercise the Windows compatibility path natively.
- Release test coordination gates before fatal assertions so a persistence regression fails immediately instead of timing out.

## Impact

Registry file contents, private permissions, no-replace commits, optimistic revisions, and public CLI behavior remain unchanged. Windows registry writes stop reporting a false access-denied failure after a successful file flush and commit.

## Why

All four transfer directions need identical copy, synchronization, metadata, cancellation, and cleanup behavior. A transport-neutral transactional engine provides those guarantees once and lets local and future SFTP backends share the same implementation.

## What Changes

- Add a filesystem backend contract and a native local backend.
- Recursively copy regular files, directories, and symlinks without deleting the source.
- Build destination content in private partial paths and atomically commit where rename semantics permit.
- Replace existing destination trees so objects absent from source are removed.
- Make interrupted and failed operations clean up staging and remain safely retryable.
- Create and verify `<source-name>.tar.gz` archives using the Go standard library.
- Emit structured stage, byte, speed, total, and elapsed progress events.

## Capabilities

### New Capabilities

- `transactional-transfer`: transport-neutral recursive copy, synchronization, metadata, staging, rollback, and cleanup.
- `archive-mode`: safe built-in tar.gz creation and integrity verification.
- `progress-events`: structured operation progress independent of terminal rendering.

### Modified Capabilities

None.

## Impact

Adds `fsx`, `transfer`, `archive`, and `progress` packages. Transfers can now run against the local backend; SSH integration follows in a separate change.

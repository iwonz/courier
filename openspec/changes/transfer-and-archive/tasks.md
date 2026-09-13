## 1. Filesystem and transfer

- [x] 1.1 Define the backend interface and local implementation
- [x] 1.2 Implement scan, recursive copy, metadata preservation, and symlinks
- [x] 1.3 Implement private staging, replacement, rollback, and cleanup
- [x] 1.4 Implement cancellation and stable transfer errors

## 2. Archive and progress

- [x] 2.1 Implement built-in tar.gz creation across backends
- [x] 2.2 Verify gzip/tar integrity and safe entry names
- [x] 2.3 Implement structured progress tracking

## 3. Verification

- [x] 3.1 Cover local synchronization, archive, interruption, rollback, and cleanup
- [x] 3.2 Run race, coverage, and strict OpenSpec validation
- [x] 3.3 Create the conventional commit on `feat/003-transfer-and-archive`

# transactional-transfer Specification

## Purpose
Define recursive staged transfer, metadata preservation, retry, and cleanup behavior.

## Requirements

### Requirement: Recursive non-destructive copy

Courier SHALL copy regular files, directories, and symlinks recursively while leaving source objects unchanged.

#### Scenario: Nested directory tree

- **WHEN** the source contains nested directories, files, and a symlink
- **THEN** destination contains the same supported object types, content, and structure, and source still exists

### Requirement: Destination synchronization

Courier SHALL make an existing target directory match source, including removal of target objects absent from source.

#### Scenario: Stale target object

- **WHEN** destination contains an object not present in source
- **THEN** the committed destination no longer contains that object

### Requirement: Metadata preservation

Courier SHALL preserve permission bits and modification timestamps where supported by the destination backend.

#### Scenario: Regular file metadata

- **WHEN** a source file has non-default permissions and timestamp
- **THEN** destination receives those permissions and timestamp

### Requirement: Transactional commit and cleanup

Courier SHALL write through a private partial path, atomically rename where supported, restore the old target after a failed commit, and attempt cleanup on success, failure, or interruption.

#### Scenario: Cancellation during copy

- **WHEN** context cancellation interrupts a transfer
- **THEN** the old destination remains intact and the partial path is removed

#### Scenario: Safe retry

- **WHEN** a transfer is rerun after an interruption
- **THEN** it starts from a consistent source and destination state and can complete without manual cleanup

### Requirement: Unsupported object protection

Courier SHALL reject unsupported special filesystem objects rather than silently altering their meaning.

#### Scenario: Named pipe source

- **WHEN** a source entry is a named pipe
- **THEN** transfer fails before commit with a transfer-stage error

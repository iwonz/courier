## ADDED Requirements

### Requirement: Destination collision rejection

Courier SHALL reject an existing resolved final path before staging and SHALL not overwrite or merge it if a late collision appears before commit.

#### Scenario: Child collision in directory container

- **WHEN** a destination directory contains an entry matching the computed source name
- **THEN** preflight fails and every existing destination entry remains unchanged

#### Scenario: Late collision

- **WHEN** another actor creates the final path after preflight
- **THEN** commit fails without replacing the actor's object

## MODIFIED Requirements

### Requirement: Recursive non-destructive copy

Courier SHALL copy regular files, directories, and symlinks recursively into an absent final path while leaving source objects and unrelated destination entries unchanged.

#### Scenario: Nested directory tree

- **WHEN** the source contains nested directories, files, and a symlink and the final path is absent
- **THEN** destination contains the same supported object types, content, and structure, source still exists, and sibling destination entries remain unchanged

### Requirement: Transactional commit and cleanup

Courier SHALL write through a private partial path, commit only to an absent final path using no-replace behavior where supported, and attempt owned-partial cleanup on success, failure, or interruption.

#### Scenario: Cancellation during copy

- **WHEN** context cancellation interrupts a transfer
- **THEN** no final destination is committed and the partial path is removed

#### Scenario: Safe retry

- **WHEN** a transfer is rerun after an interruption and the final path is still absent
- **THEN** it starts from a consistent source and can complete without manual cleanup

## REMOVED Requirements

### Requirement: Destination synchronization

**Reason**: Mirror deletion and implicit merging are outside Courier's accepted safe-copy contract.

**Migration**: Use an absent exact destination or an existing container with an absent computed child; Courier preserves every unrelated destination entry.

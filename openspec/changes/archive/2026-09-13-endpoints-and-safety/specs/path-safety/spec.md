## ADDED Requirements

### Requirement: Self-copy prevention

Courier SHALL reject a transfer when source and actual destination resolve to the same location.

#### Scenario: Local symlink alias

- **WHEN** source and destination resolve through symlinks to the same local path
- **THEN** preflight fails before destination is modified

#### Scenario: Equivalent remote paths

- **WHEN** normalized paths on the same remote identity are equal
- **THEN** preflight fails before opening a remote writer

### Requirement: Descendant-copy prevention

Courier SHALL reject copying a directory into its own descendant on the same filesystem identity.

#### Scenario: Nested destination

- **WHEN** source directory is `/data` and actual destination is `/data/backup`
- **THEN** preflight returns an unsafe-path error

### Requirement: Traversal prevention

Courier SHALL reject absolute, parent-traversing, empty, NUL-containing, and platform-escape archive entry names before joining them to an extraction root.

#### Scenario: Parent archive entry

- **WHEN** archive entry is `../../secret`
- **THEN** safe join fails and no path outside the extraction root is returned

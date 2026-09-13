# path-safety Specification

## Purpose
Define filesystem and archive containment rules that prevent unsafe path access.

## Requirements

### Requirement: Self-copy prevention

Courier SHALL treat a plain path-to-path operation whose source and actual destination resolve to the same location as a successful no-op, and SHALL reject the same identity when archive or extraction transformation is active.

#### Scenario: Local symlink alias

- **WHEN** plain source and destination resolve through symlinks to the same local path
- **THEN** preflight returns a successful no-op and performs no copy

#### Scenario: Equivalent remote paths

- **WHEN** normalized paths on the same resolved remote identity are equal for a plain transfer
- **THEN** preflight returns a successful no-op before opening a remote writer

#### Scenario: Transformed identity

- **WHEN** source and destination resolve to the same location with archive or extraction active
- **THEN** preflight fails with an unsafe collision

### Requirement: Descendant-copy prevention

Courier SHALL reject copying a directory into its own descendant on the same filesystem identity.

#### Scenario: Nested destination

- **WHEN** source directory is `/data` and actual destination is `/data/backup`
- **THEN** preflight returns an unsafe-path error

### Requirement: Traversal prevention

Courier SHALL combine canonical preflight checks with `os.Root`-bounded local operations and strict archive entry and symlink-target validation, rejecting absolute, parent-traversing, symlink-escaping, empty-component, dot-component, NUL-containing, and platform-escape paths.

#### Scenario: Local symlink escapes operation root

- **WHEN** a local source or destination path reaches outside its opened `os.Root` through a symlink
- **THEN** the filesystem operation fails without accessing the outside object

#### Scenario: Parent archive entry

- **WHEN** an archive entry is `folder/../secret` or `../../secret`
- **THEN** validation fails and no path outside or ambiguously within the extraction root is returned

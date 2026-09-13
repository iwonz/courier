# path-safety Specification

## Purpose
Define filesystem and archive containment rules that prevent unsafe path access.

## Requirements

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

Courier SHALL combine canonical preflight checks with `os.Root`-bounded local filesystem operations and safe archive entry handling, rejecting absolute, parent-traversing, symlink-escaping, empty, NUL-containing, and platform-escape paths.

#### Scenario: Local symlink escapes operation root

- **WHEN** a local source or destination path reaches outside its opened `os.Root` through a symlink
- **THEN** the filesystem operation fails without accessing the outside object

#### Scenario: Parent archive entry

- **WHEN** archive entry is `../../secret`
- **THEN** safe join fails and no path outside the extraction root is returned

## MODIFIED Requirements

### Requirement: Traversal prevention

Courier SHALL combine canonical preflight checks with `os.Root`-bounded local filesystem operations and safe archive entry handling, rejecting absolute, parent-traversing, symlink-escaping, empty, NUL-containing, and platform-escape paths.

#### Scenario: Local symlink escapes operation root

- **WHEN** a local source or destination path reaches outside its opened `os.Root` through a symlink
- **THEN** the filesystem operation fails without accessing the outside object

#### Scenario: Parent archive entry

- **WHEN** archive entry is `../../secret`
- **THEN** safe join fails and no path outside the extraction root is returned

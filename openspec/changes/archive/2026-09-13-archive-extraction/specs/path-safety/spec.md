## MODIFIED Requirements

### Requirement: Traversal prevention

Courier SHALL combine canonical preflight checks with `os.Root`-bounded local operations and strict archive entry and symlink-target validation, rejecting absolute, parent-traversing, symlink-escaping, empty-component, dot-component, NUL-containing, and platform-escape paths.

#### Scenario: Local symlink escapes operation root

- **WHEN** a local source or destination path reaches outside its opened `os.Root` through a symlink
- **THEN** the filesystem operation fails without accessing the outside object

#### Scenario: Parent archive entry

- **WHEN** an archive entry is `folder/../secret` or `../../secret`
- **THEN** validation fails and no path outside or ambiguously within the extraction root is returned

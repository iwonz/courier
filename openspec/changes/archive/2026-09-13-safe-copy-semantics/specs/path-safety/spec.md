## MODIFIED Requirements

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

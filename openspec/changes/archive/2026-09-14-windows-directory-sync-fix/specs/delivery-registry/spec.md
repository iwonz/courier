## MODIFIED Requirements

### Requirement: Private atomic registry revisions

Courier SHALL write registry state through private temporary files, flush each completed state file, and atomically commit immutable revision files under a rooted private application directory. Courier SHALL synchronize committed directory metadata where the operating system exposes supported directory-sync semantics and SHALL not fail a flushed Windows state commit solely because portable directory `Sync` is unsupported.

#### Scenario: Registry write is interrupted

- **WHEN** a new temporary revision is incomplete or fails to commit
- **THEN** readers continue to load the previous complete revision and the owned temporary file is removed

#### Scenario: Windows commits a registry revision

- **WHEN** a complete private state file is flushed and committed without replacement on Windows
- **THEN** the update succeeds without attempting unsupported directory-handle synchronization

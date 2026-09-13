## ADDED Requirements

### Requirement: Endpoint parsing

Courier SHALL parse `[user@]host:/path` as a remote endpoint and all other valid paths, including Windows drive paths, as local endpoints.

#### Scenario: SSH alias without user

- **WHEN** `source-server:/opt/data` is parsed
- **THEN** host is `source-server`, path is `/opt/data`, and user remains unset for SSH config resolution

#### Scenario: Windows path is local

- **WHEN** `C:\Users\me\data` is parsed
- **THEN** it is a local endpoint and not host `C`

#### Scenario: Unicode and spaces

- **WHEN** a valid endpoint contains Unicode or spaces in its path
- **THEN** the path is retained byte-for-byte

### Requirement: Destination semantics

Courier SHALL place the result beneath destination using the source name when destination exists as a directory or has a trailing separator; otherwise it SHALL treat destination as the exact final path.

#### Scenario: Directory destination

- **WHEN** source is `/data/photos` and destination is `/backup/`
- **THEN** actual destination is `/backup/photos`

#### Scenario: Exact destination

- **WHEN** destination `/backup/current` is not a directory and has no trailing separator
- **THEN** actual destination is `/backup/current`

#### Scenario: Archive output name

- **WHEN** archive mode supplies output name `photos.tar.gz` for directory destination `/backup/`
- **THEN** actual destination is `/backup/photos.tar.gz`

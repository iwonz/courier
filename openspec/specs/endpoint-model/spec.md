# endpoint-model Specification

## Purpose
Define unambiguous local, SSH, service, and HTTP endpoint parsing plus path destination resolution.

## Requirements

### Requirement: Endpoint parsing

Courier SHALL classify local paths, SSH paths, exact `web://` and `webhook://` service endpoints, and HTTP(S) URLs into explicit endpoint kinds while preserving their source text and rejecting unknown URI schemes.

#### Scenario: SSH alias without user

- **WHEN** `source-server:/opt/data` is parsed
- **THEN** its kind is SSH, host is `source-server`, path is `/opt/data`, and user remains unset for SSH config resolution

#### Scenario: Windows path is local

- **WHEN** `C:\Users\me\data` or `C:/Users/me/data` is parsed
- **THEN** its kind is local and it is not interpreted as host `C` or as an URI scheme

#### Scenario: Remote Windows path

- **WHEN** `user@server:C:/Users/me/data` is parsed
- **THEN** its kind is SSH and its remote path is `C:/Users/me/data`

#### Scenario: Bracketed IPv6 SSH host

- **WHEN** `root@[2001:db8::1]:/opt/data` is parsed
- **THEN** its kind is SSH and the bracketed IPv6 host is retained

#### Scenario: Exact service endpoint

- **WHEN** `web://` or `webhook://` is parsed
- **THEN** it is classified as its Courier service kind without an authority or path

#### Scenario: Unknown scheme

- **WHEN** an endpoint begins with an unrecognized `name://` scheme
- **THEN** parsing fails instead of treating it as a local or SSH path

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

### Requirement: Extraction-root destination

Courier SHALL treat the destination of `--extract` as the extraction root itself rather than applying ordinary source-name container resolution.

#### Scenario: Existing root directory

- **WHEN** `archive.tar.gz` is extracted to existing directory `/restore/`
- **THEN** archive top-level entries are placed beneath `/restore/` and no additional `archive.tar.gz` directory is introduced

### Requirement: Missing path role classification

Courier SHALL treat a missing flexible source path as an empty directory; a missing destination as a directory when its source is a directory or its spelling ends in `/` or `\`; and incoming browser, incoming webhook, and extraction roots as directories. Courier SHALL preserve an absent non-directory-hinted destination as the exact final file when its source is a file.

#### Scenario: Missing flexible source

- **WHEN** a local or SSH source path is missing in a route that accepts directories
- **THEN** Courier classifies it as an empty directory eligible for authorized creation

#### Scenario: Existing file has an exact missing destination

- **WHEN** `report.pdf` is transferred to absent `renamed.pdf` without a trailing separator
- **THEN** Courier preserves `renamed.pdf` as the exact file destination and does not request directory creation

### Requirement: File-only path absence

Courier SHALL fail clearly when a missing path occupies a file-only role and SHALL never create that path as a directory. File-only roles SHALL include an extraction archive and an unarchived outgoing webhook payload.

#### Scenario: Extraction archive is missing

- **WHEN** `--extract` names a missing source archive
- **THEN** Courier returns a path-specific preflight error without creating the source path

#### Scenario: Outgoing webhook file is missing

- **WHEN** an outgoing HTTP route names a missing source without `--archive`
- **THEN** Courier returns a path-specific preflight error before opening an HTTP request

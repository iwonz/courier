## ADDED Requirements

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


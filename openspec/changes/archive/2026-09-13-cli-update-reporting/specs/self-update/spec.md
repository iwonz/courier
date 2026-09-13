## ADDED Requirements

### Requirement: Release discovery

Courier SHALL query `https://api.github.com/repos/iwonz/courier/releases/latest`, compare semantic versions, and report when the installed version is current.

#### Scenario: New version exists

- **WHEN** GitHub reports a higher stable version
- **THEN** Courier selects the asset matching runtime OS and architecture

### Requirement: Verified replacement

Courier SHALL download the release checksum file and matching archive into a private temporary directory, verify SHA-256 before extraction, and replace the running executable through same-directory staging where platform permissions allow.

#### Scenario: Checksum mismatch

- **WHEN** downloaded asset hash differs from the published release checksum manifest
- **THEN** update fails and leaves the installed executable unchanged

### Requirement: Update cleanup

Courier SHALL remove downloaded archives, extracted files, and partial replacements after success or failure.

#### Scenario: Replacement fails

- **WHEN** the executable cannot be replaced
- **THEN** the original remains available and private temporary files are removed

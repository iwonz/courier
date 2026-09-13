# self-update Specification

## Purpose
Define verified in-place updates from immutable GitHub Release artifacts.

## Requirements

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

### Requirement: Reliable Windows self-update handoff

Courier SHALL replace a running Windows executable through a verified staged handoff and SHALL remove handoff artifacts after replacement.

#### Scenario: Windows update is available

- **WHEN** a verified newer Windows binary has been staged
- **THEN** Courier starts the staged binary in hidden handoff mode, which waits for the original process to exit before replacing the target

#### Scenario: Replacement completes

- **WHEN** the handoff process replaces the target
- **THEN** it starts the installed binary in hidden cleanup mode, and cleanup removes the handoff executable after its process exits

#### Scenario: Unsafe internal arguments

- **WHEN** an internal handoff or cleanup invocation references an invalid PID, a non-Courier staging name, or a staging file outside the target directory
- **THEN** Courier rejects the operation without replacing or deleting the target

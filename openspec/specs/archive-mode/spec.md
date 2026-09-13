# archive-mode Specification

## Purpose
Define dependency-free archive creation, verification, and safe entry handling.

## Requirements

### Requirement: Built-in archive creation

Courier SHALL create a gzip-compressed POSIX tar archive named `<source-name>.tar.gz` from selected source objects using its built-in implementation, without requiring an external archiver.

#### Scenario: Directory archive

- **WHEN** archive mode is requested for source directory `photos` with selection rules
- **THEN** a private temporary `photos.tar.gz` contains the source root and only its selected supported descendants

### Requirement: Archive integrity validation

Courier SHALL fully read and validate gzip and tar structure, entry paths, and payloads before reporting archive creation complete.

#### Scenario: Corrupted archive

- **WHEN** archive bytes are truncated or fail checksums
- **THEN** verification fails and the temporary archive is cleaned up

### Requirement: Safe archive entries

Courier SHALL reject path traversal and unsupported source objects while producing or validating an archive.

#### Scenario: Unsafe entry name

- **WHEN** a tar header contains an absolute or parent-traversing name
- **THEN** verification rejects the archive

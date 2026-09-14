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

Courier SHALL reject path traversal, duplicate normalized paths, unsafe symlink targets, unsupported types, and structurally conflicting source objects while producing, validating, or extracting an archive. A safe relative symlink read from a local Windows filesystem SHALL be converted from native separators to POSIX tar separators before the same validation is applied.

#### Scenario: Unsafe entry name

- **WHEN** a tar header contains an absolute, empty-component, dot-component, or parent-traversing name
- **THEN** verification rejects the archive

#### Scenario: Unsafe symlink target

- **WHEN** a symlink target resolves outside its top-level archive entry
- **THEN** inspection fails before destination staging

#### Scenario: Windows relative symlink target

- **WHEN** a local Windows source symlink safely targets a sibling below the same top-level source entry
- **THEN** Courier stores the equivalent slash-separated relative link in the POSIX tar header and verification succeeds

### Requirement: Extensible archive codec registry

Courier SHALL resolve archive operations through a duplicate-safe codec registry and SHALL ship tar.gz as the only supported extraction codec.

#### Scenario: Unsupported extension

- **WHEN** extraction is requested for a file not matched by a registered codec
- **THEN** preflight fails without opening a destination writer

### Requirement: Verified archive extraction

Courier SHALL inspect a tar.gz archive completely before staging selected entries into the destination extraction root.

#### Scenario: Existing extraction root

- **WHEN** the extraction root is an existing directory with unrelated entries and no selected entry conflicts
- **THEN** selected top-level entries are committed and unrelated entries remain unchanged

#### Scenario: Archive collision

- **WHEN** any selected archive entry conflicts with an existing destination path
- **THEN** preflight fails before staging or committing any entry

### Requirement: Archive bomb limits

Courier SHALL enforce a configurable expanded-size limit with a 100 GiB default plus fixed limits of 100,000 entries, depth 64, and expansion ratio 100:1.

#### Scenario: Explicit unlimited size

- **WHEN** max extracted size is unlimited
- **THEN** entry-count, depth, and expansion-ratio limits remain enforced

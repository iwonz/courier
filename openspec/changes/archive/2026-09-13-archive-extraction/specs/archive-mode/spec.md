## ADDED Requirements

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

## MODIFIED Requirements

### Requirement: Safe archive entries

Courier SHALL reject path traversal, duplicate normalized paths, unsafe symlink targets, unsupported types, and structurally conflicting source objects while producing, validating, or extracting an archive.

#### Scenario: Unsafe entry name

- **WHEN** a tar header contains an absolute, empty-component, dot-component, or parent-traversing name
- **THEN** verification rejects the archive

#### Scenario: Unsafe symlink target

- **WHEN** a symlink target resolves outside its top-level archive entry
- **THEN** inspection fails before destination staging

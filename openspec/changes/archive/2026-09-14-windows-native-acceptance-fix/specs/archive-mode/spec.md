## MODIFIED Requirements

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

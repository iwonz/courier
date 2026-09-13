## MODIFIED Requirements

### Requirement: Built-in archive creation

Courier SHALL create a gzip-compressed POSIX tar archive named `<source-name>.tar.gz` from selected source objects using its built-in implementation, without requiring an external archiver.

#### Scenario: Directory archive

- **WHEN** archive mode is requested for source directory `photos` with selection rules
- **THEN** a private temporary `photos.tar.gz` contains the source root and only its selected supported descendants

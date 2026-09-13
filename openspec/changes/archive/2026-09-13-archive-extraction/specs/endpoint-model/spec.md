## ADDED Requirements

### Requirement: Extraction-root destination

Courier SHALL treat the destination of `--extract` as the extraction root itself rather than applying ordinary source-name container resolution.

#### Scenario: Existing root directory

- **WHEN** `archive.tar.gz` is extracted to existing directory `/restore/`
- **THEN** archive top-level entries are placed beneath `/restore/` and no additional `archive.tar.gz` directory is introduced

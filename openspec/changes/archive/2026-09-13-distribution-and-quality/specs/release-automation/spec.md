## ADDED Requirements

### Requirement: Mandatory local dry run

Courier SHALL provide one local verification command that runs formatting checks, vet, race tests, exact first-party statement coverage, installer/npm checks, and `goreleaser release --snapshot --clean` before a tag can be pushed by the release command.

#### Scenario: Broken cross-build

- **WHEN** any release target fails to compile in the snapshot
- **THEN** the local release command stops before creating a tag

### Requirement: Ordered CI publication

Courier SHALL rerun quality gates in GitHub Actions, publish GitHub and catalog artifacts only after they pass, publish npm only after GitHub Release assets exist, and report verification failures.

#### Scenario: Test failure

- **WHEN** a test or coverage gate fails for a release tag
- **THEN** no release or npm publication job starts

### Requirement: One-command release

Courier SHALL provide one command that accepts or interactively requests a semantic version, validates repository and credentials/configuration, runs the local dry run, creates one annotated tag, and pushes it to start publication with generated release notes.

#### Scenario: Missing catalog setup

- **WHEN** required GitHub secrets or external repositories are not configured
- **THEN** documentation and workflow preflight identify the exact missing connection without embedding credentials in source

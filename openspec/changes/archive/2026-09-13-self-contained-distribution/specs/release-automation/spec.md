## MODIFIED Requirements

### Requirement: Ordered CI publication

Courier SHALL rerun quality gates in GitHub Actions, publish the GitHub Release and in-repository package manifests only after they pass, publish npm only after GitHub Release assets exist, and report verification failures.

#### Scenario: Test failure

- **WHEN** a test or coverage gate fails for a release tag
- **THEN** no release, manifest, or npm publication job starts

### Requirement: One-command release

Courier SHALL provide one command that accepts or interactively requests a semantic version, validates the single Courier repository and npm credential configuration, runs the local dry run, creates one annotated tag, and pushes it to start publication with generated release notes.

#### Scenario: Missing npm setup

- **WHEN** the npm publication secret is not configured
- **THEN** the local and workflow preflight identify the missing connection without embedding credentials in source

#### Scenario: Missing catalog setup

- **WHEN** the single Courier repository or in-repository manifest configuration is invalid
- **THEN** documentation and workflow preflight identify the exact problem without requiring an external catalog repository or embedding credentials in source

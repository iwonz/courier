## MODIFIED Requirements

### Requirement: Mandatory local dry run

Courier SHALL provide one local verification command that runs formatting checks, vet, race tests, exact first-party Go and TypeScript coverage, real-browser checks, workflow and strict OpenSpec validation, installer/npm checks, `goreleaser release --snapshot --clean`, artifact verification, and isolated Linux package installation before a tag can be pushed by the release command.

#### Scenario: Broken cross-build

- **WHEN** any release target fails to compile in the snapshot
- **THEN** the local release command stops before creating a tag

## ADDED Requirements

### Requirement: Acceptance-gated publication

Courier SHALL publish a tagged release only after the repeated Linux release dry run, real-browser suite, distribution package suite, macOS runtime suite, Windows exact-coverage suite, installer tests, contract freshness, workflow validation, and strict OpenSpec validation succeed.

#### Scenario: Any acceptance job fails

- **WHEN** a release candidate fails one required platform, browser, package, installer, contract, workflow, or specification check
- **THEN** GitHub Release and npm publication jobs do not start

## REMOVED Requirements

### Requirement: Exact Winget pull-request verification

**Reason**: Courier intentionally has no Winget catalog publisher, external fork, or upstream pull-request dependency; Windows remains available through Scoop, npm, PowerShell, and direct release binaries.

**Migration**: Release acceptance verifies only GitHub Releases and the package channels maintained from the single `iwonz/courier` repository.

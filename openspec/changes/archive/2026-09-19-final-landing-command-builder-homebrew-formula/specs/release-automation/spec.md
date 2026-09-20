## MODIFIED Requirements

### Requirement: Mandatory local dry run

Courier SHALL provide one local verification command that runs all existing quality gates plus a deterministic GoReleaser source snapshot, Formula rendering, Formula source-build acceptance where Homebrew is available, artifact verification, and isolated package installation before a tag can be pushed.

#### Scenario: Formula rendering drifts

- **WHEN** the snapshot source checksum and rendered Formula disagree or the Formula cannot build and run Courier
- **THEN** local release verification fails before a tag is created

#### Scenario: Broken cross-build

- **WHEN** any release target fails to compile in the snapshot
- **THEN** the local release command stops before creating a tag

#### Scenario: ShellCheck is absent locally

- **WHEN** workflow validation runs on a supported clean developer host without ShellCheck in `PATH` or the repository cache
- **THEN** Courier downloads the pinned official archive, verifies its trusted digest, stores only the executable in the ignored tool cache, and validates every embedded workflow shell script

### Requirement: Ordered CI publication

Courier SHALL publish the GitHub Release and Scoop manifest after quality gates, then run a separate macOS Formula job that reads the released source checksum, renders and audits the Formula, installs it from source, runs Courier, and commits the exact verified Formula to `main`. npm publication SHALL begin only after GitHub Release assets exist.

#### Scenario: Formula verification fails

- **WHEN** the released source archive cannot be audited, built, installed, or executed by Homebrew
- **THEN** no Formula update is committed to `main` and the release workflow reports failure

#### Scenario: Release contents are verified

- **WHEN** release verification inspects a published version
- **THEN** it requires the deterministic source archive, its checksum, the committed Formula, and Scoop manifest and does not expect a cask

#### Scenario: Test failure

- **WHEN** a test or coverage gate fails for a release tag
- **THEN** no release, manifest, Formula, or npm publication job starts

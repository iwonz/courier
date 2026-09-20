## MODIFIED Requirements

### Requirement: Mandatory local dry run

Courier SHALL provide one local verification command that runs all existing quality gates plus a deterministic GoReleaser source snapshot, Formula rendering, Formula source-build acceptance where Homebrew is available, artifact verification, and isolated package installation before a tag can be pushed. Go source used by the release candidate SHALL remain canonically formatted under the pinned supported CI Go toolchain and supported newer local toolchains.

#### Scenario: Formula rendering drifts

- **WHEN** the snapshot source checksum and rendered Formula disagree or the Formula cannot build and run Courier
- **THEN** local release verification fails before a tag is created

#### Scenario: Broken cross-build

- **WHEN** any release target fails to compile in the snapshot
- **THEN** the local release command stops before creating a tag

#### Scenario: ShellCheck is absent locally

- **WHEN** workflow validation runs on a supported clean developer host without ShellCheck in `PATH` or the repository cache
- **THEN** Courier downloads the pinned official archive, verifies its trusted digest, stores only the executable in the ignored tool cache, and validates every embedded workflow shell script

#### Scenario: Go formatter versions disagree

- **WHEN** the pinned CI Go formatter and a supported newer local Go formatter inspect release source
- **THEN** both report a clean canonical representation before commit or tag publication

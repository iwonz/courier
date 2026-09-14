## MODIFIED Requirements

### Requirement: Mandatory local dry run

Courier SHALL provide one local verification command that runs formatting checks, vet, race tests, exact first-party Go and TypeScript coverage, real-browser checks, workflow and strict OpenSpec validation, installer/npm checks, `goreleaser release --snapshot --clean`, artifact verification, and isolated Linux package installation before a tag can be pushed by the release command. Workflow validation SHALL run actionlint with a pinned checksum-verified ShellCheck executable on supported macOS and Linux developer hosts so embedded shell checks cannot silently depend on the ambient environment.

#### Scenario: Broken cross-build

- **WHEN** any release target fails to compile in the snapshot
- **THEN** the local release command stops before creating a tag

#### Scenario: ShellCheck is absent locally

- **WHEN** workflow validation runs on a supported clean developer host without ShellCheck in `PATH` or the repository cache
- **THEN** Courier downloads the pinned official archive, verifies its trusted digest, stores only the executable in the ignored tool cache, and validates every embedded workflow shell script

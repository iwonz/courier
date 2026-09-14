# release-automation Specification

## Purpose
Define local verification and ordered one-command release publication.

## Requirements

### Requirement: Mandatory local dry run

Courier SHALL provide one local verification command that runs formatting checks, vet, race tests, exact first-party Go and TypeScript coverage, real-browser checks, workflow and strict OpenSpec validation, installer/npm checks, `goreleaser release --snapshot --clean`, artifact verification, and isolated Linux package installation before a tag can be pushed by the release command. Workflow validation SHALL run actionlint with a pinned checksum-verified ShellCheck executable on supported macOS and Linux developer hosts so embedded shell checks cannot silently depend on the ambient environment.

#### Scenario: Broken cross-build

- **WHEN** any release target fails to compile in the snapshot
- **THEN** the local release command stops before creating a tag

#### Scenario: ShellCheck is absent locally

- **WHEN** workflow validation runs on a supported clean developer host without ShellCheck in `PATH` or the repository cache
- **THEN** Courier downloads the pinned official archive, verifies its trusted digest, stores only the executable in the ignored tool cache, and validates every embedded workflow shell script

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

### Requirement: Native Windows release gate

Courier SHALL run the complete first-party Go test suite on a native Windows runner, require exact 100% statement coverage, run the PowerShell installer acceptance test, and prevent publication if any native command fails.

#### Scenario: Windows test failure

- **WHEN** a Go test exits unsuccessfully on the Windows runner
- **THEN** the gate reports that test failure, skips coverage analysis and installer acceptance, and blocks publication

#### Scenario: Platform-specific representation

- **WHEN** a test observes a local path, symlink target, or filesystem metadata on Windows
- **THEN** it validates the Windows-supported representation without weakening content, integrity, safety, or cleanup assertions

#### Scenario: Corrective branch validation

- **WHEN** a branch matching `fix/**` is pushed
- **THEN** the normal quality and native Windows gates run before that change advances to `main`

### Requirement: Publication-capable npm credential

Courier SHALL require npm CI publication credentials to have package read/write authority and non-interactive 2FA-bypass capability, and SHALL keep credential values outside source, command arguments, and logs.

#### Scenario: Interactive login token

- **WHEN** a token authenticates `npm whoami` but requires an OTP for package writes
- **THEN** release documentation identifies it as unsuitable for `NPM_TOKEN` and directs the maintainer to replace the encrypted secret

### Requirement: Acceptance-gated publication

Courier SHALL publish a tagged release only after the repeated Linux release dry run, real-browser suite, distribution package suite, macOS runtime suite, Windows exact-coverage suite, installer tests, contract freshness, workflow validation, and strict OpenSpec validation succeed.

#### Scenario: Any acceptance job fails

- **WHEN** a release candidate fails one required platform, browser, package, installer, contract, workflow, or specification check
- **THEN** GitHub Release and npm publication jobs do not start

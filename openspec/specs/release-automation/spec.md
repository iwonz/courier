# release-automation Specification

## Purpose
Define local verification and ordered one-command release publication.

## Requirements

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

### Requirement: One-command release

Courier SHALL provide one end-to-end ship command that accepts a semantic version, a completed or archived OpenSpec change, and a conventional commit message. The command SHALL validate repository and npm credential configuration; require strict OpenSpec validation, completed tasks, archival, and no other active changes; synchronize the target release and generated contract artifacts; run the full local dry run; commit and push the reviewed worktree to `main`; create and push one annotated release tag; wait for GitHub Release, Scoop, Homebrew Formula, and npm verification; fast-forward to the verified Formula commit; and publish and wait for the final GitHub Pages deployment. The existing release command SHALL remain available for an already committed clean `main` and SHALL also wait for publication and Pages completion.

#### Scenario: A completed change is shipped

- **WHEN** a maintainer invokes the ship command with a valid version, completed change, and reviewed worktree
- **THEN** Courier archives the change if necessary, updates generated version artifacts, verifies, commits, pushes, publishes every configured release channel, deploys Pages, and reports the public release and site URLs

#### Scenario: Another change remains active

- **WHEN** an OpenSpec task is incomplete or another active change remains
- **THEN** the ship command stops before version mutation, commit, tag, or external publication

#### Scenario: A publication stage fails

- **WHEN** local verification, push, release workflow, Formula synchronization, package publication, or Pages deployment fails
- **THEN** the command stops at that boundary and does not report the lifecycle as complete

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

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

Courier SHALL publish the GitHub Release and Scoop manifest after quality gates, then run a separate macOS Formula job that reads the released source checksum, renders and audits the Formula, installs it from source, runs Courier, and commits the exact verified Formula to `main`. npm publication SHALL begin only after GitHub Release assets exist. Post-publication verification SHALL revalidate npm metadata online for a bounded interval that accommodates normal asynchronous registry processing before reporting failure.

#### Scenario: Formula verification fails

- **WHEN** the released source archive cannot be audited, built, installed, or executed by Homebrew
- **THEN** no Formula update is committed to `main` and the release workflow reports failure

#### Scenario: Release contents are verified

- **WHEN** release verification inspects a published version
- **THEN** it requires the deterministic source archive, its checksum, the committed Formula, and Scoop manifest and does not expect a cask

#### Scenario: npm processing is delayed

- **WHEN** npm accepts and signs the immutable version but public registry processing remains incomplete after one minute
- **THEN** verification continues online probes for up to six minutes, stops immediately when the exact version becomes visible, and does not republish the version

#### Scenario: npm never becomes public

- **WHEN** the exact published version remains unavailable after the bounded online retry window
- **THEN** verification fails with the package and version while leaving the accepted tag and other published channels unchanged

#### Scenario: Test failure

- **WHEN** a test or coverage gate fails for a release tag
- **THEN** no release, manifest, Formula, or npm publication job starts

### Requirement: One-command release

Courier SHALL provide one end-to-end ship command that runs from a feature branch containing exactly one completed OpenSpec change. It SHALL validate clean synchronized ancestry and publication configuration; select the next patch by default or accept an explicit version; archive the change, update target-release and generated artifacts, and run the complete local gate; create one conventional commit; require `origin/main` to remain at the recorded base; fast-forward local `main`; push and create one annotated release tag through the protected workflow; wait for GitHub Release, Scoop, Homebrew Formula, npm, and final GitHub Pages deployment; delete the local feature branch; and prove clean `main == origin/main`. Low-level release and Pages commands SHALL remain available only as explicitly authorized recovery operations.

#### Scenario: A completed feature change is shipped

- **WHEN** a maintainer invokes ship with a completed change and conventional message and omits the version
- **THEN** Courier selects the next patch, archives, generates, verifies, commits once, fast-forwards unchanged main, publishes every channel, deploys Pages, removes the feature branch, and reports the public URLs from clean synchronized main

#### Scenario: A completed change is shipped

- **WHEN** a maintainer invokes the ship command with a valid completed change, conventional message, and optional version override
- **THEN** Courier performs the feature-branch lifecycle and publishes every configured surface before reporting completion

#### Scenario: Another change remains active

- **WHEN** another OpenSpec change is active or the selected change has an incomplete task
- **THEN** shipping stops before version mutation, commit, merge, tag, or external publication

#### Scenario: Main moves before merge

- **WHEN** `origin/main` differs from the base recorded at change start or ship preflight
- **THEN** the command refuses the merge, push, and tag and leaves the feature commit recoverable

#### Scenario: Verification fails

- **WHEN** generation, focused checks, or `make verify` fails
- **THEN** no feature commit, main update, release tag, or publication is created

#### Scenario: Publication fails

- **WHEN** a post-tag channel or final Pages deployment fails
- **THEN** shipping reports the exact boundary, preserves published immutable state, and does not claim clean completion

#### Scenario: A publication stage fails

- **WHEN** push, release workflow, Formula synchronization, package publication, or Pages deployment fails
- **THEN** shipping stops at that boundary and does not report the lifecycle as complete

#### Scenario: Missing npm setup

- **WHEN** the npm publication secret is not configured
- **THEN** preflight identifies the missing connection without embedding credentials in source

#### Scenario: Missing catalog setup

- **WHEN** the Courier repository or in-repository manifest configuration is invalid
- **THEN** preflight identifies the exact problem without requiring an external catalog repository or embedding credentials in source

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

### Requirement: Release-gated Pages publication

Courier SHALL publish GitHub Pages only after the release matrix has succeeded, except for a documented explicitly authorized recovery dispatch. A normal push to `main` SHALL NOT trigger Pages deployment.

#### Scenario: Main receives a non-release push

- **WHEN** a commit reaches `main` without a completed release matrix
- **THEN** no Pages workflow starts and the public site cannot advertise unpublished artifacts

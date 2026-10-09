## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: Release-gated Pages publication

Courier SHALL publish GitHub Pages only after the release matrix has succeeded, except for a documented explicitly authorized recovery dispatch. A normal push to `main` SHALL NOT trigger Pages deployment.

#### Scenario: Main receives a non-release push

- **WHEN** a commit reaches `main` without a completed release matrix
- **THEN** no Pages workflow starts and the public site cannot advertise unpublished artifacts

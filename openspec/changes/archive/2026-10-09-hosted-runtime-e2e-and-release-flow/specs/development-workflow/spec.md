## ADDED Requirements

### Requirement: Tracked changes start from synchronized main

Every task that will mutate tracked files SHALL begin through one command that requires a clean local `main` equal to `origin/main`, rejects another active OpenSpec change, creates `TYPE/CHANGE`, and scaffolds that change. Read-only review, diagnosis, and planning are exempt.

#### Scenario: Main is dirty or behind origin

- **WHEN** a maintainer attempts to start a tracked change
- **THEN** the command fails before creating a branch or OpenSpec artifact

### Requirement: One completed change per feature branch

A shippable branch SHALL contain exactly one active completed OpenSpec change and no unrelated work. It SHALL pass focused tests and the complete verification gate before archival and publication.

#### Scenario: OpenSpec tasks are incomplete

- **WHEN** shipping inspects an unchecked task or invalid artifact
- **THEN** it fails before version generation, commit, merge, push, or tag creation

### Requirement: Clean published completion

Shipping SHALL create one conventional feature commit, fast-forward an unchanged local `main`, push and tag through the guarded workflow, wait for GitHub Release, Scoop, Homebrew, npm, and final Pages publication, delete the local feature branch, and prove clean `main == origin/main`.

#### Scenario: Origin main advances during verification

- **WHEN** the recorded base no longer equals `origin/main` before merge
- **THEN** shipping fails without a non-fast-forward merge, push, or tag

#### Scenario: Shipping succeeds

- **WHEN** every verification and publication surface succeeds
- **THEN** the repository ends on clean synchronized `main` with no active OpenSpec change or local feature branch

### Requirement: Protected recovery boundary

Direct pushes of `main` or release tags SHALL be rejected outside ship or an explicitly authorized recovery flow. Low-level release and Pages targets SHALL be documented and executable only as recovery operations.

#### Scenario: Developer pushes main manually

- **WHEN** the pre-push hook sees a main or release-tag update without an allowed workflow token
- **THEN** it rejects the push with the supported ship or recovery instruction

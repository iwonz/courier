## MODIFIED Requirements

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

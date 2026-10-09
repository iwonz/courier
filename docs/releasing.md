# Release runbook

Courier uses GoReleaser Community v2.18.1. A semantic Git tag is the only publication trigger. GitHub Actions repeats all quality gates before creating a GitHub Release and Scoop manifest, source-building and committing the in-repository Homebrew Formula, and publishing npm.

## One-time external setup

No catalog repository or upstream package-manager pull request is required. `iwonz/courier` is both the source repository and the custom Homebrew tap/Scoop bucket. Ensure its GitHub Actions setting permits workflows to request read/write permissions. Create the public npm package scope access needed to publish `@iwonz/courier`.

Add these GitHub Actions repository secrets to `iwonz/courier`:

| Secret | Required access |
|---|---|
| `NPM_TOKEN` | Publish `@iwonz/courier` on npmjs; use a granular access token with package read/write permission and bypass 2FA enabled |

`GITHUB_TOKEN` is supplied automatically by Actions for the Courier GitHub Release and manifest commits. Never commit or pass any token as a CLI argument.

An interactive `npm login` token is not a CI publication credential: npm accepts it for account queries but requires a one-time password for package writes. Create `NPM_TOKEN` in npm's granular access-token settings, scope it as narrowly as npm permits, enable package read/write and bypass 2FA, and send it to GitHub through the repository secret UI or standard input. Never paste it into source, workflow YAML, a command argument, or an issue.

The local preflight checks repository visibility and secret names with GitHub CLI, but GitHub does not expose secret values. A workflow preflight checks that values are non-empty before publication jobs start.

## Local dry run

Install the tracked hook once:

```sh
make hooks
```

Every commit then runs the same local gate. It can also be invoked directly:

```sh
make verify
```

The gate runs:

1. `gofmt`, `go vet`, command-contract freshness, GitHub Actions syntax, and strict OpenSpec validation;
2. all Go tests with the race detector on Linux plus exact 100% first-party statement coverage and compiled runtime checks on Linux, macOS, and Windows;
3. exact TypeScript coverage, deterministic UI builds, embedded-asset freshness, and real Chromium acceptance for all three browser surfaces;
4. npm wrapper and POSIX/PowerShell installer acceptance;
5. GoReleaser configuration validation and `goreleaser release --snapshot --clean` with the deterministic source archive;
6. one Formula renderer, Ruby syntax/contract checks, and a Homebrew source build when Homebrew is available;
7. checksum and complete primary/BSD/source artifact-matrix verification;
8. labeled native-package installation in Ubuntu, Debian, Arch, Manjaro, Fedora, Red Hat UBI, and Alpine containers, followed by ownership-scoped cleanup assertions.

The artifact matrix includes the six primary macOS/Linux/Windows targets plus exact-platform BSD helper archives. Helper archives use the same `courier_<version>_<os>_<arch>.tar.gz` convention and are verified by the snapshot gate and post-publication workflow.

The pinned GoReleaser binary is downloaded into `.cache/tools` only when absent and is verified against the upstream release checksum. No global GoReleaser installation is required.

## Start, complete, and publish one change

Every task that changes tracked files uses the same end-to-end path. Begin only from clean synchronized `main`:

```sh
make change-start TYPE=feat CHANGE=my-openspec-change
```

The command fetches `origin/main`, rejects dirty or stale state and any existing active OpenSpec change, creates `feat/my-openspec-change`, and scaffolds `openspec/changes/my-openspec-change`. Supported conventional types are `feat`, `fix`, `chore`, `docs`, `refactor`, `test`, `build`, `ci`, `perf`, and `style`.

Complete the OpenSpec artifacts and task list, implement the work, run focused tests, and run `make verify`. Then hand the entire finished worktree to the only normal publication command:

```sh
make ship CHANGE=my-openspec-change MESSAGE="feat: describe the finished change"
```

`VERSION` is optional. By default, `ship` increments the patch component of the newest stable semantic tag; use `VERSION=1.2.3` only when an explicit release number is required. `ship` requires the matching feature branch to have started exactly at unchanged `origin/main`, exactly one completed and strictly valid OpenSpec change, no earlier feature commits, and a configured publication environment. It archives the change, updates `target_release`, regenerates contract and landing assets, runs `make verify`, creates the one conventional commit, rechecks that local and remote `main` did not move, and performs a local `--ff-only` merge. It then pushes `main` through the protected flow, creates the annotated tag, waits for GitHub Release, Scoop, Homebrew, npm, and their final verification, and only then dispatches and waits for Pages. Finally it deletes the local feature branch and proves the repository is on clean `main == origin/main`.

Any failed gate stops immediately. Before the commit, the dirty feature worktree remains available for inspection; after the commit, the feature branch preserves that commit; after a tag is published, the immutable tag is never moved or recreated. Fix the failed prerequisite or publication job and follow the recovery boundary below. Use `COURIER_SHIP_YES=1` only for deliberate non-interactive execution; it bypasses the prompt, not a preflight or quality gate.

The tracked pre-push hook rejects manual pushes to `main` and semantic release tags. Install it once with `make hooks`. GitHub Pages has no automatic `main` push trigger, so the landing cannot advertise a version before its release matrix succeeds.

## Recovery after a partially completed ship

`make release` and `make pages-publish` are deliberately disabled during normal development. They are low-level recovery commands for a ship that already produced the correct clean synchronized `main` or completed release matrix:

```sh
make release VERSION=1.2.3 RECOVERY=1
make pages-publish RECOVERY=1
```

Use release recovery only when the release commit is already on `origin/main`, its `target_release` matches, no OpenSpec change is active, and the semantic tag does not yet exist. It repeats repository, authentication, verification, ancestry, tag-uniqueness, publication, Formula synchronization, and Pages gates. Use Pages recovery only after the complete release matrix is known to have succeeded. For a narrowly required manual push during incident recovery, set `COURIER_RECOVERY=1` for that one push and record why; this disables only the local hook, never GitHub protection or release validation.

`COURIER_RELEASE_VERIFIED_COMMIT`, `COURIER_RELEASE_FLOW`, `COURIER_PUSH_FLOW`, and `COURIER_PAGES_FLOW` are internal handoff variables reserved for `make ship`. They are not user-facing bypasses.

GoReleaser generates release notes from conventional commits between tags, publishes `courier_<version>_source.tar.gz`, and updates `bucket/courier.json`. A separate macOS job reads the published source checksum, renders `Formula/courier.rb`, runs `brew audit`, installs with `--build-from-source`, runs `courier version`, and commits only that verified Formula to `main` with the workflow's short-lived `GITHUB_TOKEN`. No personal GitHub token or additional repository is involved. Never move or recreate a published tag. If npm publication fails before that version is published, repair the secret and rerun the failed workflow for the existing tag. npm versions are immutable, so confirm publication status before rerunning a failed npm job. A successful npm publish can remain in asynchronous registry processing for several minutes; final verification revalidates online every ten seconds for up to six minutes and never republishes an accepted version.

## Publication order

The release workflow enforces this sequence:

```text
credentials + Linux acceptance + macOS/Windows runtime acceptance
                              |
                              v
        GitHub Release + Scoop manifest
                    |
          +---------+---------+
          |                   |
          v                   v
Homebrew Formula audit,   npmjs publication
source build, and commit       |
          |                   |
          +---------+---------+
                    v
       GitHub/npm/Formula/Scoop verification
```

See [Acceptance and release-candidate verification](acceptance.md) for prerequisites, focused commands, resource ownership, and the browser/package boundaries.

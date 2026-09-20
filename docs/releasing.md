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

## Complete and publish one change

The preferred end-to-end command owns the complete handoff from a finished worktree to public distribution:

```sh
make ship VERSION=1.2.3 CHANGE=my-openspec-change MESSAGE="feat: describe the finished change"
```

`make ship` requires `main`, a configured GitHub/npm publication environment, completed OpenSpec tasks, strict validation, and no other active OpenSpec changes. It archives the named change when necessary, updates `target_release`, regenerates contract artifacts, runs `make verify`, stages the complete reviewed worktree, creates the conventional commit, pushes `main`, creates and pushes the annotated release tag, waits for GitHub Release/Scoop/Homebrew/npm verification, fast-forwards to the workflow's verified Formula commit, and dispatches and waits for the final GitHub Pages deployment. Any failed gate stops the sequence. Use `COURIER_SHIP_YES=1` only for deliberate non-interactive execution.

## Publish one already committed release

Start from a clean, synchronized `main` branch and run either form:

```sh
make release VERSION=1.2.3
./scripts/release.sh 1.2.3
```

Omit the version to be prompted:

```sh
./scripts/release.sh
```

The command checks GitHub authentication, the Courier repository, GitHub secret names, origin, branch, clean state, tag uniqueness, and synchronization with `origin/main`. It then runs the full dry run, asks for final confirmation, creates one annotated `vMAJOR.MINOR.PATCH` tag, pushes only that tag, waits for the complete release workflow, fast-forwards local `main` to the verified Formula commit, and publishes the final landing through GitHub Pages.

For deliberate non-interactive automation, set `COURIER_RELEASE_YES=1` and provide the version. This does not bypass any quality or repository preflight. `COURIER_RELEASE_VERIFIED_COMMIT` is reserved for `make ship`; it skips the duplicate local gate only when it exactly equals the commit just verified by the ship command.

GoReleaser generates release notes from conventional commits between tags, publishes `courier_<version>_source.tar.gz`, and updates `bucket/courier.json`. A separate macOS job reads the published source checksum, renders `Formula/courier.rb`, runs `brew audit`, installs with `--build-from-source`, runs `courier version`, and commits only that verified Formula to `main` with the workflow's short-lived `GITHUB_TOKEN`. No personal GitHub token or additional repository is involved. Never move or recreate a published tag. If npm publication fails before that version is published, repair the secret and rerun the failed workflow for the existing tag. npm versions are immutable, so confirm publication status before rerunning a failed npm job.

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

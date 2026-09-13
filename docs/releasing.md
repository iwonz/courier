# Release runbook

Courier uses GoReleaser Community v2.18.1. A semantic Git tag is the only publication trigger. GitHub Actions repeats all quality gates before creating a GitHub Release, updating Homebrew and Scoop, opening the Winget pull request, and publishing npm.

## One-time external setup

Create these public repositories under the `iwonz` account:

- `iwonz/homebrew-tap`, with default branch `main`;
- `iwonz/scoop-bucket`, with default branch `main`;
- `iwonz/winget-pkgs`, as a fork of `microsoft/winget-pkgs`.

Ensure the GitHub repository Actions setting permits workflows to request read/write permissions. Create the public npm package scope access needed to publish `@iwonz/courier`.

Add these GitHub Actions repository secrets to `iwonz/courier`:

| Secret | Required access |
|---|---|
| `NPM_TOKEN` | Publish `@iwonz/courier` on npmjs; use a granular access token with package read/write permission and bypass 2FA enabled |
| `HOMEBREW_TAP_GITHUB_TOKEN` | Contents read/write on `iwonz/homebrew-tap` |
| `SCOOP_BUCKET_GITHUB_TOKEN` | Contents read/write on `iwonz/scoop-bucket` |
| `WINGET_GITHUB_TOKEN` | Contents read/write on the `iwonz/winget-pkgs` fork and permission to open the upstream pull request |

`GITHUB_TOKEN` is supplied automatically by Actions for the Courier GitHub Release. Never commit or pass any token as a CLI argument.

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

1. `gofmt` verification and `go vet`;
2. all Go tests with the race detector on the primary quality runner;
3. exact 100% first-party Go statement coverage on both Linux and Windows, including platform-specific agent and update code;
4. local POSIX installer acceptance tests and npm package tests;
5. GoReleaser configuration validation;
6. `goreleaser release --snapshot --clean`;
7. checksum and artifact-matrix verification.

The artifact matrix includes the six primary macOS/Linux/Windows targets plus exact-platform BSD helper archives. Helper archives use the same `courier_<version>_<os>_<arch>.tar.gz` convention and are verified by the snapshot gate and post-publication workflow.

The pinned GoReleaser binary is downloaded into `.cache/tools` only when absent and is verified against the upstream release checksum. No global GoReleaser installation is required.

## Publish one release

Start from a clean, synchronized `main` branch and run either form:

```sh
make release VERSION=1.2.3
./scripts/release.sh 1.2.3
```

Omit the version to be prompted:

```sh
./scripts/release.sh
```

The command checks GitHub authentication, external repositories, GitHub secret names, origin, branch, clean state, tag uniqueness, and synchronization with `origin/main`. It then runs the full dry run, asks for final confirmation, creates one annotated `vMAJOR.MINOR.PATCH` tag, and pushes only that tag.

For deliberate non-interactive automation, set `COURIER_RELEASE_YES=1` and provide the version. This does not bypass any quality or repository preflight.

GoReleaser generates release notes from conventional commits between tags. Never move or recreate a published tag. If a catalog or npm credential fails before that version is published, repair the secret and rerun the failed workflow for the existing tag. npm versions are immutable, so confirm publication status before rerunning a failed npm job.

## Publication order

The release workflow enforces this sequence:

```text
credentials + quality + PowerShell acceptance
                    |
                    v
GitHub Release + Homebrew + Scoop + Winget PR
                    |
                    v
             npmjs publication
                    |
                    v
       GitHub/npm visibility verification
```

Winget availability remains asynchronous because Microsoft reviews and merges the generated upstream pull request.

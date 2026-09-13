# Change: Keep package distribution in the Courier repository

## Why

The first production release used dedicated Homebrew and Scoop repositories and an upstream Winget pull request. Those channels work, but they create repositories and external review obligations that are disproportionate for Courier. Homebrew custom taps and Scoop custom buckets can use the main project repository directly, while the public Winget catalog cannot operate without the external `microsoft/winget-pkgs` workflow.

## What Changes

- Store the generated Homebrew cask under `Casks/` and the Scoop manifest under `bucket/` in `iwonz/courier`.
- Publish both manifests with the release workflow's short-lived `GITHUB_TOKEN` instead of personal cross-repository tokens.
- Remove Winget catalog publication and verification while retaining Windows installation through Scoop, npm, PowerShell, and direct binaries.
- Remove all release preflight requirements for auxiliary repositories and GitHub credentials.

## Non-goals

- Operating native apt, rpm, Alpine, or Arch package repositories.
- Replacing the self-contained GitHub Release artifacts or npm distribution.
- Moving or recreating the existing `v0.1.0` tag.

## Impact

- Affected specs: `install-channels`, `release-automation`.
- Affected automation: GoReleaser publishers, release preflight, and publication verification.
- Affected documentation: installation and release runbooks.
- External cleanup: close the obsolete Winget pull request and remove the three Courier-only catalog repositories after the in-repository channels are verified.

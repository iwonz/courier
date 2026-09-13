## Why

Courier is only useful as a standalone tool when users can install a verified native binary from their existing ecosystem. Releases also need one reproducible quality gate so publication cannot bypass tests, coverage, or cross-platform compilation.

## What Changes

- Make GoReleaser Community the source of release archives, checksums, Linux packages, GitHub Releases, Homebrew, Scoop, and Winget metadata.
- Publish standalone macOS, Linux, and Windows binaries for amd64 and arm64.
- Add dependency-free POSIX and PowerShell bootstrap installers that verify the release checksum before installation.
- Add the `@iwonz/courier` npm package for npm, npx, yarn, and pnpm; keep this small publisher beside GoReleaser because GoReleaser's native npm publisher is a Pro feature.
- Add GitHub Actions quality and release pipelines with ordered test, build, GitHub/package-manager publication, npm publication, and verification stages.
- Add one local release command that runs the complete dry run, creates an annotated semantic-version tag, and pushes it to start the release workflow.
- Document repositories, secrets, platform coverage, direct downloads, and recovery procedures in English.

## Capabilities

### New Capabilities

- `release-artifacts`: reproducible standalone binaries, archives, checksums, and Linux packages.
- `install-channels`: verified script and package-manager installation paths.
- `release-automation`: local dry runs, CI quality gates, tags, notes, and publication.

## Impact

Adds GoReleaser configuration, GitHub Actions, npm and installer sources, release tooling, integration checks, and distribution documentation. It introduces no runtime dependency into the Courier binary.

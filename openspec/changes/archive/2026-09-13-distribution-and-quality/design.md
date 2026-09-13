## Context

The Go application is cross-buildable and fully tested, but no published artifact or installation path exists. GoReleaser Community already owns the difficult artifact and catalog metadata; custom code is limited to channels that Community does not publish, installer bootstrap, and release orchestration.

## Goals / Non-Goals

**Goals:** six standalone target binaries, deterministic archives and checksums, deb/rpm/apk/Arch packages, GitHub Releases, npm, Homebrew, Scoop, Winget, verified install scripts, CI gates, release notes, semantic tags, and a one-command local release.

**Non-Goals:** paid GoReleaser Pro features, operating a distribution package mirror, auto-merging third-party catalog pull requests, or requiring package managers at Courier runtime.

## Decisions

### GoReleaser Community owns artifacts and supported catalog publishers

Two CGO-disabled build definitions cover Unix and Windows. A tar.gz archive exists for every target because `courier update` and npm consume one stable shape; raw binaries support minimal bootstrap installers and direct downloads; and an additional Windows zip supports Scoop and Winget. nFPM emits deb, rpm, apk, and Arch Linux packages from the Linux binaries. GoReleaser publishes the GitHub Release and updates dedicated GitHub repositories for Homebrew and Scoop; its Winget publisher opens a pull request from the configured fork.

### npm remains a narrow adjacent publisher

GoReleaser's native npm publisher is paid. The Community pipeline therefore publishes a small, dependency-free JavaScript package after the GitHub assets exist. Its postinstall selects the native tar.gz by platform, verifies `checksums.txt`, safely extracts only the Courier executable, and installs it privately inside the package. npm, npx, yarn, and pnpm all consume the same package.

### Installers never trust an unverified binary

The POSIX installer supports curl or wget and installs below the user's home by default. The PowerShell installer uses built-in web APIs and the current-user application directory. Both resolve a release, download its raw binary and checksum manifest into a private temporary directory, compare SHA-256, install atomically where practical, and always remove temporary state.

### A tag is the only publication trigger

`make verify` runs formatting checks, vet, race tests, exact coverage enforcement, npm tests, installer acceptance checks, and `goreleaser release --snapshot --clean`. `make release` runs that gate before creating and pushing one annotated `vMAJOR.MINOR.PATCH` tag. GitHub Actions repeats quality checks, publishes with GoReleaser, publishes npm only after GitHub assets succeed, and verifies the release surface.

## External Repository Contract

- `iwonz/homebrew-tap` receives `Formula/courier.rb`.
- `iwonz/scoop-bucket` receives `bucket/courier.json`.
- `iwonz/winget-pkgs` is a fork used to open pull requests to `microsoft/winget-pkgs`.

The release workflow documents the fine-grained tokens required for these repositories and npm. Missing external setup fails before a partial catalog publication is presented as complete.

## Risks / Trade-offs

- Third-party Winget review and merge are asynchronous even though manifest creation and pull-request submission are automated.
- npm installation necessarily uses the package manager's Node runtime to bootstrap and launch the packaged native executable; the Courier executable itself remains standalone.
- Homebrew casks cover supported macOS hosts; Linux users use native packages or the POSIX installer. Apple signing and notarization remain a future credentialed hardening step.
- Unsigned macOS binaries may require future Apple signing/notarization credentials; the project does not disable quarantine or instruct users to bypass host security.

## Migration Plan

Create the tap, bucket, and Winget fork once, add the documented GitHub secrets, merge this change to the default branch, then run the local release command. If a downstream catalog publication fails, rerun the GitHub workflow for the same tag after repairing credentials; never move or recreate a published version tag.

## Open Questions

None.

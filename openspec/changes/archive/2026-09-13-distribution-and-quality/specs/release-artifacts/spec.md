## ADDED Requirements

### Requirement: Standalone target matrix

Courier SHALL publish CGO-disabled binaries for macOS, Linux, and Windows on amd64 and arm64, with embedded semantic version, commit, and build date.

#### Scenario: Windows arm64 release

- **WHEN** GoReleaser builds a tagged release
- **THEN** the release contains a standalone Windows arm64 Courier executable

### Requirement: Deterministic verified artifacts

Courier SHALL publish versioned tar.gz archives and raw binaries for every target, Windows zip archives, and one SHA-256 manifest covering release artifacts.

#### Scenario: Self-update asset selection

- **WHEN** Courier version `1.2.3` runs on Linux amd64
- **THEN** `courier_1.2.3_linux_amd64.tar.gz` and its `checksums.txt` entry exist

### Requirement: Native Linux packages

Courier SHALL generate deb, rpm, apk, and Arch Linux packages for Linux amd64 and arm64 with no external runtime dependencies.

#### Scenario: Minimal distribution installation

- **WHEN** a supported native package is installed on Ubuntu, Debian, Fedora, RHEL, Alpine, Arch Linux, or Manjaro
- **THEN** `/usr/bin/courier` runs without rsync, OpenSSH, or an archiver package

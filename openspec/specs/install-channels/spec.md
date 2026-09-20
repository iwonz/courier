# install-channels Specification

## Purpose
Define verified installation paths that remain within the Courier repository.

## Requirements

### Requirement: Verified bootstrap installers

Courier SHALL provide POSIX and PowerShell installers that resolve a requested or latest version, select the runtime asset, verify SHA-256, install without administrator privileges by default, and remove all temporary state.

#### Scenario: Corrupted bootstrap download

- **WHEN** an installer's downloaded asset does not match `checksums.txt`
- **THEN** installation fails without replacing an existing Courier executable

### Requirement: JavaScript package managers

Courier SHALL publish `@iwonz/courier` so npm, npx, yarn, and pnpm install and invoke the verified native executable for the host platform.

#### Scenario: npm postinstall

- **WHEN** npm installs `@iwonz/courier@1.2.3` on macOS arm64
- **THEN** postinstall verifies and safely extracts the matching GitHub Release asset before the command becomes available

### Requirement: Native package catalogs

Courier SHALL publish a source-built Homebrew Formula and Scoop metadata inside `iwonz/courier` and SHALL NOT require an additional project-owned repository, Apple Developer ID, Gatekeeper bypass, or upstream catalog pull request for an advertised installation channel.

#### Scenario: Homebrew installation

- **WHEN** a user registers `iwonz/courier` as a custom tap with its explicit GitHub URL and installs `courier`
- **THEN** Homebrew verifies the deterministic source archive, installs Go as a build dependency, builds Courier with CGO disabled and release metadata, and installs the resulting executable

#### Scenario: Existing cask installation is migrated

- **WHEN** a user already installed the retired Courier cask
- **THEN** the installation guide removes that cask before installing the Formula and never instructs the user to disable quarantine or Gatekeeper

#### Scenario: Scoop installation

- **WHEN** a user registers `iwonz/courier` as a custom Scoop bucket and installs `courier`
- **THEN** Scoop installs the architecture-specific checksummed Windows archive from the in-repository manifest

### Requirement: Direct downloads

Courier SHALL document direct archive and checksum downloads for every supported OS and architecture.

#### Scenario: No package manager

- **WHEN** a user has no supported package manager
- **THEN** the installation guide provides a verified script or direct binary procedure

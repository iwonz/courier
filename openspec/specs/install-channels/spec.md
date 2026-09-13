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

Courier SHALL publish Homebrew and Scoop metadata inside `iwonz/courier` and SHALL NOT require an additional project-owned repository or upstream catalog pull request for any advertised installation channel.

#### Scenario: Homebrew installation

- **WHEN** a user registers `iwonz/courier` as a custom tap with its explicit GitHub URL and installs `courier`
- **THEN** Homebrew installs a checksummed GitHub Release binary from the in-repository cask

#### Scenario: Scoop installation

- **WHEN** a user registers `iwonz/courier` as a custom Scoop bucket and installs `courier`
- **THEN** Scoop installs the architecture-specific checksummed Windows archive from the in-repository manifest

### Requirement: Direct downloads

Courier SHALL document direct archive and checksum downloads for every supported OS and architecture.

#### Scenario: No package manager

- **WHEN** a user has no supported package manager
- **THEN** the installation guide provides a verified script or direct binary procedure

### Requirement: Homebrew cask trust

Courier SHALL document an item-scoped Homebrew trust step before registering its non-official custom tap and SHALL NOT require users to disable Homebrew's tap trust policy.

#### Scenario: Fresh Homebrew 6 installation

- **WHEN** a user has no Courier tap or trust entry and follows the documented Homebrew sequence
- **THEN** Homebrew trusts only `iwonz/courier/courier`, registers `iwonz/courier` from its explicit URL, and installs the checksummed cask

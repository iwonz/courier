## ADDED Requirements

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

Courier SHALL publish Homebrew and Scoop metadata to project-owned GitHub catalog repositories and SHALL submit versioned manifests to Winget through its official repository workflow.

#### Scenario: Homebrew installation

- **WHEN** a user taps `iwonz/tap` and installs `courier`
- **THEN** Homebrew installs a checksummed GitHub Release binary

### Requirement: Direct downloads

Courier SHALL document direct archive and checksum downloads for every supported OS and architecture.

#### Scenario: No package manager

- **WHEN** a user has no supported package manager
- **THEN** the installation guide provides a verified script or direct binary procedure

## MODIFIED Requirements

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

## REMOVED Requirements

### Requirement: Homebrew cask trust

**Reason**: The source-built Formula is not a quarantined prebuilt cask and does not require item-scoped cask trust.

**Migration**: Uninstall the prior cask, register the explicit repository tap, and install `iwonz/courier/courier` as a Formula.

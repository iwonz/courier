## MODIFIED Requirements

### Requirement: Native package catalogs

Courier SHALL publish Homebrew and Scoop metadata inside `iwonz/courier` and SHALL NOT require an additional project-owned repository or upstream catalog pull request for any advertised installation channel.

#### Scenario: Homebrew installation

- **WHEN** a user registers `iwonz/courier` as a custom tap with its explicit GitHub URL and installs `courier`
- **THEN** Homebrew installs a checksummed GitHub Release binary from the in-repository cask

#### Scenario: Scoop installation

- **WHEN** a user registers `iwonz/courier` as a custom Scoop bucket and installs `courier`
- **THEN** Scoop installs the architecture-specific checksummed Windows archive from the in-repository manifest

## ADDED Requirements

### Requirement: Typed self-update progress

The updater SHALL preserve `Updater.Run`, SHALL offer `RunWithProgress` with a typed optional event sink, and SHALL report release check, archive download, checksum download, verification, extraction, installation, and completion. Download events SHALL report current bytes and report total bytes when Content-Length is known.

#### Scenario: Content length is unknown

- **WHEN** a release server streams an update asset without a known length
- **THEN** progress reports downloaded bytes without inventing a percentage and verification proceeds normally

### Requirement: Canonical release destination

Update results SHALL include the canonical GitHub Release URL. Interactive update completion SHALL show the installed/current version and explicit URL after release notes, while plain update and version output SHALL append `release: <URL>`. Development builds SHALL use the repository releases index and released versions SHALL use their exact tag URL.

#### Scenario: A released build reports its version

- **WHEN** `courier version` runs for `vX.Y.Z`
- **THEN** output includes `https://github.com/iwonz/courier/releases/tag/vX.Y.Z` without an OSC-8 hyperlink sequence

## MODIFIED Requirements

### Requirement: Ordered CI publication

Courier SHALL publish the GitHub Release and Scoop manifest after quality gates, then run a separate macOS Formula job that reads the released source checksum, renders and audits the Formula, installs it from source, runs Courier, and commits the exact verified Formula to `main`. npm publication SHALL begin only after GitHub Release assets exist. Post-publication verification SHALL revalidate npm metadata online for a bounded interval that accommodates normal asynchronous registry processing before reporting failure.

#### Scenario: Formula verification fails

- **WHEN** the released source archive cannot be audited, built, installed, or executed by Homebrew
- **THEN** no Formula update is committed to `main` and the release workflow reports failure

#### Scenario: Release contents are verified

- **WHEN** release verification inspects a published version
- **THEN** it requires the deterministic source archive, its checksum, the committed Formula, and Scoop manifest and does not expect a cask

#### Scenario: npm processing is delayed

- **WHEN** npm accepts and signs the immutable version but public registry processing remains incomplete after one minute
- **THEN** verification continues online probes for up to six minutes, stops immediately when the exact version becomes visible, and does not republish the version

#### Scenario: npm never becomes public

- **WHEN** the exact published version remains unavailable after the bounded online retry window
- **THEN** verification fails with the package and version while leaving the accepted tag and other published channels unchanged

#### Scenario: Test failure

- **WHEN** a test or coverage gate fails for a release tag
- **THEN** no release, manifest, Formula, or npm publication job starts

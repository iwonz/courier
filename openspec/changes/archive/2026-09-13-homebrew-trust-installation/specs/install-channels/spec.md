## ADDED Requirements

### Requirement: Homebrew cask trust

Courier SHALL document an item-scoped Homebrew trust step before registering its non-official custom tap and SHALL NOT require users to disable Homebrew's tap trust policy.

#### Scenario: Fresh Homebrew 6 installation

- **WHEN** a user has no Courier tap or trust entry and follows the documented Homebrew sequence
- **THEN** Homebrew trusts only `iwonz/courier/courier`, registers `iwonz/courier` from its explicit URL, and installs the checksummed cask

## ADDED Requirements

### Requirement: Publication-capable npm credential

Courier SHALL require npm CI publication credentials to have package read/write authority and non-interactive 2FA-bypass capability, and SHALL keep credential values outside source, command arguments, and logs.

#### Scenario: Interactive login token

- **WHEN** a token authenticates `npm whoami` but requires an OTP for package writes
- **THEN** release documentation identifies it as unsuitable for `NPM_TOKEN` and directs the maintainer to replace the encrypted secret

### Requirement: Exact Winget pull-request verification

Courier SHALL verify the open upstream Winget pull request using the fork owner and release branch returned by GitHub's pull-request API.

#### Scenario: GoReleaser-created pull request

- **WHEN** GoReleaser opens the upstream pull request from `iwonz:winget-pkgs:courier-<version>`
- **THEN** the verification job finds its `iwonz:courier-<version>` head label through the REST filter

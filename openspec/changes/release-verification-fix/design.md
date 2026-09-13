# Design: Deterministic downstream verification

## Context

Release systems expose two distinct facts: a secret exists, and that credential has the authority required for publication. GitHub only reveals secret names, so Courier cannot determine npm write or 2FA-bypass capability during local preflight. The runbook must state the credential contract precisely and the pipeline must surface registry rejection without exposing the value.

GoReleaser creates the Winget branch `courier-<version>` in `iwonz/winget-pkgs` and opens an upstream pull request whose REST head label is `iwonz:courier-<version>`. The GitHub CLI list command does not reliably resolve that owner-qualified filter even though the pull request exists.

## Decisions

### npm credential contract

`NPM_TOKEN` remains an encrypted Actions secret supplied through `NODE_AUTH_TOKEN`. Documentation requires a granular access token with package read/write authority and bypass 2FA enabled. Interactive login tokens are explicitly excluded. Token values never enter workflow source, git history, command arguments, or logs.

### Winget REST verification

The verification job queries `repos/microsoft/winget-pkgs/pulls` with `state=open` and the exact `head=iwonz:courier-<version>` filter, then requires at least one result. This is the same label returned by the pull-request API and is independent of GitHub CLI's higher-level list filtering behavior.

### Immutable recovery

Published Git tags and npm versions are never moved or overwritten. A failed npm publish may be rerun for the same tag only after confirming that npm did not accept that version. Later workflow corrections ship in a new patch version.

## Testing

- Validate the workflow with actionlint.
- Validate the OpenSpec change strictly.
- Query the live `v0.1.0` Winget pull request with the exact REST filter.
- Retain the complete local GoReleaser dry run and exact coverage gate.

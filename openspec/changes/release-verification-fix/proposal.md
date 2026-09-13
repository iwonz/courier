# Change: Harden release credentials and Winget verification

## Why

The first production release proved that an interactive npm login token can authenticate account queries while still failing package publication with `EOTP`. It also exposed a GitHub CLI filtering mismatch: GoReleaser created the correct upstream Winget pull request, but `gh pr list --head owner:branch` did not find it.

## What Changes

- Document the exact npm granular access-token properties required by the CI publisher and safe secret handling.
- Verify the Winget pull request through GitHub's REST `head` filter, which matches the fork owner and release branch exactly.
- Preserve the immutable-tag recovery model: repair credentials and rerun only while the npm version remains unpublished.

## Non-goals

- Storing credentials in source or accepting them as release command arguments.
- Auto-merging the Microsoft Winget pull request.
- Moving or recreating `v0.1.0` after partial downstream publication.

## Impact

- Affected spec: `release-automation`.
- Affected automation: post-publication verification.
- Affected documentation: npm release credential setup and recovery.

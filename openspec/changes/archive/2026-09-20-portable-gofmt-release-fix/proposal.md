# Change: Portable Go formatting release fix

## Why

The v0.3.0 publication gate exposed a formatter difference between the supported CI Go 1.25 toolchain and a newer local Go 1.27 toolchain. Two compact map literals were accepted locally but rejected by the pinned release runner, correctly preventing publication before any release artifacts were created.

## What Changes

- Rewrite the affected updater test fixtures into a multiline form with stable formatting under Go 1.25 and Go 1.27.
- Record cross-version formatter stability as part of release-candidate acceptance.
- Publish the corrected immutable release as v0.3.1 without moving the failed v0.3.0 tag.

## Impact

Only test-source formatting and release acceptance are affected. Runtime behavior is unchanged.

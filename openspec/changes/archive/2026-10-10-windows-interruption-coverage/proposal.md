## Why

The `v0.3.15` release candidate passed Windows tests but stopped before publication because two defensive branches added by planned-interruption reporting were not executed, reducing the platform's aggregate statement coverage from the required `100.0%` to `99.9%`.

## What Changes

- Exercise an empty multi-error unwrap so cancellation classification rejects a tree without leaves.
- Exercise the command-error traversal depth guard so malformed or excessively deep wrappers terminate safely.
- Publish the unchanged runtime behavior in the next patch rather than moving the failed immutable `v0.3.15` tag.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. This is test-only completion of coverage for existing defensive behavior.

## Impact

This change affects Go unit tests and release verification only. Runtime code, public output, protocols, state, and command syntax remain unchanged.

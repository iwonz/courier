# Change: Enforce safe copy semantics

## Why

The current transfer engine replaces existing targets and synchronizes directories by deleting destination-only entries. That behavior conflicts with Courier's accepted no-overwrite contract and makes an ordinary copy unexpectedly destructive.

## What Changes

- Return success without transfer work when plain path endpoints resolve to the same location.
- Reject identical transformed operations and directory-descendant targets.
- Reject every existing final-path collision during preflight.
- Preserve unrelated destination entries and remove mirror/replacement commit behavior.
- Commit only into an absent final path, using no-replace backend primitives where available.

## Impact

Existing targets are never overwritten or merged. Callers must choose an absent exact path or a directory whose computed child path is absent.

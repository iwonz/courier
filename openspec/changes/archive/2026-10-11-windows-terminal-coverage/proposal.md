# Change: Restore exact Windows terminal renderer coverage

## Why

The theme-neutral terminal renderer added for `v0.3.18` left its wide-layout multiline continuation statement unexercised. POSIX aggregate coverage still rounded to 100.0 percent, but the native Windows release job reported 99.9 percent and correctly stopped publication. The immutable `v0.3.18` tag therefore remains unpublished and the coverage gap must be corrected in a new patch release.

## What Changes

- Add a deterministic renderer test whose labeled value wraps in the wide panel layout and exercises every continuation line.
- Assert continuation indentation, complete content, terminal-width containment, and theme-neutral output.
- Keep runtime presentation and all public behavior unchanged.

## Impact

- Affected spec: `cross-platform-acceptance`
- Affected code: shared `internal/terminalui` tests only
- No CLI, IPC, registry, worker, administration-state, contract, or product-behavior changes

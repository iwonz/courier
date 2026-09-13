# Change: Add the canonical CLI contract and provider registry

## Why

Courier's public commands are currently defined only by Cobra code and prose, which allows help, documentation, and future route combinations to drift.

## What Changes

- Add an English, machine-readable CLI contract and deterministic generated reference.
- Validate the shipped Cobra tree against the contract.
- Register commands through independent providers and reject duplicate registrations.
- Keep `help`, `version`, and `update`, while disabling Cobra's unsolicited `completion` command.

## Impact

The public surface becomes explicitly versioned. Existing shipped commands remain compatible except that the undocumented `completion` command is removed.

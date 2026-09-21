## Why

The unified landing still has a card-tinted hero and Relay area, repeated `Command` labels, and two installation actions without a leading icon. The masthead wordmark also needs an explicit vertical alignment with the Relay mark.

## What Changes

- Put the hero, Relay, and command builder directly on the page canvas, and remove the command-list and readout labels.
- Center the masthead wordmark on the mark's horizontal axis.
- Add a functional download icon to the direct-binaries link and reuse the official local npm mark for npx.
- Extend unit and browser acceptance, documentation, and OpenSpec for the revised presentation.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-landing`: simplify the command surface and make both install destinations visually identifiable.
- `ui-kit`: allow npx to reuse npm's official mark without inventing an npx brand.

## Impact

Presentation, localized catalog keys, and tests change. CLI behavior, copied commands, API contracts, brand assets, install URLs, and other product surfaces remain unchanged. The work ships as the next patch release through the established ship flow.

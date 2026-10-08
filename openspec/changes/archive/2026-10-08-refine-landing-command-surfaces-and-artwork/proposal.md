# Change: Refine landing command surfaces and artwork

## Why

The landing's command output and installation command use excessive empty padding and separate the shell selector from the command it controls. The masthead also reads as an opaque strip, while the small Relay mark is horizontally stretched and the route hero presents an isolated mascot instead of the transfer story.

## What Changes

- Group shell selection, generated commands, copy feedback, and installation commands into compact low-contrast surfaces with Copy directly below the command.
- Let the masthead share the page canvas without an opaque background fill.
- Recompose the Relay mark as a compact near-square icon that remains legible at header and favicon sizes.
- Replace the isolated route mascot with one transparent amorphous pixel-art journey connecting local, remote, browser, and delivery destinations.
- Refresh asset provenance, the Relay QA contact sheet, embedded browser assets, and regression coverage.

## Impact

This change affects the shared React command readout, the project landing composition, Relay mark and route artwork, generated embedded assets, and visual documentation. It does not change the CLI contract, transfer behavior, browser APIs, security boundaries, or publication workflow.

# Change: Simplify the landing into an interactive delivery guide

## Why

The current landing repeats the route contract as cards, separates commands from their options, and carries several narrative sections that compete with installation and route selection. Its preference controls still spend header space on visible labels and words, long installation commands are clipped inside narrow cards, and the square arrow mark does not identify Relay as Courier's mascot.

## What Changes

- Replace visible theme and locale text controls with compact icon-only selectors that retain accessible names, radio semantics, persisted state, and complete keyboard operation.
- Replace the abstract arrow mark with a compact Relay mascot mark across the shared brand component, favicon, and light/dark lockups.
- Reduce the landing to the hero, a compact illustrated installation chooser, one contract-backed interactive route explorer, and one combined commands/options reference.
- Derive route availability, endpoint labels, flags, and examples from the generated shipped CLI contract rather than introducing a second capability registry.
- Add an icon link to the Courier GitHub repository in the masthead and remove the hero eyebrow, statistics, safety cards, examples, documentation block, closing call to action, and footer.
- Make every installation command occupy the available card width and remain scrollable or wrappable without page overflow at supported viewports.

## Impact

This change affects the shared UI kit, SVG identity assets, the landing application, asset provenance, English documentation, and browser/unit tests. It does not change CLI behavior, the canonical command contract, release packaging, runtime APIs, authentication, transfer safety, or the README illustration.

# Change: Present the landing as four immersive product sections

## Why

The landing currently behaves like a conventional long page: the hero competes with installation actions, the illustrated panels are boxed away from the canvas, and dense installation and CLI content does not consistently use the available viewport. The masthead also carries unnecessary separators and the generated vector mark does not match the authored Relay artwork closely enough at compact sizes.

## What Changes

- Recompose the landing into four viewport-sized snap sections: product promise, routing, installation, and commands/options.
- Keep masthead navigation on the left in the same order as the three navigable sections, while placing GitHub, theme, and locale actions together on the right without divider rules.
- Remove the hero installation command and actions, the public full-contract label, and all remaining header rules.
- Blend the hero Relay scene into the page canvas and add pointer-, click-, and keyboard-responsive route motion that remains decorative and reduced-motion safe.
- Generate separately composed wide and portrait backgrounds for dispatch, routing, installation, and CLI verification so every scene fills its slide, survives mobile layout, and reserves environmental space for its overlaid controls.
- Scale the existing contract-backed route explorer, installation chooser, and combined command registry to the usable dimensions of their section, and add reduced-motion-safe scene parallax without duplicating functional state.
- Replace locale glyphs with native flag emoji and replace the compact product mark with a dedicated transparent Relay image generated for Courier.

## Impact

This change affects the shared brand and preference components, generated identity provenance, the landing layout and interaction model, English/Russian catalogs, documentation, generated embedded assets, and browser/unit tests. It does not change the CLI contract, transfer behavior, package channels, release automation, runtime APIs, or protected-data handling.

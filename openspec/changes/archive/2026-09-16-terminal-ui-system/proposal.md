# Change: Redesign browser surfaces as terminal workspaces

## Why

Courier's browser surfaces use the shared identity, but their large filled panels compete with the illustrated environments, the landing pointer lens can lose hit testing to its own decorative layer, and command-oriented product behavior is presented as static cards instead of an operational transcript. The landing also retains explanatory labels that duplicate visible state. Courier needs one restrained terminal grammar that remains accessible, makes simulated behavior explicit, and does not introduce a browser shell.

## What Changes

- Add shared terminal surface, transcript, and read-only command-demo primitives for every browser application.
- Track scene pointers on a stable host and replace the elliptical lens with a smooth multi-lobed mask while keeping the base image stationary.
- Replace distorted SVG route circles with code-native circular terminals.
- Recompose landing routing, installation, and CLI reference as contract-backed terminal demonstrations with copy and activation controls but no editable command input or execution.
- Recompose delivery and administration applications as terminal workspaces while preserving their existing guarded APIs and semantic form controls.
- Replace landing, delivery, and administration scenes with purpose-composed responsive cinematic editorial artwork and auditable local provenance.

## Impact

This change affects the shared Lit UI kit, all three browser applications, localized presentation copy, browser tests, documentation, and raster identity assets. It does not alter the public CLI contract, HTTP APIs, transfer semantics, README banner, release process, or credential boundaries.

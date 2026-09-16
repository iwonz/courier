## Why

The landing currently behaves as four isolated viewport slides with runnable browser previews, which fragments Relay's visual journey and spends runtime work on demonstrations that do not execute Courier. Courier needs one continuous, naturally scrolling landing whose immutable command surfaces stay useful while artwork and effects load only when they are likely to be seen.

## What Changes

- Remove visible hero tagline copy and all landing Run/Replay simulations, staged transcripts, timers, and preview messaging while retaining selectors, immutable commands, Copy feedback, route details, options, and release links.
- Replace fixed slide snapping with natural section heights, smooth header-aware anchor navigation, and active-section tracking suited to variable-height content.
- Present GitHub through the same shared control frame as locale and theme controls.
- Replace separate landing scenes with one responsive Relay journey split into continuous wide and portrait panorama segments with masked overlap between sections.
- Load landing scenes once near the viewport, create pointer refraction only on demand, and enforce explicit image and bundle budgets without changing delivery or administration behavior.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-landing`: Replace snapped slides and simulated command demonstrations with continuous natural scrolling, read-only command surfaces, continuous artwork, and bounded loading behavior.
- `ui-kit`: Replace the shared command-demo primitive with a copy-only command readout, add shared icon-link control chrome, and make responsive scene loading and refraction demand-driven.
- `cross-platform-acceptance`: Add observable browser performance, loading-order, natural-scroll, panorama-continuity, and removed-demo acceptance guarantees.

## Impact

This change affects the shared Lit UI kit, the landing application, responsive landing artwork and provenance, English documentation, localized runtime catalogs, unit/browser acceptance, and bundle/asset gates. It does not change the public CLI contract, transfer behavior, HTTP APIs, delivery authentication, administration controls, release flow, or README banner.

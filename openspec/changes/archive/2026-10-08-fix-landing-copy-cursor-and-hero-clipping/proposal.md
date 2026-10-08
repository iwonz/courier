## Why

The Copy action looks interactive but does not expose a pointer cursor, while the enlarged route illustration is washed out by its alpha treatment and vertically clipped at the bottom. These regressions weaken affordance and make the primary visual appear unfinished.

## What Changes

- Give enabled Copy actions an explicit pointer cursor while retaining their existing text-link hover treatment and button semantics.
- Remove the hero illustration's alpha mask and restore its terminal palette.
- Scale wide-screen artwork by available hero height so the complete square remains inside the composition.
- Stack the complete illustration below the headline on narrow screens to preserve both artwork and title readability.
- Update regression coverage, documentation, and generated embedded CSS.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-landing`: Require clear Copy pointer affordance and complete, non-washed-out responsive hero rendering.

## Impact

This change affects the shared command readout styling, landing hero layout, shared CSS, generated embedded browser CSS, landing documentation, and browser coverage. It does not alter copy behavior, command generation, the route-v3 raster, CLI protocols, or publication behavior.

# Change: Refine landing interaction and reference surfaces

## Why

The immersive landing now fills each viewport, but its fixed masthead does not account for mobile safe areas, its dark glass treatment muddies illustrated scenes, and pointer movement shifts whole background images. Route and installation choices also change on hover, variable content moves their internal layouts, brand-like glyphs are not authoritative marks, and the CLI registry cannot relate commands to compatible options.

## What Changes

- Keep the masthead permanently visible while measuring its rendered height, respecting safe areas, lightening its glass treatment, and identifying the active section.
- Replace scene parallax with a stationary, pointer-tracked amorphous spotlight and restrained refractive lens that degrades safely for touch and reduced motion.
- Make the hero non-interactive and describe its transfer only through canonical `Source` and `Destination` roles joined by a cubic Bezier route.
- Make route and installation selection activation-only, add hover feedback without state changes, draw the selected route through a measured Bezier connector, and keep every panel boundary invariant.
- Stretch route, installation, and CLI surfaces across the available shell width.
- Add repository-local, monochrome official brand marks with recorded provenance and no runtime asset requests.
- Make CLI commands selectable and add a checked-by-default compatibility filter driven only by the generated command-to-flag registry.

## Impact

This change affects shared UI primitives, landing layout and interaction state, brand notices, English/Russian landing catalogs, documentation, and unit/browser acceptance. It does not alter the public CLI contract, transfer runtime, release channels, generated raster scenes, or release version.

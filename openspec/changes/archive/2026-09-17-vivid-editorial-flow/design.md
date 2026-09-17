# Design

## Identity direction

Courier uses a kinetic postal-modernist palette: deep ink, clear cobalt, warm cream, safety coral, and restrained mint. A mature aerodynamic swift in a compact courier harness replaces the engraved moth. The mascot is friendly and energetic without becoming childish, toy-like, or ornamental.

The compact mark and hero composition are ImageGen-authored raster assets with local provenance. They contain no text, interface, credential, trademark, or required product information.

## Natural landing flow

The landing retains exactly three semantic sections, but none is sized as a screen. Each section receives only content-driven block padding. The first section combines headline, local hero art, route controls, and command readout in a responsive editorial composition. Installation and CLI follow as ordinary blocks in the same document flow.

There is no full-page image, panorama segmentation, section background image, `100svh` hero, or `78svh`/`88svh` minimum. Local gradients and rounded color fields provide continuity without raster seams. On narrow viewports, content expands normally and the hero illustration moves within the reading order.

## Shared product surfaces

Shared tokens and workbenches use soft radii, clear borders, high-contrast type, compact controls, and lightweight shadows. Delivery and administration inherit the refreshed palette and workbench grammar without changing their APIs or protected-data boundaries. Commands and operational values remain monospace; headings and controls use system sans-serif.

## Performance

The compact mark is loaded wherever product chrome needs it. Only the hero owns a responsive editorial image and it loads eagerly because it appears above the fold. Later landing sections make no artwork request. Stable aspect ratios prevent layout shift, and all assets remain local and budgeted.

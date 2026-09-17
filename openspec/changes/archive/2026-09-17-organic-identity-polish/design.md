# Design

## Continuous composition

The landing owns one absolute responsive panorama below all three semantic sections. A single `<picture>` selects either the desktop or mobile master and covers the complete document flow. Section content does not instantiate its own scene, mask, or background, so no image boundary can become visible during scrolling.

## Readability

Headings sit inside broad feathered contrast fields implemented as pseudo-element gradients. Text color is selected for the artwork rather than inherited from page theme. Interactive surfaces use a bounded translucent fill, generous blur, rounded corners, and a subtle shadow; they do not become opaque rectangular bands.

## Identity mark

The compact Vector mark is an ImageGen-authored transparent raster with rounded organic wings and one restrained orange waypoint. It is optimized locally, recorded in provenance, and reused by landing, delivery, administration, favicons, and wordmarks. No SVG approximation is shipped as the primary mark.

## Performance

Only one responsive panorama source is requested. The selected source loads eagerly because it is the sole full-page visual and must not pop between sections. The mark is small, locally bundled, and decoded asynchronously. Static imagery remains decorative and non-interactive.

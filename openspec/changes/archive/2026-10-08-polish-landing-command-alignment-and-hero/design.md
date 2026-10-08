## Context

The landing consumes a shared `CommandReadout` and one square transparent route-v3 sprite. The command builder nests that readout inside the shell-selection surface, while installation uses it directly. The shared stylesheet serves landing, administration, and delivery builds, and generated embedded CSS must remain synchronized.

## Goals / Non-Goals

**Goals:**

- Keep command text and actions on a stable visual axis across nested and direct readouts.
- Preserve semantic button and accessible-label behavior while making Copy read visually like a link.
- Recompose the existing route-v3 asset through responsive layout and masking without adding raster weight.

**Non-Goals:**

- Regenerating or changing the route-v3 source asset.
- Changing copy, shell-selection, command-building, or route-selection behavior.
- Adding animated or pointer-reactive hero effects.

## Decisions

The shared readout owns one symmetric padding system and renders Copy as a semantic button with transparent normal, hover, and active backgrounds. This preserves keyboard and disabled behavior while avoiding a button-shaped visual. The nested builder readout explicitly removes its own padding so the outer shell surface owns the single inset; installation keeps the shared inset.

The visible shell-syntax label is removed, but the localized string remains the selector's accessible name. This removes redundant copy without weakening assistive semantics.

The route-v3 file remains unchanged. A responsive absolute wrapper scales it beyond the previous fixed maximum, places it behind the headline layer, and clips it within the hero. A soft radial alpha mask plus theme- and viewport-specific opacity reduces visual mass. This is preferable to another generated asset because the requested change is compositional and the current transparent scene already contains the required story.

## Risks / Trade-offs

- [Large pixel artwork can dominate small screens] → Lower its mobile opacity, shift it toward the edge, and keep the headline in a higher stacking layer.
- [Nested padding can regress through utility precedence] → Use an explicit padding override on the nested readout and cover the rendered structure in tests.
- [Shared CSS changes regenerate embedded application CSS] → Rebuild and verify both embedded asset sets even though the hero class is only consumed by the landing.

## Migration Plan

Ship the responsive and shared-component changes together with regenerated embedded CSS. Rollback is a single source revision with no data or protocol migration.

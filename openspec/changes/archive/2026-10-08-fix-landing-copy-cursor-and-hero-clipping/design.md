## Context

The route-v3 source is a square transparent image. The current wide hero sizes its image element from width inside a shorter clipped container and then lowers its opacity behind a radial mask, which cuts the vertical silhouette and washes out the palette. Copy remains a semantic button styled as text but inherits the default cursor.

## Goals / Non-Goals

**Goals:**

- Preserve the complete route-v3 silhouette and intended palette at supported widths.
- Keep the headline readable without relying on a fading image mask.
- Make the Copy affordance explicit without changing its semantics or state behavior.

**Non-Goals:**

- Regenerating route-v3 or changing its encoded bytes.
- Changing command copying, status feedback, or keyboard operation.
- Adding animation or pointer-driven hero effects.

## Decisions

Enabled Copy remains a real button and receives an explicit pointer cursor in the shared readout. This retains disabled, focus, clipboard, and accessibility behavior while matching its link-like visual treatment.

At wide breakpoints, the image element is sized from the hero's available height and keeps automatic width, so its box remains square and cannot exceed the hero vertically. At narrow widths, it returns to document flow below the headline and is width-constrained to the available canvas. This is preferable to permitting overflow because the complete illustration remains visible without increasing page width.

The radial mask is removed. Narrow layouts use the asset at full opacity; wider overlapping layouts use only a slight uniform opacity reduction and restore full opacity on desktop. This keeps the authored terminal colors intact without reintroducing a hard background rectangle.

## Risks / Trade-offs

- [Stacking artwork increases narrow-page height] → Keep the scene width-bound and compact while preserving the more important title and artwork readability.
- [Fuller color can compete with the title at intermediate widths] → Retain a small uniform opacity reduction only where the composition overlaps.
- [Shared CSS changes regenerate embedded application CSS] → Rebuild and verify both embedded browser asset sets.

## Migration Plan

Deploy the styling and layout changes with regenerated embedded CSS. Rollback requires only reverting the source revision; no state or protocol migration exists.

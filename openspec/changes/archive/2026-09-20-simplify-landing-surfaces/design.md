## Context

See `proposal.md` for motivation. The landing is a single React composition backed by shared terminal controls. Its functional state is already covered by exact unit and Chromium acceptance tests; this change is limited to the surrounding presentation hierarchy.

## Goals / Non-Goals

**Goals:**

- Remove decorative information and structural rules without changing document order.
- Make installation selection visually stable and non-redundant.
- Keep command construction, validation, focus visibility, and copying intact while removing registry dividers.

**Non-Goals:**

- Changing the CLI contract, commands, localization catalog, assets, or shared control semantics.
- Removing necessary focus indicators or the native visual boundary of editable controls.
- Reworking route selection, responsive breakpoints, or release channels.

## Decisions

- Remove the metric strip from the DOM instead of hiding it. This avoids shipping redundant accessible content and leaves the hero layout content-driven.
- Remove boundary utilities at each landing composition site rather than weakening the shared border token. Other Courier surfaces still use the token for meaningful data tables and editable controls.
- Render every installation tab with the same base variant and select its calm fill through `aria-selected`. This keeps hover independent from selection and avoids variant replacement flashes.
- Omit installation `details` entirely because the tab already exposes the selected channel. The exact command and footer destinations remain unchanged.
- Remove header, row, and readout divider utilities from the CLI registry while retaining spacing, scroll regions, input boundaries, and two-pixel keyboard focus indicators.

## Risks / Trade-offs

- [Risk] Fewer rules may weaken grouping at narrow widths. → Preserve existing headings, whitespace, responsive columns, and semantic regions, then inspect 320, 390, 1024, and 1440 pixel layouts.
- [Risk] Tab state could become ambiguous in one theme. → Verify active and hover colors in both light and dark themes, with `aria-selected` remaining the semantic source of truth.
- [Risk] Removing the metrics changes hero height. → Keep the existing responsive grid and validate that route controls remain visible without clipping or overflow.

## Migration Plan

Ship as a presentation-only patch release. Rollback is a normal Git revert because no persisted data or public API changes.

## Context

Source/Destination presently owns an immutable sample command while the separate CLI builder owns independently editable arguments. The contract already defines `from`, its arguments, flags, and route pairs; no schema change is needed.

## Decisions

- Render the headline and Relay above a single command workspace inside the first top-level section. Keep `#cli` on the nested workspace for deep links. Installation becomes the second top-level section.
- Initialize the builder to `from` with the initial local-to-SSH route examples. Selecting a route template selects `from`, replaces both endpoint argument examples, and clears option values. Retain the selected route while another command is active; returning to `from` restores its examples. Selecting the already-active command does not deselect it.
- Treat route controls as example templates, not input-type constraints. Manual endpoint edits do not alter the selected template or contract-based option compatibility. Hide route controls for non-`from` commands.
- Keep one contract-built output and Copy action. Remove the immutable route output and applicable-flag badges. The existing builder remains the only place that validates, quotes, and copies values.
- Make command selection a horizontal overflow strip on small screens and a vertical list on wide screens; keep the options region unframed. Keep the binaries link outside the horizontally scrolling tab list, fixed to its right at every viewport.
- Preserve the current GitHub Releases URL but remove the redundant Linux-packages link and icon from the installation readout.

## Risks / Trade-offs

- The first section grows when the builder moves into it. Verify natural document flow and no horizontal overflow at 320, 390, 1024, and 1440 pixels.
- Route template selection replaces typed endpoint values. This is intentional and covered by tests; manual edits remain available immediately afterward.
- Other CLI commands have no route selector. Their required fields and Copy validation remain visible in the same workspace.

## Migration

No persisted state or API migration. Release as v0.3.4 after strict OpenSpec, full verification, package publication, and Pages deployment.

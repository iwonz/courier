## Context

The route and command builder already share one result. Its hero wrapper and Relay cell currently paint `--card`, while the command list and both readouts still print `Command`. The npx tab has no image, although npx is distributed under the npm identity. The direct-binaries link has no functional glyph.

## Decisions

- Remove only the hero wrapper and Relay cell fills. Keep the body's theme-specific canvas and editable controls unchanged.
- Make the shared readout heading optional and omit it for both landing instances. Keep its Copy action, command, and live status intact. Delete the now-unused landing catalog key and command-list header.
- Render the existing local npm PNG to the left of npx text; it is a shared npm identity, not a new npx asset. Keep wget text-only. Place the existing first-party pixel download icon before direct-binaries text, preserving the link's accessible name and URL.
- Give the Brand wrapper and wordmark the mark's explicit height and center alignment. Do not alter the Relay asset or other consumers' wordmark text.

## Risks / Trade-offs

- The download icon reduces space available to the mobile tab scroller. Verify 320 and 390 pixel widths still keep the link visible without document overflow.
- Shared `CommandReadout` consumers may still need a heading. Keep the prop supported and test both render paths.
- Removing light fills makes spacing more important. Inspect light and dark layouts at 320, 390, 1024, and 1440 pixels.

## Migration

No persisted state or API migration. Publish v0.3.5 after strict OpenSpec and the full release gates.

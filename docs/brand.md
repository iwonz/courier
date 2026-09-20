# Courier brand system

Courier is a dependable route for files and directories across local, SSH, browser, and webhook boundaries. The identity is precise, compact, and technically mature without looking institutional.

## Promise

**From here to anywhere.**

Courier makes Source, Destination, policy, progress, and the verified result explicit. It does not present transfer as magic and does not claim guarantees the runtime cannot prove.

## Relay pigeon

Relay is one consistent 8-bit courier pigeon, not a collection of loosely related birds. Its canonical anatomy is compact and nearly square: a cool-neutral head, one large yellow-ringed pigeon eye, short neutral beak, pale folded wings, compact dark tail, yellow feet, one small left-side teal earpiece, and one teal courier satchel. Neutral-v3 is the identity authority for every other shipped sprite.

Across the mark, route, delivery, and administration roles, Relay keeps the reviewed silhouette, physiology, body-to-head proportions, eye and beak geometry, plumage structure, earpiece, satchel, pose, and working object. The v3 family is a controlled palette edit: blue becomes teal, orange/coral becomes yellow, and blue-gray plumage becomes a cool black/white ramp. A palette refresh must not become a redraw. Relay never gains armor, a helmet, visor, glowing face panel, exoskeleton, metallic chest plate, police or military equipment, photorealism, smooth 3D shading, scenery, text, or a watermark.

Use [`courier-relay-pixel-mark-v3.webp`](../web/ui/assets/courier-relay-pixel-mark-v3.webp) in headers and favicons. Use neutral-v3 as the canonical identity authority and route/delivery/admin-v3 only in their declared roles. The [v3 QA contact sheet](assets/courier-relay-pixel-v3-contact-sheet.png) compares the complete family on light, dark, and checkerboard backgrounds. Product meaning remains complete when artwork is unavailable; a sprite never replaces a label, status, or security warning. Every sprite is transparent, square WebP rendered with pixelated sampling.

## Pixel visual language

The shared terminal palette uses these semantic colors and opacity-derived roles:

| Token | Value | Role |
|---|---|---|
| Black | `#000000` | Dark canvas and light-theme text |
| Surface 1 | `#0D1015` | Dark primary panel |
| Surface 2 | `#191C20` | Dark secondary panel |
| White | `#FFFFFF` | Light canvas and dark-theme text |
| Action teal | `#71FFF6` | Action/success fill and dark-theme text accent |
| Accessible teal | `#006B67` | Light-theme teal text |
| Selection yellow | `#FAD14F` | Selection, warning, and route signal |
| Destructive red | `#C94A55` | Failure and destructive action |

Dark mode uses white at 70% for secondary text and 55% for rules. Light mode inverts those roles with 70% and 55% black and uses 4%/9% black for its two surfaces. Bright teal and yellow fills always use black text. Controls use square corners, one-pixel semantic rules, two-pixel focus, and instant step states without chamfers, offset shadows, blur, or soft shadows. Flat surface bands and ruled rows provide hierarchy without decorative card islands.

Pixelify Sans is pinned locally at upstream commit `39df74aba80df8157546034b878e8be1eb565ced` and is used only for the wordmark and H1/H2. Overpass Mono 400/600 Latin and Cyrillic subsets come from `@fontsource/overpass-mono@5.3.0` with upstream revision `c580d28bfab7f39013568e684a65eeb23eff588d`; they serve interface copy, labels, commands, paths, identifiers, and counters. Both families ship with OFL 1.1 notices, and the four Overpass subsets total no more than 48 KiB.

First-party actions and endpoint concepts use shared code-native pixel-grid icons. Package-manager, platform, distribution, and GitHub marks are normal transparent PNG images in official geometry and color; they never use pixelated rendering. GitHub has explicit light and dark files. npx and wget have no invented mark and render as text. Locale buttons use local 16×12 pixel flags rather than platform emoji. Third-party attribution and trademark notices live in [`web/ui/NOTICE.md`](../web/ui/NOTICE.md).

## Interaction

Theme and locale are each one icon button. Theme cycles `system → light → dark`; locale cycles `English → Russian` and visibly shows the active local pixel flag. System and browser language remain defaults until an explicit stored choice overrides them. Every control exposes its current and next value in a localized accessible name.

Route geometry starts from measured Source and Destination bounds and renders one cubic Bézier connector. Its four-by-four-pixel signal uses a centered offset anchor and linear motion on that exact path. Reduced motion fixes the signal at the midpoint without hiding route state.

## Voice

1. Start with the current state or required action.
2. Name the object: Source, Destination, route, archive, server, or delivery.
3. Use concrete verbs: prepare, connect, send, receive, verify, stop, retry.
4. Report only confirmed outcomes.
5. State what remains safe after failure when known.
6. Never include credentials, query secrets, or protected metadata.

Prefer “route,” “delivery,” “verifying,” “stopped,” “collision,” and “confirmed bytes.” Avoid magic-infrastructure language, fake urgency, blame, jokes in failure states, or unsupported superlatives.

## Generated asset record

The v3 family was produced with built-in ImageGen `precise-object-edit` operations for each reviewed neutral/mark/route/delivery/admin source. The palette pass retained identity, anatomy, pose, crop, objects, and alpha while remapping the palette to teal/yellow/cool-neutral. A second controlled shape pass replaced rounded head, body, wing, tail, bag, and role-object contours with flat stepped planes and square right-angle pixel clusters without changing the role pose or functional details. Accepted PNG outputs were reviewed, nearest-neighbor resized to the established 256/384/512 canvases, and encoded as transparent WebP. The generated originals remain in the ImageGen output store; superseded v1/v2 WebPs are not shipped.

[`provenance.json`](../web/ui/assets/provenance.json) records identity invariants, controlled palette and square-silhouette variation, edit sources, prompts, ImageGen lineage, consumers, role, dimensions, byte count, SHA-256, contact sheet, font revisions, and licenses. The asset gate verifies exactly five declared v3 sprites, alpha transparency, square dimensions, the square-body identity rule, the 24/64/64/48/48 KiB per-role limits, the 248 KiB combined limit, both local font families and OFL notices, and the absence of undeclared or legacy role assets.

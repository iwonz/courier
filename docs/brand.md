# Courier brand system

Courier is a dependable route for files and directories across local, SSH, browser, and webhook boundaries. The identity is precise, compact, and technically mature without looking institutional.

## Promise

**From here to anywhere.**

Courier makes Source, Destination, policy, progress, and the verified result explicit. It does not present transfer as magic and does not claim guarantees the runtime cannot prove.

## Relay pigeon

Relay is one consistent 8-bit courier pigeon, not a collection of loosely related birds. Its canonical anatomy is compact and nearly square: a blue-gray head, one large orange-ringed pigeon eye, short ivory beak, pale folded wings, compact dark tail, coral feet, one small left-side earpiece, and one cobalt courier satchel. The neutral mascot is the identity reference for every other sprite.

Across the mark, route, delivery, and administration roles, Relay keeps the same silhouette, physiology, body-to-head proportions, eye and beak geometry, plumage, earpiece, and satchel. Only pose, role equipment or clothing, and a carried or attached object may change. Every new role starts as an identity-preserving derivative of the canonical neutral sprite; an independent mascot redraw is not an acceptable variation. Relay never gains armor, a helmet, visor, glowing face panel, exoskeleton, metallic chest plate, police or military equipment, photorealism, smooth 3D shading, scenery, text, or a watermark.

Use [`courier-relay-pixel-mark-v1.webp`](../web/ui/assets/courier-relay-pixel-mark-v1.webp) in headers and favicons. Use the neutral, route, delivery, and administration sprites only in their declared roles. Product meaning remains complete when artwork is unavailable; a sprite never replaces a label, status, or security warning. Every sprite is transparent, square, lossless WebP rendered with pixelated sampling.

## Pixel visual language

The shared palette contains exactly these twelve brand and semantic colors:

| Token | Value | Role |
|---|---|---|
| Ink | `#0B1020` | Dark canvas and primary text |
| Navy | `#17233B` | Dark secondary surface and hard shadow |
| Slate | `#3F506B` | Muted dark text and dividers |
| Steel | `#7C8DA5` | Quiet lines and disabled structure |
| Ice | `#DCECF7` | Light secondary surface |
| Paper | `#F8FBFF` | Light canvas and selected text |
| Cobalt | `#2556C7` | Selection and product action |
| Cyan | `#35B6D4` | Route signal and accent |
| Rust | `#A63F14` | Light-theme focus and warning |
| Orange | `#F07A32` | Dark-theme focus and waypoint |
| Green | `#23845D` | Success |
| Red | `#C83B4E` | Failure and destructive action |

Light and dark themes remap those tokens without introducing decorative colors. The canvas uses a restrained hard-edged eight-pixel dither. Controls follow a four-pixel unit with four-pixel chamfers, one-pixel semantic dividers, two-pixel focus, and hard two-to-four-pixel interaction shadows. Ordinary layout remains transparent: pixel frames belong to controls, fields, alerts, status, focus, and boundaries that explain independent behavior—not to decorative card islands.

Pixelify Sans is pinned locally at upstream commit `39df74aba80df8157546034b878e8be1eb565ced` and is used only for the wordmark, H1/H2 headings, and compact display labels. Its local WOFF2 contains Cyrillic and ships with the OFL 1.1 notice. Body copy stays system sans; commands, paths, identifiers, and counters stay system monospace.

First-party actions and endpoint concepts use shared code-native pixel-grid icons. Package-manager, platform, and GitHub marks are reviewed monochrome pixel-grid derivatives of pinned official geometry. Locale buttons use local 16×12 pixel flags rather than platform emoji. Third-party attribution and trademark notices live in [`web/ui/NOTICE.md`](../web/ui/NOTICE.md).

## Interaction

Theme and locale are each one icon button. Theme cycles `system → light → dark`; locale cycles `English → Russian` and visibly shows the active local pixel flag. System and browser language remain defaults until an explicit stored choice overrides them. Every control exposes its current and next value in a localized accessible name.

Route geometry starts from measured cubic Bézier coordinates, samples the curve, snaps samples to the four-pixel grid, and renders a crisp stepped path. A decorative packet may move with `steps()` timing. Reduced motion removes travel and repeated decorative animation without hiding route state.

## Voice

1. Start with the current state or required action.
2. Name the object: Source, Destination, route, archive, server, or delivery.
3. Use concrete verbs: prepare, connect, send, receive, verify, stop, retry.
4. Report only confirmed outcomes.
5. State what remains safe after failure when known.
6. Never include credentials, query secrets, or protected metadata.

Prefer “route,” “delivery,” “verifying,” “stopped,” “collision,” and “confirmed bytes.” Avoid magic-infrastructure language, fake urgency, blame, jokes in failure states, or unsupported superlatives.

## Generated asset record

The canonical neutral Relay was generated first with the built-in ImageGen workflow. The mark and role sprites were generated or cleaned with the neutral asset as their explicit identity reference. Accepted PNG sources were reviewed, nearest-neighbor resized, reduced to a shared hard 32-color raster treatment, and encoded as lossless transparent WebP. Temporary generations are not shipped.

[`provenance.json`](../web/ui/assets/provenance.json) records the identity invariants, allowed variations, prompt, reference lineage, role, dimensions, byte count, SHA-256, font revision, and license. The asset gate verifies exactly five declared Relay sprites, alpha transparency, square dimensions, the 24/64/64/48/48 KiB per-role limits, the 248 KiB combined limit, the pinned local font and OFL notice, and the absence of undeclared or legacy armored assets.

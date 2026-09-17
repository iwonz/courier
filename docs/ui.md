# Browser UI architecture

Courier's landing, delivery, and administration applications are independent React roots that share the repository-owned shadcn package in `web/ui`. The package owns Tailwind tokens, identity, localization, preferences, icons, route geometry, semantic controls, and immutable command readouts. Radix supplies behavior for checkbox, select, tabs, tooltip, progress, separator, and scroll-area primitives; Courier owns their local React source and styling.

## Shared primitives

- `Brand` renders the 256×256 pixel Relay mark with a Pixelify Sans `COURIER CLI` wordmark. `RelaySprite` accepts the typed `neutral`, `route`, `delivery`, or `admin` role and resolves only that role's local asset.
- shadcn `Button`, `Card`, `Badge`, `Checkbox`, `Input`, `Select`, `Tabs`, `Tooltip`, `Progress`, `Separator`, `ScrollArea`, and `Alert` provide one accessible control grammar. `Card` is a transparent structural wrapper, not a visual panel.
- `CommandReadout` presents immutable command text, Copy, reserved localized feedback, details, and footer actions directly on the document canvas.
- `ThemeSelector` and `LocaleSelector` are single cyclic icon buttons backed by one shared React preference provider and browser controller; locale visibly identifies its active language.
- `RouteDisplay`, `PixelIcon`, and `PixelBrandIcon` retain native semantics and render reviewed grid data with `currentColor` and crisp edges. Compatibility exports keep existing consumers source-compatible.

There is no browser shell, free command input, fake terminal prompt, pointer refraction layer, or arbitrary command execution.

## Themes and localization

English is the fallback locale and Russian is the second runtime catalog. The theme default is `system`; the locale default follows browser languages. Explicit choices persist locally. The shared controller owns document theme updates, storage failure handling, system-theme listeners, and synchronized selector state.

Product copy and source documentation remain English. `Source` and `Destination` are protocol-role terms and remain literal in both runtime locales. Protected delivery metadata is never rendered before successful authorization.

## Surface composition

The landing combines its headline and route instrument, then presents installation and the contract-generated command registry. Delivery is a destination workbench around a real protected password form and authorized manifest. Administration is an operations workbench with counters, a keyboard-operable delivery navigator, and a selected policy inspector whose UUID selection survives SSE snapshots while valid.

All applications use the same fixed twelve-color light/dark mapping, eight-pixel dither canvas, four-pixel geometry unit, chamfered controls, crisp lines, cobalt selection, and hard interaction shadows. The document background is the only page-scale surface: ordinary wrappers add no card fill, outer radius, shadow, border, or blur. Spacing, type, alignment, and occasional hairline boundaries between independently scrolling regions provide hierarchy. Inputs, alerts, focus, selected controls, status, and destructive actions retain explicit chrome because it communicates state. Pixelify Sans is limited to brand and H1/H2 display type; system sans serves body copy and system monospace serves operational values. Layouts expand on narrow screens, preserve semantic actions, and avoid horizontal overflow.

## Assets and loading

`assets.ts` exports the 256×256 mark, the typed `RelayRole`, and one resolver for the 512×512 neutral/route and 384×384 delivery/admin sprites. The neutral sprite is the immutable identity reference: every role is derived from it, keeps its anatomy, proportions, plumage, eye, beak, earpiece, and satchel, and varies only pose, equipment or clothing, and carried or attached objects. Independent mascot redraws are rejected by the family provenance gate. The local Pixelify Sans WOFF2 and OFL notice are part of the audited asset set. No responsive scene family, generated background, panorama, pointer duplicate, emoji flag, Lucide import, or runtime asset request is shipped.

## Verification

Vitest with React Testing Library enforces exactly 100% statements, branches, functions, and lines for first-party TypeScript and TSX. Asset validation checks role resolution, identity-reference lineage, alpha, square geometry, hashes, licenses, and budgets and rejects undeclared assets, legacy Relay names, emoji flags, and browser Lucide imports. Playwright covers localization, theme/locale cycling and persistence, keyboard operation, route selection, Copy feedback, per-surface role sprites, pixel rendering, security isolation, API controls, and overflow across supported viewports.

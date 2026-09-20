# Browser UI architecture

Courier's landing, delivery, and administration applications are independent React roots that share the repository-owned shadcn package in `web/ui`. The package owns Tailwind tokens, identity, localization, preferences, icons, route geometry, semantic controls, and immutable command readouts. Radix supplies behavior for checkbox, select, tabs, tooltip, progress, separator, and scroll-area primitives; Courier owns their local React source and styling.

## Shared primitives

- `Brand` renders the 256×256 pixel Relay mark with a Pixelify Sans `COURIER CLI` wordmark. `RelaySprite` accepts the typed `neutral`, `route`, `delivery`, or `admin` role and resolves only that role's local asset.
- shadcn `Button`, `Card`, `Badge`, `Checkbox`, `Input`, `Select`, `Tabs`, `Tooltip`, `Progress`, `Separator`, `ScrollArea`, and `Alert` provide one accessible control grammar. `Card` is a transparent structural wrapper, not a visual panel.
- `CommandReadout` presents immutable command text, Copy, reserved localized feedback, details, and footer actions directly on the document canvas.
- `ThemeSelector` and `LocaleSelector` are single cyclic icon buttons backed by one shared React preference provider and browser controller; locale visibly identifies its active language.
- `RouteDisplay` and `PixelIcon` render first-party functional grid data with `currentColor` and crisp edges. `BrandIcon` renders official local transparent PNG geometry and color as a normal image; compatibility exports keep existing consumers source-compatible.

There is no browser shell, free command input, fake terminal prompt, pointer refraction layer, or arbitrary command execution.

## Themes and localization

English is the fallback locale and Russian is the second runtime catalog. The theme default is `system`; the locale default follows browser languages. Explicit choices persist locally. The shared controller owns document theme updates, storage failure handling, system-theme listeners, and synchronized selector state.

Product copy and source documentation remain English. `Source` and `Destination` are protocol-role terms and remain literal in both runtime locales. Protected delivery metadata is never rendered before successful authorization.

## Surface composition

The landing is a route console followed by a calm installation chooser and an unruled command builder. Its sections have no boundary rules or decorative counters, installation tabs do not animate or repeat the selected channel below the command, and editable controls retain their own boundaries and focus treatment. Delivery is a transfer console around a real protected password form and authorized ruled manifest. Administration is an operations console with API-derived counters, a keyboard-operable delivery navigator, and a selected policy inspector whose UUID selection survives SSE snapshots while valid.

All applications use the same terminal token system: black/white canvas, `#0D1015`/`#191C20` dark surfaces, opacity-derived light surfaces and rules, teal action/success, yellow selection/warning, and red destructive state. Square controls use one-pixel rules, two-pixel focus, and instant step states without chamfer, offset shadow, blur, or soft shadow. Flat bands and row rules provide hierarchy. Pixelify Sans is limited to the wordmark and H1/H2; locally bundled Overpass Mono 400/600 serves all interface copy and operational values. Layouts expand on narrow screens, preserve semantic actions, and avoid horizontal overflow.

## Assets and loading

`assets.ts` exports the 256×256 mark-v3, the typed `RelayRole`, and one resolver for neutral-v3/route-v3 at 512×512 and delivery/admin-v3 at 384×384. Neutral-v3 is the shipped identity authority; the v3 family retains reviewed anatomy, poses, and objects while using the terminal teal/yellow/cool-neutral palette. Pixelify Sans, four Overpass Mono subsets, both OFL notices, local official brand PNGs, provenance, and the v3 contact sheet are audited. No responsive scene family, generated background, panorama, pointer duplicate, emoji flag, Lucide import, or remote runtime asset is shipped.

## Verification

Vitest with React Testing Library enforces exactly 100% statements, branches, functions, and lines for first-party TypeScript and TSX. Asset validation checks role resolution, identity-reference lineage, alpha, square geometry, hashes, licenses, and budgets and rejects undeclared assets, legacy Relay names, emoji flags, and browser Lucide imports. Playwright covers localization, theme/locale cycling and persistence, keyboard operation, route selection, Copy feedback, per-surface role sprites, pixel rendering, security isolation, API controls, and overflow across supported viewports.

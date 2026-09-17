# Browser UI architecture

Courier's landing, delivery, and administration applications are independent React roots that share the repository-owned shadcn package in `web/ui`. The package owns Tailwind tokens, identity, localization, preferences, icons, route geometry, semantic controls, and immutable command readouts. Radix supplies behavior for checkbox, select, tabs, tooltip, progress, separator, and scroll-area primitives; Courier owns their local React source and styling.

## Shared primitives

- `Brand` renders the dedicated square Relay product mark with `COURIER CLI`; `Mascot` renders the full local ImageGen-authored Relay illustration.
- shadcn `Button`, `Card`, `Badge`, `Checkbox`, `Input`, `Select`, `Tabs`, `Tooltip`, `Progress`, `Separator`, `ScrollArea`, and `Alert` provide one accessible control grammar. `Card` is a transparent structural wrapper, not a visual panel.
- `CommandReadout` presents immutable command text, Copy, reserved localized feedback, details, and footer actions directly on the document canvas.
- `ThemeSelector` and `LocaleSelector` are single cyclic icon buttons backed by one shared React preference provider and browser controller; locale visibly identifies its active language.
- `RouteDisplay`, `Icon`, and `BrandIcon` retain native semantics and bundle code-native or pinned official geometry locally.

There is no browser shell, free command input, fake terminal prompt, pointer refraction layer, or arbitrary command execution.

## Themes and localization

English is the fallback locale and Russian is the second runtime catalog. The theme default is `system`; the locale default follows browser languages. Explicit choices persist locally. The shared controller owns document theme updates, storage failure handling, system-theme listeners, and synchronized selector state.

Product copy and source documentation remain English. `Source` and `Destination` are protocol-role terms and remain literal in both runtime locales. Protected delivery metadata is never rendered before successful authorization.

## Surface composition

The landing combines its headline and route instrument, then presents installation and the contract-generated command registry. Delivery is a destination workbench around a real protected password form and authorized manifest. Administration is an operations workbench with counters, a keyboard-operable delivery navigator, and a selected policy inspector whose UUID selection survives SSE snapshots while valid.

All applications use the same ink/cobalt/off-white token system with restrained cyan routing light and orange waypoint state. The document background is the only page-scale surface: ordinary wrappers add no card fill, outer radius, shadow, border, or blur. Spacing, type, alignment, and occasional hairline boundaries between independently scrolling regions provide hierarchy. Inputs, alerts, focus, selected controls, status, and destructive actions retain explicit chrome because it communicates state. Monospace is limited to operational values. Layouts expand on narrow screens, preserve semantic actions, and avoid horizontal overflow.

## Assets and loading

`assets.ts` exports the transparent 768×768 full Relay mascot and the separate transparent 512×512 compact Relay mark. The full mascot is editorial; the mark is used by product chrome and favicons. No responsive scene family, generated background, panorama, or pointer-effect duplicate is shipped. All artwork, system fonts, icons, styles, and scripts are local.

## Verification

Vitest with React Testing Library enforces exactly 100% statements, branches, functions, and lines for first-party TypeScript and TSX. Asset validation checks dimensions, square geometry, hash, lineage, licenses, and budget. Playwright covers localization, theme/locale cycling and persistence, keyboard operation, route selection, Copy feedback, compact mascot geometry, security isolation, API controls, and overflow across supported viewports.

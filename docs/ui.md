# Browser UI architecture

Courier's landing, delivery, and administration applications are independent React roots that share the repository-owned shadcn package in `web/ui`. The package owns Tailwind tokens, identity, localization, preferences, icons, route geometry, semantic controls, and immutable command readouts. Radix supplies behavior for checkbox, select, tabs, tooltip, progress, separator, and scroll-area primitives; Courier owns their local React source and styling.

## Shared primitives

- `Brand` and `Mascot` render the local ImageGen-authored square Relay pigeon and Courier wordmark.
- shadcn `Button`, `Card`, `Badge`, `Checkbox`, `Input`, `Select`, `Tabs`, `Tooltip`, `Progress`, `Separator`, `ScrollArea`, and `Alert` provide one accessible control grammar.
- `CommandReadout` presents immutable command text, Copy, reserved localized feedback, details, and footer actions.
- `ThemeSelector` and `LocaleSelector` are single cyclic icon buttons backed by one shared React preference provider and browser controller.
- `RouteDisplay`, `Icon`, and `BrandIcon` retain native semantics and bundle code-native or pinned official geometry locally.

There is no browser shell, free command input, fake terminal prompt, pointer refraction layer, or arbitrary command execution.

## Themes and localization

English is the fallback locale and Russian is the second runtime catalog. The theme default is `system`; the locale default follows browser languages. Explicit choices persist locally. The shared controller owns document theme updates, storage failure handling, system-theme listeners, and synchronized selector state.

Product copy and source documentation remain English. `Source` and `Destination` are protocol-role terms and remain literal in both runtime locales. Protected delivery metadata is never rendered before successful authorization.

## Surface composition

The landing combines its headline and route instrument, then presents installation and the contract-generated command registry. Delivery is a destination workbench around a real protected password form and authorized manifest. Administration is an operations workbench with counters, a keyboard-operable delivery navigator, and a selected policy inspector whose UUID selection survives SSE snapshots while valid.

All applications use the same ink/cobalt/off-white token system with restrained cyan routing light and orange waypoint state. Shared shadcn surfaces use generous rounded boundaries, high-contrast translucent layers, and lightweight shadows. Monospace is limited to operational values. Layouts expand on narrow screens, preserve semantic actions, and avoid horizontal overflow.

## Assets and loading

`assets.ts` exports one transparent 768×768 WebP. The same compact Relay asset is reused for product chrome and optional above-fold decoration; no responsive scene family, generated background, panorama, or pointer-effect duplicate is shipped. All artwork, system fonts, icons, styles, and scripts are local.

## Verification

Vitest with React Testing Library enforces exactly 100% statements, branches, functions, and lines for first-party TypeScript and TSX. Asset validation checks dimensions, square geometry, hash, lineage, licenses, and budget. Playwright covers localization, theme/locale cycling and persistence, keyboard operation, route selection, Copy feedback, compact mascot geometry, security isolation, API controls, and overflow across supported viewports.

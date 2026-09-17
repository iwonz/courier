# Browser UI architecture

Courier's landing, delivery, and administration applications share the Lit package in `web/ui`. The package owns tokens, identity, localization, preferences, icons, route geometry, static scenes, form styling, workbenches, and immutable command readouts.

## Shared primitives

- `courier-brand` renders the Vector mark and Courier wordmark without a product descriptor.
- `courier-scene` performs responsive eager or one-shot proximity loading and static compositing.
- `courier-workbench` provides a neutral structured surface with heading, body, actions, status, and footer slots.
- `courier-command-readout` presents immutable command text, Copy, reserved localized feedback, details, and footer actions.
- `courier-theme-selector` and `courier-locale-selector` are single cyclic icon buttons backed by one shared browser preference controller.
- `courier-icon-link`, form controls, status, progress, route, checkbox, and brand-icon components retain native semantics and normalized Courier chrome.

There is no browser shell, free command input, fake terminal prompt, pointer refraction layer, or arbitrary command execution.

## Themes and localization

English is the fallback locale and Russian is the second runtime catalog. The theme default is `system`; the locale default follows browser languages. Explicit choices persist locally. The shared controller owns document theme updates, storage failure handling, system-theme listeners, and synchronized selector state.

Product copy and source documentation remain English. `Source` and `Destination` are protocol-role terms and remain literal in both runtime locales. Protected delivery metadata is never rendered before successful authorization.

## Surface composition

The landing combines its headline and route instrument, then presents installation and the contract-generated command registry. Delivery is a destination workbench around a real protected password form and authorized manifest. Administration is an operations workbench with counters, a keyboard-operable delivery navigator, and a selected policy inspector whose UUID selection survives SSE snapshots while valid.

All applications use the same carbon/chalk/warm-line/orange token system. Monospace is limited to operational values. Layouts expand on narrow screens, preserve semantic actions, and avoid horizontal overflow.

## Assets and loading

Surface modules export role-based scene URLs:

- `landing-scenes.ts`;
- `delivery-scenes.ts`;
- `admin-scenes.ts`;
- `identity-assets.ts`.

The first visible scene is eager. Below-fold scenes create no image request until the proximity boundary; browsers without `IntersectionObserver` fall back to native lazy loading. Every observer is disconnected on activation or element teardown. All assets, fonts, icons, and scripts are local.

## Verification

Vitest enforces exactly 100% statements, branches, functions, and lines for first-party TypeScript. Asset validation checks dimensions, hashes, lineage, responsive sequence, licenses, and budgets. Playwright covers localization, theme/locale cycling and persistence, keyboard operation, route selection, Copy feedback, responsive artwork, lazy loading, security isolation, API controls, and overflow across supported viewports.

# Change: Add the shared Courier UI kit

## Why

Browser delivery pages, the administration interface, and the project landing page need one accessible visual and localization foundation. Copying reference-page markup or CSS into each application would create divergent themes, inaccessible controls, and duplicated identity assets.

## What Changes

- Extract only approved Courier identity assets and attribution notices from the supplied local references.
- Add a reusable Lit and TypeScript UI package built with Vite.
- Define shared graphite and signal-lime design tokens, typography, icons, layout primitives, and controls.
- Add persistent `system`, `light`, and `dark` theme selection with reduced-motion and accessible-contrast behavior.
- Add typed English and Russian catalogs with English fallback and browser-language negotiation.
- Add deterministic asset provenance checks, build output checks, and exact 100% TypeScript coverage to `make verify`.
- Delete the untracked reference directory after approved assets and notices are verified.

## Impact

This change adds reusable presentation infrastructure only. It does not expose an HTTP server, a delivery route, an administration API, or a public CLI command.

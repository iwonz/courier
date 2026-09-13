# Courier UI architecture

Courier browser surfaces use one private workspace package, [`@courier/ui`](../web/ui), for visual tokens, identity assets, icons, localization, theme state, and reusable Lit components. Delivery pages, administration pages, and the static project landing import this package; they must not copy its component implementations or maintain parallel token sets.

## Package boundary

The public TypeScript entry point exports:

- `defineCourierElements` for idempotent custom-element registration;
- button, panel, progress, icon, theme-selector, and locale-selector components;
- typed theme state and browser adapters;
- typed English and Russian catalogs, browser-language negotiation, and translation helpers.

The Vite library build emits an ES module, declarations, and `courier-ui.css`. Consumers load that stylesheet once so document-level tokens and reduced-motion behavior apply consistently; component shadow styles inherit the same custom properties.

## Themes and accessibility

The stored preference is one of `system`, `light`, or `dark` and defaults to `system`. The resolved light or dark theme is written separately to the document root, so operating-system changes update a system-selected page without losing the user's preference. Storage failures fall back safely and never prevent a page from rendering.

Graphite, paper, signal lime, and beak orange are the core identity colors. Controls use native keyboard semantics, visible focus treatment, explicit labels, accessible light/dark contrast, and a reduced-motion media query.

## Localization

English is the fallback locale. Initial selection checks persisted preference and then the browser language list. Catalog shape derives from the English keys, so TypeScript rejects missing keys in Russian or future locale modules. Runtime pages can perform their initial negotiation independently; adding a locale requires no backend API change.

## Assets and provenance

Retained Courier-owned assets live in [`web/ui/assets`](../web/ui/assets) beside `provenance.json`. The manifest records the original local identity reference digest, every asset's media type, byte count, SHA-256 digest, and intended use. The verification script rejects missing, extra, changed, duplicate, or executable asset content.

No reference HTML or JavaScript is shipped. Browser-injected AdGuard resources, obsolete mirror artwork and messaging, and remote resources were excluded. [`NOTICE.md`](../web/ui/NOTICE.md) records the asset license context.

## Verification

Run the UI gate directly with:

```sh
npm ci --prefix web --ignore-scripts --no-audit --no-fund
npm run verify --prefix web
```

The gate checks asset integrity, requires 100% TypeScript statements, branches, functions, and lines through Vitest/V8, and produces a deterministic Vite library build. `make verify` includes this gate before the GoReleaser snapshot.

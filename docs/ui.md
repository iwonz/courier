# Courier UI architecture

Courier browser surfaces use one private workspace package, [`@courier/ui`](../web/ui), for visual tokens, identity assets, icons, localization, theme state, and reusable Lit components. Delivery pages, administration pages, and the static project landing import this package; they must not copy its component implementations or maintain parallel token sets.

The visual and verbal rules are defined in the [Courier brand system](brand.md). Its central idea is a field manual connected to a network control room: clear routes, calm status language, disciplined information density, and a small amount of warmth from Relay, the courier-pigeon field operator.

## Package boundary

The public TypeScript entry point exports:

- `defineCourierElements` for idempotent custom-element registration;
- brand, responsive contextual mascot and stationary scene, status, route, brand-icon, checkbox, button, panel, progress, icon, segmented-control, theme-selector, and locale-selector components;
- deterministic cubic Bezier route geometry shared by route presentations;
- shared form-control styles for normalized inputs, selects, switches, and file actions;
- typed theme state and browser adapters;
- typed English and Russian catalogs, browser-language negotiation, and translation helpers.

The Vite library build emits an ES module, declarations, and `courier-ui.css`. Consumers load that stylesheet once so document-level tokens and reduced-motion behavior apply consistently; component shadow styles inherit the same custom properties.

## Themes and accessibility

The stored preference is one of `system`, `light`, or `dark` and defaults to `system`. The resolved light or dark theme is written separately to the document root, so operating-system changes update a system-selected page without losing the user's preference. Storage failures fall back safely and never prevent a page from rendering.

Graphite, paper, signal lime, and beak orange are the core identity colors. Signal lime identifies a route, active state, or decisive action rather than covering large surfaces. The shared route and status components pair color with text. Theme and locale use icon-only shared segmented radiogroups: visible `Theme` and `Language` legends are suppressed, while localized button and group names, selected state, one roving tab stop, and arrow/Home/End keyboard behavior remain. Inputs, selects, switches, password fields, file actions, and the shared custom checkbox retain dependable form semantics but receive one normalized Courier appearance, explicit programmatic labels, visible focus treatment, accessible light/dark contrast, and reduced-motion behavior. The checkbox preserves native checked/disabled behavior and emits one composed boolean change event without exposing native chrome.

## Localization

English is the fallback locale. Initial selection checks persisted preference and then the browser language list. Catalog shape derives from the English keys, so TypeScript rejects missing keys in Russian or future locale modules. Runtime pages can perform their initial negotiation independently; adding a locale requires no backend API change.

## Assets and provenance

Retained Courier-owned assets live in [`web/ui/assets`](../web/ui/assets) beside `provenance.json`. The manifest records the source context, every asset's media type, role, dimensions, byte count, SHA-256 digest, and prompt summary. The verification script reads PNG and WebP headers and rejects missing, extra, changed, duplicate, dimension-mismatched, incomplete, or executable asset content. The compact transparent PNG mark is a generated profile of Relay, not the former source-arrow symbol. The canonical raster Relay reference anchors a contextual family for dispatch, installation, routing, verification, access control, administration operations, and the README panorama. Landing artwork reserves low-detail bays matching the live controls, concentrates character and props around their edges, and contains no embedded text, logos, third-party marks, or fake UI. `courier-mascot` accepts an optional portrait source and renders it through `<picture>` below the shared mobile breakpoint, allowing a scene to change composition without duplicating component logic. `courier-scene` composes that responsive picture as an immutable base with an optional pointer-local glow and refractive duplicate. Fine-pointer coordinates converge through elapsed-time interpolation and return gently to an ambient position; the loop stops after convergence or disconnection. Coarse pointers and reduced motion use stable fallbacks. All illustrations are non-essential: headings, controls, routes, and states remain complete if an image cannot load.

The shared `courier-brand-icon` renders reviewed monochrome product geometry from the pinned Simple Icons package or pinned official PowerShell and Scoop repository revisions. Those third-party files have a separate provenance manifest, digest gate, and license/trademark notice; they are not described as Courier-owned assets. The component makes no runtime network request and only supplies an accessible image name when a consumer requests one.

No reference HTML or JavaScript is shipped. Browser-injected AdGuard resources, obsolete mirror artwork and messaging, and remote resources were excluded. [`NOTICE.md`](../web/ui/NOTICE.md) records both Courier asset authorship and all third-party mark license, attribution, and trademark context.

## Verification

Run the UI gate directly with:

```sh
npm ci --prefix web --ignore-scripts --no-audit --no-fund
npm run verify --prefix web
```

The gate checks asset integrity, requires 100% TypeScript statements, branches, functions, and lines through Vitest/V8, and produces a deterministic Vite library build. `make verify` includes this gate before the GoReleaser snapshot.

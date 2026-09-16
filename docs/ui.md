# Courier UI architecture

Courier browser surfaces use one private workspace package, [`@courier/ui`](../web/ui), for visual tokens, identity assets, icons, localization, theme state, and reusable Lit components. Delivery pages, administration pages, and the static project landing import this package; they must not copy its component implementations or maintain parallel token sets.

The visual and verbal rules are defined in the [Courier brand system](brand.md). Its central idea is a field manual connected to a network control room: clear routes, calm status language, disciplined information density, and a small amount of warmth from Relay, the courier-pigeon field operator.

## Package boundary

The public TypeScript entry point exports:

- `defineCourierElements` for idempotent custom-element registration;
- brand, responsive contextual mascot and stationary scene, status, route, terminal, command-demo, brand-icon, checkbox, button, panel, progress, icon, segmented-control, theme-selector, and locale-selector components;
- typed terminal transcript rows, tones, and deterministic demo phases;
- deterministic cubic Bezier route geometry shared by route presentations;
- shared form-control styles for normalized inputs, selects, switches, and file actions;
- typed theme state and browser adapters;
- typed English and Russian catalogs, browser-language negotiation, and translation helpers.

The Vite library build emits an ES module, declarations, and `courier-ui.css`. Consumers load that stylesheet once so document-level tokens and reduced-motion behavior apply consistently; component shadow styles inherit the same custom properties.

## Terminal boundary

`courier-terminal` is an accessible transparent workspace with toolbar, transcript body, footer actions, and status slots. It supplies the hairline registration frame and scrolling boundary while consumers retain semantic buttons, links, forms, password fields, file inputs, selects, and checkboxes. `courier-command-demo` adds an immutable command, Copy and Run/Replay actions, a localized live result, and typed prompt/stage/result rows. It is a deterministic browser preview: it has no editable command field, shell bridge, installer execution, transfer side effect, or arbitrary network action. Reduced motion completes the transcript immediately; normal motion uses one bounded timer that is reset on selection change and cleared on replay or disconnection.

Delivery uses the terminal as a presentation boundary around its real authentication and versioned transfer API. Administration uses it around the existing registry, SSE snapshots, policy forms, refresh, and stop APIs. Neither product surface presents a shell prompt or changes its security boundary.

## Themes and accessibility

The stored preference is one of `system`, `light`, or `dark` and defaults to `system`. The resolved light or dark theme is written separately to the document root, so operating-system changes update a system-selected page without losing the user's preference. Storage failures fall back safely and never prevent a page from rendering.

Graphite, paper, signal lime, and beak orange are the core identity colors. Signal lime identifies a route, active state, or decisive action rather than covering large surfaces. The shared route and status components pair color with text. Theme and locale use icon-only shared segmented radiogroups: visible `Theme` and `Language` legends are suppressed, while localized button and group names, selected state, one roving tab stop, and arrow/Home/End keyboard behavior remain. Inputs, selects, switches, password fields, file actions, and the shared custom checkbox retain dependable form semantics but receive one normalized Courier appearance, explicit programmatic labels, visible focus treatment, accessible light/dark contrast, and reduced-motion behavior. The checkbox preserves native checked/disabled behavior and emits one composed boolean change event without exposing native chrome.

## Localization

English is the fallback locale. Initial selection checks persisted preference and then the browser language list. Catalog shape derives from the English keys, so TypeScript rejects missing keys in Russian or future locale modules. Runtime pages can perform their initial negotiation independently; adding a locale requires no backend API change.

## Assets and provenance

Retained Courier-owned assets live in [`web/ui/assets`](../web/ui/assets) beside `provenance.json`. The manifest records the source context, every asset's media type, role, dimensions, byte count, SHA-256 digest, and prompt summary. The verification script reads PNG and WebP headers and rejects missing, extra, changed, duplicate, dimension-mismatched, incomplete, or executable asset content. The compact transparent PNG mark is a generated profile of Relay, not the former source-arrow symbol. The canonical raster Relay reference anchors a twelve-image terminal scene family plus the unchanged README panorama. Hero, routing, installation, CLI reference, delivery, and administration each have separately composed wide and portrait cinematic editorial assets. They reserve low-detail regions matching the real terminal layouts, concentrate character and props around their edges, and contain no embedded text, terminal UI, logos, third-party marks, credentials, or required information. `courier-mascot` accepts an optional portrait source and renders it through `<picture>` below the shared mobile breakpoint, allowing a scene to change composition without duplicating component logic. `courier-scene` composes that responsive picture as an immutable base with a pointer-local multi-lobed glow and refractive duplicate. The stable host owns tracking, decorative layers ignore hit testing, fine-pointer coordinates converge through one elapsed-time loop, and the loop stops after convergence or disconnection. Coarse pointers and reduced motion use stable fallbacks. All illustrations are non-essential: headings, controls, routes, and states remain complete if an image cannot load.

The shared `courier-brand-icon` renders reviewed monochrome product geometry from the pinned Simple Icons package or pinned official PowerShell and Scoop repository revisions. Those third-party files have a separate provenance manifest, digest gate, and license/trademark notice; they are not described as Courier-owned assets. The component makes no runtime network request and only supplies an accessible image name when a consumer requests one.

No reference HTML or JavaScript is shipped. Browser-injected AdGuard resources, obsolete mirror artwork and messaging, and remote resources were excluded. [`NOTICE.md`](../web/ui/NOTICE.md) records both Courier asset authorship and all third-party mark license, attribution, and trademark context.

## Verification

Run the UI gate directly with:

```sh
npm ci --prefix web --ignore-scripts --no-audit --no-fund
npm run verify --prefix web
```

The gate checks asset integrity, requires 100% TypeScript statements, branches, functions, and lines through Vitest/V8, and produces a deterministic Vite library build. `make verify` includes this gate before the GoReleaser snapshot.

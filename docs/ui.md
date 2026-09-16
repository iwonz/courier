# Courier UI architecture

Courier browser surfaces share one private workspace package, [`@courier/ui`](../web/ui), for visual tokens, identity assets, icons, localization, theme state, and Lit components. Delivery, administration, and the static landing import this package; consumers do not copy component implementations or maintain parallel token sets.

The visual and verbal rules are defined in the [Courier brand system](brand.md): a field manual connected to a network control room, with clear routes, calm status language, disciplined density, and restrained warmth from Relay.

## Package boundary

The TypeScript entry point exports:

- idempotent `defineCourierElements` registration;
- brand, responsive mascot, stationary scene, status, route, terminal, command-readout, icon-link, brand-icon, checkbox, button, panel, progress, icon, segmented-control, theme-selector, and locale-selector components;
- deterministic cubic Bezier route geometry;
- shared normalized form-control styles;
- typed theme state and browser adapters;
- typed English/Russian catalogs, browser-language negotiation, and translation helpers.

The Vite library build emits an ES module, declarations, and `courier-ui.css`. Consumers load the stylesheet once so document tokens, shared control frames, and reduced-motion behavior remain consistent.

## Terminal and command boundaries

`courier-terminal` is an accessible transparent workspace with toolbar, content, footer, and status slots. Delivery and administration use it around real authentication, transfer, registry, SSE, policy, refresh, and stop state. Consumers retain semantic buttons, links, forms, password/file inputs, selects, and checkboxes. Neither product surface exposes a shell prompt.

`courier-command-readout` is deliberately smaller. It exposes immutable command text, localized heading and description, Copy, a reserved live result region, details, and footer actions. Changing command identity synchronously clears previous copy feedback. It has no timer, transcript, editable field, Run/Replay API, shell bridge, installer execution, transfer side effect, or arbitrary network behavior.

## Themes, controls, and localization

Theme preference is `system`, `light`, or `dark` and defaults to `system`. Resolved theme and stored preference remain separate so operating-system changes can update a system-selected page. Storage failure never prevents rendering.

Graphite, paper, signal lime, and beak orange are shared semantic roles. Theme and locale use icon-only segmented radiogroups with localized accessible names, visible selected state, roving focus, and arrow/Home/End behavior. Shared control-frame tokens define dimensions, border, surface, radius, hover, focus, and theme behavior for segmented groups and semantic icon links. `courier-icon-link` keeps anchor semantics and explicit target/relationship attributes.

Inputs, selects, switches, password fields, file actions, and the shared checkbox keep native semantics while normalizing visible chrome. English is the fallback. Catalog shape derives from English keys, so TypeScript rejects missing Russian or future keys; adding a locale requires no backend change.

## Scenes and loading

Retained assets live in [`web/ui/assets`](../web/ui/assets) beside `provenance.json`. The manifest records media type, role, dimensions, bytes, SHA-256, prompt summary, and landing sequence position. Validation rejects missing, extra, changed, duplicate, dimension-mismatched, executable, incomplete, or over-budget content.

The canonical Relay reference anchors eight landing journey segments and four delivery/administration scenes plus the unchanged README panorama. Landing segments continue one route through departure, routing junction, verification depot, and destination archive. Wide and portrait sequences have four independently composed images each with matching boundary geometry. The landing total is capped at 1.6 MiB and each segment at 225 KiB.

`courier-mascot` renders an optional portrait source through `<picture>`. `courier-scene` activates eager images immediately and non-eager images once at half-viewport proximity. Until activation, no image element or request exists. If `IntersectionObserver` is unavailable, native lazy loading is used. The stationary base is the only image before interaction.

Fine-pointer movement creates a refracted duplicate on demand; ambient convergence removes it. The stable host owns hit testing, decorative layers ignore input, one elapsed-time loop converges and stops, and disconnection cancels observers and frames. Coarse pointers and reduced motion never create moving refraction. Every illustration is optional to meaning.

The shared `courier-brand-icon` renders pinned local official geometry in Courier monochrome. Simple Icons 16.31.0 supplies reviewed marks; official pinned PowerShell and Scoop revisions supply missing geometry. Their provenance, license, attribution, and trademark context are recorded separately in [`NOTICE.md`](../web/ui/NOTICE.md). GNU Wget uses a Courier functional download glyph because this set has no distinct official product mark.

No reference HTML or JavaScript is shipped. Browser-injected AdGuard resources, obsolete mirror artwork/messages, superseded landing scenes, temporary storyboards, and remote runtime assets are excluded.

## Verification

Run:

```sh
npm ci --prefix web --ignore-scripts --no-audit --no-fund
npm run verify --prefix web
```

The gate verifies asset and license provenance, enforces exactly 100% TypeScript statements/branches/functions/lines, builds every UI, validates embedded assets, and enforces landing bundle budgets. `make verify` includes it before the GoReleaser snapshot.

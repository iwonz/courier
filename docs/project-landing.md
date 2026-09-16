# Project landing

Courier's project site is a static Lit/Vite application in `web/landing`. It is deployed to GitHub Pages from this repository and uses neither another repository nor a `gh-pages` maintenance branch.

## Sources of truth

`docs/cli-contract.yaml` is the only command inventory. Running:

```sh
go run ./cmd/contractdoc --write
```

generates `docs/cli-reference.md` and `web/landing/src/contract.generated.json`. The landing receives only shipped endpoints, routes, flags, examples, and product/system commands. Its route pairs, applicable option lists, command compatibility, and examples are projections of that data rather than a parallel capability registry. `contract-check` compares both generated files byte-for-byte and validates the live Cobra tree and operation planner.

Installation commands and repository links are reviewed static content. Download links target stable GitHub Releases pages and never embed credentials or mutable release tokens.

## Continuous application composition

The landing is a workspace peer of the embedded delivery and administration applications. It imports `@courier/ui` for identity, responsive scenes, route geometry, immutable command readouts, brand marks, themes, localization, and controls. English is the fallback; browser languages may select Russian on first use. Theme and locale preferences remain icon-only segmented radio groups with localized accessible names, visible state, roving focus, and arrow/Home/End navigation.

The page contains exactly four sections in this order:

1. dispatch hero;
2. routing;
3. installation;
4. commands and options.

Mandatory scroll snapping is absent. The hero fills at least one small viewport. Routing and installation use natural content height with stable `42rem`/`78svh` minimum geometry, and the CLI reference uses `52rem`/`88svh`; narrow and localized content may expand without clipping. Below-fold sections use `content-visibility` with stable intrinsic placeholders.

The safe-area-aware masthead publishes its real `ResizeObserver` height as the section anchor offset. Its graphite glass and gradual transparent fade have no opaque panel or lower rule. A narrow, header-aware `IntersectionObserver` reading band tracks variable-height sections and exposes the active route, installation, or CLI link through `aria-current`. In-page navigation is smooth unless reduced motion is requested. GitHub is a semantic external `courier-icon-link` that shares the same control-frame dimensions, border, surface, radius, hover, focus, and theme tokens as the preference controls. Absolute off-landing links open with `_blank`, `noopener`, and `noreferrer`.

## Continuous Relay panorama

The landing is one responsive Relay journey split into four wide and four portrait WebP segments:

1. departure and route overview;
2. routing junction;
3. verified distribution/checksum depot;
4. destination control and archive.

Adjacent segments share signal-route geometry, perspective, light, palette, and low-detail boundary bands. Each section extends its scene into a masked overlap so the transition fades into the page rather than revealing a tile or uncovered region. Portrait scenes are independently composed and selected by `<picture>`, not automatic desktop crops. The README banner and delivery/administration artwork remain unchanged. All raster scenes are decorative and contain no text, fake UI, logos, credentials, watermarks, weapons, or required product information.

The hero contains the headline but no visible tagline, installation command, CTA, or activation state. It depicts the literal roles `Source` and `Destination` as crisp circular HTML terminals joined by responsive cubic Bezier routes.

`courier-scene` keeps its base completely stationary. The hero loads eagerly with high fetch priority. Non-eager scenes create no image element until a one-shot half-viewport proximity observer activates them, then disconnect their observer. Browsers without `IntersectionObserver` receive native lazy-loaded responsive images. The boundary keeps the next segment ready while preventing later segments from joining the initial request path.

The stable scene host owns fine-pointer tracking, while every decorative descendant ignores hit testing. One elapsed-time loop converges on the latest target and returns to ambient. A refracted duplicate is created only after fine-pointer interaction and removed after ambient convergence. Touch, coarse pointers, and reduced motion never create it. Disconnection cancels pending observers and animation frames.

## Read-only command surfaces

The route explorer projects `from <source> to <destination>` from the shipped contract. `Source` and `Destination` remain untranslated protocol-role labels. Both sides present `Local`, `Remote`, `Web`, and `Web Hook`: `Remote` maps to SSH, source `Web Hook` maps to incoming `webhook://`, and destination `Web Hook` maps to outgoing HTTP(S). Direction-specific icons and descriptions explain each role. Hover and focus are previews only; pointer or keyboard activation changes selection. Invalid destinations are disabled buttons. A measured cubic Bezier connects the selected controls.

Routing exposes the exact immutable command, concise route explanation, applicable options, Copy, and a reserved localized result. Installation exposes compact icon-led channel tags, one exact immutable command, Copy, repository release actions, and locally bundled official monochrome marks. The CLI reference exposes selectable generated commands, independently scrollable options, and a designed compatibility checkbox backed by each command's generated flag list. A selected command exposes its generated usage in the same readout.

No landing surface contains Run/Replay, staged transcripts, demo timers, preview status, editable command input, shell execution, installer execution, transfer execution, or arbitrary browser network behavior. Changing a route, channel, or command synchronously clears old copy feedback without resizing its surface.

The previous statistics, route-card grid, safety section, examples, documentation section, full-contract link, closing CTA, and footer remain intentionally absent.

## Performance and verification

The UI and landing Vitest configurations require exactly 100% statements, branches, functions, and lines. Asset validation records dimensions, bytes, SHA-256, prompt, role, orientation, and sequence position; it caps each landing WebP at 225 KiB and all eight at 1.6 MiB. The deterministic Pages build caps landing JavaScript at 45 KiB gzip and CSS at 8 KiB gzip.

Pinned Playwright acceptance covers English/Russian, system/light/dark, reduced motion, keyboard, fine/coarse pointer, 320×568, 360×740, 390×844, short desktop, and 1440×900. It verifies natural scrolling, reading-band navigation, Copy and localized live status, absence of simulations/editable commands, responsive source selection, hero/next-only initial requests, panorama coverage, invariant selector geometry, circular hero terminals, stationary base scenes, external-link protection, compatibility filtering, overflow safety, delivery metadata isolation, administration controls, and zero external runtime requests.

Build the current contract-backed artifact with:

```sh
make pages-build
```

`.github/workflows/pages.yml` runs after every push to `main` and may be dispatched manually. The build job has read-only repository access and uploads only `web/landing/dist`; the deploy job alone receives `pages: write` and short-lived OIDC permission.

An authenticated maintainer can rebuild and republish synchronized `main` with:

```sh
make pages-publish
```

The command requires `gh auth login`, clean local `main` matching `origin/main`, and no token argument. It verifies the local build, enables workflow-based Pages if needed, dispatches the exact revision, waits for deployment, and reports <https://iwonz.github.io/courier/>.

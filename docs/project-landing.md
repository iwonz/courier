# Project landing

Courier's project site is a static Lit/Vite application in `web/landing`. It is deployed at GitHub Pages from this repository and does not use another repository or a `gh-pages` branch.

## Source of truth

`docs/cli-contract.yaml` remains the only command inventory. Running:

```sh
go run ./cmd/contractdoc --write
```

generates both `docs/cli-reference.md` and `web/landing/src/contract.generated.json`. The landing projection contains shipped endpoints, routes, options, product commands, system commands, and current examples; entries marked `planned` are filtered before serialization. The page expands this projection into endpoint pairs and exact route-specific option lists, so its interactive route illustration does not maintain a second capability registry. `contract-check` compares both generated files byte-for-byte and still checks the live Cobra tree and operation planner.

Installation commands and repository links are static reviewed content. Download links target the stable GitHub Releases pages instead of embedding mutable release API responses or credentials.

## Application

The landing is a workspace peer of the embedded data and administration applications. It imports `@courier/ui` for the shared terminal-workbench identity, contextual Relay scenes, route grammar, terminal surfaces and demonstrations, graphite/signal-lime tokens, persisted `system`/`light`/`dark` themes, reduced-motion handling, and English/Russian icon-only segmented selectors. Their visible legends and option words are removed, while localized accessible names, radiogroup state, roving focus, and arrow/Home/End navigation remain. English is the fallback and supported browser languages select Russian on first use.

The page has exactly four viewport-sized snap sections in this order: dispatch hero, routing, installation, and commands/options. The permanently visible masthead respects `viewport-fit=cover` safe areas and publishes its measured `ResizeObserver` height to the slide layout instead of assuming a mobile constant. It has no rule or opaque panel: translucent graphite glass and a soft lower fade merge it with every scene. Its left navigation follows the three working sections and identifies the currently intersecting section with `aria-current`, while GitHub and the icon-only preferences remain at the right. Every absolute off-landing link opens in a new context with `noopener noreferrer`; in-page section links stay in the current context.

Every section uses a purpose-composed full-bleed Relay environment instead of a mascot tile. The hero leaves headline space inside a signal-station landscape, contains no install command or CTA, exposes no activation state, and continuously depicts the literal roles `Source` and `Destination` as crisp circular HTML terminals joined by independently composed desktop and mobile cubic Bezier handoffs. The routing signal yard, release inspection depot, and command archive reserve quiet interface zones while Relay and narrative detail stay around their edges. Translucent terminal surfaces therefore read as part of the environment instead of covering unrelated artwork. Each scene has an independently generated portrait counterpart selected by the shared responsive component rather than cropping the desktop image.

The base scene remains completely stationary. The stable `courier-scene` host owns pointer tracking; its base image, refracted duplicate, glow, and veil never intercept input. On a fine pointer, raw targets and rendered coordinates are separate: one animation loop uses elapsed-time exponential interpolation, clamps long frame gaps, converges below a fixed threshold, and eases back to the ambient position on exit. The resulting coordinates drive only an irregular multi-lobed glow and a masked, slightly scaled duplicate that creates restrained local refraction. Touch and other coarse pointers retain a fixed ambient treatment. Reduced motion removes lens travel, breathing, and refraction while keeping content readable. The effect never changes functional state or participates in layout, and the loop stops after convergence or component disconnection.

The route explorer projects the exact `from <source> to <destination>` matrix from the shipped contract. `Source` and `Destination` are untranslated protocol-role labels in both catalogs. Its compact presentation model shows `Local`, `Remote`, `Web`, and `Web Hook` on both sides: `Remote` maps to SSH, source `Web Hook` maps to incoming `webhook://`, and destination `Web Hook` maps to outgoing HTTP(S). Direction-specific folder, server, browser-upload/browser-share, and webhook receiver/sender icons plus concise descriptions explain what each role does without duplicating route availability. Hover and focus provide preview styling, but only pointer or keyboard activation changes a selection. Invalid destinations are actual disabled buttons, and a measured cubic Bezier path joins the selected terminal controls. The full-width command terminal exposes Copy and Run demo; its contract-backed transcript advances through preflight, route selection, transfer preview, verification, and completion while explicitly stating that no data leaves the page. Selection changes reset the transcript without moving panel boundaries.

Installation presents compact, wrapping icon-led tags and one width-safe terminal surface. It changes selection only through activation, and its chooser, command, transcript, and release actions occupy invariant tracks. Copy writes the exact selected command and announces a localized success or failure inside a reserved status row. Run demo simulates release resolution, OS/architecture selection, artifact download, checksum verification, and installation-path reporting; the browser never runs the installer. Official package-manager, shell, operating-system, and distribution geometry is bundled locally through the shared monochrome brand component; GNU Wget intentionally uses a functional download glyph. Source revisions, licenses, attributions, and trademark constraints are recorded in the UI notice and checked with the assets.

Commands and options share bounded, independently scrollable transparent terminal panes. Command buttons start with no selection and toggle through `aria-pressed`. The shared designed checkbox starts checked but disabled while no command is selected, then filters options through the selected command's generated `flags` array; it can be unchecked without clearing the command. Selecting a command enables a separate immutable help demo with Copy and Run demo. Its transcript renders generated usage, compatible options and defaults, then reports a simulated successful help invocation; commands without flags retain the localized empty state. No editable command field or browser-side shell exists. The previous statistics, route-card grid, safety, examples, documentation, complete-contract link, closing CTA, and page footer are intentionally absent. GitHub source access remains in the masthead. The README uses the unchanged repository-local 3:1 route panorama as a full-width introduction.

The page covers curl, wget, PowerShell, npm, npx, Yarn, pnpm, Homebrew, Scoop, direct binaries, checksums, Deb, RPM, APK, and Arch packages. Each channel has a platform or package-manager mark and a bounded command region that wraps safely and remains selectable instead of clipping, widening the document, or changing the panel height. It lists the supported distributions and renders only the generated shipped CLI surface.

## Verification and deployment

The landing's Vitest configuration requires 100% statements, branches, functions, and lines. The root web gate tests the UI kit and all three consumers, verifies local asset/license provenance, produces a deterministic production build, and is part of `make verify` and every GoReleaser dry run. Pinned Playwright acceptance exercises the compiled landing in Chromium across both locales, all theme preferences, reduced motion, keyboard activation, fine/coarse pointers, 320×568, 360×740, 390×844, short desktop, and 1440×900 viewports. It runs routing, installation, and help demos through pointer and keyboard activation, verifies Copy and localized live output, checks every valid route and installation channel for invariant bounds within one pixel, proves hero terminals stay circular, checks header containment and safe-area spacing at every snap position, stationary base-scene transforms, smoothed spotlight coordinate changes, protected external links, command compatibility filtering, internal overflow, and the absence of external runtime requests.

Build the current contract-backed production artifact locally with:

```sh
make pages-build
```

`.github/workflows/pages.yml` automatically runs after every push to `main` and can also be dispatched manually. Its build job has read-only repository access, calls the same `pages-build` target, and uploads only `web/landing/dist`. The deploy job alone receives `pages: write` and short-lived OIDC permission for the protected `github-pages` environment. No publication token is written into source or static assets.

An authenticated maintainer can explicitly rebuild and republish the synchronized default branch with:

```sh
make pages-publish
```

The command requires `gh auth login`, a clean local `main` exactly matching `origin/main`, and no separate token argument. It verifies the local build, enables workflow-based Pages if necessary, dispatches the exact revision, waits for deployment, and prints the public URL: <https://iwonz.github.io/courier/>.

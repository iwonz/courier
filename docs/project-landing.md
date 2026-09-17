# Project landing

Courier's static Lit/Vite landing lives in `web/landing` and deploys to GitHub Pages from this repository without a second repository or maintenance branch.

## Contract-backed content

`docs/cli-contract.yaml` is the command source of truth. `go run ./cmd/contractdoc --write` generates both `docs/cli-reference.md` and `web/landing/src/contract.generated.json`. Route pairs, applicable options, command compatibility, and usage are projections of shipped contract entries. Installation commands and GitHub links are reviewed static content.

## Three-section journey

The page contains exactly three naturally scrolling sections:

1. the headline **From here to anywhere.** combined with the Source-to-Destination route instrument;
2. installation channels and one immutable command readout;
3. the generated commands and options reference.

The first section is at least one small viewport on capable desktops and expands on short or narrow screens. The fixed safe-area masthead measures its rendered height for anchor offsets and tracks Install/CLI through a header-aware observer band. GitHub and all other external links open in a new tab with `noopener noreferrer`.

Source and Destination controls expose Local, Remote, Web, and Web Hook using direction-specific icons and assistive descriptions. Descriptions are not visible labels. Selection changes only on click or keyboard activation; invalid destinations use real disabled semantics. A measured cubic Bézier and restrained orange pulse visualize the active contract route. Reduced motion keeps the route static.

Installation uses locally bundled official monochrome marks and compact channel controls. Commands/options use two independently scrollable workbench columns, selectable generated commands, compatible-option filtering, and Copy. There are no simulations, Run/Replay controls, transcripts, editable commands, shell execution, or browser transfer execution.

## Continuous editorial panorama

One responsive wide/portrait WebP pair spans the complete page behind all three semantic sections. The selected image depicts one uninterrupted journey from route selection through installation to the command registry; sections never create their own scene, image mask, repeated horizon, or overlap. A page-wide theme-aware veil and feathered heading fields preserve contrast without turning headings into opaque cards. Product information never depends on the decorative panorama.

Each panorama is at most 420 KiB, the pair totals at most 700 KiB, and the ImageGen-authored compact raster mark is at most 96 KiB. Exactly one browser-selected panorama source loads eagerly at high priority; the alternate responsive source is not requested.

The production bundle is capped at 45 KiB gzip JavaScript and 9 KiB gzip CSS. Build locally with `make pages-build`. `main` deploys automatically through `.github/workflows/pages.yml`; an authenticated maintainer can rebuild and dispatch the synchronized revision with `make pages-publish`.

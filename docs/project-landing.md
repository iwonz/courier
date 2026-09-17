# Project landing

Courier's static React/Vite landing lives in `web/landing`, uses repository-owned shadcn components, and deploys to GitHub Pages from this repository without a second repository or maintenance branch.

## Contract-backed content

`docs/cli-contract.yaml` is the command source of truth. `go run ./cmd/contractdoc --write` generates both `docs/cli-reference.md` and `web/landing/src/contract.generated.json`. Route pairs, applicable options, command compatibility, and usage are projections of shipped contract entries. Installation commands and GitHub links are reviewed static content.

## Three-section flow

The page contains exactly three naturally scrolling sections:

1. the headline **From here to anywhere.** combined with the Source-to-Destination route instrument;
2. installation channels and one immutable command readout;
3. the generated commands and options reference.

Every section uses natural content height and ordinary block spacing: there are no viewport-height targets, full-screen slides, or snap points. The sticky header keeps the brand, Install/CLI navigation, GitHub, theme, and locale controls available without obscuring content. GitHub and all other external links open in a new tab with `noopener noreferrer`.

Source and Destination controls expose Local, Remote, Web, and Web Hook using direction-specific icons and assistive descriptions. Descriptions are not visible labels. Selection changes only on click or keyboard activation; invalid destinations use real disabled semantics. A measured cubic Bézier and restrained orange pulse visualize the active contract route. Reduced motion keeps the route static.

Installation uses locally bundled official monochrome marks and compact channel controls. Commands/options use two independently scrollable workbench columns, selectable generated commands, compatible-option filtering, and Copy. There are no simulations, Run/Replay controls, transcripts, editable commands, shell execution, or browser transfer execution.

## Compact Relay hero

One transparent 1:1 WebP appears inside the first section beside the headline. It depicts the compact technological Relay pigeon between code-native Source and Destination waypoints. Installation and CLI do not create decorative images, background scenes, masks, repeated horizons, or stitched artwork. Product information never depends on the illustration.

The ImageGen-authored mascot is at most 80 KiB and is the only generated raster shipped by the UI package.

The React/shadcn production bundle is capped at 145 KiB gzip JavaScript and 9 KiB gzip CSS. Build locally with `make pages-build`. `main` deploys automatically through `.github/workflows/pages.yml`; an authenticated maintainer can rebuild and dispatch the synchronized revision with `make pages-publish`.

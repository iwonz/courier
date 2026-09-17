# Project landing

Courier's static React/Vite landing lives in `web/landing`, uses repository-owned shadcn components, and deploys to GitHub Pages from this repository without a second repository or maintenance branch.

## Contract-backed content

`docs/cli-contract.yaml` is the command source of truth. `go run ./cmd/contractdoc --write` generates both `docs/cli-reference.md` and `web/landing/src/contract.generated.json`. Route pairs, applicable options, command compatibility, and usage are projections of shipped contract entries. Installation commands and GitHub links are reviewed static content.

## Three-section flow

The page contains exactly three naturally scrolling sections:

1. the headline **From here to anywhere.** combined with the Source-to-Destination route instrument;
2. installation channels and one immutable command readout;
3. the generated commands and options reference.

Every section uses natural content height and compact adjacent spacing: there are no viewport-height targets, full-screen slides, snap points, stacked section padding, rounded section boards, or tinted card islands. The route hero and controls form one transparent composition, installation continues on the same canvas, and the CLI registry uses responsive columns with only a meaningful internal divider. The backgroundless header stays fixed to the viewport; document padding keeps content below it. It places the square Relay mark and `COURIER CLI` beside Install/CLI navigation, with GitHub and preference controls on the opposite edge. GitHub and all other external links open in a new tab with `noopener noreferrer`.

Source and Destination controls expose Local, Remote, Web, and Web Hook using direction-specific icons and assistive descriptions. Descriptions are not visible labels. Selection changes only on click or keyboard activation; invalid destinations use real disabled semantics. A measured cubic Bézier and restrained orange pulse visualize the active contract route. Reduced motion keeps the route static.

Installation uses locally bundled official monochrome marks and compact channel controls. Commands/options use two independently scrollable transparent regions, selectable generated commands, compatible-option filtering, and Copy. Immutable command readouts sit directly in the document flow. There are no simulations, Run/Replay controls, transcripts, editable commands, shell execution, or browser transfer execution.

## Compact Relay hero

One transparent 1:1 WebP appears inside the first section beside the headline. It depicts the compact technological Relay pigeon between code-native Source and Destination waypoints. Installation and CLI do not create decorative images, background scenes, masks, repeated horizons, or stitched artwork. Product information never depends on the illustration.

The full ImageGen-authored mascot and dedicated square header mark are each at most 80 KiB and together remain below 160 KiB. The mark is also the favicon; later sections allocate no decorative raster.

The React/shadcn production bundle is capped at 145 KiB gzip JavaScript and 9 KiB gzip CSS. Build locally with `make pages-build`. `main` deploys automatically through `.github/workflows/pages.yml`; an authenticated maintainer can rebuild and dispatch the synchronized revision with `make pages-publish`.

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

Source and Destination controls expose Local, Remote, Web, and Web Hook using direction-specific pixel icons and assistive descriptions. Descriptions are not visible labels. Selection changes only on click or keyboard activation; invalid destinations use real disabled semantics. The measured cubic Bézier is sampled and snapped to the four-pixel grid; a small packet uses stepped timing. Reduced motion keeps the route static.

Installation uses locally bundled official monochrome marks and compact channel controls. Commands/options use two independently scrollable transparent regions, selectable generated commands, compatible-option filtering, and Copy. Immutable command readouts sit directly in the document flow. There are no simulations, Run/Replay controls, transcripts, editable commands, shell execution, or browser transfer execution.

## Pixel Relay hero

The 512×512 transparent route sprite appears inside the first section beside the Pixelify Sans headline. It depicts the same canonical 8-bit Relay pigeon between code-native Source and Destination signals; its anatomy and proportions match the neutral, delivery, admin, and compact-mark family. Installation and CLI do not create decorative images, background scenes, masks, repeated horizons, or stitched artwork. Product information never depends on the illustration.

The route sprite is at most 64 KiB and the 256×256 mark is at most 24 KiB. The mark is also the favicon; later sections allocate no decorative raster or request any other role sprite. Both use pixelated sampling and lossless transparent WebP.

The React/shadcn production bundle is capped at 145 KiB gzip JavaScript and 9 KiB gzip CSS. Build locally with `make pages-build`. `main` deploys automatically through `.github/workflows/pages.yml`; an authenticated maintainer can rebuild and dispatch the synchronized revision with `make pages-publish`.

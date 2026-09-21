# Project landing

Courier's static React/Vite landing lives in `web/landing`, uses repository-owned shadcn components, and deploys to GitHub Pages from this repository without a second repository or maintenance branch.

## Contract-backed content

`docs/cli-contract.yaml` is the command source of truth. `go run ./cmd/contractdoc --write` generates both `docs/cli-reference.md` and `web/landing/src/contract.generated.json`. Route pairs, command paths, ordered arguments, typed parameter values, choices, dependencies, conflicts, and applicability are projections of shipped contract entries. Installation commands and GitHub links are reviewed static content.

## Two-section flow

The page contains exactly two naturally scrolling sections:

1. the headline **From here to anywhere.**, Relay, and one unified route-and-command workspace;
2. installation channels and one immutable command readout.

Every section uses natural content height and compact adjacent spacing: there are no viewport-height targets, full-screen slides, snap points, rounded boards, detached card islands, contract-count strips, or boundary rules between the header and sections. Installation uses horizontally scrollable channel tabs with stable, immediate active and hover states plus one unframed command area. A single direct-binaries text link remains beside the tabs, outside their scroller, including on mobile. The nested `#cli` workspace uses responsive command selection and options without an outer table frame, header rules, row rules, inter-column dividers, or a readout rule. The transparent `h-16` header is a static normal-flow row that scrolls away. It places the square Relay mark and `COURIER CLI` on the left and GitHub plus preference controls on the right, with no internal navigation or compensating main padding.

Source and Destination controls expose Local, Remote, Web, and Web Hook using direction-specific pixel icons and assistive descriptions. Descriptions are not visible labels. Selection changes only on click or keyboard activation; invalid destinations use real disabled semantics. A four-by-four-pixel signal uses a centered offset anchor and moves linearly on the measured connector path. Reduced motion places it statically at the midpoint. These controls are templates for `courier from`: the initial local-to-SSH pair pre-fills editable arguments, selecting another pair replaces both examples and clears parameters, and manual input is not restricted to the template's endpoint kind. Other CLI commands hide the route templates; returning to `from` restores the selected pair's examples. One contract-built readout and Copy action serves the entire workspace.

Installation uses locally bundled transparent PNG marks in official geometry and color; GitHub has light/dark variants, while npx and wget remain text-only. The selected package-manager identity appears only in its tab and is not repeated below the command. Commands/options use a horizontal command strip on narrow screens, a vertical list on wide screens, compatible-option filtering, typed value controls, and Copy. The builder emits path, arguments, and explicitly supplied parameters in contract order, preserves repeatable order, resolves dependencies/conflicts, validates required values, and supports exact POSIX or PowerShell quoting. Defaults are hints and are never materialized. There is no Run action or browser command execution.

## Pixel Relay hero

The 512×512 transparent route-v3 sprite appears inside the first section beside the Pixelify Sans headline. It depicts the terminal-palette Relay between code-native Source and Destination signals; its anatomy and proportions match neutral-v3, delivery-v3, admin-v3, and mark-v3. Installation does not create decorative images, background scenes, masks, repeated horizons, or stitched artwork. Product information never depends on the illustration.

The route sprite is at most 64 KiB and the 256×256 mark is at most 24 KiB. The mark is also the favicon; the installation section allocates no decorative raster or requests any other role sprite. Both use pixelated sampling and transparent WebP.

The React/shadcn production bundle is capped at 150 KiB gzip JavaScript and 9 KiB gzip CSS. Build locally with `make pages-build`. `main` deploys automatically through `.github/workflows/pages.yml`; an authenticated maintainer can rebuild and dispatch the synchronized revision with `make pages-publish`.

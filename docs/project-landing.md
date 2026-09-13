# Project landing

Courier's project site is a static Lit/Vite application in `web/landing`. It is deployed at GitHub Pages from this repository and does not use another repository or a `gh-pages` branch.

## Source of truth

`docs/cli-contract.yaml` remains the only command inventory. Running:

```sh
go run ./cmd/contractdoc --write
```

generates both `docs/cli-reference.md` and `web/landing/src/contract.generated.json`. The landing projection contains shipped endpoints, routes, options, product commands, system commands, and current examples; entries marked `planned` are filtered before serialization. `contract-check` compares both generated files byte-for-byte and still checks the live Cobra tree and operation planner.

Installation commands and repository links are static reviewed content. Download links target the stable GitHub Releases pages instead of embedding mutable release API responses or credentials.

## Application

The landing is a workspace peer of the embedded data and administration applications. It imports `@courier/ui` for graphite/signal-lime tokens, shared elements, identity conventions, persisted `system`/`light`/`dark` themes, reduced-motion handling, and the English/Russian locale selector. English is the fallback and supported browser languages select Russian on first use.

The page covers curl, wget, PowerShell, npm, npx, Yarn, pnpm, Homebrew, Scoop, direct binaries, checksums, Deb, RPM, APK, and Arch packages. It lists the supported distributions and renders only the generated shipped CLI surface.

## Verification and deployment

The landing's Vitest configuration requires 100% statements, branches, functions, and lines. The root web gate tests the UI kit and all three consumers, produces a deterministic production build, and is part of `make verify` and every GoReleaser dry run. Pinned Playwright acceptance also exercises the compiled landing in Chromium across both locales, all theme preferences, keyboard navigation, persisted selectors, and narrow/wide viewports.

`.github/workflows/pages.yml` runs independently on `main` and by manual dispatch. Its build job has read-only repository access, verifies generated contract data, tests and builds the landing, and uploads only `web/landing/dist`. The deploy job alone receives `pages: write` and short-lived OIDC permission for the protected `github-pages` environment. No publication token is written into source or static assets.

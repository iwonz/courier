# Change: Publish the Courier project landing

## Why

Courier's installation, command, and release information currently lives across repository Markdown files. Users need one accessible, localized entry point that never drifts from the machine-readable CLI contract and can be published without creating another repository or maintenance branch.

## What Changes

- Add a static Lit/Vite landing application that consumes the shared Courier UI kit, themes, identity assets, and typed English/Russian catalogs.
- Generate deterministic landing command data from `docs/cli-contract.yaml`, exposing only entries whose status is `shipped` or `system` as appropriate.
- Present install channels, shipped routes and commands, practical examples, GitHub Release downloads, source, documentation, and security links.
- Add freshness, exact TypeScript coverage, build, responsive, accessibility-oriented, locale, and theme tests to `make verify`.
- Publish the built artifact from `main` through GitHub Pages Actions with least-privilege workflow permissions and no `gh-pages` branch.

## Impact

CLI behavior and release packaging do not change. GitHub Pages becomes an independently deployable repository surface. No token is embedded in the site, no runtime API is introduced, and no external package or Pages repository is required.

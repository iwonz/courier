# Context

The repository already contains a canonical YAML contract, deterministic Markdown generation, a shared Lit UI kit, two embedded applications, and release workflows. The landing page should reuse these assets and gates instead of becoming a manually maintained documentation copy.

# Decisions

## Generated public data

`contractdoc` generates `web/landing/src/contract.generated.json` together with `docs/cli-reference.md`. The JSON contains contract/release versions, shipped endpoint kinds, shipped routes, shipped product commands, system commands, and examples. Planned items are filtered before serialization. Both `--write` and `--check` handle the two outputs atomically from the operator's perspective, and tests reject stale or accidentally exposed planned entries.

The landing application imports this static JSON at build time. It never fetches mutable metadata or credentials at runtime. Release download buttons point to the stable GitHub Releases pages rather than guessing an asset for the visitor.

## Shared static application

`web/landing` is a workspace peer of `ui`, `data`, and `admin`. It aliases `@courier/ui` to source for deterministic builds and uses the same persisted `system`/`light`/`dark` selector, locale selector, components, tokens, identity assets, reduced-motion behavior, and keyboard semantics. English is the fallback; browser language selects Russian when supported.

The layout is a responsive single-page guide with semantic navigation, install cards, command/route tables, copyable code blocks, and repository links. User-controlled HTML is never rendered.

## Pages deployment

A dedicated workflow runs on successful pushes to `main` and manual dispatch. It installs pinned Node dependencies, verifies generated contract freshness, runs the landing test/build gate, uploads `web/landing/dist`, and deploys through the official Pages artifact action. The build job has read-only contents access; only the deploy job receives Pages and OIDC permissions. Concurrency cancels superseded deployments.

# Risks / Trade-offs

The static landing shows the contract version from the checked-in main branch rather than querying releases at runtime, which makes builds reproducible and avoids API rate limits. Release links remain useful between releases, while the target release and contract versions clearly describe the documented surface.

# Migration Plan

Add the OpenSpec delta, extend deterministic contract generation, build and test the shared-kit landing, add the Pages workflow and documentation, run all release gates, archive the change, and commit once.

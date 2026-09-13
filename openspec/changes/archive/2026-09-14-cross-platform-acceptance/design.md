# Context

Courier already has deterministic contract tests, 100% first-party statement coverage, in-process native SSH/SFTP tests, secure HTTP handler tests, package artifacts, and shared UI unit tests. Acceptance should compose those boundaries and avoid maintaining a second transfer engine or external infrastructure repository.

# Decisions

## Layered acceptance

Fast Go/TypeScript tests remain the exhaustive logic layer. A black-box Go acceptance package invokes the compiled command composition with isolated mapped endpoint backends, exercising all four path directions and transformation/error semantics without network flakiness. Native SSH/SFTP cryptography, authentication, ProxyJump, known-hosts, and helper cleanup remain covered by the existing in-process protocol integration suite.

Large streams are generated, never checked into the repository. Bounded writers measure maximum request size and block/release to prove backpressure and cancellation. Race-enabled tests cover every shared mutable domain in the normal Go gate.

## Browser acceptance

Pinned Playwright drives Chromium against three local Vite servers. Requests for data/admin APIs are intercepted with controlled secret-bearing fixtures so tests prove protected values do not render. Tests exercise locale persistence, all theme preferences, keyboard focus, semantic landmarks, and narrow/wide viewports. Vitest remains responsible for exact first-party TypeScript coverage; browser tooling is not counted as application source.

## Distribution containers

The package script chooses the architecture reported by Docker and installs the matching GoReleaser artifact in seven representative images on the required amd64 CI runner. The official Arch image is amd64-only; an ARM development host verifies the arm64 Arch package through Manjaro and leaves the native Arch image check to CI. Every container carries a unique Courier ownership label and `--rm`; an EXIT/HUP/INT/TERM trap forcibly removes leftovers. Final assertions query containers, networks, and volumes only by that label and fail if any remain. The script never prunes or deletes unrelated Docker resources.

## CI and release gates

Ubuntu runs `make verify`, including browser and package acceptance. macOS and Windows run their compiled runtime suites, with exact first-party coverage on each OS path. GoReleaser remains the only release builder and cross-builds the BSD helper matrix. Workflow syntax is checked locally and in CI with pinned actionlint. Release publication depends on the same acceptance jobs and continues to use only repository GitHub Releases/npm/Homebrew/Scoop outputs.

# Risks / Trade-offs

Container and browser checks add time to the full gate, so fast package-local test commands remain available during development. Distribution images are pinned by stable tags rather than immutable digests to receive security fixes; package behavior is deterministic because the mounted Courier artifact and verification command are fixed.

# Migration Plan

Add the OpenSpec delta, fix real-browser event compatibility, add browser and black-box acceptance, add labeled package tests and platform scripts, wire CI/release gates, validate cleanup and performance locally, then archive and commit once.

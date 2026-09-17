# Courier CLI implementation plan

Work is split into stacked branches. Each branch is based on the previous branch, owns an OpenSpec change, and ends with one conventional commit.

| # | Branch | OpenSpec change | Outcome |
|---|---|---|---|
| 1 | `feat/001-project-foundation` | `project-foundation` | Go/CLI foundation, MIT, architecture contract |
| 2 | `feat/002-endpoints-and-safety` | `endpoints-and-safety` | Local/remote parsing, destination semantics, path safety |
| 3 | `feat/003-transfer-and-archive` | `transfer-and-archive` | Filesystem engine, sync, atomic staging, archives, progress events |
| 4 | `feat/004-native-ssh-transport` | `native-ssh-transport` | SSH config/auth/known_hosts/ProxyJump, SFTP, remote OS detection |
| 5 | `feat/005-cli-update-reporting` | `cli-update-reporting` | `from … to …`, reporting, exit codes, self-update |
| 6 | `feat/006-distribution-and-quality` | `distribution-and-quality` | GoReleaser, installers, npm/Brew, CI, integration and cleanup tests |
| 7 | `feat/007-optional-remote-helper` | `optional-remote-helper` | Consent-gated helper fallback for genuine remote capability gaps and Windows runtime hardening |
| 8 | `fix/008-windows-ci` | `windows-ci` | Portable Windows CI and coverage handling |
| 9 | `fix/009-release-verification` | `release-verification` | Release and package publication verification |
| 10 | `fix/010-self-contained-distribution` | `self-contained-distribution` | Single-repository package manifests and self-contained release channels |
| 11 | `fix/011-homebrew-trust` | `homebrew-trust` | Homebrew cask trust guidance |
| 12 | `chore/012-openspec-baseline` | `openspec-baseline` | Consolidated OpenSpec baseline for completed changes 001–011 |
| 13 | `feat/013-cli-contract-registry` | `cli-contract-registry` | Canonical CLI contract, generated reference, contract gate, and command providers |
| 14 | `feat/014-operation-routing` | `operation-routing` | Typed endpoints and routes, option matrix, and preflight taxonomy |
| 15 | `fix/015-safe-copy-semantics` | `safe-copy-semantics` | No-op identity, collision rejection, no-replace commits, and destination preservation |
| 16 | `feat/016-selection-engine` | `selection-engine` | Shared ordered gitignore and regular-expression selection |
| 17 | `feat/017-archive-extraction` | `archive-extraction` | Tar.gz codec registry, verified extraction, extraction-root semantics, and bomb limits |
| 18 | `feat/018-courier-ui-kit` | `courier-ui-kit` | Sanitized identity assets and shared Lit/TypeScript/Vite UI kit |
| 19 | `feat/019-delivery-registry` | `delivery-registry` | Server and delivery domain, policies, counters, tombstones, and private persistence |
| 20 | `feat/020-worker-lifecycle` | `worker-lifecycle` | Shared-bind workers, versioned IPC, leases, stop lifecycle, and stale cleanup |
| 21 | `feat/021-delivery-policy-core` | `delivery-policy-core` | Authentication, sessions, IP policy, reservations, and aggregate size/rate enforcement |
| 22 | `feat/022-web-deliveries` | `web-deliveries` | Browser upload/download delivery routes and opaque resources |
| 23 | `feat/023-webhook-deliveries` | `webhook-deliveries` | Incoming multipart and outgoing single-file webhooks |
| 24 | `feat/024-server-control` | `server-control` | Server listing and UUID-scoped stop commands |
| 25 | `feat/025-admin-ui` | `admin-ui` | UI singleton, guarded admin API, live progress, and optimistic policy editing |
| 26 | `feat/026-operational-reporting` | `operational-reporting` | Unified stages, counters, redacted diagnostics, history, and final exit-code mapping |
| 27 | `feat/027-project-landing` | `project-landing` | Static GitHub Pages landing generated from shipped contract data |
| 28 | `test/028-cross-platform-acceptance` | `cross-platform-acceptance` | Route, security, platform, installer, performance, race, and browser acceptance suites |
| 29 | `feat/029-pages-publication` | `pages-publication` | Automatic GitHub Pages deployment and guarded one-command republishing |
| 30 | `fix/030-workflow-shellcheck` | `workflow-shellcheck-fix` | Pinned local ShellCheck parity and release-workflow lint correction |
| 31 | `fix/031-windows-directory-sync` | `windows-directory-sync-fix` | Platform-correct registry durability and fail-fast Windows tests |
| 32 | `fix/032-windows-native-acceptance` | `windows-native-acceptance-fix` | Windows named-pipe identity, rooted paths, archive links, and permission expectations |
| 33 | `fix/033-windows-pipe-test-lifecycle` | `windows-pipe-test-lifecycle` | Unique bounded Windows named-pipe acceptance and failure-safe cleanup |
| 34 | `feat/034-brand-identity-system` | `brand-identity-system` | Mature brand manual, Relay mascot, shared visual system, and redesigned browser surfaces |
| 35 | `feat/035-illustrated-interface-system` | `illustrated-interface-system` | Branded control grammar, contextual Relay scenes, and full-width README panorama |
| 36 | `feat/036-interactive-landing-system` | `interactive-landing-system` | Icon-only preferences, Relay product mark, compact install channels, interactive route explorer, and unified CLI reference |
| 37 | `feat/037-immersive-landing-sections` | `immersive-landing-sections` | Full-bleed responsive Relay scenes, four snap sections, blended masthead, and scene-aware interactions |
| 38 | `feat/038-landing-interaction-polish` | `landing-interaction-polish` | Safe-area masthead, stationary refractive scenes, click-only selectors, official brand marks, invariant panels, and command-compatible option filtering |
| 39 | `feat/039-landing-scene-integration` | `landing-scene-integration` | Continuous scene integration, pointer spotlight refinement, and responsive route/install composition |
| 40 | `feat/040-terminal-ui-system` | `terminal-ui-system` | Read-only terminal workspaces, product UI integration, and cinematic browser scenes |
| 41 | `feat/041-continuous-landing-flow` | `continuous-landing-flow` | Natural three-stage panorama loading, copy-only command readouts, and page-performance budgets |
| 42 | `feat/042-vector-identity-system` | `vector-identity-system` | Vector moth identity, three-section landing, cyclic preferences, neutral workbenches, and static responsive scenes |

## Definition of Done

- `go test ./...` and the race detector pass;
- first-party Go package statement coverage is 100%;
- `openspec validate --all --strict` passes;
- primary binaries build for macOS, Linux, and Windows on amd64 and arm64, with exact-platform BSD helper archives;
- tests use local fixtures or containers only and always register cleanup;
- local and remote temporary artifacts are removed after success, failure, and cancellation;
- installation and security documentation match the shipped CLI;
- generated command and web assets are current and contract-gated;
- first-party TypeScript statements, branches, functions, and lines are 100% covered once web sources are introduced;
- one release command creates and pushes the tag, builds with GoReleaser Community, publishes the GitHub Release with notes, and publishes `@iwonz/courier` when credentials are available.

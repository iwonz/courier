# Acceptance and release-candidate verification

Courier uses layered acceptance so production implementations, packages, browser surfaces, and release configuration are exercised without test-only engines.

## Local gates

`make test` runs formatting, vetting, the race detector, exact first-party Go statement coverage, and compiled runtime checks. `make verify` adds npm wrapper checks, exact TypeScript coverage, embedded web assets, real Chromium acceptance, workflow validation, strict OpenSpec validation, GoReleaser Community snapshot builds, installers, checksums, and native Linux package installation.

Local prerequisites are Go 1.25+, Node.js 24+, Docker, OpenSpec 1.11.0, and the pinned Playwright Chromium runtime. Repository scripts bootstrap checksum-verified release tools into ignored caches.

## Browser acceptance

Playwright runs landing, delivery, and administration in English/Russian and system/light/dark themes. It covers:

- exactly two landing sections with a unified headline/route/command workspace and installation below;
- pointer and keyboard selection for every valid route pair;
- one exact contract-built POSIX/PowerShell Copy action in the route/command workspace and Copy for every installation channel;
- single-button theme/locale cycles, active-locale flag icons, browser defaults, persistence, and reduced motion;
- 320×568, 390×844, 1024-wide, short desktop, and 1440×900 layouts in light and dark themes without clipping or horizontal overflow;
- one compact transparent pixel mark plus the exact route, delivery, or administration role sprite required by the current surface, stable identity anatomy across every pose, no generated background or unrelated role request, and no external runtime requests;
- the exact terminal light/dark tokens, local Pixelify Sans display scope, local Overpass Mono 400/600 interface scope, pixelated Relay v3 sprites, normal local PNG third-party marks, square controls, and a four-pixel centered linear route signal with a reduced-motion midpoint;
- route console, calm channel chooser, unruled command builder, transfer console, manifest rows, operations metrics, registry, and policy composition with no chamfers, outer radii, offset or soft shadows, blur, fictional data, or decorative card islands;
- compact landing section spacing, a static transparent masthead that scrolls away without reserved content offset, aligned header centers, no top-level boundary rules or contract counters, and visible selected Source/Destination contrast;
- absence of terminal/refraction code, editable commands, Run/Replay, transcripts, emoji flags, Lucide imports, remote fonts, scripts, images, or analytics;
- protected delivery metadata isolation, authentication, upload/download behavior, admin policy updates, stop actions, selection persistence, and SSE fallback.

Asset validation requires exactly five square transparent Relay v3 WebPs: mark-v3 ≤24 KiB, neutral-v3/route-v3 ≤64 KiB each, delivery/admin-v3 ≤48 KiB each, and ≤248 KiB combined. It verifies neutral-v3 authority, controlled terminal-palette and square-silhouette lineage, identity invariants, consumers, hashes, the three-background QA sheet, local transparent 128×128 PNG brand marks and provenance, text-only npx/wget, Pixelify Sans, four Overpass Mono subsets totaling ≤48 KiB, and both OFL notices. Landing JavaScript remains below 150 KiB gzip and CSS below 9 KiB gzip. First-party TypeScript and TSX maintain exact 100% statements, branches, functions, and lines.

## Runtime and cleanup

Black-box tests cover local and SSH path directions, archive/extraction, no-op identity, collision preservation, interruption, confirmed-byte reporting, staging cleanup, and bounded streaming. Native SSH/SFTP tests use an in-process authenticated server.

Temporary files use `t.TempDir`/`t.Cleanup`. Distribution containers carry a unique Courier acceptance label and are removed by success, failure, and signal traps. Acceptance never prunes or deletes unrelated Docker resources.

## CI and publication

CI runs the release-candidate gate on Linux and compiled platform gates on macOS and Windows. It cross-builds all declared client/helper artifacts and installs packages in the documented distribution matrix. A tagged release includes the deterministic source archive; a separate macOS job audits, source-builds, tests, and commits the Formula before final verification. GitHub Releases, Formula/Scoop metadata, and npm publishing remain gated.

# Acceptance and release-candidate verification

Courier uses layered acceptance so production implementations, packages, browser surfaces, and release configuration are exercised without test-only engines.

## Local gates

`make test` runs formatting, vetting, the race detector, exact first-party Go statement coverage, and compiled runtime checks. `make verify` adds npm wrapper checks, exact TypeScript coverage, embedded web assets, real Chromium acceptance, workflow validation, strict OpenSpec validation, GoReleaser Community snapshot builds, installers, checksums, and native Linux package installation.

Local prerequisites are Go 1.25+, Node.js 24+, Docker, OpenSpec 1.11.0, and the pinned Playwright Chromium runtime. Repository scripts bootstrap checksum-verified release tools into ignored caches.

## Browser acceptance

Playwright runs landing, delivery, and administration in English/Russian and system/light/dark themes. It covers:

- exactly three landing sections and the combined headline/route composition;
- pointer and keyboard selection for every valid route pair;
- Copy for routes, every installation channel, and CLI commands;
- single-button theme/locale cycles, active-locale flag icons, browser defaults, persistence, and reduced motion;
- 320×568, 360×740, 390×844, short desktop, and 1440×900 layouts without clipping or horizontal overflow;
- one compact transparent pixel mark plus the exact route, delivery, or administration role sprite required by the current surface, stable identity anatomy across every pose, no generated background or unrelated role request, and no external runtime requests;
- the fixed twelve-color light/dark mapping, local Cyrillic Pixelify Sans display scope, eight-pixel dither, pixelated sprites, crisp icon registries, chamfered controls, and a grid-snapped route with reduced-motion-safe stepped motion;
- continuous route, CLI registry, delivery manifest, administration metrics, and administration workspace composition with no ordinary card backgrounds, outer radii, smooth shadows, blur, or unnecessary borders;
- compact landing section spacing, a transparent viewport-fixed masthead with reserved content offset, and visible selected Source/Destination contrast;
- absence of terminal/refraction code, editable commands, Run/Replay, transcripts, emoji flags, Lucide imports, remote fonts, scripts, images, or analytics;
- protected delivery metadata isolation, authentication, upload/download behavior, admin policy updates, stop actions, selection persistence, and SSE fallback.

Asset validation requires exactly five square transparent lossless Relay WebPs: mark ≤24 KiB, neutral/route ≤64 KiB each, delivery/admin ≤48 KiB each, and ≤248 KiB combined. It verifies that neutral is the canonical identity reference, role variants preserve the declared invariant physiology, the pinned Pixelify Sans WOFF2 is ≤48 KiB, the OFL notice and source revision match, and undeclared or legacy assets fail. Landing JavaScript remains below 145 KiB gzip and CSS below 9 KiB gzip. First-party TypeScript and TSX maintain exact 100% statements, branches, functions, and lines.

## Runtime and cleanup

Black-box tests cover local and SSH path directions, archive/extraction, no-op identity, collision preservation, interruption, confirmed-byte reporting, staging cleanup, and bounded streaming. Native SSH/SFTP tests use an in-process authenticated server.

Temporary files use `t.TempDir`/`t.Cleanup`. Distribution containers carry a unique Courier acceptance label and are removed by success, failure, and signal traps. Acceptance never prunes or deletes unrelated Docker resources.

## CI and publication

CI runs the release-candidate gate on Linux and compiled platform gates on macOS and Windows. It cross-builds all declared client/helper artifacts and installs packages in the documented distribution matrix. Tagged release publication, repository-owned Homebrew/Scoop manifests, GitHub Releases, and npm publishing remain blocked until platform and release-candidate gates succeed.

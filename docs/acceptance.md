# Acceptance and release-candidate verification

Courier uses layered acceptance so the same transfer, archive, HTTP, worker, and reporting implementations are exercised without introducing a second test-only engine.

## Local gates

The fast runtime gate is:

```sh
make test
```

On macOS and Linux this runs the complete Go suite, the race detector, exact 100% first-party statement coverage, and smoke checks against a newly compiled `courier` binary. Windows runs the equivalent platform coverage and compiled-binary checks through `scripts/test-windows.ps1`.

The release-candidate gate is:

```sh
make verify
```

It additionally verifies the npm wrapper, exact TypeScript coverage, deterministic embedded UI assets, real Chromium behavior, GitHub Actions syntax, strict OpenSpec state, the complete GoReleaser snapshot, installers, artifact checksums, and native Linux packages. The CI package checks install the snapshot in Ubuntu, Debian, Arch Linux, Manjaro, Fedora, Red Hat UBI, and Alpine containers and execute the installed binary. The official Arch image is amd64-only, so an ARM Docker host verifies the arm64 Arch package through Manjaro while the required amd64 CI job performs the native Arch check.

Local prerequisites are Go 1.25 or newer, Node.js 24 or newer, Docker, OpenSpec 1.11.0, and Playwright's pinned Chromium runtime. Prepare the non-Go tools with:

```sh
npm install --global @fission-ai/openspec@1.11.0
npm ci --prefix web --ignore-scripts --no-audit --no-fund
cd web && npx playwright install chromium
```

On Linux, use `npx playwright install --with-deps chromium` when the host also needs Chromium system libraries. GoReleaser, actionlint, and ShellCheck are version-pinned by the repository. Go downloads actionlint on first use; checksum-verified GoReleaser and ShellCheck bootstraps store their binaries under ignored `.cache/tools`. This keeps embedded workflow-shell validation identical on developer macOS/Linux hosts and GitHub's Linux runner.

Native Windows acceptance uses the current-process pseudo token to scope named-pipe ACLs, assigns each pipe handshake a unique valid endpoint, bounds both sides of that handshake, and registers listener cleanup before connecting. It also converts rooted `fs.FS` directory names to slash form, converts local relative symlink targets to portable tar form, and checks POSIX permission bits only on filesystems that expose them. These adaptations do not weaken endpoint identity, archive traversal checks, or private-file policy on supported platforms.

## Runtime and browser boundaries

Black-box command tests compose the public Cobra tree with isolated local and mapped SSH backends. They cover all four path directions, Unicode and spaces, archive/extraction round trips, successful identity no-ops, collision preservation, interruption, confirmed-byte reporting, and staging cleanup. Native SSH/SFTP tests continue to use an in-process authenticated server so cryptography, host keys, SFTP, and helper behavior stay within the same Go gate.

Large input is generated as a synthetic stream. A deliberately blocked destination proves that reads and writes never exceed the configured buffer, cancellation releases backpressure, confirmed staging is reported, and neither a final path nor a Courier partial path survives.

Playwright starts the three Vite applications and drives a real Chromium instance. The suite covers the data UI, administration UI, and landing in English and Russian; system, light, and dark themes; persisted icon-only selectors; keyboard operation; natural four-section scrolling; measured-header anchor offsets and active navigation; wide/portrait source switching; hero-first and next-only image request order; on-demand fine-pointer refraction; coarse/reduced-motion fallbacks; masked panorama coverage; exact route/install/CLI Copy behavior; absence of Run/Replay, transcripts, and editable command surfaces; width-safe installation commands; narrow and wide overflow; external-request isolation; and protected-metadata non-disclosure. Browser artifacts are written below a private temporary directory and removed by exit and signal traps.

Application code remains subject to Vitest's exact 100% statements, branches, functions, and lines threshold; Playwright and generated bundles are not counted as first-party executable source. Asset validation caps every landing WebP at 225 KiB and their aggregate at 1.6 MiB. The deterministic production build caps landing JavaScript at 45 KiB gzip and CSS at 8 KiB gzip.

## Resource ownership and cleanup

Go tests register local cleanup immediately with `t.Cleanup` or `t.TempDir`. Distribution containers receive a unique `com.iwonz.courier.acceptance=<uuid>` label and use `--rm`. Success, failure, and signal traps remove and then query only resources carrying that exact label. The acceptance scripts never prune Docker and never inspect or delete unrelated containers, networks, volumes, or files.

## CI and publication

CI uses Linux for the complete release-candidate gate, macOS for compiled runtime and exact platform coverage, and Windows for compiled runtime, exact platform coverage, and PowerShell installer acceptance. The Linux snapshot cross-builds every declared macOS, Linux, Windows, and BSD helper artifact before package installation begins.

The tagged release workflow repeats all three platform gates. GitHub Release creation, in-repository Homebrew/Scoop manifest updates, and npm publication cannot start unless Linux browser/package/workflow/OpenSpec acceptance and the macOS/Windows jobs all succeed.

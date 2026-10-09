# cross-platform-acceptance Specification

## Purpose
Define the release-candidate acceptance layers for runtime routes, bounded streams, browsers, platforms, distribution packages, and owned-resource cleanup.

## Requirements

### Requirement: Route and transformation acceptance

Courier SHALL exercise all four filesystem directions, browser/webhook route handlers, archive creation/extraction, no-op, collisions, interruption, path safety, Unicode, special paths, partial reporting, and cleanup through isolated acceptance fixtures.

#### Scenario: An acceptance operation is interrupted

- **WHEN** a transfer is canceled after staging data
- **THEN** Courier returns exit code 130, leaves the previous final namespace unchanged, and removes only its owned partial resources

### Requirement: Bounded performance acceptance

Courier SHALL verify large synthetic streams with bounded buffers, observable backpressure, cancellation, and no payload-sized memory allocation.

#### Scenario: A slow destination blocks

- **WHEN** a source produces data faster than the destination accepts it
- **THEN** transfer reads remain bounded by the configured buffer and cancellation releases the operation and its staging resources

### Requirement: Cross-platform release acceptance

Courier SHALL run compiled runtime suites on Linux, macOS, and Windows, cross-build the declared primary and BSD helper artifacts, and install matching release packages in isolated Ubuntu, Debian, Arch Linux, Manjaro, Fedora, RHEL-compatible, and Alpine containers. Native Windows acceptance SHALL exercise owner-scoped named pipes using the current-process identity, a unique valid endpoint per test run, bounded connection contexts, cleanup registered before the handshake, rooted filesystem operations across the native-path and `fs.FS` path dialects, portable archive symlink targets, and platform-supported metadata semantics.

#### Scenario: A distribution package check exits

- **WHEN** installation succeeds, fails, or is interrupted
- **THEN** its cleanup trap removes every uniquely labeled Courier container and proves that no labeled container, network, volume, or temporary file remains

#### Scenario: Windows uses a rooted directory path

- **WHEN** a native Windows path produced by `filepath` crosses into an `fs.FS` operation
- **THEN** the adapter supplies slash-separated syntax without weakening the rooted confinement boundary

#### Scenario: Windows creates a private control pipe

- **WHEN** Courier resolves the named-pipe security principal
- **THEN** it queries the supported current-process token and grants pipe access only to that user SID

#### Scenario: Windows pipe acceptance stalls

- **WHEN** a native named-pipe test cannot complete its client-server handshake
- **THEN** its bounded context terminates the attempt and pre-registered cleanup releases the listener and accept goroutine without waiting for the package timeout

### Requirement: Real-browser acceptance

Courier SHALL test landing, delivery, and administration in Chromium for the existing product behaviors plus a static normal-flow masthead, local official PNG third-party marks, aligned four-pixel route signaling, one unframed route-and-command workspace, exact POSIX and PowerShell clipboard output, and visual layouts at 320, 390, 1024, and 1440 pixels in light and dark themes. Landing acceptance SHALL verify two top-level sections, route-template prefill and command switching, and one direct-binaries link beside scrollable installation tabs.

#### Scenario: Landing polish is inspected

- **WHEN** browser acceptance measures the masthead, brand media, route signal, and command workspace
- **THEN** the header scrolls away, control centers align within one pixel, no third-party mark is pixelated or remotely loaded, the signal center follows the path, no CLI separator exists, and both shell commands copy exactly

#### Scenario: Reduced motion is requested

- **WHEN** the route signal renders under reduced motion
- **THEN** its four-pixel center remains at the measured connector midpoint without animation

#### Scenario: The landing first renders

- **WHEN** the initial landing viewport loads at a supported wide or portrait size
- **THEN** it presents the localized headline, Relay, and one route-and-command workspace without a separate route command, viewport-sized section constraint, clipping, or unreadable artwork overlap

#### Scenario: Natural section geometry is inspected

- **WHEN** browser acceptance measures the unified route/command section and installation at every supported viewport
- **THEN** both bounds are content-driven, follow one another without scroll snapping or hard border seams, and no stylesheet rule imposes `svh`-based section height

#### Scenario: Static artwork is inspected

- **WHEN** the page loads or a pointer crosses the hero composition
- **THEN** the local Relay mascot and compact mark remain stable, installation allocates no decorative raster, and no external image, font, script, or analytics request occurs

#### Scenario: Preferences are activated

- **WHEN** theme and locale buttons are activated by pointer or keyboard and the page reloads
- **THEN** each control advances one cyclic value, the locale icon identifies the active locale, accessible names announce current and next values, and the valid preference persists without a radiogroup

#### Scenario: Product applications are exercised

- **WHEN** browser acceptance authenticates, transfers data, receives administration snapshots, changes policy, or stops a target
- **THEN** the lightweight React/shadcn applications preserve existing guarded API behavior, secret isolation, authoritative state, and current selection where valid

#### Scenario: Protected delivery returns an authentication error

- **WHEN** the data page receives an unauthenticated response containing an adversarial secret marker
- **THEN** the marker is not rendered in either locale and the accessible generic authentication state is shown

#### Scenario: A visitor uses landing command surfaces

- **WHEN** route templates, installation, or generated CLI Copy are activated in English or Russian
- **THEN** the exact immutable or shell-quoted command is offered to the clipboard and localized status is announced without simulation, API mutation, or layout shift

#### Scenario: A visitor uses natural navigation

- **WHEN** the visitor scrolls between the route/command workspace and installation
- **THEN** the page does not snap, the static masthead follows document flow, and reduced motion removes smooth scrolling

#### Scenario: A visitor crosses a section boundary

- **WHEN** the visitor moves from the unified route/command section to installation
- **THEN** no panorama or bordered section boundary exists and the shared document canvas remains continuous without an image seam or layout shift

#### Scenario: A visitor crosses a panorama boundary

- **WHEN** the visitor moves from the unified route/command section to installation
- **THEN** no panorama or bordered boundary exists and the canvas remains continuous without a seam or layout shift

#### Scenario: Browser surfaces render their primary workflows

- **WHEN** landing, delivery, and administration are inspected at supported desktop and mobile viewports
- **THEN** their primary content follows the continuous document canvas without outer card islands, clipping, horizontal overflow, or loss of interaction state

#### Scenario: Landing chrome and spacing are inspected

- **WHEN** the landing is measured before and after scrolling at supported desktop and mobile sizes
- **THEN** its masthead is static and transparent, content requires no header offset, adjacent content blocks have compact spacing, and the selected Source and Destination templates remain visibly contrasted

#### Scenario: A visitor configures a command

- **WHEN** Chromium selects a route template, edits endpoint values, switches CLI commands, and copies in either shell mode
- **THEN** only the active generated command is copied and the browser creates no horizontal document overflow

#### Scenario: A visitor installs or downloads

- **WHEN** Chromium scrolls the installation channel tabs on a narrow viewport
- **THEN** the direct-binaries link remains visible beside them and only the selected installation command is copied

### Requirement: Exact browser coverage

Courier SHALL include first-party `.ts` and `.tsx` sources in exact statement, branch, function, and line coverage gates.

#### Scenario: Unit coverage is measured

- **WHEN** shared UI, landing, delivery, and administration tests complete
- **THEN** every first-party TypeScript and TSX metric is exactly 100 percent without excluding React component sources

### Requirement: Pixel browser acceptance

Courier SHALL verify landing, delivery, and administration in English/Russian and system/light/dark at 320×568, 390×844, 1024-wide, and 1440×900. Acceptance SHALL assert exact terminal tokens, local Overpass Mono and Pixelify Sans scope, five local Relay v3 assets and correct per-surface resolution, unmodified local brand PNGs, square controls, one-pixel rules, two-pixel focus, absence of shadows/chamfers/card islands/external requests/overflow, centered route signal and reduced-motion midpoint, truthful contract/API data, and all existing Copy, builder, auth, upload/download, administration mutation, and preference flows. Landing JavaScript SHALL remain at most 150 KiB gzip and CSS at most 9 KiB gzip.

#### Scenario: Browser acceptance runs

- **WHEN** the production surfaces are exercised at every required locale, theme, motion setting, and viewport
- **THEN** visual structure, asset resolution, fonts, colors, overflow, interaction, API behavior, and bundle limits satisfy the terminal identity contract

#### Scenario: Browser surfaces render the pixel identity

- **WHEN** each production surface renders in a supported theme and viewport
- **THEN** the expected v3 assets, fonts, tokens, rules, square controls, truthful values, and no external requests or overflow are observed

#### Scenario: Product workflows are exercised

- **WHEN** Playwright performs route, Copy, builder, auth, upload/download, administration, and preference scenarios
- **THEN** the redesign preserves the exact established functional outcomes in English and Russian

### Requirement: Compiled hosted-route browser acceptance

Courier SHALL build the real executable before browser acceptance and SHALL serially exercise embedded hosted applications and HTTP routes in isolated homes with unique ports and mandatory process/resource cleanup.

#### Scenario: Browser uploads sequential files

- **WHEN** Playwright uploads Unicode text and binary files through a compiled `from web://` delivery
- **THEN** destination bytes match exactly, each committed name is present, and external stop makes the foreground command exit `0`

#### Scenario: Browser uploads an extraction archive

- **WHEN** Playwright uploads a safe tar.gz through `from web:// --extract`
- **THEN** the exact safe tree is committed without retaining the archive

#### Scenario: Browser downloads a directory

- **WHEN** Playwright opens `from <directory> to web://` and downloads the shared directory
- **THEN** the returned archive contains the expected nested byte-identical tree

#### Scenario: Administration runtime is exercised

- **WHEN** Playwright opens the embedded UI from a real foreground `ui start` and a second process runs `ui stop`
- **THEN** the page was reachable, the initiating process exits `0`, and administration state is removed

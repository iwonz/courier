## MODIFIED Requirements

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

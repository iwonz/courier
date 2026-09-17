## MODIFIED Requirements

### Requirement: Real-browser acceptance

Courier SHALL test the React/shadcn landing, delivery, and administration surfaces in Chromium for English/Russian catalogs, active-locale icons, cyclic system/light/dark preferences, reduced motion, touch and fine-pointer input, keyboard operation, responsive layouts, protected-metadata isolation, exactly three natural-height landing sections, local Relay rasters, external-request isolation, and layout stability. Acceptance SHALL verify the continuous route composition, unified CLI registry, continuous administration metrics/workspace, and the absence of unnecessary masthead and section-edge rules. The gate SHALL cap each Relay raster at 80 KiB, all Courier identity rasters at 160 KiB, landing JavaScript at 145 KiB gzip, and landing CSS at 9 KiB gzip.

#### Scenario: The landing first renders

- **WHEN** the initial landing viewport loads at a supported wide or portrait size
- **THEN** it presents the localized headline and route instrument as one composition without a full-page panorama, separate route card, viewport-sized section constraint, clipped content, or unreadable artwork overlap

#### Scenario: Natural section geometry is inspected

- **WHEN** browser acceptance measures hero, installation, and CLI at every supported viewport
- **THEN** their bounds are content-driven, they follow one another without scroll snapping or hard border seams, and no stylesheet rule imposes `svh`-based section height

#### Scenario: Static artwork is inspected

- **WHEN** the page loads or a pointer crosses the hero composition
- **THEN** the local Relay mascot and compact mark remain stable, later sections allocate no decorative raster, and no external image, font, script, or analytics request occurs

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

- **WHEN** route, installation, or CLI Copy is activated in English or Russian
- **THEN** the exact immutable command is offered to the clipboard and localized status is announced without simulation, API mutation, or layout shift

#### Scenario: A visitor uses natural navigation

- **WHEN** the visitor scrolls or activates installation or CLI navigation
- **THEN** the page does not snap, active navigation follows the current content block, and reduced motion removes smooth scrolling

#### Scenario: A visitor crosses a panorama boundary

- **WHEN** the visitor moves from one landing section to the next
- **THEN** no panorama or bordered section boundary exists and the shared document canvas remains continuous without an image seam or layout shift

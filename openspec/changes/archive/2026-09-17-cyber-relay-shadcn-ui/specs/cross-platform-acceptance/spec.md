## MODIFIED Requirements

### Requirement: Real-browser acceptance

Courier SHALL test the React/shadcn landing, delivery, and administration surfaces in Chromium for English/Russian catalogs, cyclic system/light/dark preferences, reduced motion, touch and fine-pointer input, keyboard operation, responsive layouts, protected-metadata isolation, exactly three natural-height landing sections, one compact local Relay raster, external-request isolation, and layout stability. The gate SHALL cap the Relay raster at 80 KiB, landing JavaScript at 145 KiB gzip, and landing CSS at 9 KiB gzip.

#### Scenario: The landing first renders

- **WHEN** the initial landing viewport loads at a supported wide or portrait size
- **THEN** it presents the localized headline and route instrument without a full-page panorama, viewport-sized section constraints, clipped content, or unreadable artwork overlap

#### Scenario: Natural section geometry is inspected

- **WHEN** browser acceptance measures hero, installation, and CLI at every supported viewport
- **THEN** their bounds are content-driven, they follow one another without scroll snapping, and no stylesheet rule imposes `svh`-based section height

#### Scenario: Static artwork is inspected

- **WHEN** the page loads or a pointer crosses the hero composition
- **THEN** one local Relay raster remains stable, later sections allocate no decorative raster, and no external image, font, script, or analytics request occurs

#### Scenario: Preferences are activated

- **WHEN** theme and locale buttons are activated by pointer or keyboard and the page reloads
- **THEN** each control advances one cyclic value, announces the current and next value, and persists the valid preference without a radiogroup

#### Scenario: Product applications are exercised

- **WHEN** browser acceptance authenticates, transfers data, receives administration snapshots, changes policy, or stops a target
- **THEN** the React/shadcn applications preserve existing guarded API behavior, secret isolation, authoritative state, and current selection where valid

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
- **THEN** no panorama boundary exists and the shared document canvas remains continuous without an image seam or layout shift

## ADDED Requirements

### Requirement: Exact browser coverage

Courier SHALL include first-party `.ts` and `.tsx` sources in exact statement, branch, function, and line coverage gates.

#### Scenario: Unit coverage is measured

- **WHEN** shared UI, landing, delivery, and administration tests complete
- **THEN** every first-party TypeScript and TSX metric is exactly 100 percent without excluding React component sources

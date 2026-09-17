## MODIFIED Requirements

### Requirement: Real-browser acceptance

Courier SHALL test landing, delivery, and administration in Chromium for English/Russian catalogs, cyclic system/light/dark preferences, reduced motion, touch and fine-pointer input, keyboard operation, responsive layouts, protected-metadata isolation, three-section natural landing flow, one responsive full-page panorama, raster identity use, heading contrast, external-request isolation, and layout stability. The gate SHALL enforce a maximum 420 KiB per selected landing panorama, 700 KiB for both landing panorama sources, 96 KiB for the compact raster mark, 200 KiB per delivery/admin scene, 45 KiB gzip landing JavaScript, and 9 KiB gzip landing CSS.

#### Scenario: The landing first renders

- **WHEN** the initial landing viewport loads at a supported wide or portrait size
- **THEN** it presents the localized Vector headline and route instrument, contains exactly three sections, requests exactly one selected full-page panorama, and retains readable heading contrast without segmented scene elements

#### Scenario: Preferences are activated

- **WHEN** theme and locale buttons are activated by pointer or keyboard and the page reloads
- **THEN** each control advances one cyclic value, announces the current and next value, and persists the valid preference without a radiogroup

#### Scenario: Product applications are exercised

- **WHEN** browser acceptance authenticates, transfers data, receives administration snapshots, changes policy, or stops a target
- **THEN** the softened workbench UI preserves existing guarded API behavior, secret isolation, and authoritative state

#### Scenario: Static artwork is inspected

- **WHEN** the pointer moves across the panorama in any motion mode
- **THEN** the base image remains single and stationary with no pointer listener outcome, duplicate image, refraction, or external request

#### Scenario: Protected delivery returns an authentication error

- **WHEN** the data page receives an unauthenticated response containing an adversarial secret marker
- **THEN** the marker is not rendered in either locale and the accessible generic authentication state is shown

#### Scenario: A visitor uses natural navigation

- **WHEN** the visitor scrolls or activates installation or CLI navigation
- **THEN** the page does not snap, the masthead identifies the section at its reading band, and reduced motion removes smooth scrolling

#### Scenario: A visitor crosses a panorama boundary

- **WHEN** adjacent landing sections are visible together at a supported viewport
- **THEN** one uninterrupted responsive panorama remains visible with no seam, repeated stripe, geometry shift, or inactive-source request

#### Scenario: A visitor uses landing command surfaces

- **WHEN** route, installation, or CLI Copy is activated in English or Russian
- **THEN** the exact immutable command is offered to the clipboard and localized status is announced without simulation, API mutation, or layout shift

## MODIFIED Requirements

### Requirement: Real-browser acceptance

Courier SHALL test data UI, admin UI, and landing behavior in Chromium for English/Russian catalogs, system/light/dark preferences, reduced motion, touch and fine-pointer input, keyboard navigation, responsive viewports, protected-metadata non-disclosure, natural landing scrolling, responsive panorama loading, external-request isolation, and layout stability. The acceptance gate SHALL enforce a maximum of 225 KiB per landing WebP, 1.6 MiB for all eight landing panorama assets, 45 KiB gzip for landing JavaScript, and 8 KiB gzip for landing CSS.

#### Scenario: Protected delivery returns an authentication error

- **WHEN** the data page receives an unauthenticated response containing an adversarial secret marker
- **THEN** the marker is not rendered in either locale and the accessible generic authentication state is shown

#### Scenario: The landing first renders

- **WHEN** the initial landing viewport loads at a supported wide or portrait size
- **THEN** it contains no visible tagline, Run or Replay action, transcript stage, editable command, horizontal overflow, or external runtime request and it requests only the selected hero asset plus at most the next panorama segment

#### Scenario: A visitor uses natural navigation

- **WHEN** the visitor scrolls or activates route, installation, or CLI navigation
- **THEN** the page does not snap between sections, the masthead identifies the section at its header-aware reading band, and reduced-motion preference removes smooth scrolling

#### Scenario: A visitor crosses a panorama boundary

- **WHEN** adjacent landing sections are visible together at a supported viewport
- **THEN** their matching responsive segments overlap without uncovered background, visible tiling, geometry shift, or loading the inactive responsive source

#### Scenario: A visitor uses landing command surfaces

- **WHEN** route, installation, or CLI Copy is activated in English or Russian
- **THEN** the exact immutable command is offered to the clipboard and localized status is announced without a simulation, API mutation, or layout shift

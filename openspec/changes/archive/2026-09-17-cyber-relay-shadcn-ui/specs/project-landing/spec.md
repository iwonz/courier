## MODIFIED Requirements

### Requirement: Shared-kit static landing

Courier SHALL provide a responsive static React application built from repository-owned shadcn components with English fallback, Russian localization, accessible system/light/dark controls, reduced-motion support, and semantic keyboard navigation.

#### Scenario: A Russian-language browser visits for the first time

- **WHEN** the browser language resolves to Russian and no selector preference exists
- **THEN** the landing renders the Russian catalog while the same shipped contract data and installation commands remain available

### Requirement: Brand-led product narrative

The landing SHALL use React and shared shadcn components in exactly three naturally scrolling, content-height sections ordered as a combined hero and route instrument, installation, and combined shipped command/option reference. No section SHALL use viewport-relative minimum height or behave as a mandatory screen. The first section SHALL present `From here to anywhere.` in English and `Отсюда — куда угодно.` in Russian and contain the contract-backed Source-to-Destination interaction without a separate routing introduction, visible endpoint explanation, installation command, or execution simulation.

#### Scenario: A visitor reads the complete landing

- **WHEN** the visitor scrolls from route selection through installation and the CLI reference
- **THEN** each section occupies only its content-driven height, retains smooth anchor navigation, and exposes no snap behavior, artificial viewport filler, clipping, or horizontal overflow

#### Scenario: A new visitor opens Courier

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor sees the localized promise, Relay identity, Source and Destination controls, selected route, immutable command, Copy action, and applicable options in the first content-height section

#### Scenario: A new visitor moves through the landing

- **WHEN** the visitor scrolls through the landing
- **THEN** route, installation, and CLI content follow natural document flow without any section being expanded to a viewport target

#### Scenario: A new visitor arrives

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor can identify Courier, explore valid routes, inspect installation channels, inspect commands and options, and reach source or releases without encountering unshipped claims

#### Scenario: The first viewport is constrained

- **WHEN** the first section cannot fit within the initial viewport
- **THEN** it expands in normal document flow without clipping controls, commands, artwork, or horizontal overflow

### Requirement: Illustrated operational landing

The route hero SHALL compose the canonical transparent technological Relay pigeon with responsive code-native Source-to-Destination route geometry. It SHALL NOT use a generated scene, cyberpunk city, page background, panorama, or responsive duplicate. The mascot SHALL be decorative, stable, local, and free of text, fake UI, logos, credentials, and required information.

#### Scenario: The landing loads

- **WHEN** a supported browser opens the landing
- **THEN** it loads one transparent Relay source while semantic shadcn controls and code-native geometry present every route and command function

#### Scenario: Artwork is unavailable

- **WHEN** the Relay mascot cannot load
- **THEN** the headline, route choices, command, options, and every interaction remain complete and readable

#### Scenario: The first section loads

- **WHEN** a supported browser opens the landing
- **THEN** it requests exactly one local transparent Relay source and no generated wide, portrait, or background scene

#### Scenario: A pointer moves over artwork

- **WHEN** a visitor moves any pointer across the Relay mascot
- **THEN** the artwork remains stable and decorative without translation, refraction, duplicate images, or animation work

#### Scenario: A visitor moves across a landing scene

- **WHEN** any fine pointer moves over the first section
- **THEN** functional state responds only to explicit semantic controls

#### Scenario: Motion or precise pointer input is unavailable

- **WHEN** reduced motion is requested or the device uses coarse pointer input
- **THEN** Relay remains readable and route signaling becomes static without alternate image behavior

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor clicks or keyboard-navigates the first section
- **THEN** only semantic route controls respond while Relay remains inert

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the route hero through installation and CLI reference
- **THEN** code-native color fields continue the identity without a full-page image, image seam, repeated horizon, uncovered region, or hard poster boundary

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** the locally versioned transparent Relay composition remains legible in light and dark GitHub themes

### Requirement: Contract-backed interactive route illustration

Courier SHALL project the shipped route matrix into compact Source and Destination endpoint controls integrated with the first-section Relay composition, a measured cubic Bezier connector, a reduced-motion-safe route signal, an immutable generated command, and applicable options. Endpoint descriptions SHALL remain available to assistive technology but SHALL not render as visible explanatory copy. Route endpoints, installation channels, command rows, compatible-option filtering, copy actions, and external links SHALL use shared shadcn controls while remaining generated from the CLI contract and installation registry.

#### Scenario: A route is selected

- **WHEN** a visitor activates a valid Source or Destination
- **THEN** the selected controls, connector, command, and applicable options update from generated contract data without layout movement

#### Scenario: A visitor compares endpoint roles

- **WHEN** the route instrument renders its choices
- **THEN** both sides show Local, Remote, Web, and Web Hook with directional iconography and assistive role descriptions while visible copy remains limited to endpoint names

#### Scenario: A visitor previews an endpoint

- **WHEN** the visitor hovers or focuses an endpoint control
- **THEN** the control receives visual feedback without changing the selected route

#### Scenario: A visitor activates a source or destination

- **WHEN** the visitor clicks or keyboard-activates a valid endpoint
- **THEN** the selected pair, Bezier connector, example syntax, and applicable options update from the generated contract while bounds remain unchanged

#### Scenario: A visitor explores a source endpoint

- **WHEN** the visitor hovers or focuses a source endpoint
- **THEN** only affordance changes until activation highlights valid destinations and updates contract-backed details

#### Scenario: A visitor explores a destination endpoint

- **WHEN** the visitor hovers or focuses a valid destination
- **THEN** only affordance changes until activation updates the route and the control remains operable without hover

#### Scenario: A visitor hovers an endpoint

- **WHEN** a pointer or keyboard focus enters an endpoint
- **THEN** visual affordance changes without changing route selection or exposing visible explanatory prose

#### Scenario: A visitor operates the landing

- **WHEN** the visitor selects a route, install channel, command, or compatibility filter with pointer or keyboard
- **THEN** the relevant immutable command and options update without layout instability and Copy reports localized success or failure

## REMOVED Requirements

### Requirement: Terminal-composed illustrated landing

**Reason**: Courier no longer ships a responsive generated hero scene or terminal composition.

**Migration**: Compose the single compact transparent Relay raster with semantic shadcn controls and code-native route geometry.

### Requirement: Proximity-loaded landing panorama

**Reason**: The landing has no panorama, responsive scene pair, or below-fold decorative image loading.

**Migration**: Load the single compact Relay raster with stable dimensions in the first section only.

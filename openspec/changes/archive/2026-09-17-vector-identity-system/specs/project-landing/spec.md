## MODIFIED Requirements

### Requirement: Brand-led product narrative

The landing SHALL consist of exactly three naturally scrolling sections ordered as a combined hero and route instrument, installation, and combined shipped command/option reference. The first section SHALL present `From here to anywhere.` in English and `Отсюда — куда угодно.` in Russian, occupy at least one small viewport on capable desktop layouts, expand without clipping on constrained layouts, and contain the contract-backed Source-to-Destination interaction without a separate routing introduction, visible endpoint explanation, installation command, or execution simulation.

#### Scenario: A new visitor opens Courier

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor sees the localized promise, Vector identity, Source and Destination controls, selected route, immutable command, Copy action, and applicable options in the first section

#### Scenario: A new visitor moves through the landing

- **WHEN** the landing loads at a supported viewport and the visitor scrolls through its sections
- **THEN** the combined hero and route instrument occupies at least the initial capable desktop viewport, subsequent content follows natural document flow, and every claim remains backed by shipped contract or repository-owned distribution data

#### Scenario: A new visitor arrives

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor can identify Courier from the localized headline and Vector route instrument, explore valid routes, inspect installation channels, inspect commands and options, and reach source or releases without encountering unshipped claims

#### Scenario: The first viewport is constrained

- **WHEN** the first section cannot fit safely within one viewport
- **THEN** it expands in natural flow without clipping controls, command content, or horizontal overflow

### Requirement: Contract-backed interactive route illustration

Courier SHALL project the shipped route matrix into compact Source and Destination endpoint controls integrated with the first-section artwork, a measured cubic Bezier connector, a reduced-motion-safe route signal, an immutable generated command, and applicable options. Endpoint descriptions SHALL remain available to assistive technology but SHALL not render as visible explanatory copy. Selection SHALL change only through activation and invalid destinations SHALL use disabled button semantics.

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

### Requirement: Compact landing composition

The landing SHALL omit a separate hero tagline, route section, route heading, route introduction, section indices, endpoint prose, examples, documentation, footer, execution preview, and simulated transcript. Its fixed masthead SHALL link the Courier identity to the first section and expose only Installation and CLI in-page navigation plus GitHub, theme, and locale actions.

#### Scenario: A visitor scans the page

- **WHEN** the visitor moves from the combined first section through installation and CLI reference
- **THEN** exactly three sections form one continuous composition with no removed standalone content or duplicated route introduction

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor scrolls through the landing
- **THEN** the combined route hero, installation, and CLI reference follow one another without snapping, clipping, disconnected tiles, or removed standalone content

### Requirement: Illustrated operational landing

The landing SHALL present one optional Vector journey as three separately loadable wide and portrait aerospace-editorial segments with matching transition geometry, carbon/chalk composition, and sparse burnt-orange route signals. Artwork SHALL remain static, contain no required information, and use no pointer-following, refraction, duplicated scene image, remote asset, text, fake UI, logo, credential, or watermark.

#### Scenario: The first section loads

- **WHEN** a supported browser opens the landing
- **THEN** it requests the selected first-section source eagerly and at most the next section within the preload boundary while farther artwork remains unrequested

#### Scenario: A pointer moves over artwork

- **WHEN** a visitor moves any pointer across a scene
- **THEN** no scene image translates, refracts, duplicates, or schedules pointer animation work

#### Scenario: A visitor moves across a landing scene

- **WHEN** any fine pointer moves over an active scene
- **THEN** the single static base illustration remains unchanged and functional state responds only to explicit controls

#### Scenario: Motion or precise pointer input is unavailable

- **WHEN** reduced motion is requested or the device uses coarse pointer input
- **THEN** artwork remains readable and route signaling becomes static without allocating alternate scene behavior

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor clicks or keyboard-navigates the first section
- **THEN** only semantic route controls respond while decorative Vector artwork remains inert

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the combined route hero through installation and CLI reference
- **THEN** the responsive panorama depicts route selection, verified distribution, and arrival as one optional environment with no uncovered boundary or remote asset request

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** the locally versioned text-free Vector route banner remains legible in light and dark GitHub themes

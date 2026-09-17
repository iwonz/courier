## MODIFIED Requirements

### Requirement: Terminal-composed illustrated landing

The project landing SHALL use one purpose-composed wide/portrait inline hero pair inside the first content-height section. Installation and CLI SHALL remain image-free, and functional route, command, and reference information SHALL be rendered as semantic interface content rather than embedded artwork. The hero artwork SHALL remain non-interactive.

#### Scenario: A fine pointer crosses a scene

- **WHEN** it moves over the inline hero artwork
- **THEN** the static image remains unchanged without pointer tracking, refraction, translation, functional state change, or continuous work

### Requirement: Transparent masthead integration

Courier SHALL keep the safe-area-aware measured masthead permanently visible using a compact translucent canvas surface, backdrop blur, shared control chrome, and no lower rule. Its measured height SHALL drive exact section anchor offsets.

#### Scenario: A visitor scrolls across contrasting scenes

- **WHEN** any natural-height section reaches the masthead
- **THEN** navigation and actions remain legible while the masthead remains fully within the visual viewport and does not expose a strip of the preceding section at the target anchor

### Requirement: Brand-led product narrative

The landing SHALL consist of exactly three naturally scrolling, content-height sections ordered as a combined hero and route instrument, installation, and combined shipped command/option reference. No section SHALL use viewport-relative minimum height or behave as a mandatory screen. The first section SHALL present `From here to anywhere.` in English and `Отсюда — куда угодно.` in Russian and contain the contract-backed Source-to-Destination interaction without a separate routing introduction, visible endpoint explanation, installation command, or execution simulation.

#### Scenario: A visitor reads the complete landing

- **WHEN** the visitor scrolls from the route hero through installation and CLI reference
- **THEN** every section occupies only the height required by its content and responsive spacing without snap behavior, viewport-height filler, clipping, or artificial blank regions

#### Scenario: A new visitor opens Courier

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor sees the localized promise, swift identity, Source and Destination controls, selected route, immutable command, Copy action, and applicable options in the first content-height section

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

The landing SHALL use a repository-local responsive wide/portrait hero illustration inside the first section rather than a full-page or per-section background. The artwork SHALL present the ImageGen-authored swift courier in a vivid editorial route composition, remain decorative, contain no required information, and create no seam, repeated tile, pointer-following effect, remote request, text, fake UI, logo lettering, credential, or watermark.

#### Scenario: The landing loads

- **WHEN** a supported browser opens the landing
- **THEN** it requests only the browser-selected hero illustration and compact mark while installation and CLI remain image-free content sections

#### Scenario: Artwork is unavailable

- **WHEN** the hero illustration cannot load
- **THEN** the headline, route choices, command, options, and every interaction remain complete and readable

#### Scenario: The first section loads

- **WHEN** a supported browser opens the landing
- **THEN** it requests exactly one browser-selected local hero source eagerly while the alternate responsive source remains unrequested

#### Scenario: A pointer moves over artwork

- **WHEN** a visitor moves any pointer across the hero illustration
- **THEN** the artwork remains stable and decorative without translation, refraction, duplicate images, or animation work

#### Scenario: A visitor moves across a landing scene

- **WHEN** any fine pointer moves over the first section
- **THEN** functional state responds only to explicit semantic controls

#### Scenario: Motion or precise pointer input is unavailable

- **WHEN** reduced motion is requested or the device uses coarse pointer input
- **THEN** artwork remains readable and route signaling becomes static without alternate image behavior

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor clicks or keyboard-navigates the first section
- **THEN** only semantic route controls respond while decorative artwork remains inert

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the route hero through installation and CLI reference
- **THEN** local color fields continue the identity without a full-page image, image seam, repeated horizon, uncovered region, or hard poster boundary

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** the locally versioned text-free Courier banner remains legible in light and dark GitHub themes

### Requirement: Compact landing composition

The landing SHALL use one continuous document canvas with content-driven spacing, local color fields, and no full-page panorama. Its fixed masthead SHALL link the Courier identity to the first section and expose only Installation and CLI in-page navigation plus GitHub, theme, and locale actions.

#### Scenario: A visitor scans the page

- **WHEN** the visitor moves between the three sections
- **THEN** the sections follow one another as ordinary document blocks without snapping, full-screen sizing, image stacking, background seams, or disconnected poster composition

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor scrolls through the landing
- **THEN** the combined route hero, installation, and CLI reference follow one another without snapping, clipping, viewport fillers, or removed standalone content

### Requirement: Proximity-loaded landing panorama

Courier SHALL NOT render a full-page or segmented landing panorama. It SHALL render one responsive hero picture with stable dimensions inside the first section, eagerly request only the selected wide or portrait source, and keep later sections independent of image decoding.

#### Scenario: The landing first loads

- **WHEN** the first section becomes interactive
- **THEN** exactly one responsive hero source is requested and no full-page panorama or later-section artwork is allocated

#### Scenario: A visitor approaches a later section

- **WHEN** installation or CLI enters the viewport
- **THEN** no additional landing artwork is allocated or requested and section geometry remains stable

#### Scenario: Intersection observation is unavailable

- **WHEN** the browser lacks IntersectionObserver
- **THEN** the eager responsive hero remains available without observation or a change in accessible content order

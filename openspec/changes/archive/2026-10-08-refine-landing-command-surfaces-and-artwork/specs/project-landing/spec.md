## ADDED Requirements

### Requirement: Quiet command surfaces

The landing SHALL group the generated CLI command with its shell selector and present installation commands in compact, flat, low-contrast surface bands. Each Copy action SHALL appear immediately below its command and align to the inline start. The masthead SHALL share the continuous page canvas without an opaque background fill. These treatments SHALL retain square corners, visible focus, localized status feedback, command wrapping, and keyboard access.

#### Scenario: A visitor reads a generated command

- **WHEN** the command builder renders at narrow or wide width
- **THEN** shell selection, the generated command, and its Copy action remain in one compact surface with Copy directly below the command at the inline start

#### Scenario: A visitor chooses an installation channel

- **WHEN** the selected installation command changes
- **THEN** the exact command wraps within one compact surface and its Copy action remains immediately below it without excess horizontal padding

#### Scenario: Landing chrome is rendered

- **WHEN** the masthead and command surfaces appear in light or dark theme
- **THEN** the masthead has no opaque strip while the command surfaces remain distinguishable through a restrained theme-derived fill

## MODIFIED Requirements

### Requirement: Illustrated operational landing

The route hero SHALL compose one local transparent amorphous pixel-art journey containing the recognizable Relay pigeon, a cyan route, and symbolic local-folder, remote-server, browser-portal, and parcel-delivery destinations. The generated artwork SHALL remain decorative, inert, stable, and free of text, credentials, logos, fake interactive UI, an opaque rectangular background, or required information. Semantic shadcn controls and code-native Source-to-Destination geometry SHALL continue presenting every route and command function.

#### Scenario: The landing loads

- **WHEN** a supported browser opens the landing
- **THEN** it loads one transparent thematic Relay source while semantic shadcn controls and code-native geometry present every route and command function

#### Scenario: Artwork is unavailable

- **WHEN** the Relay illustration cannot load
- **THEN** the headline, route choices, command, options, and every interaction remain complete and readable

#### Scenario: The first section loads

- **WHEN** a supported browser opens the landing
- **THEN** it requests exactly one local transparent route illustration with an irregular silhouette and no full-width or opaque background scene

#### Scenario: A pointer moves over artwork

- **WHEN** a visitor moves any pointer across the Relay illustration
- **THEN** the artwork remains stable and decorative without translation, refraction, duplicate images, or animation work

#### Scenario: A visitor moves across a landing scene

- **WHEN** any fine pointer moves over the first section
- **THEN** functional state responds only to explicit semantic controls

#### Scenario: Motion or precise pointer input is unavailable

- **WHEN** reduced motion is requested or the device uses coarse pointer input
- **THEN** the illustration remains readable and route signaling becomes static without alternate image behavior

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor clicks or keyboard-navigates the first section
- **THEN** only semantic route controls respond while the illustration remains inert

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the route hero through installation and CLI reference
- **THEN** code-native color fields and the transparent amorphous illustration continue the identity without a page background, image seam, repeated horizon, uncovered region, or hard poster boundary

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** the locally versioned transparent Relay composition remains legible in light and dark GitHub themes

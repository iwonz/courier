## MODIFIED Requirements

### Requirement: Illustrated operational landing

The route hero SHALL compose one local transparent amorphous pixel-art journey containing the recognizable Relay pigeon, a cyan route, and symbolic local-folder, remote-server, browser-portal, and parcel-delivery destinations. The illustration SHALL preserve its intended terminal palette without a fading mask, occupy approximately half of the available wide-screen hero, scale responsively, and remain fully visible within the hero's vertical bounds. Narrow layouts SHALL place the complete square illustration below the localized headline, while wider layouts MAY place it near or partially behind the headline without reducing contrast or readability. The generated artwork SHALL remain decorative, inert, stable, and free of text, credentials, logos, fake interactive UI, an opaque rectangular background, or required information. Semantic shadcn controls and code-native Source-to-Destination geometry SHALL continue presenting every route and command function.

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

#### Scenario: The hero reflows

- **WHEN** the hero renders at a supported narrow or wide viewport in either theme
- **THEN** the illustration retains its full palette and complete vertical silhouette without clipping the headline, overflowing horizontally, or losing its lower edge

### Requirement: Quiet command surfaces

The landing SHALL group the generated CLI command with its shell selector and present installation commands in compact, flat, low-contrast surface bands. The shell selector SHALL retain a localized accessible name without rendering that name as a visible label. Shell selection, command text, and Copy SHALL share a consistent inline-start axis and even surface padding. Each Copy action SHALL appear immediately below its command as an icon-led text action whose pointer hover changes its link treatment without adding a background fill, and each enabled Copy action SHALL expose a pointer cursor. The masthead SHALL share the continuous page canvas without an opaque background fill. These treatments SHALL retain square corners, visible focus, localized status feedback, command wrapping, and keyboard access.

#### Scenario: A visitor reads a generated command

- **WHEN** the command builder renders at narrow or wide width
- **THEN** shell selection, the generated command, and its Copy action remain in one evenly padded surface on one inline-start axis with no visible shell-syntax label

#### Scenario: A visitor uses assistive technology

- **WHEN** the POSIX and PowerShell selector is exposed to the accessibility tree
- **THEN** it retains its localized shell-syntax name even though that name is not rendered visually

#### Scenario: A pointer hovers Copy

- **WHEN** a visitor hovers an enabled Copy action
- **THEN** its icon and text gain the link emphasis, its cursor is a pointer, and no button-shaped background fill appears

#### Scenario: A visitor chooses an installation channel

- **WHEN** the selected installation command changes
- **THEN** the exact command wraps within one evenly padded surface and its Copy action remains aligned immediately below it

#### Scenario: Landing chrome is rendered

- **WHEN** the masthead and command surfaces appear in light or dark theme
- **THEN** the masthead has no opaque strip while the command surfaces remain distinguishable through a restrained theme-derived fill

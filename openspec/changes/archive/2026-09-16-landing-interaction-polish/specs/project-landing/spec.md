## MODIFIED Requirements

### Requirement: Illustrated operational landing

The project landing and repository README SHALL use responsive, context-specific Relay scenes to orient users at major narrative moments without obscuring interactive route data, installation commands, or the combined CLI reference. Each landing scene SHALL fill its slide as a thematic background rather than appear inside a separate illustration tile. Landing roles SHALL provide separately composed wide and portrait sources selected for the active viewport. Landing backgrounds SHALL remain positionally stationary while a non-essential, reduced-motion-safe amorphous spotlight and restrained refractive lens MAY follow a fine pointer. The landing hero SHALL be non-interactive and SHALL depict a decorative handoff from `Source` to `Destination`.

#### Scenario: A visitor moves across a landing scene

- **WHEN** a fine pointer moves over any landing slide
- **THEN** a local irregular highlight and restrained lens MAY follow it without translating the base illustration, changing functional state, or becoming necessary to understand the page

#### Scenario: Motion or precise pointer input is unavailable

- **WHEN** reduced motion is requested or the device uses coarse pointer input
- **THEN** the scene remains readable with a fixed ambient treatment and without pointer-driven travel or refraction animation

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor moves a pointer across, clicks, or keyboard-navigates the hero illustration
- **THEN** only the non-essential local spotlight may respond to a fine pointer while the decorative `Source` to `Destination` handoff remains non-interactive and unchanged

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the landing hero through routing and installation
- **THEN** each stationary full-slide background depicts a distinct matching Courier scenario, frames its overlaid controls as part of the environment, selects a purpose-composed wide or portrait source, preserves readable content hierarchy, and remains optional to understanding the section

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** a locally versioned, full-width, text-free Courier route panorama introduces the project and remains legible in light and dark GitHub themes

### Requirement: Contract-backed interactive route illustration

Courier SHALL project the shipped route matrix into one keyboard- and pointer-operable `from <source> to <destination>` illustration without maintaining a second route capability list. `Source` and `Destination` SHALL identify endpoint roles, while `Local`, `SSH`, `Web`, `Webhook`, and `HTTP(S)` SHALL identify endpoint types. Selection SHALL change only through activation, and a measured cubic Bezier connector SHALL join the selected endpoints without affecting layout.

#### Scenario: A visitor previews an endpoint

- **WHEN** the visitor hovers or focuses an endpoint control
- **THEN** the control receives visual feedback without changing the selected route

#### Scenario: A visitor activates a source or destination

- **WHEN** the visitor clicks or keyboard-activates a valid endpoint
- **THEN** the selected pair, Bezier connector, example syntax, and applicable options update from the generated contract while panel bounds remain unchanged

#### Scenario: A visitor explores a source endpoint

- **WHEN** the visitor hovers or focuses a source endpoint
- **THEN** the control indicates its affordance without changing the selected route, and only activation highlights its valid destinations and updates contract-backed details

#### Scenario: A visitor explores a destination endpoint

- **WHEN** the visitor hovers or focuses a valid destination
- **THEN** the control indicates its affordance without changing the selected route, and only activation updates the route while remaining operable without hover

### Requirement: Viewport-aware installation chooser

Courier SHALL present every supported installation channel through an icon-led interactive chooser with one exact, width-safe command readout and compact repository-owned package/download actions. Selection SHALL change only through activation, official third-party mark geometry SHALL be bundled locally and rendered in Courier monochrome, and every channel SHALL occupy invariant panel geometry.

#### Scenario: A user inspects an installation channel

- **WHEN** the user clicks or keyboard-activates a channel
- **THEN** the section exposes that channel's complete command without document overflow, hover-dependent state, layout movement, or a remote asset request

### Requirement: Minimal masthead composition

Courier SHALL align section navigation beside the product identity and group GitHub, theme, and locale actions at the opposite edge without bottom rules or separators between actions. Courier SHALL keep the landing masthead fully visible across scroll positions, browser safe areas, supported zoom levels, locales, and viewports. Its actual rendered height SHALL determine slide spacing, its background SHALL use a light scene-blending glass treatment in every theme, and its navigation SHALL identify the currently visible working section.

#### Scenario: A visitor scrolls through snapped sections

- **WHEN** any landing section aligns with the viewport
- **THEN** the complete masthead remains inside the visual viewport without clipping or document overflow and the matching route, installation, or CLI link exposes `aria-current`

#### Scenario: The masthead is rendered

- **WHEN** the landing loads at desktop or mobile width
- **THEN** route, installation, and CLI links follow page order on the left and the GitHub and preference controls remain reachable on the right without decorative divider lines

## ADDED Requirements

### Requirement: Full-width compatible CLI reference

Courier SHALL present generated commands and options in a full-width, invariant two-column reference. Commands SHALL be selectable without displaying repeated product/system badges. A checked-by-default custom checkbox SHALL filter options by the selected command's generated flag registry; no command selection SHALL display every option, and commands without flags SHALL display a localized empty state.

#### Scenario: A visitor filters options by command

- **WHEN** the visitor selects a command while compatibility filtering is checked
- **THEN** only flags registered for that command are shown without changing reference bounds

#### Scenario: No command is selected

- **WHEN** the reference first opens or the selected command is activated again
- **THEN** all options are shown and the checked compatibility control is disabled until another command is selected

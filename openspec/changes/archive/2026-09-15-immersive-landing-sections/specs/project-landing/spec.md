## MODIFIED Requirements

### Requirement: Brand-led product narrative

The landing SHALL introduce Courier through a viewport-sized product-promise scene without installation commands or action buttons, then present routing, installation, and the combined shipped command/option registry as three viewport-sized sections in that order.

#### Scenario: A new visitor moves through the landing

- **WHEN** the landing loads at a supported viewport and the visitor advances through its sections
- **THEN** each scroll-aligned section presents one primary concept, the hero remains free of installation calls to action, and every product claim remains backed by shipped contract or repository-owned distribution data

#### Scenario: A new visitor arrives

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor can identify what Courier moves, explore valid local, remote, web, webhook, and HTTP(S) route endpoints, inspect every installation channel, inspect commands and options, and reach source or releases without encountering unshipped claims

### Requirement: Illustrated operational landing

The project landing and repository README SHALL use responsive, context-specific Relay scenes to orient users at major narrative moments without obscuring interactive route data, installation commands, or the combined CLI reference. Each landing scene SHALL fill its slide as a thematic background rather than appear inside a separate illustration tile. Landing roles SHALL provide separately composed wide and portrait sources selected for the active viewport. The landing hero scene SHALL blend into its canvas and offer optional pointer, activation, and keyboard-responsive route feedback; other slide scenes MAY provide non-essential reduced-motion-safe pointer parallax while functional controls retain their own state.

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor moves a pointer across, activates, or keyboard-operates the hero illustration
- **THEN** the scene responds with non-essential route feedback, remains understandable without the effect, and honors reduced-motion preferences

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the landing hero through routing and installation
- **THEN** each full-slide background depicts a distinct matching Courier scenario, frames its overlaid controls as part of the environment, selects a purpose-composed wide or portrait source, preserves readable content hierarchy, and remains optional to understanding the section

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** a locally versioned, full-width, text-free Courier route panorama introduces the project and remains legible in light and dark GitHub themes

### Requirement: Compact landing composition

The landing SHALL consist of exactly four viewport-sized sections ordered hero, route explorer, installation chooser, and combined CLI reference. It SHALL omit a separate route-card grid, safety narrative, statistics strip, examples section, documentation section, full-contract link, closing call to action, and footer while keeping GitHub source access in the masthead.

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor moves through the landing
- **THEN** scroll snapping aligns one primary section at a time, the masthead navigation follows route, installation, and CLI order, and no removed standalone content interrupts the sequence

## ADDED Requirements

### Requirement: Viewport-aware installation chooser

Courier SHALL present every supported installation channel through an icon-led interactive chooser with one exact, width-safe command readout and compact repository-owned package/download actions.

#### Scenario: A user inspects an installation channel

- **WHEN** the user hovers, focuses, or activates a channel at a supported viewport
- **THEN** the section exposes that channel's complete command without document overflow and remains fully operable without pointer hover

### Requirement: Minimal masthead composition

Courier SHALL align section navigation beside the product identity and group GitHub, theme, and locale actions at the opposite edge without bottom rules or separators between actions.

#### Scenario: The masthead is rendered

- **WHEN** the landing loads at desktop or mobile width
- **THEN** route, installation, and CLI links follow page order on the left and the GitHub and preference controls remain reachable on the right without decorative divider lines

## MODIFIED Requirements

### Requirement: Brand-led product narrative

The landing SHALL introduce Courier through a hero with a minimum height of one small viewport and no visible tagline, installation command, or action button, then present routing, installation, and the combined shipped command/option registry as naturally sized sections in that order. Routing and installation SHALL have a minimum height of `42rem` or `78svh`, whichever is greater for the active viewport, while the CLI reference SHALL have a minimum height of `52rem` or `88svh`, whichever is greater; every section SHALL expand when its content requires more space.

#### Scenario: A new visitor moves through the landing

- **WHEN** the landing loads at a supported viewport and the visitor scrolls through its sections
- **THEN** the hero occupies at least the initial viewport, subsequent content follows natural document flow without clipping, and every product claim remains backed by shipped contract or repository-owned distribution data

#### Scenario: A new visitor arrives

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor can identify Courier from its headline and decorative Source-to-Destination story without a visible tagline, explore valid endpoint routes, inspect every installation channel, inspect commands and options, and reach source or releases without encountering unshipped claims

### Requirement: Illustrated operational landing

The project landing and repository README SHALL use responsive, context-specific Relay scenes without obscuring interactive route data, installation commands, or the combined CLI reference. The landing SHALL present one continuous journey as four adjacent wide and portrait panorama segments with matching transition geometry, lighting, perspective, and palette; masked overlap SHALL prevent uncovered or visibly tiled section boundaries. Landing backgrounds SHALL remain positionally stationary while a non-essential, reduced-motion-safe amorphous spotlight and demand-created restrained refractive lens MAY follow a fine pointer. The landing hero SHALL remain non-interactive and SHALL depict a decorative handoff between crisp `Source` and `Destination` route points.

#### Scenario: A visitor moves across a landing scene

- **WHEN** a fine pointer first moves over an active landing scene
- **THEN** a local irregular highlight and on-demand restrained lens converge smoothly without snapping, translating the base illustration, changing functional state, or remaining allocated after the pointer returns to ambient

#### Scenario: Motion or precise pointer input is unavailable

- **WHEN** reduced motion is requested or the device uses coarse pointer input
- **THEN** the scene remains readable with a fixed ambient treatment and does not create a pointer-driven refracted image

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor moves a pointer across, clicks, or keyboard-navigates the hero illustration
- **THEN** only the non-essential local spotlight may respond to a fine pointer while the decorative `Source` to `Destination` handoff remains non-interactive and unchanged

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the landing hero through routing, installation, and CLI reference
- **THEN** the responsive panorama depicts departure, route selection, verified distribution, and destination control as one continuous optional environment with no uncovered boundary or remote asset request

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** the existing locally versioned, full-width, text-free Courier route panorama remains unchanged and legible in light and dark GitHub themes

### Requirement: Compact landing composition

The landing SHALL consist of exactly four sections ordered hero, route explorer, installation chooser, and combined CLI reference. It SHALL use natural document scrolling without mandatory scroll snapping and SHALL omit a separate route-card grid, safety narrative, statistics strip, examples section, documentation section, full-contract link, closing call to action, and footer while keeping GitHub source access in the masthead.

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor scrolls through the landing
- **THEN** each naturally sized section follows the previous section without snapping, clipping, or a disconnected visual tile and no removed standalone content interrupts the sequence

### Requirement: Minimal masthead composition

Courier SHALL align section navigation beside the product identity and group GitHub, theme, and locale actions at the opposite edge without bottom rules or separators. GitHub SHALL use the same shared control-frame dimensions, border, padding, radius, surface, hover, focus, and theme tokens as the preference controls while preserving link semantics. The safe-area-aware masthead SHALL remain fully visible and its measured height SHALL define section anchor offsets. Its active navigation SHALL follow a header-aware reading band through variable-height sections. Every link whose destination leaves the landing page SHALL open a new browsing context without retaining opener access; in-page navigation SHALL remain in the current context and SHALL scroll smoothly unless reduced motion is requested.

#### Scenario: A visitor scrolls through natural sections

- **WHEN** a landing section crosses the reading band below the masthead
- **THEN** the complete masthead remains inside the visual viewport and the matching route, installation, or CLI link exposes `aria-current`

#### Scenario: A visitor scrolls through snapped sections

- **WHEN** the visitor scrolls through the landing after mandatory snapping has been removed
- **THEN** the complete masthead remains inside the visual viewport and active navigation follows the header-aware reading band rather than a snap position

#### Scenario: The masthead is rendered

- **WHEN** the landing loads at desktop or mobile width
- **THEN** route, installation, and CLI links follow page order on the left and GitHub and preference controls share one visual frame grammar on the right without divider lines

#### Scenario: A visitor follows an off-landing link

- **WHEN** the visitor activates the GitHub action, a release link, or another absolute off-landing destination
- **THEN** the destination opens in a new tab or window with `noopener` and `noreferrer` protection while the landing remains open

#### Scenario: A visitor follows section navigation

- **WHEN** an in-page link is activated
- **THEN** its target is offset by the measured masthead height and scrolling is smooth unless the visitor requests reduced motion

### Requirement: Terminal-composed illustrated landing

The project landing SHALL use purpose-composed wide and portrait panorama segments for hero, routing, installation, and CLI reference. The stationary responsive backgrounds SHALL form one continuous Relay journey and reserve quiet space for transparent read-only command and reference surfaces. The hero SHALL remain non-interactive and SHALL depict a decorative handoff between circular `Source` and `Destination` terminals.

#### Scenario: A fine pointer crosses a scene

- **WHEN** it moves over the base, optional refracted, glow, or veil layers
- **THEN** the local effect follows smoothly without hit-test jitter, base translation, functional state change, or continuous work after convergence

## ADDED Requirements

### Requirement: Read-only landing command surfaces

Courier SHALL present the selected route and installation commands as immutable command readouts with localized heading, description, Copy action, reserved live copy status, details, and optional repository actions. The selected generated CLI command SHALL provide the same copy affordance when a command is selected. The landing SHALL expose no Run or Replay control, staged transcript, demo timer, editable command field, preview status, shell, installer, transfer, or browser execution surface.

#### Scenario: A visitor copies a landing command

- **WHEN** Copy is activated for a route, installation channel, or selected CLI command
- **THEN** Courier copies the exact immutable command and announces localized success or failure without moving the surrounding layout

#### Scenario: A visitor inspects a read-only command

- **WHEN** any landing command surface is rendered
- **THEN** it contains no Run or Replay action, transcript stage, editable field, or claim that a browser-side operation ran

### Requirement: Proximity-loaded landing panorama

Courier SHALL load the responsive hero panorama eagerly with high fetch priority, activate the following segment no earlier than one viewport before it is needed, and omit farther image elements and requests until their one-shot preload boundaries are crossed. Browsers without intersection observation SHALL receive native lazy-loaded responsive images. Below-fold sections SHALL use stable intrinsic placeholders and deferred rendering without changing accessible content order.

#### Scenario: The landing first loads

- **WHEN** the initial viewport becomes interactive
- **THEN** only the selected hero source and at most the next responsive panorama segment have been requested

#### Scenario: A visitor approaches a later section

- **WHEN** a non-eager scene enters its preload boundary
- **THEN** its matching wide or portrait source is activated once without shifting section geometry or requesting the alternate source

#### Scenario: Intersection observation is unavailable

- **WHEN** the browser lacks IntersectionObserver
- **THEN** non-hero scenes render responsive images with native lazy loading and remain fully usable

## REMOVED Requirements

### Requirement: Terminal route demonstration

**Reason**: Landing commands are read-only copy surfaces and no longer simulate execution or transfer transcripts.

**Migration**: Use the selected route's immutable command, route description, applicable options, and Copy action.

### Requirement: Terminal installation demonstration

**Reason**: Installation commands are read-only copy surfaces and no longer simulate resolution, download, verification, or installation.

**Migration**: Use the selected channel's immutable command, Copy action, and repository-owned release links.

### Requirement: Terminal CLI demonstration

**Reason**: CLI reference selection no longer simulates a help invocation.

**Migration**: Select a generated command to inspect compatible options and copy its immutable usage string.

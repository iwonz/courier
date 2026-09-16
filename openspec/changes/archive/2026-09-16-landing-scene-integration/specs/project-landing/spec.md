## MODIFIED Requirements

### Requirement: Illustrated operational landing

The project landing and repository README SHALL use responsive, context-specific Relay scenes to orient users at major narrative moments without obscuring interactive route data, installation commands, or the combined CLI reference. Each landing scene SHALL fill its slide as a thematic background and SHALL reserve purpose-composed quiet regions that visually integrate its live controls rather than placing artwork and panels in competing layers. Landing roles SHALL provide separately composed wide and portrait sources selected for the active viewport. Landing backgrounds SHALL remain positionally stationary while a non-essential, reduced-motion-safe amorphous spotlight and smoothly interpolated restrained refractive lens MAY follow a fine pointer. The landing hero SHALL be non-interactive and SHALL depict a decorative handoff between crisp, intentional `Source` and `Destination` route points.

#### Scenario: A visitor moves across a landing scene

- **WHEN** a fine pointer moves over any landing slide
- **THEN** a local irregular highlight and restrained lens converge smoothly on it without snapping, translating the base illustration, changing functional state, or becoming necessary to understand the page

#### Scenario: Motion or precise pointer input is unavailable

- **WHEN** reduced motion is requested or the device uses coarse pointer input
- **THEN** the scene remains readable with a fixed ambient treatment and without pointer-driven travel or refraction animation

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor moves a pointer across, clicks, or keyboard-navigates the hero illustration
- **THEN** only the non-essential local spotlight may respond to a fine pointer while the decorative `Source` to `Destination` handoff remains non-interactive and unchanged

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the landing hero through routing, installation, and CLI reference
- **THEN** each full-slide background frames its corresponding live controls through deliberate quiet space and edge detail so the illustration and interface read as one responsive composition

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** a locally versioned, full-width, text-free Courier route panorama introduces the project and remains legible in light and dark GitHub themes

### Requirement: Contract-backed interactive route illustration

Courier SHALL project the shipped route matrix into one keyboard- and pointer-operable `from <source> to <destination>` illustration without maintaining a second capability list. `Source` and `Destination` SHALL identify endpoint roles. The landing SHALL present the endpoint types as `Local`, `Remote`, `Web`, and `Web Hook`, mapping `Remote` to SSH, source `Web Hook` to incoming `webhook://`, and destination `Web Hook` to outgoing HTTP(S). Each role SHALL use semantic, directional iconography and concise role-specific explanation. Selection SHALL change only through activation, and a measured cubic Bezier connector SHALL join compact selected endpoint controls without affecting layout.

#### Scenario: A visitor compares endpoint roles

- **WHEN** the route explorer renders its source and destination choices
- **THEN** both sides show Local, Remote, Web, and Web Hook while Web and Web Hook communicate upload, serving, receiving, or sending according to the selected side

#### Scenario: A visitor previews an endpoint

- **WHEN** the visitor hovers or focuses an endpoint control
- **THEN** the control receives visual feedback without changing the selected route

#### Scenario: A visitor activates a source or destination

- **WHEN** the visitor clicks or keyboard-activates a valid compact endpoint control
- **THEN** the selected pair, Bezier connector, example syntax, and applicable options update from the generated contract while panel bounds remain unchanged

#### Scenario: A visitor explores a source endpoint

- **WHEN** the visitor hovers or focuses a source endpoint
- **THEN** the control indicates its affordance without changing the selected route, and only activation highlights its valid destinations and updates contract-backed details

#### Scenario: A visitor explores a destination endpoint

- **WHEN** the visitor hovers or focuses a valid destination
- **THEN** the control indicates its affordance without changing the selected route, and only activation updates the route while remaining operable without hover

### Requirement: Viewport-aware installation chooser

Courier SHALL present every supported installation channel through compact, icon-led tag controls with one exact, width-safe command readout, an accessible copy action with bounded status feedback, and compact repository-owned package/download actions. Selection SHALL change only through activation, official third-party mark geometry SHALL be bundled locally and rendered in Courier monochrome, and every channel SHALL occupy invariant panel geometry.

#### Scenario: A user inspects an installation channel

- **WHEN** the user clicks or keyboard-activates a channel tag
- **THEN** the section exposes that channel's complete command without expanding the tag grid, causing document overflow, depending on hover state, moving panel bounds, or requesting a remote asset

#### Scenario: A user copies an installation command

- **WHEN** the visitor activates the copy action and clipboard access succeeds or fails
- **THEN** Courier attempts to copy the exact visible command and announces a localized result without moving or resizing the installation surface

### Requirement: Minimal masthead composition

Courier SHALL align section navigation beside the product identity and group GitHub, theme, and locale actions at the opposite edge without bottom rules or separators between actions. Courier SHALL keep the landing masthead fully visible across scroll positions, browser safe areas, supported zoom levels, locales, and viewports. Its actual rendered height SHALL determine slide spacing, its background SHALL use dark scene-integrated translucent glass in every theme, and its navigation SHALL identify the currently visible working section. Every link whose destination leaves the landing page SHALL open a new browsing context without retaining opener access; in-page navigation SHALL remain in the current context.

#### Scenario: A visitor scrolls through snapped sections

- **WHEN** any landing section aligns with the viewport
- **THEN** the complete masthead remains inside the visual viewport without clipping or document overflow and the matching route, installation, or CLI link exposes `aria-current`

#### Scenario: The masthead is rendered

- **WHEN** the landing loads at desktop or mobile width
- **THEN** route, installation, and CLI links follow page order on the left and the GitHub and preference controls remain reachable on the right without decorative divider lines

#### Scenario: A visitor follows an off-landing link

- **WHEN** the visitor activates the GitHub action, a release link, or another absolute off-landing destination
- **THEN** the destination opens in a new tab or window with `noopener` and `noreferrer` protection while the landing remains open

## MODIFIED Requirements

### Requirement: Contract-backed interactive route illustration

Courier SHALL project the shipped route matrix into Source and Destination controls, a measured cubic Bezier connector, a four-pixel centered signal that moves linearly on the connector, a reduced-motion midpoint state, an immutable generated command, and applicable options. Endpoint descriptions SHALL remain available to assistive technology but SHALL not render as visible explanatory copy.

#### Scenario: Route signaling is inspected

- **WHEN** the selected route connector is rendered with or without reduced motion
- **THEN** the signal is exactly four by four pixels and its center remains on the connector, moving linearly in normal motion and resting at the midpoint for reduced motion

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
- **THEN** the selected pair, connector, example syntax, and applicable options update from the generated contract while bounds remain unchanged

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

### Requirement: Compact landing composition

The landing SHALL use one continuous document canvas with compact content-driven spacing. Its masthead SHALL be a static normal-flow block with one `h-16` horizontal row containing the brand on the left and GitHub, theme, and locale actions on the right. It SHALL expose no internal navigation, reserve no fixed-header document offset, and SHALL leave the viewport when the document scrolls. The CLI registry SHALL use responsive columns without outer or internal divider lines.

#### Scenario: Landing chrome is inspected

- **WHEN** the landing is measured before and after scrolling at a supported viewport
- **THEN** the masthead has static positioning, no navigation, aligned control centers within one pixel, no compensating main padding, and scrolls out of view

#### Scenario: CLI registry is inspected

- **WHEN** commands, parameters, and the generated command render in one or two columns
- **THEN** no vertical, horizontal, or readout separator line is visible

#### Scenario: A visitor scans the page

- **WHEN** the visitor moves between the three sections
- **THEN** sections and related content remain in ordinary document flow without card backgrounds, outer radii, nested islands, hard section seams, full-screen sizing, or disconnected poster composition

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor scrolls through the landing
- **THEN** the combined route hero, installation, and CLI reference remain visually continuous while preserving their semantic headings and independent interaction regions

#### Scenario: A visitor uses a command surface

- **WHEN** the route, installation, or CLI command changes
- **THEN** the immutable command, Copy action, status, and details remain aligned directly on the page without gaining a containing card

### Requirement: Viewport-aware installation chooser

Courier SHALL present every supported installation channel through compact controls and one exact command readout. Official third-party marks SHALL be transparent local PNG images in official geometry and color without pixel rendering; GitHub SHALL provide light and dark variants. npx and wget SHALL render as text because they have no independent official mark.

#### Scenario: Third-party branding is inspected

- **WHEN** GitHub, operating system, distribution, and package-manager identities render
- **THEN** every available official mark is a local `img` PNG with no external request, SVG DOM, crisp-edge hint, or pixelated rendering, while npx and wget contain no invented icon

#### Scenario: A user inspects an installation channel

- **WHEN** the user clicks or keyboard-activates a channel control
- **THEN** the section exposes that channel's complete command without document overflow, hover-dependent state, layout movement, or a remote asset request

#### Scenario: A user copies an installation command

- **WHEN** the visitor activates the copy action and clipboard access succeeds or fails
- **THEN** Courier attempts to copy the exact visible command and announces a localized result without moving or resizing the installation surface

### Requirement: Minimal masthead composition

Courier SHALL render the compact Relay mark beside the literal wordmark `COURIER CLI` on the left and GitHub, theme, and locale actions on the right of one transparent `h-16` row. It SHALL contain no internal navigation, use static normal-flow positioning, reserve no document offset, and leave the viewport when the page scrolls. External links SHALL still open a new browsing context without retaining opener access.

#### Scenario: The masthead is rendered

- **WHEN** the landing loads at a supported width
- **THEN** the brand and right-side controls have visual centers aligned within one pixel on a transparent static masthead with no navigation or divider

#### Scenario: A visitor scrolls through natural sections

- **WHEN** a landing section crosses the reading band
- **THEN** the masthead scrolls away with the document and does not cover content

#### Scenario: A visitor scrolls through snapped sections

- **WHEN** the visitor scrolls through the landing after mandatory snapping has been removed
- **THEN** the page follows natural document position without a snap target, sticky header, or active navigation state

#### Scenario: A visitor follows an off-landing link

- **WHEN** the visitor activates the GitHub action, a release link, or another absolute off-landing destination
- **THEN** the destination opens in a new tab or window with `noopener` and `noreferrer` protection while the landing remains open

#### Scenario: A visitor follows section navigation

- **WHEN** the in-page brand link is activated
- **THEN** its route target remains in the current context without a fixed-header offset

### Requirement: Transparent masthead integration

Courier SHALL keep the transparent masthead in normal document flow as one compact `h-16` row with shared control chrome and no fill, blur, lower rule, measured overlay, or section-anchor offset.

#### Scenario: A visitor scrolls across contrasting scenes

- **WHEN** any natural-height section scrolls toward the top of the viewport
- **THEN** the masthead leaves the viewport with the preceding document content and never overlays the section

### Requirement: Pixel landing composition

The landing SHALL retain exactly three natural-height borderless sections on one continuous canvas while applying the shared modern 8-bit grammar to first-party identity and controls. The static transparent masthead SHALL use the Relay pixel mark, display wordmark, official local GitHub PNG, theme action, and active local locale flag. The route hero SHALL use only the transparent route-v2 sprite, code-native pixel waypoints, a grid-snapped measured connector, a centered four-pixel signal with linear path motion and a reduced-motion midpoint, contract-backed endpoint controls, immutable command, Copy action, and applicable options. Installation SHALL use only declared local official brand PNGs, and CLI SHALL remain free of decorative raster art.

#### Scenario: A visitor scans and operates the page

- **WHEN** the visitor scrolls, changes route, installation channel, command, theme, or locale, or copies a command
- **THEN** the three compact sections retain their geometry and behavior without card islands, smooth decorative gradients, legacy artwork, emoji flags, external requests, or browser execution

### Requirement: Pixel masthead composition

Courier SHALL keep the masthead static and transparent while rendering the generated square pixel Relay mark beside a local Pixelify Sans `COURIER CLI` wordmark. It SHALL contain no section navigation. GitHub, theme, and locale actions SHALL remain aligned in the same row, keyboard-operable, and protected according to existing external-link rules.

#### Scenario: A visitor scrolls the landing

- **WHEN** the document scrolls at a supported viewport
- **THEN** the pixel identity and controls leave the viewport with the masthead, no background appears, and no reserved document offset is present

## ADDED Requirements

### Requirement: Structured landing command builder

Courier SHALL build copy-ready CLI commands from ordered generated contract fields rather than usage-string replacement. It SHALL support command arguments, boolean toggles, enum selects, scalar typed fields, ordered repeatable values, explicit conflicts and dependencies, and POSIX or PowerShell quoting without executing a command.

#### Scenario: A visitor constructs a command

- **WHEN** required arguments and explicitly selected parameter values are valid
- **THEN** the readout contains command path, ordered arguments, and ordered parameters with exact shell-safe quoting and Copy offers that exact value

#### Scenario: Builder input is incomplete

- **WHEN** a required argument, UUID, number, enum, or dependency is invalid
- **THEN** validation is exposed and Copy remains disabled

#### Scenario: A visitor changes commands

- **WHEN** another command is selected
- **THEN** argument values, parameter values, and prior copy status are cleared while the shell preference remains selected

## MODIFIED Requirements

### Requirement: Compact landing composition

The landing SHALL use one continuous document canvas with compact content-driven spacing and no full-page panorama. Adjacent sections SHALL NOT accumulate large top and bottom padding. The route headline, Relay illustration, endpoint controls, command, and options SHALL form one composition without a rounded or tinted outer container. Installation SHALL flow directly on the same canvas, and the CLI registry SHALL use transparent responsive columns with only a meaningful internal divider rather than an outer surface. The fixed masthead SHALL link the Courier identity to the first section and expose Installation and CLI navigation plus GitHub, theme, and locale actions.

#### Scenario: A visitor scans the page

- **WHEN** the visitor moves between the three sections
- **THEN** sections and related content remain in ordinary document flow without card backgrounds, outer radii, nested islands, hard section seams, full-screen sizing, or disconnected poster composition

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor scrolls through the landing
- **THEN** the combined route hero, installation, and CLI reference remain visually continuous while preserving their semantic headings and independent interaction regions

#### Scenario: A visitor uses a command surface

- **WHEN** the route, installation, or CLI command changes
- **THEN** the immutable command, Copy action, status, and details remain aligned directly on the page without gaining a containing card

### Requirement: Minimal masthead composition

Courier SHALL render the compact Relay mark beside the literal wordmark `COURIER CLI`, place Installation and CLI navigation directly after that brand cluster, and align GitHub, theme, and locale actions at the opposite edge. The masthead SHALL remain fixed to the visual viewport without a background fill, backdrop blur, lower rule, or separator, while document padding SHALL prevent it from covering content. GitHub and preference actions SHALL share the same quiet control treatment while preserving link and button semantics. Every link whose destination leaves the landing page SHALL open a new browsing context without retaining opener access; in-page navigation SHALL remain in the current context and SHALL scroll smoothly unless reduced motion is requested.

#### Scenario: The masthead is rendered

- **WHEN** the landing loads at desktop width
- **THEN** the brand and in-page navigation form one left-aligned cluster while GitHub and preference controls form the right-aligned cluster on a transparent fixed masthead without a full-width divider

#### Scenario: A visitor scrolls through natural sections

- **WHEN** a landing section crosses the reading band below the masthead
- **THEN** the complete masthead remains inside the visual viewport and the corresponding in-page link exposes `aria-current`

#### Scenario: A visitor scrolls through snapped sections

- **WHEN** the visitor scrolls through the landing after mandatory snapping has been removed
- **THEN** the complete masthead remains inside the visual viewport and navigation follows natural document position rather than a snap target

#### Scenario: A visitor follows an off-landing link

- **WHEN** the visitor activates the GitHub action, a release link, or another absolute off-landing destination
- **THEN** the destination opens in a new tab or window with `noopener` and `noreferrer` protection while the landing remains open

#### Scenario: A visitor follows section navigation

- **WHEN** an in-page link is activated
- **THEN** its target remains offset below the fixed masthead and scrolling is smooth unless the visitor requests reduced motion

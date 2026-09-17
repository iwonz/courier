## MODIFIED Requirements

### Requirement: Compact landing composition

The landing SHALL use one continuous document canvas with content-driven spacing, restrained tonal fields, and no full-page panorama. The route headline, Relay illustration, endpoint controls, command, and options SHALL form one composition rather than separate bordered islands. Installation SHALL not introduce section-edge rules, and the CLI registry SHALL use one responsive shared surface instead of separate outlined columns. The fixed masthead SHALL link the Courier identity to the first section and expose Installation and CLI navigation plus GitHub, theme, and locale actions.

#### Scenario: A visitor scans the page

- **WHEN** the visitor moves between the three sections
- **THEN** sections and related controls follow one another as ordinary document content with no nested card-on-card framing, hard section seam, full-screen sizing, image stacking, or disconnected poster composition

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor scrolls through the landing
- **THEN** the combined route hero, installation, and CLI reference remain visually continuous while preserving their semantic headings and independent interaction regions

### Requirement: Minimal masthead composition

Courier SHALL render the compact Relay mark beside the literal wordmark `COURIER CLI`, place Installation and CLI navigation directly after that brand cluster, and align GitHub, theme, and locale actions at the opposite edge without bottom rules or separators. GitHub and preference actions SHALL share the same quiet control treatment while preserving link and button semantics. Every link whose destination leaves the landing page SHALL open a new browsing context without retaining opener access; in-page navigation SHALL remain in the current context and SHALL scroll smoothly unless reduced motion is requested.

#### Scenario: The masthead is rendered

- **WHEN** the landing loads at desktop width
- **THEN** the brand and in-page navigation form one left-aligned cluster while GitHub and preference controls form the right-aligned cluster without a full-width divider

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
- **THEN** its target remains offset below the sticky masthead and scrolling is smooth unless the visitor requests reduced motion

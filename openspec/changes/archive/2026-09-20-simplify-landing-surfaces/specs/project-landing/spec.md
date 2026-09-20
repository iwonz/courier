## MODIFIED Requirements

### Requirement: Brand-led product narrative

The landing SHALL retain exactly three naturally scrolling sections: route console, installation console, and command builder. The route console SHALL combine the localized promise, Relay route-v3, interactive endpoints, immutable command, and applicable options without endpoint, route, command, telemetry, or monitoring counters.

#### Scenario: A visitor reads the complete landing

- **WHEN** the visitor scrolls through the page
- **THEN** the three content-height terminal sections retain all route, installation, and command-builder behavior without viewport filler or horizontal overflow

#### Scenario: A new visitor opens Courier

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor sees the localized promise, route-v3, endpoint controls, generated route command, Copy, and applicable options without a contract-count strip

#### Scenario: A new visitor moves through the landing

- **WHEN** the visitor scrolls from route through installation and the command builder
- **THEN** all three sections follow normal document flow and the static header scrolls away

#### Scenario: A new visitor arrives

- **WHEN** the first route console is visible
- **THEN** product promise, route identity, and primary route controls remain readable without counters or unsupported monitoring content

#### Scenario: The first viewport is constrained

- **WHEN** the landing is 320×568 or 390×844
- **THEN** the route console reflows without clipping, overflow, or inaccessible controls

### Requirement: Compact landing composition

The landing SHALL use one continuous document canvas with compact content-driven spacing. Its masthead SHALL be a static normal-flow block with one `h-16` horizontal row containing the brand on the left and GitHub, theme, and locale actions on the right. It SHALL expose no internal navigation, reserve no fixed-header document offset, and SHALL leave the viewport when the document scrolls. The masthead, top-level sections, and route-console regions SHALL not be separated by border rules. The CLI registry SHALL use responsive columns without outer, header, row, column, or readout divider lines.

#### Scenario: Landing chrome is inspected

- **WHEN** the landing is measured before and after scrolling at a supported viewport
- **THEN** the masthead has static positioning, no navigation, aligned control centers within one pixel, no lower rule, no compensating main padding, and scrolls out of view

#### Scenario: CLI registry is inspected

- **WHEN** commands, parameters, and the generated command render in one or two columns
- **THEN** no vertical, horizontal, outer, header, row, or readout separator line is visible while editable controls and focus indicators remain perceivable

#### Scenario: A visitor scans the page

- **WHEN** the visitor moves between the three sections
- **THEN** sections and route-console regions remain in ordinary document flow without boundary rules, outer radii, nested islands, hard section seams, full-screen sizing, or disconnected poster composition

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor scrolls through the landing
- **THEN** the combined route hero, installation, and CLI reference remain visually continuous while preserving their semantic headings and independent interaction regions

#### Scenario: A visitor uses a command surface

- **WHEN** the route, installation, or CLI command changes
- **THEN** the immutable command, Copy action, status, and details remain aligned directly on the page without gaining a containing card or separator rule

### Requirement: Viewport-aware installation chooser

Installation SHALL use keyboard-operable horizontally scrollable channel tabs and an unframed command area. The selected tab SHALL have an immediate stable active fill, unselected tabs SHALL have a calm hover state, and neither state SHALL animate or flash. The selected channel identity SHALL appear in the tab only and SHALL not be repeated below the command. Official local PNG marks remain unmodified and non-pixelated; npx and wget remain text-only.

#### Scenario: Installation is narrow

- **WHEN** installation renders at 320 or 390 pixels
- **THEN** every channel remains reachable through the tab row and the exact selected command remains readable and copyable without document overflow

#### Scenario: Third-party branding is inspected

- **WHEN** installation and GitHub marks render
- **THEN** each supported brand is a local official-color PNG without pixelated image rendering while npx and wget remain text-only

#### Scenario: A user inspects an installation channel

- **WHEN** a channel tab receives pointer, hover, or keyboard activation
- **THEN** its exact reviewed command replaces the prior selection with a stable active state, a distinct calm hover state, no animation or flash, and no duplicate channel label below the command

#### Scenario: A user copies an installation command

- **WHEN** Copy is activated for the selected channel
- **THEN** the exact immutable installation command is written to the clipboard and localized status is announced

### Requirement: Full-width compatible CLI reference

The command builder SHALL render as a flat responsive composition without an outer table frame or any header, row, column, or readout divider rules. Necessary editable-control boundaries and keyboard focus indicators SHALL remain available. All structured argument, option, dependency, conflict, validation, shell quoting, reset, and exact Copy behavior SHALL remain unchanged.

#### Scenario: A command is configured

- **WHEN** the visitor supplies valid arguments and parameters
- **THEN** the unframed builder produces and copies the same exact POSIX or PowerShell command from the canonical contract

#### Scenario: A visitor filters options by command

- **WHEN** a command is selected and compatible filtering is active
- **THEN** only its contract-declared arguments and parameters remain enabled without adding any registry divider

#### Scenario: No command is selected

- **WHEN** the builder has no active command
- **THEN** it explains the empty state, disables Copy, and preserves access to the contract command list without adding an empty-state rule

### Requirement: Pixel landing composition

The landing SHALL use terminal palette tokens, Overpass Mono UI type, Pixelify Sans wordmark/H1/H2, route-v3 and mark-v3 assets, a teal measured path, and an exactly four-by-four yellow signal whose center follows the path or rests at its midpoint under reduced motion.

#### Scenario: The route console renders

- **WHEN** the route endpoints are measured
- **THEN** Relay v3, the teal connector, and the centered yellow signal render without clipping or a contract-count strip in light or dark mode

#### Scenario: A visitor scans and operates the page

- **WHEN** the visitor selects endpoints, channels, commands, arguments, options, shell syntax, theme, or locale
- **THEN** every state uses the terminal palette and retains the existing semantic behavior without fictional or decorative data

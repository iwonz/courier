## MODIFIED Requirements

### Requirement: Brand-led product narrative

The landing SHALL retain exactly three naturally scrolling sections: route console, installation console, and command data-grid. The route console SHALL combine the localized promise, Relay route-v3, interactive endpoints, immutable command, and a contract-derived strip containing endpoint, route, and command counts. It SHALL NOT invent telemetry or monitoring values.

#### Scenario: A visitor reads the complete landing

- **WHEN** the visitor scrolls through the page
- **THEN** the three content-height terminal sections retain all route, installation, and command-builder behavior without viewport filler or horizontal overflow

#### Scenario: A new visitor opens Courier

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor sees the localized promise, route-v3, truthful contract metrics, endpoint controls, generated route command, Copy, and applicable options

#### Scenario: A new visitor moves through the landing

- **WHEN** the visitor scrolls from route through installation and the command builder
- **THEN** all three sections follow normal document flow and the static header scrolls away

#### Scenario: A new visitor arrives

- **WHEN** the first route console is visible
- **THEN** product promise, route identity, and primary route controls remain readable without unsupported monitoring content

#### Scenario: The first viewport is constrained

- **WHEN** the landing is 320×568 or 390×844
- **THEN** the route console reflows without clipping, overflow, or inaccessible controls

### Requirement: Viewport-aware installation chooser

Installation SHALL use keyboard-operable horizontally scrollable channel tabs and a ruled command panel. Official local PNG marks remain unmodified and non-pixelated; npx and wget remain text-only.

#### Scenario: Installation is narrow

- **WHEN** installation renders at 320 or 390 pixels
- **THEN** every channel remains reachable through the tab row and the exact selected command remains readable and copyable without document overflow

#### Scenario: Third-party branding is inspected

- **WHEN** installation and GitHub marks render
- **THEN** each supported brand is a local official-color PNG without pixelated image rendering while npx and wget remain text-only

#### Scenario: A user inspects an installation channel

- **WHEN** a channel tab receives pointer or keyboard activation
- **THEN** its exact reviewed command and matching channel identity replace the prior selection

#### Scenario: A user copies an installation command

- **WHEN** Copy is activated for the selected channel
- **THEN** the exact immutable installation command is written to the clipboard and localized status is announced

### Requirement: Full-width compatible CLI reference

The command builder SHALL render as a flat responsive data-grid with horizontal row rules, no outer table frame, and no vertical divider between command and parameter columns. All structured argument, option, dependency, conflict, validation, shell quoting, reset, and exact Copy behavior SHALL remain unchanged.

#### Scenario: A command is configured

- **WHEN** the visitor supplies valid arguments and parameters
- **THEN** the ruled grid produces and copies the same exact POSIX or PowerShell command from the canonical contract

#### Scenario: A visitor filters options by command

- **WHEN** a command is selected and compatible filtering is active
- **THEN** only its contract-declared arguments and parameters remain enabled without adding an inter-column divider

#### Scenario: No command is selected

- **WHEN** the builder has no active command
- **THEN** it explains the empty state, disables Copy, and preserves access to the contract command list

### Requirement: Pixel landing composition

The landing SHALL use terminal palette tokens, Overpass Mono UI type, Pixelify Sans wordmark/H1/H2, route-v3 and mark-v3 assets, a teal measured path, and an exactly four-by-four yellow signal whose center follows the path or rests at its midpoint under reduced motion.

#### Scenario: The route console renders

- **WHEN** the route endpoints are measured
- **THEN** Relay v3, the teal connector, the centered yellow signal, and contract-derived counts render without clipping in light or dark mode

#### Scenario: A visitor scans and operates the page

- **WHEN** the visitor selects endpoints, channels, commands, arguments, options, shell syntax, theme, or locale
- **THEN** every state uses the terminal palette and retains the existing semantic behavior without fictional data

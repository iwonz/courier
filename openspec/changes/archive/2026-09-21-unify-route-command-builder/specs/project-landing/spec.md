## MODIFIED Requirements

### Requirement: Brand-led product narrative

The landing SHALL render two naturally scrolling top-level sections: a hero containing one unified route-and-command workspace, followed by installation. The hero SHALL retain the localized promise and route-v3 Relay without a visible route-console eyebrow or duplicate command output. The nested command workspace SHALL retain the `#cli` anchor.

#### Scenario: A visitor opens the landing

- **WHEN** the first section renders
- **THEN** `courier from` is selected, the initial Source/Destination route template fills editable arguments, and exactly one generated command and Copy action are shown in the workspace

#### Scenario: A visitor reads the complete landing

- **WHEN** the visitor scrolls through the page
- **THEN** the two natural-height sections retain the route, command-builder, and installation workflows without viewport filler or horizontal overflow

#### Scenario: A new visitor opens Courier

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor sees the localized promise, route-v3, route templates, one generated CLI command and Copy, and installation below without a contract-count strip

#### Scenario: A new visitor moves through the landing

- **WHEN** the visitor scrolls from the unified hero/workspace to installation
- **THEN** both sections follow normal document flow and the static header scrolls away

#### Scenario: A new visitor arrives

- **WHEN** the first route/command section is visible
- **THEN** product promise, route identity, and the primary controls remain readable without counters or unsupported monitoring content

#### Scenario: The first viewport is constrained

- **WHEN** the landing is 320×568 or 390×844
- **THEN** the combined section reflows without clipping, overflow, or inaccessible controls

#### Scenario: A visitor switches commands

- **WHEN** a non-`from` command is selected
- **THEN** route controls are hidden while that command's arguments, options, validation, and one readout remain available; returning to `from` restores the saved route template examples

#### Scenario: A visitor selects a route template

- **WHEN** a valid Source or Destination is activated
- **THEN** `from` becomes active, both endpoint examples update, parameters and Copy status reset, and the single generated command reflects the selected template

#### Scenario: A visitor edits an endpoint

- **WHEN** a template-filled endpoint argument is edited manually
- **THEN** the edited value is used in the generated command without constraining it to the selected template's endpoint kind

### Requirement: Contract-backed interactive route illustration

Courier SHALL project the shipped route matrix into Source and Destination template controls, a measured cubic Bezier connector, and a four-pixel centered signal that moves linearly on the connector or rests at its midpoint for reduced motion. The templates SHALL populate editable `from` arguments in the unified builder, not create another command output. Endpoint descriptions SHALL remain available to assistive technology without visible explanatory copy.

#### Scenario: A route is selected

- **WHEN** a visitor activates a valid Source or Destination template
- **THEN** the selected controls, connector, editable endpoint examples, and single generated command update from contract data

#### Scenario: Route signaling is inspected

- **WHEN** the selected template connector renders with or without reduced motion
- **THEN** the signal is four by four pixels with its center on the path, moving linearly or resting at the midpoint respectively

#### Scenario: A different command is selected

- **WHEN** a command other than `from` is active
- **THEN** the route controls and connector are hidden until `from` is selected again

#### Scenario: A visitor compares endpoint roles

- **WHEN** the route templates render
- **THEN** Source and Destination present Local, Remote, Web, and Web Hook with directional icons and assistive descriptions while visible copy stays limited to endpoint names

#### Scenario: A visitor previews an endpoint

- **WHEN** an endpoint receives hover or focus
- **THEN** only its visual affordance changes, without changing the selected template or generated command

#### Scenario: A visitor activates a source or destination

- **WHEN** a valid endpoint is activated by click or keyboard
- **THEN** the selected pair, connector, editable examples, and generated command update without changing the workspace bounds unexpectedly

#### Scenario: A visitor explores a source endpoint

- **WHEN** a source endpoint receives hover or focus
- **THEN** selection stays unchanged until activation and valid destinations remain distinguishable

#### Scenario: A visitor explores a destination endpoint

- **WHEN** a valid destination receives hover or focus
- **THEN** selection stays unchanged until activation and the control remains operable without hover

#### Scenario: A visitor hovers an endpoint

- **WHEN** a pointer or keyboard focus enters an endpoint
- **THEN** affordance changes without adding visible explanatory prose

#### Scenario: A visitor operates the landing

- **WHEN** the visitor selects a route template, install channel, command, or compatibility filter with pointer or keyboard
- **THEN** the relevant single readout and parameters update without layout instability and Copy reports localized success or failure

### Requirement: Compact landing composition

The landing SHALL use one continuous document canvas with two compact content-driven top-level sections. Its masthead SHALL be a static normal-flow `h-16` row with brand left and GitHub/theme/locale right, no navigation, and no fixed-header offset. The hero and nested command workspace SHALL form one section before installation. No section, workspace, or command region SHALL add decorative borders, row/column dividers, outer radii, or card islands; editable controls and focus indicators remain visible.

#### Scenario: Landing chrome is inspected

- **WHEN** the landing is measured before and after scrolling
- **THEN** the aligned masthead scrolls away, the route/command workspace remains in the first section, and installation follows without a boundary rule or horizontal overflow

#### Scenario: Command workspace is inspected

- **WHEN** command selection, route templates, options, and readout render at narrow or wide width
- **THEN** they form one unframed workspace with one generated command and one Copy action

#### Scenario: CLI registry is inspected

- **WHEN** commands, parameters, and the generated command render at narrow or wide width
- **THEN** no outer, header, row, column, or readout separator appears while editable controls and focus indicators remain perceivable

#### Scenario: A visitor scans the page

- **WHEN** the visitor moves between the two top-level sections
- **THEN** content stays in normal flow without boundary rules, outer radii, card islands, or full-screen sizing

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor scrolls through the landing
- **THEN** the combined hero/CLI workspace and installation remain visually continuous with distinct semantic headings

#### Scenario: A visitor uses a command surface

- **WHEN** a route template, CLI command, or installation channel changes
- **THEN** the relevant immutable command, Copy action, and status remain aligned directly on the page without a containing card or separator

### Requirement: Read-only landing command surfaces

Courier SHALL present the generated CLI command and selected installation command as separate immutable readouts with localized Copy actions and status. The route templates SHALL update the CLI readout rather than expose a duplicate route readout. The landing SHALL expose no Run, Replay, staged transcript, demo timer, editable command output, installer execution, transfer execution, or browser command execution.

#### Scenario: A visitor copies a landing command

- **WHEN** Copy is activated in the builder or installation section
- **THEN** Courier copies that section's exact command and announces localized success or failure without layout movement

#### Scenario: A visitor inspects command surfaces

- **WHEN** the landing renders
- **THEN** the builder has one readout and Copy action, installation has one readout and Copy action, and no route-only readout or execution action exists

#### Scenario: A visitor inspects a read-only command

- **WHEN** either command readout is rendered
- **THEN** it contains no Run or Replay action, transcript stage, editable output field, or claim of browser-side execution

### Requirement: Viewport-aware installation chooser

Installation SHALL keep horizontally scrollable channel tabs and one immutable command readout. A single text link to the latest direct binaries SHALL remain outside the tab scroller to its right, including at narrow widths. The redundant Linux-packages link SHALL not render.

#### Scenario: Installation is narrow

- **WHEN** installation renders at 320 or 390 pixels
- **THEN** channels remain scrollable while the direct-binaries link stays visible to their right without document overflow

#### Scenario: A visitor follows a direct download

- **WHEN** the direct-binaries link is activated
- **THEN** it opens the current GitHub Releases latest destination in a new protected tab

#### Scenario: Third-party branding is inspected

- **WHEN** installation and GitHub marks render
- **THEN** supported brands use local official-color PNGs without pixelated sampling while npx and wget stay text-only

#### Scenario: A user inspects an installation channel

- **WHEN** a channel tab receives pointer, hover, or keyboard activation
- **THEN** its exact command replaces the prior selection with stable active and hover states, no flash, and no repeated channel label below

#### Scenario: A user copies an installation command

- **WHEN** Copy is activated for the selected channel
- **THEN** the exact immutable command is written to the clipboard and localized status is announced

### Requirement: Full-width compatible CLI reference

The contract-backed command workspace SHALL use one generated readout and Copy action, horizontal command selection on narrow screens and vertical selection on wide screens, and no decorative outer, row, column, or readout divider rules. Existing typed values, dependency/conflict handling, shell quoting, and validation SHALL remain intact.

#### Scenario: A command is configured

- **WHEN** a visitor enters arguments and explicitly enables parameters
- **THEN** the single readout emits and copies the exact POSIX or PowerShell command in contract order

#### Scenario: A visitor switches commands

- **WHEN** the selected command changes
- **THEN** its arguments, parameters, compatible-only filter, and Copy status reset, with route examples populated only for `from`

#### Scenario: A visitor filters options by command

- **WHEN** a command is selected and compatible filtering is active
- **THEN** only its contract-declared arguments and parameters remain enabled without registry dividers

#### Scenario: No command is selected

- **WHEN** the builder first loads or the active command is clicked again
- **THEN** `from` remains selected initially or the current selection persists, so an unselected empty state cannot interrupt the workspace

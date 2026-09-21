# project-landing Specification

## Purpose
Define the contract-backed, shared-kit Courier project site and its least-privilege publication as a repository-owned GitHub Pages artifact.

## Requirements

### Requirement: Shared-kit static landing

Courier SHALL provide a responsive static React application built from repository-owned shadcn components with English fallback, Russian localization, accessible system/light/dark controls, reduced-motion support, and semantic keyboard navigation.

#### Scenario: A Russian-language browser visits for the first time

- **WHEN** the browser language resolves to Russian and no selector preference exists
- **THEN** the landing renders the Russian catalog while the same shipped contract data and installation commands remain available

### Requirement: Installation and release guidance

The landing SHALL present supported one-command installers, npm-compatible clients, Homebrew, Scoop, direct GitHub Releases, and Linux package formats as compact icon-led channels without clipping their exact commands or embedding tokens or mutable release credentials.

#### Scenario: A user chooses a distribution channel

- **WHEN** the installation section is opened on any supported viewport
- **THEN** every repository-owned command or GitHub Release destination uses the available width, remains fully reachable by keyboard, and does not force document overflow

### Requirement: Repository-owned Pages deployment

Courier SHALL automatically publish the deterministic landing artifact after every push to `main` through GitHub Actions Pages deployment with least-privilege permissions, no external repository, and no `gh-pages` maintenance branch. Courier SHALL also provide one Make command that builds a clean synchronized `main`, dispatches the same workflow, waits for completion, and reports the public URL.

#### Scenario: Main passes the landing gate

- **WHEN** the Pages workflow builds a synchronized contract and tested production landing
- **THEN** it uploads and deploys only the static artifact through the GitHub Pages environment

#### Scenario: Main changes

- **WHEN** any revision is pushed to `main`
- **THEN** the Pages workflow verifies the synchronized contract, tests and builds the production landing, and deploys only its static artifact

#### Scenario: A maintainer explicitly republishes

- **WHEN** an authenticated maintainer runs the publication Make command from clean local `main` matching `origin/main`
- **THEN** Courier builds locally, enables workflow-based Pages if absent, dispatches the exact main revision, waits for the new run to succeed, and reports the Pages URL

#### Scenario: Local state is not publishable

- **WHEN** the publication command runs from a feature branch, dirty worktree, or main revision different from `origin/main`
- **THEN** it fails before changing Pages configuration or dispatching a workflow

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

### Requirement: Illustrated operational landing

The route hero SHALL compose the canonical transparent technological Relay pigeon with responsive code-native Source-to-Destination route geometry. It SHALL NOT use a generated scene, cyberpunk city, page background, panorama, or responsive duplicate. The mascot SHALL be decorative, stable, local, and free of text, fake UI, logos, credentials, and required information.

#### Scenario: The landing loads

- **WHEN** a supported browser opens the landing
- **THEN** it loads one transparent Relay source while semantic shadcn controls and code-native geometry present every route and command function

#### Scenario: Artwork is unavailable

- **WHEN** the Relay mascot cannot load
- **THEN** the headline, route choices, command, options, and every interaction remain complete and readable

#### Scenario: The first section loads

- **WHEN** a supported browser opens the landing
- **THEN** it requests exactly one local transparent Relay source and no generated wide, portrait, or background scene

#### Scenario: A pointer moves over artwork

- **WHEN** a visitor moves any pointer across the Relay mascot
- **THEN** the artwork remains stable and decorative without translation, refraction, duplicate images, or animation work

#### Scenario: A visitor moves across a landing scene

- **WHEN** any fine pointer moves over the first section
- **THEN** functional state responds only to explicit semantic controls

#### Scenario: Motion or precise pointer input is unavailable

- **WHEN** reduced motion is requested or the device uses coarse pointer input
- **THEN** Relay remains readable and route signaling becomes static without alternate image behavior

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor clicks or keyboard-navigates the first section
- **THEN** only semantic route controls respond while Relay remains inert

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the route hero through installation and CLI reference
- **THEN** code-native color fields continue the identity without a full-page image, image seam, repeated horizon, uncovered region, or hard poster boundary

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** the locally versioned transparent Relay composition remains legible in light and dark GitHub themes

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

### Requirement: Viewport-aware installation chooser

Installation SHALL keep horizontally scrollable channel tabs and one immutable command readout. A single direct-binaries link with a leading first-party download icon SHALL remain outside the tab scroller to its right, including at narrow widths. The redundant Linux-packages link SHALL not render. The npx tab SHALL show the official local npm mark beside its npx label; wget remains text-only.

#### Scenario: Installation is narrow

- **WHEN** installation renders at 320 or 390 pixels
- **THEN** channels remain scrollable while the direct-binaries link and its icon stay visible to their right without document overflow

#### Scenario: A visitor follows a direct download

- **WHEN** the direct-binaries link is activated
- **THEN** it opens the current GitHub Releases latest destination in a new protected tab

#### Scenario: Third-party branding is inspected

- **WHEN** installation and GitHub marks render
- **THEN** supported brands use local official-color PNGs without pixelated sampling, npx shares the npm mark, and wget stays text-only

#### Scenario: A user inspects an installation channel

- **WHEN** a channel tab receives pointer, hover, or keyboard activation
- **THEN** its exact command replaces the prior selection with stable active and hover states, no flash, and no repeated channel label below

#### Scenario: A user copies an installation command

- **WHEN** Copy is activated for the selected channel
- **THEN** the exact immutable command is written to the clipboard and localized status is announced

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

### Requirement: Transparent masthead integration

Courier SHALL keep the transparent masthead in normal document flow as one compact `h-16` row with shared control chrome and no fill, blur, lower rule, measured overlay, or section-anchor offset.

#### Scenario: A visitor scrolls across contrasting scenes

- **WHEN** any natural-height section scrolls toward the top of the viewport
- **THEN** the masthead leaves the viewport with the preceding document content and never overlays the section

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

### Requirement: Pixel landing composition

The landing SHALL use terminal palette tokens, Overpass Mono UI type, Pixelify Sans wordmark/H1/H2, route-v3 and mark-v3 assets, a teal measured path, and an exactly four-by-four yellow signal whose center follows the path or rests at its midpoint under reduced motion.

#### Scenario: The route console renders

- **WHEN** the route endpoints are measured
- **THEN** Relay v3, the teal connector, and the centered yellow signal render without clipping or a contract-count strip in light or dark mode

#### Scenario: A visitor scans and operates the page

- **WHEN** the visitor selects endpoints, channels, commands, arguments, options, shell syntax, theme, or locale
- **THEN** every state uses the terminal palette and retains the existing semantic behavior without fictional or decorative data

### Requirement: Pixel masthead composition

Courier SHALL keep the masthead static and transparent while rendering the generated square pixel Relay mark beside a local Pixelify Sans `COURIER CLI` wordmark. It SHALL contain no section navigation. GitHub, theme, and locale actions SHALL remain aligned in the same row, keyboard-operable, and protected according to existing external-link rules.

#### Scenario: A visitor scrolls the landing

- **WHEN** the document scrolls at a supported viewport
- **THEN** the pixel identity and controls leave the viewport with the masthead, no background appears, and no reserved document offset is present

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

### Requirement: Transparent labeled-only landing workspace

The landing SHALL place the headline, Relay, and unified command workspace directly on the theme canvas without a tinted wrapper or sprite-cell fill. Its command selector and both immutable readouts SHALL omit the redundant visible `Command` label while retaining one Copy action per readout. The masthead SHALL vertically center the `COURIER CLI` wordmark on the Relay mark's axis.

#### Scenario: A visitor scans the command workspace

- **WHEN** the hero, command selector, and command readout render in either theme
- **THEN** the hero and Relay cell have transparent backgrounds, the command form remains functional, and no `Command` heading precedes the selector or appears beside Copy

#### Scenario: A visitor scans the installation readout

- **WHEN** an installation channel is selected
- **THEN** its exact command and Copy action render without a repeated visible `Command` label

#### Scenario: A visitor scans the masthead

- **WHEN** the compact brand renders
- **THEN** the mark and wordmark share a horizontal center axis within one pixel

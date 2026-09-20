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

### Requirement: Transparent masthead integration

Courier SHALL keep the transparent masthead in normal document flow as one compact `h-16` row with shared control chrome and no fill, blur, lower rule, measured overlay, or section-anchor offset.

#### Scenario: A visitor scrolls across contrasting scenes

- **WHEN** any natural-height section scrolls toward the top of the viewport
- **THEN** the masthead leaves the viewport with the preceding document content and never overlays the section

### Requirement: Read-only landing command surfaces

Courier SHALL present the selected route and installation commands as immutable command readouts with localized heading, description, Copy action, reserved live copy status, details, and optional repository actions. The selected generated CLI command SHALL provide the same copy affordance when a command is selected. The landing SHALL expose no Run or Replay control, staged transcript, demo timer, editable command field, preview status, shell, installer, transfer, or browser execution surface.

#### Scenario: A visitor copies a landing command

- **WHEN** Copy is activated for a route, installation channel, or selected CLI command
- **THEN** Courier copies the exact immutable command and announces localized success or failure without moving the surrounding layout

#### Scenario: A visitor inspects a read-only command

- **WHEN** any landing command surface is rendered
- **THEN** it contains no Run or Replay action, transcript stage, editable field, or claim that a browser-side operation ran

### Requirement: Pixel landing composition

The landing SHALL use terminal palette tokens, Overpass Mono UI type, Pixelify Sans wordmark/H1/H2, route-v3 and mark-v3 assets, a teal measured path, and an exactly four-by-four yellow signal whose center follows the path or rests at its midpoint under reduced motion.

#### Scenario: The route console renders

- **WHEN** the route endpoints are measured
- **THEN** Relay v3, the teal connector, the centered yellow signal, and contract-derived counts render without clipping in light or dark mode

#### Scenario: A visitor scans and operates the page

- **WHEN** the visitor selects endpoints, channels, commands, arguments, options, shell syntax, theme, or locale
- **THEN** every state uses the terminal palette and retains the existing semantic behavior without fictional data

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

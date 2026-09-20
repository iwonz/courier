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

The landing SHALL use React and shared shadcn components in exactly three naturally scrolling, content-height sections ordered as a combined hero and route instrument, installation, and combined shipped command/option reference. No section SHALL use viewport-relative minimum height or behave as a mandatory screen. The first section SHALL present `From here to anywhere.` in English and `Отсюда — куда угодно.` in Russian and contain the contract-backed Source-to-Destination interaction without a separate routing introduction, visible endpoint explanation, installation command, or execution simulation.

#### Scenario: A visitor reads the complete landing

- **WHEN** the visitor scrolls from route selection through installation and the CLI reference
- **THEN** each section occupies only its content-driven height, retains smooth anchor navigation, and exposes no snap behavior, artificial viewport filler, clipping, or horizontal overflow

#### Scenario: A new visitor opens Courier

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor sees the localized promise, Relay identity, Source and Destination controls, selected route, immutable command, Copy action, and applicable options in the first content-height section

#### Scenario: A new visitor moves through the landing

- **WHEN** the visitor scrolls through the landing
- **THEN** route, installation, and CLI content follow natural document flow without any section being expanded to a viewport target

#### Scenario: A new visitor arrives

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor can identify Courier, explore valid routes, inspect installation channels, inspect commands and options, and reach source or releases without encountering unshipped claims

#### Scenario: The first viewport is constrained

- **WHEN** the first section cannot fit within the initial viewport
- **THEN** it expands in normal document flow without clipping controls, commands, artwork, or horizontal overflow

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

### Requirement: Full-width compatible CLI reference

Courier SHALL present generated commands and options in a full-width, invariant two-column reference. Commands SHALL be selectable without displaying repeated product/system badges. A checked-by-default custom checkbox SHALL filter options by the selected command's generated flag registry; no command selection SHALL display every option, and commands without flags SHALL display a localized empty state.

#### Scenario: A visitor filters options by command

- **WHEN** the visitor selects a command while compatibility filtering is checked
- **THEN** only flags registered for that command are shown without changing reference bounds

#### Scenario: No command is selected

- **WHEN** the reference first opens or the selected command is activated again
- **THEN** all options are shown and the checked compatibility control is disabled until another command is selected

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

The landing SHALL retain exactly three natural-height borderless sections on one continuous canvas while applying the shared modern 8-bit grammar to first-party identity and controls. The static transparent masthead SHALL use the Relay pixel mark, display wordmark, official local GitHub PNG, theme action, and active local locale flag. The route hero SHALL use only the transparent route-v2 sprite, code-native pixel waypoints, a grid-snapped measured connector, a centered four-pixel signal with linear path motion and a reduced-motion midpoint, contract-backed endpoint controls, immutable command, Copy action, and applicable options. Installation SHALL use only declared local official brand PNGs, and CLI SHALL remain free of decorative raster art.

#### Scenario: A visitor scans and operates the page

- **WHEN** the visitor scrolls, changes route, installation channel, command, theme, or locale, or copies a command
- **THEN** the three compact sections retain their geometry and behavior without card islands, smooth decorative gradients, legacy artwork, emoji flags, external requests, or browser execution

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

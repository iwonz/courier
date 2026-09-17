# project-landing Specification

## Purpose
Define the contract-backed, shared-kit Courier project site and its least-privilege publication as a repository-owned GitHub Pages artifact.

## Requirements

### Requirement: Shared-kit static landing

Courier SHALL provide a responsive static Lit application built from the shared UI kit with English fallback, Russian localization, accessible system/light/dark controls, reduced-motion support, and semantic keyboard navigation.

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

The landing SHALL consist of exactly three naturally scrolling, content-height sections ordered as a combined hero and route instrument, installation, and combined shipped command/option reference. No section SHALL use viewport-relative minimum height or behave as a mandatory screen. The first section SHALL present `From here to anywhere.` in English and `Отсюда — куда угодно.` in Russian and contain the contract-backed Source-to-Destination interaction without a separate routing introduction, visible endpoint explanation, installation command, or execution simulation.

#### Scenario: A visitor reads the complete landing

- **WHEN** the visitor scrolls from the route hero through installation and CLI reference
- **THEN** every section occupies only the height required by its content and responsive spacing without snap behavior, viewport-height filler, clipping, or artificial blank regions

#### Scenario: A new visitor opens Courier

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor sees the localized promise, swift identity, Source and Destination controls, selected route, immutable command, Copy action, and applicable options in the first content-height section

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

The landing SHALL use a repository-local responsive wide/portrait hero illustration inside the first section rather than a full-page or per-section background. The artwork SHALL present the ImageGen-authored swift courier in a vivid editorial route composition, remain decorative, contain no required information, and create no seam, repeated tile, pointer-following effect, remote request, text, fake UI, logo lettering, credential, or watermark.

#### Scenario: The landing loads

- **WHEN** a supported browser opens the landing
- **THEN** it requests only the browser-selected hero illustration and compact mark while installation and CLI remain image-free content sections

#### Scenario: Artwork is unavailable

- **WHEN** the hero illustration cannot load
- **THEN** the headline, route choices, command, options, and every interaction remain complete and readable

#### Scenario: The first section loads

- **WHEN** a supported browser opens the landing
- **THEN** it requests exactly one browser-selected local hero source eagerly while the alternate responsive source remains unrequested

#### Scenario: A pointer moves over artwork

- **WHEN** a visitor moves any pointer across the hero illustration
- **THEN** the artwork remains stable and decorative without translation, refraction, duplicate images, or animation work

#### Scenario: A visitor moves across a landing scene

- **WHEN** any fine pointer moves over the first section
- **THEN** functional state responds only to explicit semantic controls

#### Scenario: Motion or precise pointer input is unavailable

- **WHEN** reduced motion is requested or the device uses coarse pointer input
- **THEN** artwork remains readable and route signaling becomes static without alternate image behavior

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor clicks or keyboard-navigates the first section
- **THEN** only semantic route controls respond while decorative artwork remains inert

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the route hero through installation and CLI reference
- **THEN** local color fields continue the identity without a full-page image, image seam, repeated horizon, uncovered region, or hard poster boundary

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** the locally versioned text-free Courier banner remains legible in light and dark GitHub themes

### Requirement: Contract-backed interactive route illustration

Courier SHALL project the shipped route matrix into compact Source and Destination endpoint controls integrated with the first-section artwork, a measured cubic Bezier connector, a reduced-motion-safe route signal, an immutable generated command, and applicable options. Endpoint descriptions SHALL remain available to assistive technology but SHALL not render as visible explanatory copy. Selection SHALL change only through activation and invalid destinations SHALL use disabled button semantics.

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
- **THEN** the selected pair, Bezier connector, example syntax, and applicable options update from the generated contract while bounds remain unchanged

#### Scenario: A visitor explores a source endpoint

- **WHEN** the visitor hovers or focuses a source endpoint
- **THEN** only affordance changes until activation highlights valid destinations and updates contract-backed details

#### Scenario: A visitor explores a destination endpoint

- **WHEN** the visitor hovers or focuses a valid destination
- **THEN** only affordance changes until activation updates the route and the control remains operable without hover

#### Scenario: A visitor hovers an endpoint

- **WHEN** a pointer or keyboard focus enters an endpoint
- **THEN** visual affordance changes without changing route selection or exposing visible explanatory prose

### Requirement: Compact landing composition

The landing SHALL use one continuous document canvas with content-driven spacing, local color fields, and no full-page panorama. Its fixed masthead SHALL link the Courier identity to the first section and expose only Installation and CLI in-page navigation plus GitHub, theme, and locale actions.

#### Scenario: A visitor scans the page

- **WHEN** the visitor moves between the three sections
- **THEN** the sections follow one another as ordinary document blocks without snapping, full-screen sizing, image stacking, background seams, or disconnected poster composition

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor scrolls through the landing
- **THEN** the combined route hero, installation, and CLI reference follow one another without snapping, clipping, viewport fillers, or removed standalone content

### Requirement: Viewport-aware installation chooser

Courier SHALL present every supported installation channel through compact, icon-led tag controls with one exact, width-safe command readout, an accessible copy action with bounded status feedback, and compact repository-owned package/download actions. Selection SHALL change only through activation, official third-party mark geometry SHALL be bundled locally and rendered in Courier monochrome, and every channel SHALL occupy invariant panel geometry.

#### Scenario: A user inspects an installation channel

- **WHEN** the user clicks or keyboard-activates a channel tag
- **THEN** the section exposes that channel's complete command without document overflow, hover-dependent state, layout movement, or a remote asset request

#### Scenario: A user copies an installation command

- **WHEN** the visitor activates the copy action and clipboard access succeeds or fails
- **THEN** Courier attempts to copy the exact visible command and announces a localized result without moving or resizing the installation surface

### Requirement: Minimal masthead composition

Courier SHALL align section navigation beside the product identity and group GitHub, theme, and locale actions at the opposite edge without bottom rules or separators. GitHub SHALL use the same shared control-frame dimensions, border, padding, radius, surface, hover, focus, and theme tokens as the preference controls while preserving link semantics. The safe-area-aware masthead SHALL remain fully visible and its measured height SHALL define section anchor offsets. Its active navigation SHALL follow a header-aware reading band through variable-height sections. Every link whose destination leaves the landing page SHALL open a new browsing context without retaining opener access; in-page navigation SHALL remain in the current context and SHALL scroll smoothly unless reduced motion is requested.

#### Scenario: A visitor scrolls through natural sections

- **WHEN** a landing section crosses the reading band below the masthead
- **THEN** the complete masthead remains inside the visual viewport and the matching route, installation, or CLI link exposes `aria-current`

#### Scenario: A visitor scrolls through snapped sections

- **WHEN** the visitor scrolls through the landing after mandatory snapping has been removed
- **THEN** the complete masthead remains inside the visual viewport and active navigation follows the header-aware reading band rather than a snap position

#### Scenario: The masthead is rendered

- **WHEN** the landing loads at desktop or mobile width
- **THEN** route, installation, and CLI links follow page order on the left and GitHub and preference controls share one visual frame grammar on the right without divider lines

#### Scenario: A visitor follows an off-landing link

- **WHEN** the visitor activates the GitHub action, a release link, or another absolute off-landing destination
- **THEN** the destination opens in a new tab or window with `noopener` and `noreferrer` protection while the landing remains open

#### Scenario: A visitor follows section navigation

- **WHEN** an in-page link is activated
- **THEN** its target is offset by the measured masthead height and scrolling is smooth unless the visitor requests reduced motion

### Requirement: Full-width compatible CLI reference

Courier SHALL present generated commands and options in a full-width, invariant two-column reference. Commands SHALL be selectable without displaying repeated product/system badges. A checked-by-default custom checkbox SHALL filter options by the selected command's generated flag registry; no command selection SHALL display every option, and commands without flags SHALL display a localized empty state.

#### Scenario: A visitor filters options by command

- **WHEN** the visitor selects a command while compatibility filtering is checked
- **THEN** only flags registered for that command are shown without changing reference bounds

#### Scenario: No command is selected

- **WHEN** the reference first opens or the selected command is activated again
- **THEN** all options are shown and the checked compatibility control is disabled until another command is selected

### Requirement: Terminal-composed illustrated landing

The project landing SHALL use one purpose-composed wide/portrait inline hero pair inside the first content-height section. Installation and CLI SHALL remain image-free, and functional route, command, and reference information SHALL be rendered as semantic interface content rather than embedded artwork. The hero artwork SHALL remain non-interactive.

#### Scenario: A fine pointer crosses a scene

- **WHEN** it moves over the inline hero artwork
- **THEN** the static image remains unchanged without pointer tracking, refraction, translation, functional state change, or continuous work

### Requirement: Transparent masthead integration

Courier SHALL keep the safe-area-aware measured masthead permanently visible using a compact translucent canvas surface, backdrop blur, shared control chrome, and no lower rule. Its measured height SHALL drive exact section anchor offsets.

#### Scenario: A visitor scrolls across contrasting scenes

- **WHEN** any natural-height section reaches the masthead
- **THEN** navigation and actions remain legible while the masthead remains fully within the visual viewport and does not expose a strip of the preceding section at the target anchor

### Requirement: Read-only landing command surfaces

Courier SHALL present the selected route and installation commands as immutable command readouts with localized heading, description, Copy action, reserved live copy status, details, and optional repository actions. The selected generated CLI command SHALL provide the same copy affordance when a command is selected. The landing SHALL expose no Run or Replay control, staged transcript, demo timer, editable command field, preview status, shell, installer, transfer, or browser execution surface.

#### Scenario: A visitor copies a landing command

- **WHEN** Copy is activated for a route, installation channel, or selected CLI command
- **THEN** Courier copies the exact immutable command and announces localized success or failure without moving the surrounding layout

#### Scenario: A visitor inspects a read-only command

- **WHEN** any landing command surface is rendered
- **THEN** it contains no Run or Replay action, transcript stage, editable field, or claim that a browser-side operation ran

### Requirement: Proximity-loaded landing panorama

Courier SHALL NOT render a full-page or segmented landing panorama. It SHALL render one responsive hero picture with stable dimensions inside the first section, eagerly request only the selected wide or portrait source, and keep later sections independent of image decoding.

#### Scenario: The landing first loads

- **WHEN** the first section becomes interactive
- **THEN** exactly one responsive hero source is requested and no full-page panorama or later-section artwork is allocated

#### Scenario: A visitor approaches a later section

- **WHEN** installation or CLI enters the viewport
- **THEN** no additional landing artwork is allocated or requested and section geometry remains stable

#### Scenario: Intersection observation is unavailable

- **WHEN** the browser lacks IntersectionObserver
- **THEN** the eager responsive hero remains available without observation or a change in accessible content order

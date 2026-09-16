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

The landing SHALL introduce Courier through a viewport-sized product-promise scene without installation commands or action buttons, then present routing, installation, and the combined shipped command/option registry as three viewport-sized sections in that order.

#### Scenario: A new visitor moves through the landing

- **WHEN** the landing loads at a supported viewport and the visitor advances through its sections
- **THEN** each scroll-aligned section presents one primary concept, the hero remains free of installation calls to action, and every product claim remains backed by shipped contract or repository-owned distribution data

#### Scenario: A new visitor arrives

- **WHEN** the landing loads at a supported viewport
- **THEN** the visitor can identify what Courier moves, explore valid local, remote, web, webhook, and HTTP(S) route endpoints, inspect every installation channel, inspect commands and options, and reach source or releases without encountering unshipped claims

### Requirement: Illustrated operational landing

The project landing and repository README SHALL use responsive, context-specific Relay scenes to orient users at major narrative moments without obscuring interactive route data, installation commands, or the combined CLI reference. Each landing scene SHALL fill its slide as a thematic background and SHALL reserve purpose-composed quiet regions that visually integrate its live controls rather than placing artwork and panels in competing layers. Landing roles SHALL provide separately composed wide and portrait sources selected for the active viewport. Landing backgrounds SHALL remain positionally stationary while a non-essential, reduced-motion-safe amorphous spotlight and smoothly interpolated restrained refractive lens MAY follow a fine pointer. The landing hero SHALL be non-interactive and SHALL depict a decorative handoff between crisp, intentional `Source` and `Destination` route points.

#### Scenario: A visitor moves across a landing scene

- **WHEN** a fine pointer moves over any landing slide
- **THEN** a local irregular highlight and restrained lens converge smoothly on it without snapping, translating the base illustration, changing functional state, or becoming necessary to understand the page

#### Scenario: Motion or precise pointer input is unavailable

- **WHEN** reduced motion is requested or the device uses coarse pointer input
- **THEN** the scene remains readable with a fixed ambient treatment and without pointer-driven travel or refraction animation

#### Scenario: A visitor interacts with the hero scene

- **WHEN** the visitor moves a pointer across, clicks, or keyboard-navigates the hero illustration
- **THEN** only the non-essential local spotlight may respond to a fine pointer while the decorative `Source` to `Destination` handoff remains non-interactive and unchanged

#### Scenario: A visitor scans the project story

- **WHEN** the visitor moves from the landing hero through routing and installation
- **THEN** each stationary full-slide background depicts a distinct matching Courier scenario, frames its overlaid controls through deliberate quiet space and edge detail as part of one responsive environment, selects a purpose-composed wide or portrait source, preserves readable content hierarchy, and remains optional to understanding the section

#### Scenario: A visitor opens the repository README

- **WHEN** GitHub renders the README
- **THEN** a locally versioned, full-width, text-free Courier route panorama introduces the project and remains legible in light and dark GitHub themes

### Requirement: Contract-backed interactive route illustration

Courier SHALL project the shipped route matrix into one keyboard- and pointer-operable `from <source> to <destination>` illustration without maintaining a second route capability list. `Source` and `Destination` SHALL identify endpoint roles. The landing SHALL present endpoint types as `Local`, `Remote`, `Web`, and `Web Hook`, mapping `Remote` to SSH, source `Web Hook` to incoming `webhook://`, and destination `Web Hook` to outgoing HTTP(S). Each role SHALL use semantic, directional iconography and concise role-specific explanation. Selection SHALL change only through activation, and a measured cubic Bezier connector SHALL join the compact selected endpoint controls without affecting layout.

#### Scenario: A visitor compares endpoint roles

- **WHEN** the route explorer renders its source and destination choices
- **THEN** both sides show Local, Remote, Web, and Web Hook while Web and Web Hook communicate upload, serving, receiving, or sending according to the selected side

#### Scenario: A visitor previews an endpoint

- **WHEN** the visitor hovers or focuses an endpoint control
- **THEN** the control receives visual feedback without changing the selected route

#### Scenario: A visitor activates a source or destination

- **WHEN** the visitor clicks or keyboard-activates a valid endpoint
- **THEN** the selected pair, Bezier connector, example syntax, and applicable options update from the generated contract while panel bounds remain unchanged

#### Scenario: A visitor explores a source endpoint

- **WHEN** the visitor hovers or focuses a source endpoint
- **THEN** the control indicates its affordance without changing the selected route, and only activation highlights its valid destinations and updates contract-backed details

#### Scenario: A visitor explores a destination endpoint

- **WHEN** the visitor hovers or focuses a valid destination
- **THEN** the control indicates its affordance without changing the selected route, and only activation updates the route while remaining operable without hover

### Requirement: Compact landing composition

The landing SHALL consist of exactly four viewport-sized sections ordered hero, route explorer, installation chooser, and combined CLI reference. It SHALL omit a separate route-card grid, safety narrative, statistics strip, examples section, documentation section, full-contract link, closing call to action, and footer while keeping GitHub source access in the masthead.

#### Scenario: A visitor scans the complete page

- **WHEN** the visitor moves through the landing
- **THEN** scroll snapping aligns one primary section at a time, the masthead navigation follows route, installation, and CLI order, and no removed standalone content interrupts the sequence

### Requirement: Viewport-aware installation chooser

Courier SHALL present every supported installation channel through compact, icon-led tag controls with one exact, width-safe command readout, an accessible copy action with bounded status feedback, and compact repository-owned package/download actions. Selection SHALL change only through activation, official third-party mark geometry SHALL be bundled locally and rendered in Courier monochrome, and every channel SHALL occupy invariant panel geometry.

#### Scenario: A user inspects an installation channel

- **WHEN** the user clicks or keyboard-activates a channel tag
- **THEN** the section exposes that channel's complete command without document overflow, hover-dependent state, layout movement, or a remote asset request

#### Scenario: A user copies an installation command

- **WHEN** the visitor activates the copy action and clipboard access succeeds or fails
- **THEN** Courier attempts to copy the exact visible command and announces a localized result without moving or resizing the installation surface

### Requirement: Minimal masthead composition

Courier SHALL align section navigation beside the product identity and group GitHub, theme, and locale actions at the opposite edge without bottom rules or separators between actions. Courier SHALL keep the landing masthead fully visible across scroll positions, browser safe areas, supported zoom levels, locales, and viewports. Its actual rendered height SHALL determine slide spacing, its background SHALL use dark scene-integrated translucent glass in every theme, and its navigation SHALL identify the currently visible working section. Every link whose destination leaves the landing page SHALL open a new browsing context without retaining opener access; in-page navigation SHALL remain in the current context.

#### Scenario: A visitor scrolls through snapped sections

- **WHEN** any landing section aligns with the viewport
- **THEN** the complete masthead remains inside the visual viewport without clipping or document overflow and the matching route, installation, or CLI link exposes `aria-current`

#### Scenario: The masthead is rendered

- **WHEN** the landing loads at desktop or mobile width
- **THEN** route, installation, and CLI links follow page order on the left and the GitHub and preference controls remain reachable on the right without decorative divider lines

#### Scenario: A visitor follows an off-landing link

- **WHEN** the visitor activates the GitHub action, a release link, or another absolute off-landing destination
- **THEN** the destination opens in a new tab or window with `noopener` and `noreferrer` protection while the landing remains open

### Requirement: Full-width compatible CLI reference

Courier SHALL present generated commands and options in a full-width, invariant two-column reference. Commands SHALL be selectable without displaying repeated product/system badges. A checked-by-default custom checkbox SHALL filter options by the selected command's generated flag registry; no command selection SHALL display every option, and commands without flags SHALL display a localized empty state.

#### Scenario: A visitor filters options by command

- **WHEN** the visitor selects a command while compatibility filtering is checked
- **THEN** only flags registered for that command are shown without changing reference bounds

#### Scenario: No command is selected

- **WHEN** the reference first opens or the selected command is activated again
- **THEN** all options are shown and the checked compatibility control is disabled until another command is selected

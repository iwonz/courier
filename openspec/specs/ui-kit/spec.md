# ui-kit Specification

## Purpose
Define the single reusable visual, accessibility, theme, localization, component, and identity-asset foundation shared by every Courier browser surface.

## Requirements

### Requirement: Shared UI package

Courier SHALL provide one React and TypeScript UI package for delivery, administration, and the project landing, composed from repository-owned shadcn source, Radix primitives, Tailwind utilities, CVA variants, and a shared class-merging helper. Production browser code SHALL NOT retain Lit custom elements or copied component implementations between consumers.

#### Scenario: Consumer imports a control

- **WHEN** a Courier web application needs a button, card, badge, checkbox, input, select, tabs, tooltip, progress indicator, separator, scroll area, alert, command readout, brand, mascot, icon, or preference control
- **THEN** it imports the shared shadcn component and preserves native semantics, keyboard behavior, focus visibility, disabled state, localized labeling, and accessible state

#### Scenario: A browser application is built

- **WHEN** Vite compiles a Courier browser application
- **THEN** React, component styles, icons, system fonts, and identity artwork are bundled locally without a runtime request to shadcn, a CDN, or a remote asset host

### Requirement: Auditable identity assets

Courier SHALL retain exactly five transparent pixel Relay v3 WebP assets, one pinned local Pixelify Sans display font, four pinned local Overpass Mono Latin/Cyrillic 400/600 subsets totaling at most 48 KiB, their OFL notices, and the declared local transparent PNG third-party marks. Relay SHALL use a compact square-bodied silhouette with a flat stepped crown and angular right-angle pixel clusters in every role. Neutral-v3 SHALL be the identity authority for mark, route, delivery, and administration v3 assets. Provenance SHALL record role, dimensions, byte count, generation prompt or upstream source, lineage or revision, license, consumers, and SHA-256 digest. Legacy Relay v1/v2 files SHALL NOT ship. Existing per-sprite and combined budgets remain unchanged.

#### Scenario: Asset integrity check

- **WHEN** repository verification runs
- **THEN** undeclared assets, legacy Relay files, missing provenance or licenses, incorrect dimensions, broken v3 neutral lineage, altered third-party branding, excess font bytes, or an exceeded sprite budget fail the gate

#### Scenario: Identity assets are validated

- **WHEN** the asset gate inspects Relay, fonts, and third-party media
- **THEN** it verifies alpha, dimensions, hashes, neutral-v3 linkage, consumers, font licenses and sizes, official brand provenance, and the absence of third-party pixel rendering

### Requirement: User-selectable themes

Courier SHALL expose theme through one semantic React icon button that cycles `system`, `light`, and `dark`, defaults to browser-resolved `system`, persists an explicit valid preference when storage is available, meets accessible contrast, and respects reduced-motion preferences. The system choice SHALL continue following browser color-scheme changes until the user stores a different choice.

#### Scenario: Theme cycles from its default

- **WHEN** a visitor with no saved theme repeatedly activates the theme control
- **THEN** the shared preference provider cycles system to light to dark to system while system resolution continues following the browser media query

#### Scenario: Storage is unavailable

- **WHEN** reading or writing the theme preference fails
- **THEN** the current document state still updates without crashing or installing a duplicate listener

#### Scenario: System theme changes

- **WHEN** the stored preference is `system` and the operating-system color scheme changes
- **THEN** the resolved theme updates through the shared provider without replacing the stored preference

### Requirement: Extensible localization

Courier SHALL provide typed English and Russian React catalogs, use English as the fallback, negotiate the initial locale from browser languages, and expose locale through one semantic icon button that cycles English and Russian. A valid explicit choice SHALL persist when storage is available, and additional catalog modules SHALL require no backend changes.

#### Scenario: Locale cycles from its default

- **WHEN** a visitor activates the locale control
- **THEN** the shared preference provider alternates English and Russian, persists the valid choice when storage is available, and updates visible and accessible content once

#### Scenario: Unsupported browser locale

- **WHEN** no supported locale matches the browser language list
- **THEN** English is selected and every requested message key resolves

### Requirement: Exact UI source coverage

Courier SHALL enforce 100% statements, branches, functions, and lines for first-party TypeScript sources while excluding generated output and third-party code.

#### Scenario: Untested branch is introduced

- **WHEN** a first-party TypeScript branch is not executed by the test suite
- **THEN** the repository verification gate fails

### Requirement: Coherent brand identity system

Courier SHALL use one repository-local terminal 8-bit Relay v3 family across landing, delivery, administration, README, favicon, and product chrome. Every v3 role SHALL use the same compact square-bodied silhouette, flat stepped crown, angular feather planes, and right-angle pixel clusters while preserving the established anatomy, fundamental proportions, eye, beak, earpiece, satchel, pose, and working object. The family SHALL use teal in place of cobalt, yellow in place of orange/coral, and a cool black/white neutral plumage ramp. Assets SHALL retain true alpha and hard pixel clusters.

#### Scenario: A browser surface presents Courier

- **WHEN** landing, delivery, or administration renders
- **THEN** its header uses mark-v3 beside `COURIER CLI`, its optional v3 role sprite matches the surface, and product meaning remains complete without the image

#### Scenario: Relay changes roles

- **WHEN** all five v3 assets are compared on light, dark, and checkerboard backgrounds
- **THEN** the same recognizable square-bodied Relay identity, angular silhouette grammar, role poses, object count, dimensions, transparent background, and terminal palette remain visible

### Requirement: Non-essential mascot guidance

Courier SHALL use Relay only as restrained orientation and identity artwork, with locally verified provenance and without making meaning depend on the image.

#### Scenario: Mascot artwork is unavailable

- **WHEN** Relay artwork cannot be loaded or is hidden from assistive technology
- **THEN** headings, status text, controls, route geometry, and progress information still communicate the complete workflow

### Requirement: Calm operational communication

Courier SHALL use concise, concrete, non-alarmist language that identifies actions and verified states without unsupported guarantees, secret disclosure, or humor during failures.

#### Scenario: A delivery changes state

- **WHEN** a user sees preparation, transfer, verification, completion, or failure feedback
- **THEN** the message names the current or stopped operation and preserves an actionable, technically accurate tone in every supported locale

### Requirement: Intentional control appearance

Courier SHALL present all three browser surfaces through one terminal data-grid grammar. Dark mode SHALL use `#000000`, `#0D1015`, `#191C20`, white, 70% white secondary text, and 55% white rules. Light mode SHALL invert those roles with white and 4%/9%/70%/55% black. Action/success SHALL use `#71FFF6`, selection/warning `#FAD14F`, destructive/error `#C94A55`, and light-theme teal text `#006B67`. Bright teal and yellow fills SHALL use black text. Controls SHALL use square corners, one-pixel rules, two-pixel focus, and instant step states without chamfer, offset shadow, blur, or soft shadow. Pixelify Sans SHALL be limited to wordmark and H1/H2; Overpass Mono 400/600 SHALL serve interface text, labels, values, and commands.

#### Scenario: A user operates a Courier control

- **WHEN** a control is rendered, focused, selected, disabled, hovered, or pressed in any theme
- **THEN** it remains contrast-safe, keyboard-operable, square, sharply ruled, and free of soft or offset effects

#### Scenario: A product surface groups related content

- **WHEN** content is grouped without its own interactive state
- **THEN** it uses a flat surface band, typography, alignment, spacing, or a meaningful one-pixel rule instead of a detached card island

### Requirement: Compact icon preference controls

Courier SHALL render theme and locale selectors as separate semantic React icon buttons without visible labels or radiogroups. The locale button SHALL display an icon that identifies the currently active language. Each button SHALL expose localized current and next values, use shared quiet shadcn control styling, preserve keyboard activation and focus indication, and delegate persistence and document updates to one shared preference provider.

#### Scenario: A user operates an icon preference selector

- **WHEN** the user points to, focuses, or activates the theme or locale button
- **THEN** its current and next states are available programmatically and activation advances exactly once through the documented cycle

#### Scenario: A user operates a locale selector

- **WHEN** the visitor activates the locale button by pointer or keyboard
- **THEN** its visible icon changes from the current English flag to the current Russian flag or vice versa while its localized accessible name identifies both current and next locale

### Requirement: Local official brand marks

Courier SHALL provide a typed first-party pixel icon registry for functional endpoint, preference, file, transfer, status, and administration icons. GitHub, installation, platform, operating-system, distribution, shell, and package-manager identities SHALL use typed local transparent PNG marks in official geometry and color without pixel sampling. Every third-party raster SHALL retain pinned official source, revision, color, license, attribution, and trademark records. npx SHALL reuse the official npm mark rather than invent an independent npx brand; channels without an applicable official mark, including wget, SHALL remain text-only. Locale controls SHALL use local first-party pixel flags rather than platform emoji.

#### Scenario: A browser surface renders an icon

- **WHEN** Courier renders a functional icon, channel mark, or locale flag
- **THEN** it uses the declared local first-party registry or official brand PNG with an accessible name where needed and makes no runtime request

#### Scenario: A branded channel is rendered

- **WHEN** a supported package manager, shell, operating system, or distribution mark appears
- **THEN** the UI uses its transparent official-color PNG through an `img`, retains its pinned source record, avoids pixelated rendering, and exposes no remote asset URL

### Requirement: Shared styled checkbox

Courier SHALL provide a shared checkbox primitive with custom visual treatment, native checked and disabled semantics, localized accessible labeling, keyboard activation, focus indication, and a composed change event.

#### Scenario: A visitor toggles a view filter

- **WHEN** the checkbox is activated by pointer or keyboard
- **THEN** its checked state and accessible state update once and consumers receive the new boolean value

### Requirement: Shared Bezier route geometry

Courier SHALL preserve measured cubic route positioning while sampling and snapping its visual points to the four-pixel grid. The landing route signal SHALL be exactly four by four pixels with a centered anchor and linear motion on the existing path; reduced motion SHALL place that center at the path midpoint without travel or repeated animation.

#### Scenario: A route is rendered

- **WHEN** finite Source and Destination positions are supplied
- **THEN** a deterministic crisp polyline connects them, its terminals remain fixed-aspect, and the signal center follows that path or rests at its midpoint under reduced motion

#### Scenario: A consumer supplies two endpoint positions

- **WHEN** a route path is requested for finite coordinates
- **THEN** the helper returns a stable cubic source path and a deterministic four-pixel-grid sample for the rendered connector

### Requirement: Fixed-aspect route terminals

Courier SHALL provide deterministic cubic Bezier paths and fixed-aspect code-native terminal markers for source-to-destination presentations without coupling consumers to a stretched SVG viewport.

#### Scenario: A route is rendered in a non-square viewport

- **WHEN** the connector SVG scales to a wide or portrait layout
- **THEN** Source and Destination markers remain circular within one rendered pixel

### Requirement: Shared command readout

Courier SHALL provide an independent immutable command strip with localized heading, an optionally disabled Copy action, reserved live feedback, optional details and footer actions, and deterministic reset when command identity changes. It SHALL not execute the represented command.

#### Scenario: Copy is unavailable

- **WHEN** a consumer supplies an incomplete or invalid generated command
- **THEN** Copy is disabled and no clipboard write occurs

#### Scenario: Copy succeeds or fails

- **WHEN** a user activates Copy and clipboard access succeeds or fails
- **THEN** the exact command is offered to the clipboard and a localized result is announced in the reserved status region without geometry movement

#### Scenario: The command identity changes

- **WHEN** a consumer replaces the displayed command selection
- **THEN** prior copy feedback resets synchronously and no timed work remains

### Requirement: Shared control frame

Courier SHALL provide shared control-frame tokens and a semantic icon-link component so external icon links and cyclic preference buttons use the same dimensions, border, padding, radius, surface, hover, focus, and theme behavior while preserving their native accessible roles.

#### Scenario: Header controls are rendered together

- **WHEN** an external icon link appears beside theme or locale buttons
- **THEN** their outer chrome matches across light, dark, and system themes while the link remains a link and each preference remains a button

### Requirement: Shared workbench primitive

Courier SHALL compose product state from flat terminal panels, horizontal row rules, semantic controls, and truthful product data. Decorative card islands, fake shell execution, fictional monitoring data, soft shadows, blur, and rounded framing SHALL NOT be introduced.

#### Scenario: A product surface presents state

- **WHEN** landing, delivery, or administration groups related information
- **THEN** hierarchy comes from typography, spacing, flat surface bands, and meaningful one-pixel rules while all displayed values originate from the CLI contract or runtime API

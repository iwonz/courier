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

Courier SHALL retain only approved Courier identity assets and notices with local digests and SHALL exclude injected scripts, credentials, remote resources, obsolete mirror messaging, reference-page executable code, and superseded swift, moth, panorama, realistic-pigeon, or cyberpunk-city artwork. Every shipped Relay raster SHALL record its role, dimensions, byte count, generation prompt, lineage, and SHA-256 digest.

#### Scenario: Asset integrity check

- **WHEN** the repository quality gate runs
- **THEN** only the declared transparent Relay raster is shipped and its manifest digest, dimensions, role, provenance, and budget match the file on disk

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

Courier SHALL use one repository-local technological Relay identity across landing, delivery, administration, README, and product chrome. Relay SHALL be an ImageGen-authored square transparent pigeon mascot built from compact graphite, cobalt, and off-white volumes, tucked segmented wing plates, an integrated courier module, restrained cyan routing light, and one orange waypoint beacon.

#### Scenario: A browser surface presents Courier

- **WHEN** landing, delivery, or administration renders
- **THEN** it uses the same compact technological Relay identity, local assets, accessible contrast, and shared shadcn interaction grammar without a swift, moth, scenic background, remote font, or remote image

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

Courier SHALL present theme, locale, form, file, and policy controls through one shared visual and interaction grammar without exposing unstyled browser-native chrome.

#### Scenario: A user operates a Courier control

- **WHEN** the control is rendered, focused, selected, disabled, or activated with a keyboard or pointer
- **THEN** it uses Courier tokens and visible state treatment while preserving the expected semantic role, accessible name, focus order, and change behavior

### Requirement: Compact icon preference controls

Courier SHALL render theme and locale selectors as separate semantic React icon buttons without visible labels or radiogroups. Each button SHALL expose localized current and next values, use shared shadcn control styling, preserve keyboard activation and focus indication, and delegate persistence and document updates to one shared preference provider.

#### Scenario: A user operates an icon preference selector

- **WHEN** the user points to, focuses, or activates the theme or locale button
- **THEN** its current and next states are available programmatically and activation advances exactly once through the documented cycle

#### Scenario: A user operates a locale selector

- **WHEN** the visitor activates the locale button by pointer or keyboard
- **THEN** its localized accessible name identifies the current and next locale while the visible icon advances once without exposing a radiogroup

### Requirement: Local official brand marks

Courier SHALL provide a shared brand-icon primitive that renders pinned, repository-bundled official geometry in Courier monochrome without runtime network requests. Third-party source, license, attribution, and trademark constraints SHALL be recorded separately from Courier-owned identity assets.

#### Scenario: A branded channel is rendered

- **WHEN** a supported package manager, shell, operating system, or distribution mark appears
- **THEN** the UI uses its registered official geometry, inherits the requested Courier color, and exposes no remote asset URL

### Requirement: Shared styled checkbox

Courier SHALL provide a shared checkbox primitive with custom visual treatment, native checked and disabled semantics, localized accessible labeling, keyboard activation, focus indication, and a composed change event.

#### Scenario: A visitor toggles a view filter

- **WHEN** the checkbox is activated by pointer or keyboard
- **THEN** its checked state and accessible state update once and consumers receive the new boolean value

### Requirement: Shared Bezier route geometry

Courier SHALL provide deterministic cubic Bezier path generation for source-to-destination route presentations without coupling consumers to a specific layout.

#### Scenario: A consumer supplies two endpoint positions

- **WHEN** a route path is requested for finite coordinates
- **THEN** the helper returns a stable horizontal cubic Bezier path spanning those positions

### Requirement: Fixed-aspect route terminals

Courier SHALL provide deterministic cubic Bezier paths and fixed-aspect code-native terminal markers for source-to-destination presentations without coupling consumers to a stretched SVG viewport.

#### Scenario: A route is rendered in a non-square viewport

- **WHEN** the connector SVG scales to a wide or portrait layout
- **THEN** Source and Destination markers remain circular within one rendered pixel

### Requirement: Shared command readout

Courier SHALL provide an independent immutable command strip with localized heading, Copy action, reserved live feedback, optional details and footer actions, and deterministic reset when command identity changes. It SHALL not depend on the workbench component or render a fake shell prompt, editable field, transcript, timer, Run action, or execution behavior.

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

Courier SHALL use repository-owned shadcn cards, scroll areas, separators, alerts, and related primitives for structured product state instead of a separate terminal or workbench component. It SHALL use high-contrast surfaces, generous rounded geometry, compact spacing, and normal sans-serif content without viewport-sized framing, shell prompts, fake execution, or arbitrary command input.

#### Scenario: A product surface presents state

- **WHEN** landing, delivery, or administration renders structured content
- **THEN** shared shadcn composition presents it with readable contrast and lightweight spatial hierarchy without copied consumer-specific implementations

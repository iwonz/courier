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

Courier SHALL retain exactly the declared five transparent pixel Relay WebP assets and one pinned local display font with role, dimensions, byte count, generation prompt, lineage or upstream revision, license, and SHA-256 digest. The mark SHALL remain at most 24 KiB, neutral and route sprites at most 64 KiB each, delivery and administration sprites at most 48 KiB each, all Relay assets at most 248 KiB combined, and the display font at most 48 KiB.

#### Scenario: Asset integrity check

- **WHEN** repository verification runs
- **THEN** undeclared assets, legacy armored Relay files, missing provenance, missing licenses, incorrect dimensions, or an exceeded budget fail the gate

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

Courier SHALL use one repository-local modern 8-bit Relay identity across landing, delivery, administration, README, favicon, and product chrome. Relay SHALL be a compact square pigeon courier with a satchel, parcel, folded wings, compact tail, short beak, visible eye, and at most one small earpiece. The transparent ImageGen-authored mark, neutral mascot, route pose, delivery pose, and administration pose SHALL use hard pixel edges and SHALL NOT depict armor, a helmet, visor, glowing face panel, exoskeleton, metallic chest plate, police or military equipment, photorealism, smooth 3D shading, scenery, text, logos, or watermarks.

The neutral sprite SHALL be the canonical identity reference. Every role sprite SHALL be an identity-preserving derivative of that canonical asset and SHALL preserve its body size, body-to-head proportions, physiology, head, eye, beak, folded wings, compact tail, base plumage, earpiece, and satchel. Only pose, role equipment, clothing, and carried or attached objects MAY vary. Independent mascot redraws SHALL NOT be accepted as role variants.

#### Scenario: A browser surface presents Courier

- **WHEN** landing, delivery, or administration renders
- **THEN** its header uses the square pixel Relay mark beside `COURIER CLI`, its optional role sprite matches the surface, and complete product meaning remains available without either image

#### Scenario: Relay changes roles

- **WHEN** the landing, delivery, administration, README, and header assets are compared
- **THEN** they depict the same recognizable mascot proportions and physiology while only the declared pose and role equipment differ

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

Courier SHALL present semantic browser controls through one shared modern pixel grammar using the declared twelve-color palette, four-pixel geometry unit, crisp dividers, chamfered interactive frames, hard offset state shadows, and visible two-pixel focus. Ordinary layout wrappers SHALL remain transparent and SHALL NOT add card backgrounds, outer radii, blur, soft shadows, or borders. Body content SHALL remain system sans, commands SHALL remain system monospace, and the locally bundled Cyrillic Pixelify Sans font SHALL be limited to brand and display headings.

#### Scenario: A user operates a Courier control

- **WHEN** a control is rendered, focused, selected, disabled, or activated
- **THEN** its pixel state remains contrast-safe, keyboard-operable, semantically native, and free of an unnecessary surrounding island

#### Scenario: A product surface groups related content

- **WHEN** content is grouped for layout without its own interactive state
- **THEN** it remains on the continuous document canvas and uses spacing, typography, alignment, or a meaningful crisp divider instead of a card island

### Requirement: Compact icon preference controls

Courier SHALL render theme and locale selectors as separate semantic React icon buttons without visible labels or radiogroups. The locale button SHALL display an icon that identifies the currently active language. Each button SHALL expose localized current and next values, use shared quiet shadcn control styling, preserve keyboard activation and focus indication, and delegate persistence and document updates to one shared preference provider.

#### Scenario: A user operates an icon preference selector

- **WHEN** the user points to, focuses, or activates the theme or locale button
- **THEN** its current and next states are available programmatically and activation advances exactly once through the documented cycle

#### Scenario: A user operates a locale selector

- **WHEN** the visitor activates the locale button by pointer or keyboard
- **THEN** its visible icon changes from the current English flag to the current Russian flag or vice versa while its localized accessible name identifies both current and next locale

### Requirement: Local official brand marks

Courier SHALL provide typed first-party and third-party pixel icon registries rendered from reviewed sixteen- or twenty-four-pixel grids with crisp edges and current color. Installation, platform, operating-system, GitHub, endpoint, preference, file, transfer, status, and administration icons SHALL be local. Third-party pixel derivatives SHALL retain pinned official source, license, attribution, and trademark records. Locale controls SHALL use local pixel flags rather than platform emoji.

#### Scenario: A browser surface renders an icon

- **WHEN** Courier renders a functional icon, channel mark, or locale flag
- **THEN** it uses the declared local pixel registry with an accessible name where needed and makes no runtime request

#### Scenario: A branded channel is rendered

- **WHEN** a supported package manager, shell, operating system, or distribution mark appears
- **THEN** the UI uses its reviewed pixel-grid derivative, inherits the requested Courier color, retains its pinned source record, and exposes no remote asset URL

### Requirement: Shared styled checkbox

Courier SHALL provide a shared checkbox primitive with custom visual treatment, native checked and disabled semantics, localized accessible labeling, keyboard activation, focus indication, and a composed change event.

#### Scenario: A visitor toggles a view filter

- **WHEN** the checkbox is activated by pointer or keyboard
- **THEN** its checked state and accessible state update once and consumers receive the new boolean value

### Requirement: Shared Bezier route geometry

Courier SHALL preserve measured cubic route positioning while sampling and snapping its visual points to the four-pixel grid. Decorative route packets and status signals SHALL use short stepped motion only, and reduced motion SHALL render the same state without travel or repeated animation.

#### Scenario: A route is rendered

- **WHEN** finite Source and Destination positions are supplied
- **THEN** a deterministic crisp polyline connects them, its terminals remain fixed-aspect, and any packet animation stops under reduced motion

#### Scenario: A consumer supplies two endpoint positions

- **WHEN** a route path is requested for finite coordinates
- **THEN** the helper returns a stable cubic source path and a deterministic four-pixel-grid sample for the rendered connector

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

Courier SHALL use repository-owned shadcn layout and control primitives for structured product state without a separate terminal or workbench component. Card SHALL remain a semantic structural wrapper with no default surface chrome; command readouts, route displays, and tab lists SHALL likewise avoid a containing panel. The system SHALL use normal sans-serif content without viewport-sized framing, shell prompts, fake execution, or arbitrary command input.

#### Scenario: A product surface presents state

- **WHEN** landing, delivery, or administration renders structured content
- **THEN** shared shadcn composition presents it directly on the document canvas with readable hierarchy and no ordinary card islands

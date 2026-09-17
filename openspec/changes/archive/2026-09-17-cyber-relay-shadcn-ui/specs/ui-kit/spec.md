## MODIFIED Requirements

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

### Requirement: Compact icon preference controls

Courier SHALL render theme and locale selectors as separate semantic React icon buttons without visible labels or radiogroups. Each button SHALL expose localized current and next values, use shared shadcn control styling, preserve keyboard activation and focus indication, and delegate persistence and document updates to one shared preference provider.

#### Scenario: A user operates an icon preference selector

- **WHEN** the user points to, focuses, or activates the theme or locale button
- **THEN** its current and next states are available programmatically and activation advances exactly once through the documented cycle

#### Scenario: A user operates a locale selector

- **WHEN** the visitor activates the locale button by pointer or keyboard
- **THEN** its localized accessible name identifies the current and next locale while the visible icon advances once without exposing a radiogroup

### Requirement: Shared workbench primitive

Courier SHALL use repository-owned shadcn cards, scroll areas, separators, alerts, and related primitives for structured product state instead of a separate terminal or workbench component. It SHALL use high-contrast surfaces, generous rounded geometry, compact spacing, and normal sans-serif content without viewport-sized framing, shell prompts, fake execution, or arbitrary command input.

#### Scenario: A product surface presents state

- **WHEN** landing, delivery, or administration renders structured content
- **THEN** shared shadcn composition presents it with readable contrast and lightweight spatial hierarchy without copied consumer-specific implementations

## REMOVED Requirements

### Requirement: Smooth stationary scene response

**Reason**: Courier no longer ships a generated scene system; one static transparent Relay raster is composed directly by React consumers.

**Migration**: Use the shared mascot component and code-native route geometry.

### Requirement: Contextual Vector illustration family

**Reason**: The Vector/scene family was superseded by one compact technological Relay mascot without generated backgrounds.

**Migration**: Use the locally verified Relay raster exported from the shared asset module.

### Requirement: Vector compact mark

**Reason**: The Vector moth mark was superseded by the compact ImageGen-authored Relay pigeon.

**Migration**: Use the shared Relay brand and mascot components.

# ui-kit Specification

## Purpose
Define the single reusable visual, accessibility, theme, localization, component, and identity-asset foundation shared by every Courier browser surface.

## Requirements

### Requirement: Shared UI package

Courier SHALL provide one Lit and TypeScript UI package for delivery pages, administration pages, and the project landing page, with no copied component implementations between consumers.

#### Scenario: Consumer imports a control

- **WHEN** a Courier web application imports a public component
- **THEN** its behavior, tokens, icons, and accessible states come from the shared package

### Requirement: Auditable identity assets

Courier SHALL retain only approved Courier identity assets and notices with local digests and SHALL exclude injected scripts, credentials, remote resources, obsolete mirror messaging, and reference-page executable code.

#### Scenario: Asset integrity check

- **WHEN** the repository quality gate runs
- **THEN** every retained identity asset matches its provenance digest and every manifest entry resolves locally

### Requirement: User-selectable themes

Courier SHALL provide visible `system`, `light`, and `dark` theme choices, default to `system`, persist explicit preference locally, meet accessible contrast, and respect reduced-motion preferences.

#### Scenario: System theme changes

- **WHEN** the stored preference is `system` and the operating-system color scheme changes
- **THEN** the resolved theme updates without replacing the stored preference

### Requirement: Extensible localization

Courier SHALL provide typed English and Russian catalogs, use English as the fallback, negotiate the initial locale from browser languages, and permit additional catalog modules without backend changes.

#### Scenario: Unsupported browser locale

- **WHEN** no supported locale matches the browser language list
- **THEN** English is selected and every requested message key resolves

### Requirement: Exact UI source coverage

Courier SHALL enforce 100% statements, branches, functions, and lines for first-party TypeScript sources while excluding generated output and third-party code.

#### Scenario: Untested branch is introduced

- **WHEN** a first-party TypeScript branch is not executed by the test suite
- **THEN** the repository verification gate fails

### Requirement: Coherent brand identity system

Courier SHALL define and apply one mature, utilitarian identity across browser surfaces, documentation, and retained assets using shared semantic tokens and components rather than consumer-specific brand implementations.

#### Scenario: A browser surface presents Courier

- **WHEN** the landing, delivery, or administration application renders
- **THEN** it uses the same positioning, wordmark, semantic color roles, typography, focus treatment, and operational visual grammar from the shared UI package

### Requirement: Non-essential mascot guidance

Courier SHALL use the original Relay courier-pigeon character only as a restrained orientation and state-communication aid, with locally verified provenance and without making meaning depend on the image.

#### Scenario: Mascot artwork is unavailable

- **WHEN** Relay cannot be loaded or is hidden from assistive technology
- **THEN** headings, status text, controls, and progress information still communicate the complete workflow

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

### Requirement: Contextual Relay illustration family

Courier SHALL use role-specific Relay illustrations that preserve the character identity while matching the operational context of each appearance.

#### Scenario: Relay appears on a product surface

- **WHEN** a landing, authentication, delivery, or administration view includes Relay
- **THEN** the illustration's wardrobe, tools, posture, and surrounding scene communicate that view's purpose without duplicating an unrelated pose or carrying required information

#### Scenario: Illustration assets are audited

- **WHEN** the asset verification gate runs
- **THEN** every shipped illustration has a declared role, local path, dimensions, digest, prompt, and authorship notice and no undeclared identity raster is bundled

### Requirement: Compact icon preference controls

Courier SHALL render theme and locale selectors as shared icon-only segmented radiogroups without visible group or option labels while retaining localized accessible names, selected state, persisted preference, roving focus, and arrow/Home/End keyboard behavior. Locale options SHALL use native flag emoji as their visible symbols.

#### Scenario: A user operates a locale selector

- **WHEN** the user points to, focuses, or navigates an English or Russian locale option
- **THEN** a flag emoji identifies the option visually while its localized name and radio state remain available programmatically

#### Scenario: A user operates an icon preference selector

- **WHEN** the user points to, focuses, or navigates a theme or locale option
- **THEN** its meaning is available programmatically, its state is visibly distinguishable, and changing it has the same persisted behavior as the labeled control

### Requirement: Relay compact mark

Courier SHALL use a recognizable, repository-local, generated Relay mascot image as its compact product mark across browser components and favicons.

#### Scenario: A compact Courier identity is rendered

- **WHEN** the wordmark has limited space or a favicon is displayed
- **THEN** the locally stored transparent mark depicts Relay with sufficient light/dark contrast and no remote image dependency

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

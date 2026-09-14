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

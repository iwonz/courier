## ADDED Requirements

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

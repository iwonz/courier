## MODIFIED Requirements

### Requirement: Shared accessible administration application

Courier SHALL embed a deterministic Lit administration application built from the shared terminal workspace and UI kit with typed English and Russian catalogs, system/light/dark themes, keyboard-operable API-backed controls, responsive live registry transcripts, and visible conflict/error states. It SHALL provide no command prompt or arbitrary execution surface.

#### Scenario: An operator inspects and changes a delivery

- **WHEN** a live snapshot arrives or a semantic policy, refresh, or stop control is activated
- **THEN** the terminal workspace updates the authoritative secret-free state through the existing guarded API without recreating a shell interface

#### Scenario: Locale and theme change

- **WHEN** an operator selects Russian and dark theme
- **THEN** the administration content localizes, the shared theme resolves to dark, and both preferences persist through the shared browser adapters

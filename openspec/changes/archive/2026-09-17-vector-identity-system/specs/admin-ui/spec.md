## MODIFIED Requirements

### Requirement: Shared accessible administration application

Courier SHALL embed a deterministic Lit operations workbench built from the shared UI kit with typed English/Russian catalogs, cyclic system/light/dark preference, overview counters, a keyboard-operable server/delivery navigator, a selected route and policy inspector, and explicit refresh, save, and stop actions backed only by the guarded API. Selection SHALL persist by UUID across authoritative SSE snapshots when possible and SHALL fall back deterministically when the selected object disappears. The application SHALL provide no terminal prompt or arbitrary execution surface.

#### Scenario: A snapshot updates the selected delivery

- **WHEN** an SSE snapshot still contains the selected UUID
- **THEN** counters, route, state, and policy update in place while selection and focused controls remain stable

#### Scenario: A selected object disappears

- **WHEN** an authoritative snapshot no longer contains the selected server or delivery
- **THEN** the navigator selects the first valid remaining object or renders the localized empty state without exposing stale data

#### Scenario: The viewport is narrow

- **WHEN** the administration application renders on a mobile viewport
- **THEN** navigator and inspector stack without hiding actions or creating horizontal overflow

#### Scenario: An operator inspects and changes a delivery

- **WHEN** a live snapshot arrives or a semantic policy, refresh, or stop control is activated
- **THEN** the operations workbench updates authoritative secret-free state through the existing guarded API without recreating a shell interface

#### Scenario: Locale and theme change

- **WHEN** an operator cycles to Russian and dark theme
- **THEN** administration content localizes, the shared theme resolves to dark, and both preferences persist through shared browser adapters

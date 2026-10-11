## MODIFIED Requirements

### Requirement: Shared accessible administration application

Courier SHALL embed a deterministic React operations application built from the repository-owned shadcn component system with typed English/Russian catalogs, cyclic system/light/dark preference, overview counters, a keyboard-operable server/delivery navigator, a selected route and policy inspector, a recent upload activity view, and explicit refresh, save, and stop actions backed only by the guarded API. Counters, navigator, inspector, activity, route facts, and policy controls SHALL render directly on one continuous canvas without rounded panel backgrounds or nested card islands. Selection SHALL persist by UUID across authoritative SSE snapshots when possible and SHALL fall back deterministically when the selected object disappears. The application SHALL provide no terminal prompt or arbitrary execution surface.

#### Scenario: A snapshot updates the selected delivery

- **WHEN** an SSE snapshot still contains the selected UUID
- **THEN** counters, route, state, policy, and matching activity update in place while selection and focused controls remain stable

#### Scenario: A selected object disappears

- **WHEN** an authoritative snapshot no longer contains the selected server or delivery
- **THEN** the navigator selects the first valid remaining object or renders the localized empty state without exposing stale data

#### Scenario: An operator inspects upload activity

- **WHEN** the selected delivery has recent incoming browser or webhook receipts
- **THEN** the UI shows localized time, actor and peer, original and committed names, human-readable and exact bytes, route/destination, and extraction outcome from authoritative secret-free records

#### Scenario: Activity contains no rows for the selection

- **WHEN** the selected delivery has no retained receipts
- **THEN** the UI renders a localized empty activity state without showing unrelated delivery records

#### Scenario: The viewport is narrow

- **WHEN** administration including upload activity renders on a mobile viewport
- **THEN** navigator, inspector, and activity rows stack without hiding values, actions, or creating horizontal overflow

#### Scenario: An operator manages a delivery

- **WHEN** the operator selects a delivery
- **THEN** lightweight shadcn navigation, activity, and policy controls present current route, counters, receipts, policy, save, and stop actions with complete keyboard operation and without nested card-on-card framing

#### Scenario: An operator inspects and changes a delivery

- **WHEN** a live snapshot arrives or a semantic policy, refresh, or stop control is activated
- **THEN** the React operations application updates authoritative secret-free state through the existing guarded API without recreating a shell interface

#### Scenario: Locale and theme change

- **WHEN** an operator cycles to Russian and dark theme
- **THEN** administration content and quantity details localize, the shared theme resolves to dark, and both preferences persist through the shared React preference provider

#### Scenario: An operator inspects a delivery

- **WHEN** the operator chooses a delivery from the server navigator
- **THEN** current route, counters, policy, activity, save, and stop actions remain keyboard-operable in one borderless document flow without changing guarded API behavior

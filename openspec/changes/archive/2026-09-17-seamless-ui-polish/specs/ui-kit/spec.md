## MODIFIED Requirements

### Requirement: Coherent brand identity system

Courier SHALL use one repository-local technological Relay identity across landing, delivery, administration, README, and product chrome. The full Relay mascot SHALL remain an ImageGen-authored square transparent pigeon illustration, while compact product chrome SHALL use a distinct ImageGen-authored head-and-shoulders Relay mark with a broad near-square silhouette, no feet, perch, tail fan, or rounded chicken-like body. Both SHALL use compact graphite, cobalt, and off-white volumes, restrained cyan routing light, and one orange waypoint beacon.

#### Scenario: A browser surface presents Courier

- **WHEN** landing, delivery, or administration renders
- **THEN** its header uses the compact Relay mark beside `COURIER CLI`, editorial usage may use the full mascot, and both remain local, decorative, and legible without a remote font or image

### Requirement: Intentional control appearance

Courier SHALL present theme, locale, form, file, policy, navigation, and command controls through one shared lightweight shadcn visual grammar. Borders SHALL be reserved for inputs, focus, destructive emphasis, or meaningful internal boundaries; ordinary grouping SHALL prefer spacing, tone, typography, and a continuous canvas over nested outlined cards.

#### Scenario: A user operates a Courier control

- **WHEN** the control is rendered, focused, selected, disabled, or activated with a keyboard or pointer
- **THEN** it uses Courier tokens and visible semantic state while preserving the expected role, accessible name, focus order, and change behavior without an unnecessary surrounding border

### Requirement: Compact icon preference controls

Courier SHALL render theme and locale selectors as separate semantic React icon buttons without visible labels or radiogroups. The locale button SHALL display an icon that identifies the currently active language. Each button SHALL expose localized current and next values, use shared quiet shadcn control styling, preserve keyboard activation and focus indication, and delegate persistence and document updates to one shared preference provider.

#### Scenario: A user operates an icon preference selector

- **WHEN** the user points to, focuses, or activates the theme or locale button
- **THEN** its current and next states are available programmatically and activation advances exactly once through the documented cycle

#### Scenario: A user operates a locale selector

- **WHEN** the visitor activates the locale button by pointer or keyboard
- **THEN** its visible icon changes from the current English flag to the current Russian flag or vice versa while its localized accessible name identifies both current and next locale

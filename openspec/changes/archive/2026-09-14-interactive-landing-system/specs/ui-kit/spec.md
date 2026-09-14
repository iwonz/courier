## ADDED Requirements

### Requirement: Compact icon preference controls

Courier SHALL render theme and locale selectors as shared icon-only segmented radiogroups without visible group or option labels while retaining localized accessible names, selected state, persisted preference, roving focus, and arrow/Home/End keyboard behavior.

#### Scenario: A user operates an icon preference selector

- **WHEN** the user points to, focuses, or navigates a theme or locale option
- **THEN** its meaning is available programmatically, its state is visibly distinguishable, and changing it has the same persisted behavior as the labeled control

### Requirement: Relay compact mark

Courier SHALL use a recognizable, repository-local Relay mascot symbol as its compact product mark across browser components, favicons, and light/dark lockups.

#### Scenario: A compact Courier identity is rendered

- **WHEN** the wordmark has limited space or a favicon is displayed
- **THEN** the mark depicts Relay with sufficient light/dark contrast rather than an abstract route arrow

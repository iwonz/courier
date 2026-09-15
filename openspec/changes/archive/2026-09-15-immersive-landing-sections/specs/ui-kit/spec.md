## MODIFIED Requirements

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

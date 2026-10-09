## ADDED Requirements

### Requirement: Localized structured parameter presentation

The landing command builder SHALL render generated parameter syntax, localized Required or Optional status, localized plain-language descriptions, human-readable applicable routes or commands, and relevant default, repeatability, dependency, and conflict metadata in English and Russian without changing controls or generated command values.

#### Scenario: Locale changes

- **WHEN** a visitor switches between English and Russian
- **THEN** parameter descriptions, status, applicability labels, and metadata labels change locale while entered values and generated command syntax remain unchanged

#### Scenario: Parameter applicability is shown

- **WHEN** a visitor inspects a route-scoped option
- **THEN** the option names the applicable routes in readable language instead of exposing internal route identifiers

#### Scenario: Parameter rows use a narrow viewport

- **WHEN** the command builder renders on a supported narrow viewport
- **THEN** syntax, description, applicability, metadata, and controls remain readable and keyboard-operable without horizontal page overflow


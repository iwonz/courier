## MODIFIED Requirements

### Requirement: Extensible command registry

Courier SHALL dispatch top-level behavior through a Cobra command tree so new commands are added as isolated constructors without changing root argument parsing.

#### Scenario: Registered command is executed

- **WHEN** a Cobra child command name is the first CLI argument
- **THEN** Cobra invokes its handler with validated remaining arguments and flags

#### Scenario: Unknown command is rejected

- **WHEN** the first CLI argument does not name a registered command
- **THEN** Courier returns the stable CLI usage exit code and does not execute another handler

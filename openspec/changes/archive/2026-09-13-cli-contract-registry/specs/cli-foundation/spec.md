## MODIFIED Requirements

### Requirement: Extensible command registry

Courier SHALL construct its Cobra tree from independently registered command providers, reject duplicate command names, and allow a provider to be added without changing parsing logic.

#### Scenario: Registered command is executed

- **WHEN** a provider contributes a unique top-level command
- **THEN** the root attaches that command and Cobra invokes its handler with the remaining arguments

#### Scenario: Duplicate command

- **WHEN** two providers contribute the same top-level command name
- **THEN** root construction fails before command execution

#### Scenario: Unknown command is rejected

- **WHEN** the first CLI argument does not name a registered command
- **THEN** Courier returns a stable usage error and does not execute another handler

### Requirement: Help output

Courier SHALL print usage and registered commands when invoked without arguments or with `help`, `-h`, or `--help`, SHALL keep `version` and `update` as system commands, and SHALL disable Cobra's default `completion` command.

#### Scenario: Root help requested

- **WHEN** the user invokes Courier without a transfer command
- **THEN** output includes `from`, `help`, `update`, and `version` but not `completion`

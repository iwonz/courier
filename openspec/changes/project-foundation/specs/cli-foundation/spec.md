## ADDED Requirements

### Requirement: Extensible command registry

Courier SHALL dispatch top-level commands through a command registry so a new command can be registered without changing the registry parser.

#### Scenario: Registered command is executed

- **WHEN** a registered command name is the first CLI argument
- **THEN** the registry invokes its handler with the remaining arguments

#### Scenario: Unknown command is rejected

- **WHEN** the first CLI argument does not name a registered command
- **THEN** Courier returns a stable usage error and does not execute another handler

### Requirement: Build identity

Courier SHALL expose its semantic version, commit and build date through the `version` command, with development-safe defaults when linker metadata is absent.

#### Scenario: Development build reports identity

- **WHEN** a locally built binary runs `courier version`
- **THEN** it prints a non-empty version, commit and build date

### Requirement: Help output

Courier SHALL print usage and registered commands when invoked without arguments or with `help`, `-h`, or `--help`.

#### Scenario: Root help requested

- **WHEN** the user invokes Courier without a transfer command
- **THEN** the output includes the `courier from <source> to <destination> [flags]` syntax

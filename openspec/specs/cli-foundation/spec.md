# cli-foundation Specification

## Purpose
Define Courier's extensible command assembly, build identity, and help behavior.

## Requirements

### Requirement: Extensible command registry

Courier SHALL dispatch top-level behavior through a Cobra command tree so new commands are added as isolated constructors without changing root argument parsing.

#### Scenario: Registered command is executed

- **WHEN** a Cobra child command name is the first CLI argument
- **THEN** Cobra invokes its handler with validated remaining arguments and flags

#### Scenario: Unknown command is rejected

- **WHEN** the first CLI argument does not name a registered command
- **THEN** Courier returns the stable CLI usage exit code and does not execute another handler

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

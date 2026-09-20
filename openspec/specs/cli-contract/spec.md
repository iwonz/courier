# cli-contract Specification

## Purpose
Define the canonical machine-readable public CLI inventory and generated-reference parity gate.

## Requirements

### Requirement: Canonical command contract

Courier SHALL keep a versioned machine-readable English contract containing every public command and option combination plus ordered command paths and arguments. Arguments SHALL declare name, kind, required state, optional literal prefix, and optional suppressing flag. Parameters SHALL declare a value kind, optional choices, placeholder, dependencies, repeatability, conflicts, and default.

#### Scenario: Contract metadata is invalid

- **WHEN** an argument or parameter has an unknown kind, invalid choices, duplicate reference, impossible dependency, unknown suppressing flag, or inconsistent repeatability
- **THEN** strict contract validation fails

#### Scenario: Planned command

- **WHEN** a command is present in the target contract but is not implemented
- **THEN** it is marked planned and does not appear in live help

### Requirement: Generated reference parity

Courier SHALL deterministically generate its human command reference and validate the shipped Cobra tree against the canonical contract.

#### Scenario: Undocumented Cobra command

- **WHEN** Cobra exposes a command or flag not present as shipped or system behavior in the contract
- **THEN** the local verification gate fails

### Requirement: Generated shipped-only landing data

Courier SHALL deterministically generate landing command-builder data from the canonical CLI contract, preserving command, argument, parameter, choice, dependency, and conflict order while omitting planned inventory.

#### Scenario: Landing data is regenerated

- **WHEN** the canonical contract changes
- **THEN** generated JSON contains the exact structured path, arguments, and typed parameter metadata or the freshness gate fails

#### Scenario: A planned command remains in the contract

- **WHEN** the landing data is generated or checked
- **THEN** the planned command is omitted while shipped product commands and system commands remain represented according to their status

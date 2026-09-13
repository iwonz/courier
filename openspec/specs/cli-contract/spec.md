# cli-contract Specification

## Purpose
Define the canonical machine-readable public CLI inventory and generated-reference parity gate.

## Requirements

### Requirement: Canonical command contract

Courier SHALL keep a versioned machine-readable English contract containing every shipped, planned, system, and explicitly unsupported public command and option combination.

#### Scenario: Planned command

- **WHEN** a command is present in the target contract but is not implemented
- **THEN** it is marked planned and does not appear in live help

### Requirement: Generated reference parity

Courier SHALL deterministically generate its human command reference and validate the shipped Cobra tree against the canonical contract.

#### Scenario: Undocumented Cobra command

- **WHEN** Cobra exposes a command or flag not present as shipped or system behavior in the contract
- **THEN** the local verification gate fails

### Requirement: Generated shipped-only landing data

Courier SHALL deterministically generate landing-page command data from the canonical CLI contract and SHALL reject stale output or any planned entry exposed as available.

#### Scenario: A planned command remains in the contract

- **WHEN** the landing data is generated or checked
- **THEN** the planned command is omitted while shipped product commands and system commands remain represented according to their status

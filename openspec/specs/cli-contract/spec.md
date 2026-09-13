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

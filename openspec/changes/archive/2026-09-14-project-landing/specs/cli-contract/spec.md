## ADDED Requirements

### Requirement: Generated shipped-only landing data

Courier SHALL deterministically generate landing-page command data from the canonical CLI contract and SHALL reject stale output or any planned entry exposed as available.

#### Scenario: A planned command remains in the contract

- **WHEN** the landing data is generated or checked
- **THEN** the planned command is omitted while shipped product commands and system commands remain represented according to their status

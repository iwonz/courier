## ADDED Requirements

### Requirement: Platform-valid Unix control endpoints

Courier SHALL create new Unix control sockets inside its private state directory using a deterministic, lossless compact representation of the full server UUID that fits the supported platform pathname limit under Courier's default state-directory layout. Courier SHALL continue using a persisted legacy control endpoint exactly as recorded, and SHALL leave Windows named-pipe generation unchanged.

#### Scenario: Administration starts from the default macOS state directory

- **WHEN** Courier generates an administration control endpoint beneath `~/Library/Application Support/courier/state`
- **THEN** the endpoint fits Darwin's Unix-socket pathname boundary and the private listener can start

#### Scenario: A web delivery starts from the default macOS state directory

- **WHEN** Courier acquires a new data worker for a web delivery
- **THEN** the worker receives the same deterministic compact endpoint that its launch configuration validates

#### Scenario: A legacy process remains live across upgrade

- **WHEN** persisted state records a live `control-<UUID>.sock` endpoint created by an older Courier version
- **THEN** the upgraded client probes, controls, and removes that exact recorded endpoint without rewriting it

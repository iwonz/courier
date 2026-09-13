## ADDED Requirements

### Requirement: Structured progress

Courier SHALL emit structured events containing operation stage, current bytes, total bytes, elapsed time, and calculated transfer speed.

#### Scenario: Data chunk copied

- **WHEN** a data chunk is confirmed written
- **THEN** current bytes increase monotonically and speed is derived from elapsed time

### Requirement: Error accounting

Courier SHALL include the failing stage and confirmed transferred byte count in transfer errors.

#### Scenario: Writer failure

- **WHEN** destination fails after accepting some bytes
- **THEN** the returned error identifies the transfer stage and the confirmed byte count

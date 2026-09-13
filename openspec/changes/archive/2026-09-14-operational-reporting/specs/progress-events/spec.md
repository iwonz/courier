## MODIFIED Requirements

### Requirement: Structured progress

Courier SHALL emit structured events from a closed stage vocabulary containing monotonic read, sent, and confirmed byte counters, total payload bytes, elapsed time, and speed derived from confirmed bytes, and SHALL enforce `confirmed <= sent <= read`.

#### Scenario: Data chunk copied

- **WHEN** a source chunk is read, submitted to a transport, and durably accepted by its target
- **THEN** the corresponding counters advance without regression and the confirmed-byte speed reflects elapsed time

### Requirement: Error accounting

Courier SHALL include the failing stage and the last known read, sent, and confirmed byte counts in typed transfer errors and final failure output.

#### Scenario: Writer failure

- **WHEN** a destination fails after accepting part of a stream
- **THEN** the returned error identifies the transfer stage and the exact confirmed count without claiming unconfirmed bytes

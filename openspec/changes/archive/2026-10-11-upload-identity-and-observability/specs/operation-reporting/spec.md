## ADDED Requirements

### Requirement: Human-readable quantity presentation

Courier SHALL format interactive byte, bit, and rate quantities with space-grouped base-unit integers, binary scaled units through EiB/Eibit, at most one fractional digit, and unambiguous `B` versus `bit` symbols. Detailed and accessible web presentation SHALL retain the grouped exact base-unit value. Redirected CLI fields and JSON numeric fields SHALL remain raw and stable.

#### Scenario: Byte quantity is rendered interactively

- **WHEN** an interactive surface displays `411789` bytes
- **THEN** it presents a human-readable binary byte quantity and makes the exact `411 789 B` value available without labeling it as bits

#### Scenario: Byte rate is rendered

- **WHEN** progress reports a per-second byte rate
- **THEN** Courier uses a byte unit ending in `/s` and never a lowercase bit unit

#### Scenario: Output is redirected

- **WHEN** an operation report is written to a pipe
- **THEN** its existing raw integer byte fields remain deterministic for machine consumers

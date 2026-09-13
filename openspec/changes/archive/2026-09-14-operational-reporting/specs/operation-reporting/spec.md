## MODIFIED Requirements

### Requirement: Interactive progress rendering

Courier SHALL render the current stage, read, sent, confirmed, total, speed, and elapsed time with mpb on an interactive terminal, and SHALL emit deterministic, cursor-free line output otherwise.

#### Scenario: Non-interactive output

- **WHEN** stderr is not a terminal
- **THEN** each progress line uses stable named fields and contains no cursor-control sequences

### Requirement: Error summary and exit codes

Courier SHALL print the failing stage, sanitized reason, and known byte counters, SHALL preserve exit codes 0, 2, 10, 20, 30, and 40, and SHALL return 130 for context cancellation or process interruption.

#### Scenario: Sensitive connection failure

- **WHEN** a connection or control error contains URL user-info, query secrets, authorization data, or credential assignments
- **THEN** Courier reports the useful failure context with the sensitive value redacted and returns the stable category code

#### Scenario: Connection failure

- **WHEN** SSH setup fails before transfer
- **THEN** Courier exits with the documented connection code and no credentials in output

## ADDED Requirements

### Requirement: Unified final accounting

Courier SHALL use confirmed payload bytes for the success result and SHALL never describe read-only or merely submitted bytes as transferred successfully.

#### Scenario: Transport outcome is uncertain

- **WHEN** bytes were read and sent but remote acceptance cannot be confirmed
- **THEN** the operation fails with separate sent and confirmed counts and no success summary

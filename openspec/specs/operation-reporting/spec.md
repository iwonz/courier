# operation-reporting Specification

## Purpose
Define safe, stable progress, success, failure, and exit-code reporting.

## Requirements

### Requirement: Interactive progress rendering

Courier SHALL render the current stage, read, sent, confirmed, total, speed, and elapsed time with mpb on an interactive terminal, and SHALL emit deterministic, cursor-free line output otherwise.

#### Scenario: Non-interactive output

- **WHEN** stderr is not a terminal
- **THEN** each progress line uses stable named fields and contains no cursor-control sequences

### Requirement: Success summary

Courier SHALL print source, actual destination, confirmed transferred bytes, elapsed time, and successful result after commit.

#### Scenario: Successful transfer

- **WHEN** commit and cleanup complete
- **THEN** the summary names the actual destination after directory/archive semantics

### Requirement: Error summary and exit codes

Courier SHALL print the failing stage, sanitized reason, and known byte counters, SHALL preserve exit codes 0, 2, 10, 20, 30, and 40, and SHALL return 130 for context cancellation or process interruption.

#### Scenario: Sensitive connection failure

- **WHEN** a connection or control error contains URL user-info, query secrets, authorization data, or credential assignments
- **THEN** Courier reports the useful failure context with the sensitive value redacted and returns the stable category code

#### Scenario: Connection failure

- **WHEN** SSH setup fails before transfer
- **THEN** Courier exits with the documented connection code and no credentials in output

### Requirement: Unified final accounting

Courier SHALL use confirmed payload bytes for the success result and SHALL never describe read-only or merely submitted bytes as transferred successfully.

#### Scenario: Transport outcome is uncertain

- **WHEN** bytes were read and sent but remote acceptance cannot be confirmed
- **THEN** the operation fails with separate sent and confirmed counts and no success summary

### Requirement: Terminal-capability runtime presentation

Courier SHALL use the shared square-panel and table presentation for runtime commands only when their target output stream is a real TTY, SHALL disable color for `NO_COLOR` or `TERM=dumb`, and SHALL sanitize dynamic terminal values. Redirected output SHALL retain deterministic plain fields without ANSI or cursor-control sequences.

#### Scenario: Runtime output is redirected

- **WHEN** transfer, hosted, UI, servers, update, or version output is piped or redirected
- **THEN** Courier emits its stable plain representation without borders, color, or cursor movement

#### Scenario: Dynamic text contains terminal controls

- **WHEN** a path, remote error, release note, or other dynamic value contains escape or unsafe control sequences
- **THEN** the interactive renderer removes those sequences before layout

### Requirement: Neutral interactive interruption

Courier SHALL retain exit code `130` for a user-cancelled foreground operation and SHALL present the interruption as stopped by the user in an interactive terminal. A concurrent cleanup or lease-release failure SHALL remain an error.

#### Scenario: A foreground server receives Ctrl+C

- **WHEN** the wait ends only because the user cancelled the command
- **THEN** the TTY presentation is neutral and identifies exit code 130 without describing the delivery as failed

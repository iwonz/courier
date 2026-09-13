# operation-reporting Specification

## Purpose
Define safe, stable progress, success, failure, and exit-code reporting.

## Requirements

### Requirement: Interactive progress rendering

Courier SHALL render stage, confirmed bytes, total bytes, speed, and elapsed time with mpb on an interactive terminal, and SHALL emit deterministic line output otherwise.

#### Scenario: Non-interactive output

- **WHEN** stderr is not a terminal
- **THEN** progress output contains no cursor-control sequences

### Requirement: Success summary

Courier SHALL print source, actual destination, confirmed transferred bytes, elapsed time, and successful result after commit.

#### Scenario: Successful transfer

- **WHEN** commit and cleanup complete
- **THEN** the summary names the actual destination after directory/archive semantics

### Requirement: Error summary and exit codes

Courier SHALL print the failing stage, safe reason, and confirmed byte count, and SHALL use stable codes for CLI, connection, transfer, and update failures.

#### Scenario: Connection failure

- **WHEN** SSH setup fails before transfer
- **THEN** Courier exits with the documented connection code and no credentials in output

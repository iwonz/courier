## ADDED Requirements

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

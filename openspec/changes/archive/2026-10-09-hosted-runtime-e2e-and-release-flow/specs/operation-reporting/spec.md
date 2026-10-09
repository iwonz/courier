## ADDED Requirements

### Requirement: Readiness QR projection

Courier SHALL render a QR code containing the exact validated readiness URL only when output is an interactive terminal with sufficient width. It SHALL retain the text URL, use a four-module quiet zone, omit ANSI when color is disabled, and omit the QR for redirected output or `TERM=dumb`.

#### Scenario: Readiness is redirected

- **WHEN** administration, browser delivery, or incoming webhook readiness is written to a pipe
- **THEN** the stable text URL is present and no QR modules or ANSI sequences are emitted

#### Scenario: Terminal is too narrow

- **WHEN** the complete QR including its quiet zone would exceed the detected terminal width
- **THEN** Courier prints the readiness fields without a partial or wrapped QR

### Requirement: Compact usage diagnostics

Courier SHALL classify every command-shape failure with exit code `2` and render only a sanitized error, an available suggestion, and the relevant help command. Usage diagnostics SHALL NOT include operational stage or byte-accounting fields.

#### Scenario: Unknown command resembles a real command

- **WHEN** a user invokes `courier services`
- **THEN** Courier reports the unknown command, suggests `courier servers`, identifies the help command, and omits `stage`, `read`, `sent`, `confirmed`, and `result`

#### Scenario: Unsupported route is supplied

- **WHEN** source and destination form no supported route
- **THEN** Courier exits with code `2` through the same compact usage renderer

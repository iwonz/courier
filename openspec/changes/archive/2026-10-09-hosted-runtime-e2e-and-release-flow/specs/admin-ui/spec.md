## ADDED Requirements

### Requirement: Administration QR readiness

Interactive `courier ui start` readiness SHALL include a QR whose payload is the exact printed local URL when the terminal can render it. The administration listener SHALL remain loopback-only and the output SHALL not claim phone reachability.

#### Scenario: Administration starts in a capable TTY

- **WHEN** the verified loopback administration URL is ready
- **THEN** Courier prints the URL and exact-payload QR while retaining loopback binding

### Requirement: Clean foreground administration stop

The foreground command that started the administration UI SHALL distinguish verified `courier ui stop` from parent-context interruption.

#### Scenario: UI is stopped from another Courier process

- **WHEN** verified control requests shutdown
- **THEN** the initiating command prints a neutral stopped result and exits `0`

#### Scenario: UI foreground receives Ctrl+C

- **WHEN** its parent context is cancelled
- **THEN** the initiating command exits `130`

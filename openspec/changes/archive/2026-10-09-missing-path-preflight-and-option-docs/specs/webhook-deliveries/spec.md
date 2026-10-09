## ADDED Requirements

### Requirement: Webhook path preflight

Courier SHALL treat incoming webhook destinations and extraction roots as directories and SHALL treat an unarchived outgoing webhook source as a file-only role. It SHALL complete local or SSH path preflight before worker acquisition or an outgoing HTTP request.

#### Scenario: Incoming destination is absent

- **WHEN** an incoming webhook targets a missing destination directory
- **THEN** Courier confirms or force-creates it before launching or registering the worker

#### Scenario: Outgoing file-only source is absent

- **WHEN** an unarchived outgoing webhook source does not exist
- **THEN** Courier fails at preflight without creating it or opening an HTTP request

### Requirement: Discoverable detached incoming webhook lifecycle

Courier SHALL document and expose `--background` as the sole detached mode for incoming webhooks. It SHALL return after printing the URL and delivery UUID, survive terminal closure, and end through delivery or server stop, configured stop behavior, fatal worker failure, or process termination.

#### Scenario: Background incoming webhook is ready

- **WHEN** an incoming webhook starts with `--background`
- **THEN** the initiating command exits after readiness and the delivery remains listed by `courier servers`


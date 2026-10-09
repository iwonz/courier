## ADDED Requirements

### Requirement: Browser delivery path preflight

Courier SHALL validate and, when authorized, create the local or SSH path used by a browser delivery before acquiring or registering a worker. Browser-upload destinations SHALL be directories, while missing browser-download sources SHALL be directory-capable.

#### Scenario: Hosted path preflight fails

- **WHEN** the browser-delivery path is missing, wrong-type, or cannot be created as authorized
- **THEN** Courier exits at preflight with code `20` and does not launch or register a worker

### Requirement: Discoverable detached browser lifecycle

Courier SHALL document and expose `--background` as the sole detached mode for browser downloads and uploads. It SHALL return after printing the URL and delivery UUID, survive terminal closure, and end through delivery or server stop, configured stop behavior, fatal worker failure, or process termination.

#### Scenario: Background browser delivery is ready

- **WHEN** a browser download or upload starts with `--background`
- **THEN** the initiating command exits after readiness and the delivery remains listed by `courier servers`


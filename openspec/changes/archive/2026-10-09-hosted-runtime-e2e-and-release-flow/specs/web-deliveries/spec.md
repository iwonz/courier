## ADDED Requirements

### Requirement: Browser delivery QR readiness

Interactive browser upload and download readiness SHALL include a scannable QR whose payload is the exact printed delivery URL when the terminal can render it. The QR SHALL NOT alter the bind address or advertise loopback as phone-reachable.

#### Scenario: Browser delivery is ready in a wide TTY

- **WHEN** `from web://` or `to web://` reaches readiness in a capable terminal
- **THEN** Courier shows the text URL and an exact-payload QR after the readiness panel

### Requirement: Externally stopped browser foreground

A foreground browser delivery SHALL exit successfully after its delivery or containing server is stopped through verified Courier control.

#### Scenario: Delivery UUID is stopped from another process

- **WHEN** `courier servers stop <delivery-uuid>` commits removal
- **THEN** the initiating browser-delivery process reports a neutral stop, exits `0`, and leaves no route or worker after the final delivery ends

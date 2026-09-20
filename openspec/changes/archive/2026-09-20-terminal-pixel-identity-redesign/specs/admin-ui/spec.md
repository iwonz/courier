## MODIFIED Requirements

### Requirement: Pixel administration identity

The administration React application SHALL render as a terminal operations console using mark-v3 and admin-v3. Its title region, API-derived metric strip, server navigator, selected delivery facts, and policy editor SHALL use flat surfaces and structural one-pixel rules rather than card islands. Every counter and identifier SHALL come from the authoritative runtime snapshot; no synthetic uptime, status chart, or telemetry SHALL appear.

#### Scenario: Live administration data renders

- **WHEN** a secret-free server snapshot is available
- **THEN** the operations console shows only its real server, delivery, confirmed-byte, route, and policy values while stop, refresh, selection, and save actions preserve existing semantics

#### Scenario: An operator manages live state

- **WHEN** the operator selects a delivery, changes policy, refreshes, or stops a target
- **THEN** the terminal console performs the existing guarded API operation and reconciles the authoritative snapshot without exposing secrets

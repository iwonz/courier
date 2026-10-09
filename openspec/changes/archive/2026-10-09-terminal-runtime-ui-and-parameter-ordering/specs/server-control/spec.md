## ADDED Requirements

### Requirement: Responsive interactive server inventory

`courier servers` SHALL render a TTY server overview with Status, UUID, Bind, PID, State, and Started/Updated values followed by a delivery table for each server with State, UUID, Route, Source/Destination, Read/Sent/Confirmed, Policy, and Created/Updated values. When the multi-column layout does not fit, Courier SHALL use field/value tables without truncating UUIDs, paths, or policy values.

#### Scenario: Inventory is wide

- **WHEN** the output terminal can display the complete multi-column layout
- **THEN** Courier prints one server overview and nested delivery tables without exposing secrets

#### Scenario: Inventory is narrow

- **WHEN** the complete multi-column layout would not fit
- **THEN** Courier switches to per-record field/value tables and preserves every identifier and value

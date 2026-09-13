# operation-routing Specification

## Purpose
Define side-effect-free route selection and option validation before Courier acquires runtime resources.

## Requirements

### Requirement: Side-effect-free route planning

Courier SHALL classify and validate an operation before opening local paths, SSH sessions, sockets, HTTP clients, or workers.

#### Scenario: Supported direction

- **WHEN** source and destination kinds form path-to-path, web-to-path, path-to-web, webhook-to-path, or path-to-HTTP
- **THEN** planning returns the corresponding typed route

#### Scenario: Unsupported direction

- **WHEN** source and destination kinds do not form a supported route
- **THEN** planning fails with a preflight route error and performs no I/O

### Requirement: Typed path direction

Courier SHALL distinguish local-to-local, local-to-SSH, SSH-to-local, and SSH-to-SSH path routes for policy applicability.

#### Scenario: Two SSH endpoints

- **WHEN** both path endpoints are SSH-backed
- **THEN** the planned direction is SSH-to-SSH

### Requirement: Option matrix validation

Courier SHALL validate option applicability, conflicts, defaults, and repeatability against the planned route using ordered option occurrences.

#### Scenario: Repeated non-repeatable option

- **WHEN** a non-repeatable option occurs more than once
- **THEN** planning fails before applying a default or starting runtime work

#### Scenario: Interleaved selection rules

- **WHEN** exclude and exclude-from options are interleaved
- **THEN** their values remain in command-line occurrence order for the selection engine

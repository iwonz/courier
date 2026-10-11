## MODIFIED Requirements

### Requirement: Localized structured parameter metadata

Courier SHALL maintain CLI contract schema version `3` and contract version `0.17.0` with English and Russian descriptions for every public positional argument and flag, explicit required or optional state, and the existing type, default, repeatability, dependency, conflict, and applicability metadata. All flags SHALL be optional, and argument requiredness SHALL remain command-specific.

#### Scenario: Upload filename policy is inspected

- **WHEN** a consumer reads an incoming browser or webhook route
- **THEN** it finds optional `--filename-policy <preserve|unique>`, default `preserve`, with `--extract` declared as a conflict

#### Scenario: Parameter metadata is incomplete

- **WHEN** a public argument or flag lacks either localized description or explicit requiredness
- **THEN** strict contract validation fails

#### Scenario: Directory creation flag is inspected

- **WHEN** a consumer reads a path-bearing `from` route
- **THEN** it finds optional boolean `--force-source-creation` with applicability to that route

## ADDED Requirements

### Requirement: Localized structured parameter metadata

Courier SHALL maintain CLI contract schema version `3` and contract version `0.16.0` with English and Russian descriptions for every public positional argument and flag, explicit required or optional state, and the existing type, default, repeatability, dependency, conflict, and applicability metadata. All flags SHALL be optional, and argument requiredness SHALL remain command-specific.

#### Scenario: Parameter metadata is incomplete

- **WHEN** a public argument or flag lacks either localized description or explicit requiredness
- **THEN** strict contract validation fails

#### Scenario: Directory creation flag is inspected

- **WHEN** a consumer reads a path-bearing `from` route
- **THEN** it finds optional boolean `--force-source-creation` with applicability to that route

### Requirement: Unified generated parameter documentation

Courier SHALL deterministically generate structured Cobra help metadata, the CLI reference, shipped landing JSON, and managed README argument and option tables from the canonical contract. Each presentation SHALL include syntax, Required or Optional status, a plain-language description, human-readable applicability, and any relevant default, repeatability, dependency, or conflict metadata.

#### Scenario: A generated consumer is stale

- **WHEN** contract metadata changes without regenerating terminal help metadata, reference documentation, landing JSON, or the managed README sections
- **THEN** the verification freshness gate fails

#### Scenario: Background help is rendered

- **WHEN** a user inspects help or generated documentation for `--background`
- **THEN** it explains supported browser downloads, browser uploads, incoming webhooks, and `ui start`, immediate return after URL and UUID output, survival after terminal closure, and supported stop conditions


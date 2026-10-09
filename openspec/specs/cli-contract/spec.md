# cli-contract Specification

## Purpose
Define the canonical machine-readable public CLI inventory and generated-reference parity gate.

## Requirements

### Requirement: Canonical command contract

Courier SHALL keep a versioned machine-readable English contract containing every public command and option combination plus ordered command paths and arguments. Arguments SHALL declare name, kind, required state, optional literal prefix, and optional suppressing flag. Parameters SHALL declare a value kind, optional choices, placeholder, dependencies, repeatability, conflicts, and default.

#### Scenario: Contract metadata is invalid

- **WHEN** an argument or parameter has an unknown kind, invalid choices, duplicate reference, impossible dependency, unknown suppressing flag, or inconsistent repeatability
- **THEN** strict contract validation fails

#### Scenario: Planned command

- **WHEN** a command is present in the target contract but is not implemented
- **THEN** it is marked planned and does not appear in live help

### Requirement: Generated reference parity

Courier SHALL deterministically generate its human command reference and validate the shipped Cobra tree against the canonical contract.

#### Scenario: Undocumented Cobra command

- **WHEN** Cobra exposes a command or flag not present as shipped or system behavior in the contract
- **THEN** the local verification gate fails

### Requirement: Generated shipped-only landing data

Courier SHALL deterministically generate landing command-builder data from the canonical CLI contract, preserving command, argument, parameter, choice, dependency, and conflict order while omitting planned inventory.

#### Scenario: Landing data is regenerated

- **WHEN** the canonical contract changes
- **THEN** generated JSON contains the exact structured path, arguments, and typed parameter metadata or the freshness gate fails

#### Scenario: A planned command remains in the contract

- **WHEN** the landing data is generated or checked
- **THEN** the planned command is omitted while shipped product commands and system commands remain represented according to their status

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

### Requirement: Stable required-first parameter ordering

Courier SHALL present required parameters before optional parameters in structured help, generated README tables, the CLI reference, and landing data while preserving canonical contract order inside each requiredness group. README and CLI reference option tables SHALL use separate Default, Repeatable, Applies to, Requires, and Conflicts columns.

#### Scenario: Mixed requiredness is generated

- **WHEN** a command contains required and optional parameters in interleaved contract order
- **THEN** every generated presentation lists the required parameters first and retains relative order within both groups

#### Scenario: Option relationships are documented

- **WHEN** a generated option has applicability, dependency, or conflict metadata
- **THEN** each metadata category appears in its own documentation column rather than a combined relationship field

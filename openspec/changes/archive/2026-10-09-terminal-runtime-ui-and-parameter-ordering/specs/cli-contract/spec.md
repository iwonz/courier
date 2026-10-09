## ADDED Requirements

### Requirement: Stable required-first parameter ordering

Courier SHALL present required parameters before optional parameters in structured help, generated README tables, the CLI reference, and landing data while preserving canonical contract order inside each requiredness group. README and CLI reference option tables SHALL use separate Default, Repeatable, Applies to, Requires, and Conflicts columns.

#### Scenario: Mixed requiredness is generated

- **WHEN** a command contains required and optional parameters in interleaved contract order
- **THEN** every generated presentation lists the required parameters first and retains relative order within both groups

#### Scenario: Option relationships are documented

- **WHEN** a generated option has applicability, dependency, or conflict metadata
- **THEN** each metadata category appears in its own documentation column rather than a combined relationship field

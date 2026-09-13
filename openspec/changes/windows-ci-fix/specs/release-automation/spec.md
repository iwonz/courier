## ADDED Requirements

### Requirement: Native Windows release gate

Courier SHALL run the complete first-party Go test suite on a native Windows runner, require exact 100% statement coverage, run the PowerShell installer acceptance test, and prevent publication if any native command fails.

#### Scenario: Windows test failure

- **WHEN** a Go test exits unsuccessfully on the Windows runner
- **THEN** the gate reports that test failure, skips coverage analysis and installer acceptance, and blocks publication

#### Scenario: Platform-specific representation

- **WHEN** a test observes a local path or filesystem metadata on Windows
- **THEN** it validates the Windows-supported representation without weakening content, integrity, safety, or cleanup assertions

#### Scenario: Corrective branch validation

- **WHEN** a branch matching `fix/**` is pushed
- **THEN** the normal quality and native Windows gates run before that change advances to `main`

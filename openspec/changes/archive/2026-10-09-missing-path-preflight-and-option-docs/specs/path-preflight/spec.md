## Purpose

Define safe, predictable handling for path endpoints that are absent before a local or hosted Courier operation starts.

## ADDED Requirements

### Requirement: Complete missing-directory inspection

Courier SHALL inspect every local or SSH path endpoint, classify its expected file or directory role, and report all missing directory endpoints before creating any of them.

#### Scenario: Source and destination directories are missing

- **WHEN** one operation has a missing directory-capable source and a missing directory destination
- **THEN** Courier requests or obtains authorization for both in source-then-destination order before creating either directory

### Requirement: Explicit directory-creation authorization

Courier SHALL ask `Source directory <path> does not exist. Create it? [y/N]` or `Destination directory <path> does not exist. Create it? [y/N]` for each missing directory unless `--force-source-creation` is present. The flag SHALL bypass only these directory-creation confirmations.

#### Scenario: Interactive creation is declined

- **WHEN** the user declines any missing-directory confirmation
- **THEN** Courier fails at preflight without creating any inspected directory

#### Scenario: Noninteractive creation lacks authorization

- **WHEN** a missing directory is found without an interactive terminal and without `--force-source-creation`
- **THEN** Courier exits at preflight with code `20` and instructs the user to rerun with `--force-source-creation`

#### Scenario: Forced creation is requested

- **WHEN** `--force-source-creation` is present and all missing path roles permit directory creation
- **THEN** Courier creates those directories without prompting and does not bypass any unrelated trust, credential, helper, authentication, collision, or overwrite protection

### Requirement: Safe ordered directory creation

Courier SHALL recursively create authorized directories in deterministic source-then-destination order with mode `0700`, revalidate each resulting object as a directory, and retain successfully created directories if a later creation or operation fails.

#### Scenario: A path appears during preflight

- **WHEN** another process creates a missing path before Courier creates or revalidates it
- **THEN** Courier continues only if the resulting object is a directory and otherwise returns a path-specific wrong-type failure

#### Scenario: A later creation fails

- **WHEN** Courier creates the source directory and then cannot create the destination directory
- **THEN** it returns a preflight failure and leaves the source directory intact

### Requirement: Local and SSH parity

Courier SHALL apply the same classification, authorization, ordering, creation, and revalidation behavior to local and SSH paths while preserving existing SSH trust and authentication requirements.

#### Scenario: Remote creation is authorized

- **WHEN** an SSH path has a missing directory role and creation is authorized
- **THEN** Courier creates it recursively with private permissions before transfer or hosted-worker acquisition


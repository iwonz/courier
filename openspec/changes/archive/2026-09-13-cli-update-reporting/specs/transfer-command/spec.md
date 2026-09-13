## ADDED Requirements

### Requirement: Transfer command grammar

Courier SHALL implement `courier from <source> to <destination> [flags]` as an isolated Cobra command handler and SHALL reject any other separator or argument count before connecting or writing.

#### Scenario: Valid grammar

- **WHEN** the arguments are `from ./data to server:/srv/data`
- **THEN** the handler parses exactly one source and one destination

#### Scenario: Invalid separator

- **WHEN** the middle argument is not `to`
- **THEN** Courier returns the CLI usage exit code without side effects

### Requirement: Four transfer directions

Courier SHALL support local-to-local, local-to-remote, remote-to-local, and remote-to-remote using the same transfer engine.

#### Scenario: Remote-to-remote direct stream

- **WHEN** both authenticated SFTP endpoints are available and archive mode is off
- **THEN** Courier streams through bounded memory from source SFTP to destination SFTP without local disk staging

#### Scenario: Remote-to-remote archive

- **WHEN** archive mode is enabled
- **THEN** Courier uses a private verified local archive stage and cleans it on every exit path

### Requirement: Independent remote setup

Courier SHALL resolve, verify, and authenticate distinct remote endpoints independently and may connect them concurrently with cancellation propagation.

#### Scenario: Second remote fails

- **WHEN** one endpoint connects and the other fails
- **THEN** the connected endpoint is closed and no transfer mutation begins

### Requirement: Archive flag

Courier SHALL accept `--archive` and transfer a verified `<source-name>.tar.gz` artifact instead of the source tree.

#### Scenario: Archive destination directory

- **WHEN** source is `photos`, destination is a directory, and `--archive` is set
- **THEN** actual destination ends with `photos.tar.gz`

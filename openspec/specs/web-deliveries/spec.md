# web-deliveries Specification

## Purpose
Define secure browser-based uploads and downloads that reuse Courier's endpoint,
selection, policy, worker, and transactional storage layers without exposing
credentials or protected metadata.

## Requirements

### Requirement: Opaque browser delivery routing

Courier SHALL route every browser delivery through a random 256-bit resource token that contains no delivery, endpoint, path, credential, or policy data and SHALL revoke that route when the delivery stops.

#### Scenario: Stopped token is requested

- **WHEN** a browser requests a resource token after its delivery stopped
- **THEN** the data host returns the same not-found response used for an unknown token

### Requirement: Protected metadata isolation

Courier SHALL complete peer admission and configured authentication before reading or returning protected source or destination metadata.

#### Scenario: Unauthenticated directory is requested

- **WHEN** a client without valid authentication requests a shared directory
- **THEN** no entry name, size, path, or archive metadata is rendered or embedded

### Requirement: Transactional browser upload

Courier SHALL reserve policy capacity before reading an upload, stage accepted bytes privately, and commit only to an absent safe destination name while preserving unrelated entries.

#### Scenario: Uploaded name conflicts

- **WHEN** the sanitized final upload name already exists
- **THEN** Courier rejects the upload without consuming or replacing that entry and removes its staging data

#### Scenario: Uploaded archive is extracted

- **WHEN** a browser upload delivery uses `--extract` and receives a safe tar.gz
- **THEN** Courier applies the shared extraction limits and selection policy and commits non-conflicting top-level entries into the destination root without retaining the uploaded archive

### Requirement: Safe browser navigation

Courier SHALL constrain requested relative paths to the selected source root, apply the shared selection engine, and reject traversal, symlink escape, and unsupported special files.

#### Scenario: Encoded traversal is requested

- **WHEN** a download or listing path resolves outside the shared root
- **THEN** Courier rejects it without opening the escaped object

### Requirement: Browser file and directory download

Courier SHALL stream selected files with bounded buffers and SHALL stream a selected directory as a verified-structure tar.gz with exactly one top-level source entry.

#### Scenario: Whole directory is downloaded

- **WHEN** an authorized client requests archive download for a directory
- **THEN** the response is a tar.gz whose safe entries are rooted under the source base name

#### Scenario: Explicit archive delivery is registered

- **WHEN** a path-to-browser delivery uses `--archive`
- **THEN** Courier creates and verifies one private `<source-name>.tar.gz` before registration and publishes only that object

### Requirement: Ephemeral SSH endpoint credentials

Courier SHALL preflight SSH-backed browser endpoints with normal interactive authentication, send only credentials actually requested and prior helper consent over private IPC, clear them after opening the worker resource, and never persist or log them.

#### Scenario: Encrypted identity is needed by a background delivery

- **WHEN** SSH preflight requests an identity passphrase and the background worker opens the same endpoint
- **THEN** the worker consumes the matching private prompt response without another terminal and the public registry remains secret-free

### Requirement: Shared browser policy enforcement

Courier SHALL use the delivery policy core for authentication attempts, bans or stop decisions, IP admission, concurrent reservations, file-size limits, CSRF, and aggregate directional rates.

#### Scenario: Concurrent uploads exceed the delivery limit

- **WHEN** an additional upload cannot reserve a configured delivery slot
- **THEN** it is rejected before its multipart body is consumed

### Requirement: UI and JSON surfaces

Courier SHALL serve the embedded React data application by default as a shared shadcn destination workspace with typed localization, cyclic theme and locale controls, protected authentication, verified route presentation, breadcrumbs, manifests, and semantic upload/download actions. Unauthenticated rendering SHALL contain only the real authentication form and generic state, and protected metadata SHALL appear only after authorization. Optional local Relay artwork SHALL remain decorative and SHALL reveal no delivery state. The application SHALL use only the existing guarded delivery API and SHALL expose only versioned JSON/download surfaces when `--no-ui` is active.

#### Scenario: Authentication is required

- **WHEN** an unauthenticated visitor opens a protected delivery
- **THEN** only a generic shadcn authentication card and real password form render without protected names, paths, sizes, routes, credentials, command input, or terminal chrome

#### Scenario: An authorized delivery opens

- **WHEN** authentication and metadata admission succeed
- **THEN** shadcn cards, buttons, badges, separators, and progress affordances present the verified route, navigation, entries, uploads, and downloads without changing API semantics

#### Scenario: An authorized browser delivery opens

- **WHEN** authentication and metadata admission succeed
- **THEN** the destination workspace presents the verified route, navigation, entries, and permitted actions without arbitrary command execution

#### Scenario: UI-disabled root is requested

- **WHEN** a browser opens a no-UI delivery root
- **THEN** Courier returns a JSON description without embedding HTML or protected metadata

### Requirement: Foreground and background ownership

Courier SHALL bind foreground browser deliveries to their initiating control lease and SHALL keep explicitly background deliveries active after the initiating CLI exits.

#### Scenario: Foreground Courier is interrupted

- **WHEN** its lease closes
- **THEN** only that browser delivery and its opened endpoint resources are removed

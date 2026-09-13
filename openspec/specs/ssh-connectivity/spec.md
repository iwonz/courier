# ssh-connectivity Specification

## Purpose
Define native SSH configuration, trust, authentication, routing, and helper fallback.

## Requirements

### Requirement: Native SSH configuration

Courier SHALL resolve `HostName`, `User`, `Port`, `IdentityFile`, agent, `known_hosts`, and `ProxyJump` from OpenSSH-compatible user configuration and explicit endpoint user values.

#### Scenario: SSH alias

- **WHEN** endpoint host matches an SSH config alias
- **THEN** Courier connects to the resolved hostname, user, port, identities, and jump chain

#### Scenario: Explicit endpoint user

- **WHEN** `[user@]` is present on the endpoint
- **THEN** that user overrides SSH config user resolution

### Requirement: Secure authentication

Courier SHALL try key and agent authentication without exposing secrets, and SHALL request a password interactively only after non-password authentication is unavailable or rejected.

#### Scenario: Password fallback

- **WHEN** key and agent authentication fail and an interactive prompt is available
- **THEN** Courier prompts without echo and retries authentication with the entered password

#### Scenario: Password flag attempt

- **WHEN** a user attempts to pass a password through CLI arguments
- **THEN** no supported flag accepts or logs it

### Requirement: Host-key verification

Courier SHALL verify every target and jump host against known_hosts and SHALL NOT provide an insecure-disable option.

#### Scenario: Unknown or changed key

- **WHEN** the presented host key is absent or differs from known_hosts
- **THEN** connection fails before authentication data or file content is transferred

### Requirement: ProxyJump

Courier SHALL create each ProxyJump hop as a verified native SSH client and independently authenticate the final endpoints used in remote-to-remote transfers.

#### Scenario: Jump-host route

- **WHEN** target configuration contains `ProxyJump bastion`
- **THEN** Courier verifies and authenticates the bastion before dialing and independently authenticating the target through it

### Requirement: Native SFTP filesystem

Courier SHALL expose remote files through the transfer backend without invoking rsync, scp, or a local ssh executable.

#### Scenario: Remote tree operation

- **WHEN** the transfer engine reads or writes a remote endpoint
- **THEN** filesystem operations use the authenticated SFTP subsystem with POSIX path handling

### Requirement: Optional remote helper consent

Courier SHALL use no remote helper when native SSH/SFTP and built-in archive behavior suffice. If an incompatible remote requires a helper, Courier SHALL explain why and require explicit interactive confirmation before downloading or placing it, then verify its checksum and clean temporary bootstrap files.

#### Scenario: SFTP is available

- **WHEN** the server provides a usable SFTP subsystem
- **THEN** Courier does not prompt for or deploy a helper

#### Scenario: Helper is required but declined

- **WHEN** Courier identifies a helper-only capability gap and the user declines
- **THEN** the operation stops without modifying the remote host

### Requirement: Verified temporary helper fallback

Courier SHALL offer a remote helper only after verified SSH authentication succeeds and native SFTP reports a capability gap. Courier SHALL require explicit interactive consent before downloading or writing helper data, verify the exact-platform helper locally and on the remote endpoint, and remove all staging data after success, error, or cancellation.

#### Scenario: Native SFTP succeeds

- **WHEN** the authenticated server starts a usable SFTP subsystem
- **THEN** Courier uses native SFTP without prompting, downloading, or staging a helper

#### Scenario: User declines the helper

- **WHEN** SFTP is unavailable and the user declines or no interactive terminal exists
- **THEN** Courier returns a connection error without downloading an artifact or modifying the remote endpoint

#### Scenario: User accepts the helper

- **WHEN** SFTP is unavailable and the user explicitly accepts the explained fallback
- **THEN** Courier selects the detected remote OS and architecture, verifies the helper, stages it privately, verifies it remotely, and serves the existing filesystem backend over the SSH channel

#### Scenario: Helper setup fails

- **WHEN** acquisition, upload, checksum verification, startup, transfer, or cancellation fails
- **THEN** Courier closes the helper process and attempts to remove every local and remote temporary artifact before closing SSH

### Requirement: Native Windows SSH agent

Courier SHALL connect to the Windows OpenSSH agent through its named pipe without invoking an SSH executable.

#### Scenario: Windows default agent

- **WHEN** Courier runs on Windows and no `SSH_AUTH_SOCK` override is present
- **THEN** it attempts the standard OpenSSH agent named pipe through the platform-native pipe client

#### Scenario: Explicit agent endpoint

- **WHEN** `SSH_AUTH_SOCK` identifies an agent endpoint
- **THEN** Courier passes that endpoint to the platform-specific native dialer

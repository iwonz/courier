## ADDED Requirements

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

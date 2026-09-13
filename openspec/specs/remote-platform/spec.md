# remote-platform Specification

## Purpose
Define remote OS, architecture, and archive capability detection.

## Requirements

### Requirement: Remote platform detection

Courier SHALL normalize remote Linux, macOS, BSD, and Windows operating systems and amd64/arm64 architectures from constant, non-user-controlled probe commands.

#### Scenario: Darwin arm64

- **WHEN** probe output reports `Darwin` and `arm64`
- **THEN** the normalized platform is `darwin/arm64`

#### Scenario: Windows amd64

- **WHEN** POSIX probing fails and the Windows probe reports AMD64
- **THEN** the normalized platform is `windows/amd64`

### Requirement: Archive capability selection

Courier SHALL select an available `tar`, `bsdtar`, `gtar`, or `tar.exe` appropriate to the detected platform, and SHALL fall back to Courier's built-in streaming archiver when none is available.

#### Scenario: No remote archiver

- **WHEN** no supported remote archiver is discovered
- **THEN** archive mode uses the built-in implementation without asking the user to configure the server

## MODIFIED Requirements

### Requirement: Cross-platform release acceptance

Courier SHALL run compiled runtime suites on Linux, macOS, and Windows, cross-build the declared primary and BSD helper artifacts, and install matching release packages in isolated Ubuntu, Debian, Arch Linux, Manjaro, Fedora, RHEL-compatible, and Alpine containers. Native Windows acceptance SHALL exercise owner-scoped named pipes using the current-process identity, a unique valid endpoint per test run, bounded connection contexts, cleanup registered before the handshake, rooted filesystem operations across the native-path and `fs.FS` path dialects, portable archive symlink targets, and platform-supported metadata semantics.

#### Scenario: A distribution package check exits

- **WHEN** installation succeeds, fails, or is interrupted
- **THEN** its cleanup trap removes every uniquely labeled Courier container and proves that no labeled container, network, volume, or temporary file remains

#### Scenario: Windows uses a rooted directory path

- **WHEN** a native Windows path produced by `filepath` crosses into an `fs.FS` operation
- **THEN** the adapter supplies slash-separated syntax without weakening the rooted confinement boundary

#### Scenario: Windows creates a private control pipe

- **WHEN** Courier resolves the named-pipe security principal
- **THEN** it queries the supported current-process token and grants pipe access only to that user SID

#### Scenario: Windows pipe acceptance stalls

- **WHEN** a native named-pipe test cannot complete its client-server handshake
- **THEN** its bounded context terminates the attempt and pre-registered cleanup releases the listener and accept goroutine without waiting for the package timeout

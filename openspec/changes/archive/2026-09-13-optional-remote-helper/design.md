# Design: Optional remote helper and Windows hardening

## Context

The authenticated `x/crypto/ssh` connection remains usable when `pkg/sftp` cannot start the server subsystem. Courier can use that verified channel to stage a matching Courier binary and run its embedded SFTP server. This preserves one filesystem protocol and one transfer engine instead of introducing a second set of copy semantics.

## Decisions

### Capability-driven fallback

`sshx.Factory` invokes a fallback callback only when native SFTP construction returns `CapabilityError{Capability: "sftp"}`. The callback serializes prompts, explains the gap, and defaults to refusal. Non-interactive execution fails without changing the remote host. No fallback code runs on a successful native SFTP connection.

### Exact helper artifact

When remote OS and architecture match the local process, Courier stages its current executable and computes SHA-256 locally. Otherwise it downloads the archive for the exact running semantic version from the corresponding GitHub Release plus `checksums.txt`, verifies the archive, safely extracts the single Courier executable, and computes the executable digest. Development builds cannot download a cross-platform helper because they have no immutable release identity.

### Private remote bootstrap

Courier uploads over an authenticated SSH session into an OS temporary directory created with private permissions. Fixed POSIX shell or PowerShell bootstrap programs accept no endpoint path as executable syntax. They compare the uploaded file digest with the locally computed digest before starting it. The helper is invoked through a hidden internal entry point and serves `pkg/sftp` over standard input/output.

The helper runtime owns its SSH session and remote staging path. Closing the endpoint closes SFTP, waits for the helper, removes the staging directory, and only then closes SSH hops. Setup failures execute the same cleanup path.

### Windows agent

Platform-specific files select a Unix socket dialer outside Windows and `go-winio` named-pipe dialing on Windows. Windows defaults to the OpenSSH agent pipe when `SSH_AUTH_SOCK` is absent. Authentication policy and `ssh/agent` use remain shared.

### Windows update handoff

On Windows, the updater stages and verifies the new executable next to the target, then launches it in handoff mode. The handoff waits for the original PID, copies itself to a second private partial, atomically replaces the target, and launches the installed target in cleanup mode. Cleanup waits for the handoff PID and deletes the original staged executable. Hidden modes validate that staging and target paths share a directory and that staging names use Courier's private prefix.

## Security notes

- Existing host-key verification and independent endpoint authentication happen before helper logic.
- Consent precedes every download or remote mutation.
- Checksums are required locally and remotely; mismatch prevents execution.
- Bootstrap commands contain only validated digests and quoted generated paths.
- The embedded SFTP implementation runs with the authenticated remote user's privileges and creates no persistent service.
- Helper and update staging is private and cleanup is attempted after success, failure, and cancellation.

## Testing

- Unit tests cover refusal, non-interactive behavior, artifact selection, checksum failures, quoting, cleanup, named-pipe selection, and handoff validation.
- An in-process SSH integration server disables its SFTP subsystem, accepts the helper upload, and exercises the fallback lifecycle without external infrastructure.
- Cross-build tests compile all six release targets.
- The repository-wide gate retains exact 100% first-party Go statement coverage and race testing.

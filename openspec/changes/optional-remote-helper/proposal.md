# Change: Add an optional remote helper and harden Windows runtime behavior

## Why

Courier's native SSH transport is intentionally SFTP-first, but some correctly authenticated SSH servers do not expose an SFTP subsystem. Those servers need an explicit, auditable fallback instead of a generic connection failure. Windows also needs native OpenSSH agent named-pipe support and a self-update handoff because a running executable cannot replace itself reliably.

## What Changes

- Offer a temporary Courier helper only after a typed SFTP capability failure and explicit interactive consent.
- Acquire the helper for the detected remote OS and architecture, verify its SHA-256 digest, upload it through the already verified SSH connection, and remove every local and remote staging artifact.
- Run the helper as an SFTP server over SSH standard streams so the existing transfer engine and public transport boundary remain unchanged.
- Add native Windows OpenSSH agent named-pipe dialing.
- Add a two-process Windows self-update handoff that replaces the executable after the original process exits and then removes staging files.
- Keep internal helper and update handoff entry points hidden from the public Cobra command surface.

## Non-goals

- Replacing native SFTP when the server already provides it.
- Installing a persistent remote daemon or modifying the remote `PATH`.
- Reimplementing SSH, SFTP, or archive formats.
- Claiming complete behavioral compatibility with an OpenSSH executable.

## Impact

- Affected specs: `ssh-connectivity`, `self-update`.
- Affected packages: `internal/sshx`, `internal/helper`, `internal/update`, `internal/app`, and `cmd/courier`.
- New dependency: `github.com/Microsoft/go-winio` for the Windows named-pipe client only.

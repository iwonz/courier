## Context

The transfer engine already consumes an abstract filesystem. Native SSH creates a secure connection and SFTP adapts the remote filesystem to that contract. No user-controlled path is placed into a shell command.

## Goals / Non-Goals

**Goals:** OpenSSH-compatible configuration, strict trust, key/agent/password auth, ProxyJump, SFTP backend, platform probes, and capability-driven helper policy.

**Non-Goals:** silently installing server software, accepting command-line passwords, or disabling host-key verification.

## Decisions

### Established Go SSH, config, and SFTP libraries

`x/crypto/ssh`, `kevinburke/ssh_config`, and `pkg/sftp` are linked into each release binary. Courier never implements the SSH protocol and never shells out to local `ssh`, `scp`, or `rsync`.

### Strict trust before password fallback

Every connection configuration receives a `knownhosts` callback. Password prompting occurs only after a key/agent attempt returns an authentication error, and prompt values are held in memory only.

### Jump hosts are real SSH clients

ProxyJump chains dial the next TCP endpoint through the previous authenticated client. Every hop resolves its own configuration and host-key callback. Both remote endpoints create independent chains.

### SFTP first, helper only by explicit consent

The normal path uses SFTP and built-in archive streaming. A typed capability error may trigger an interactive helper proposal later in command orchestration. There is no automatic installation and no helper prompt when SFTP works.

### Constant remote probes

Platform and archiver commands are fixed literals. User paths and endpoint text never enter these commands, preventing command injection.

## Risks / Trade-offs

- Some servers disable SFTP; those produce a typed capability error rather than an unsafe shell fallback.
- OpenSSH config is broad. The established parser handles matching and includes; Courier explicitly applies the required connection fields without promising perfect OpenSSH client parity.
- Remote timestamp and permission behavior may be limited by server policy or filesystem capabilities.

## Migration Plan

No server-side state is created on the default SFTP path. Connection objects close SFTP and every SSH hop in reverse order.

## Open Questions

None.

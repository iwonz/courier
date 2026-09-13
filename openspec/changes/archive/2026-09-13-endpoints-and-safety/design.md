## Context

Endpoint parsing and path safety run before SSH connection or mutation. The same contract must work on macOS, Linux, and Windows without interpreting spaces or shell metacharacters.

## Goals / Non-Goals

**Goals:** a pure endpoint model; cross-platform naming semantics; canonical safety checks; safe archive names.

**Non-Goals:** SSH config resolution, filesystem stat, and actual copying.

## Decisions

### Remote endpoints require the strict `:/` delimiter

A remote endpoint has a non-empty authority without path separators and an absolute slash path. Thus `C:\...`, `C:/...`, and local colon-containing names do not become remote accidentally. Bracketed IPv6 is supported.

### Destination receives stat as input

The pure function accepts `destinationIsDirectory`, which the appropriate backend will compute later. Product semantics stay separate from I/O.

### Canonical local paths resolve the existing ancestor

For a destination that does not exist yet, symlinks in its nearest existing ancestor are resolved. Remote comparisons happen after SSH config resolution using the effective user/host and normalized POSIX path.

## Risks / Trade-offs

- Remote Windows SFTP still uses slash paths; the server exposes native drive paths inside its absolute SFTP namespace.
- Local path case folding follows the runtime OS, matching the semantics of Courier's host.

## Migration Plan

There is no persisted data and rollback requires no migration.

## Open Questions

None.

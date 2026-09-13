## Context

Copy behavior must not depend on local OS APIs, SFTP details, or rsync. Atomic replacement and a complete staged tree naturally implement synchronization while keeping the previous destination untouched until commit.

## Goals / Non-Goals

**Goals:** one backend contract; recursive metadata-preserving copy; deterministic cleanup; verified built-in archives; renderer-neutral progress.

**Non-Goals:** SSH connectivity, interactive terminal rendering, and release packaging.

## Decisions

### Backend-neutral tree engine

`fsx.Backend` exposes only the file operations required by transfer and archive code. Paths are joined by the backend, so POSIX SFTP paths never inherit client OS separator rules.

### Stage beside the final destination

The engine writes a randomly named, mode-0700 partial sibling. If the target exists, it is renamed to a backup, the partial is renamed into place, then the backup is removed. A failed second rename attempts restoration. Deferred cleanup covers every return path.

### Complete replacement implements sync

Building a complete source snapshot and replacing the target removes stale destination entries without mutating the visible tree during copy. Re-running after interruption is safe even without offset resume.

### Built-in archive fallback is always available

The Go tar/gzip implementation reads any `fsx.Backend`, so it also serves remote sources when no suitable server archiver exists. Verification streams the resulting archive to EOF to force checksum validation.

## Risks / Trade-offs

- Complete staging requires destination space roughly equal to the source plus the previous target until commit.
- Cross-backend transfers stream through Courier; remote-to-remote optimization can choose the same engine without local disk staging for non-archive mode.
- Some backends cannot set every metadata field; unsupported metadata operations must return an explicit capability error or documented no-op.

## Migration Plan

No persisted schema exists. Partial names are namespaced and always cleaned opportunistically.

## Open Questions

None.

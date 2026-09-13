## Context

The reusable engine and transports are complete, but direct package calls are not a product. This change assembles them with established CLI/progress/concurrency libraries and keeps orchestration testable through injected interfaces.

## Goals / Non-Goals

**Goals:** exact grammar, four directions, secure prompting, destination preflight, bounded local access, terminal/non-terminal reports, stable exit codes, and verified updates.

**Non-Goals:** package-manager publication and release pipeline configuration, which remain in the distribution change.

## Decisions

### Cobra owns command and flag parsing

Root construction adds isolated child commands. The transfer handler validates its special infix `to` grammar and owns `--archive`; no password flag exists.

### Structured progress remains the source of truth

Courier counts confirmed bytes in its engine. An mpb adapter renders those events rather than owning byte accounting, which prevents UI behavior from changing correctness.

### Independent connections use errgroup

Source and destination remotes connect independently under one cancelable context. Cleanup is registered immediately and always runs before command return. Plain remote-to-remote transfer streams between SFTP backends; archive mode uses the verified private local artifact.

### Local backends are rooted

Each local endpoint opens an `os.Root` at its nearest existing parent and addresses the endpoint relatively. Kernel-assisted containment blocks symlink escapes and narrows TOCTOU exposure beyond lexical checks. Documented OS limitations remain explicit.

### Updates are archive-plus-checksum transactions

The updater parses GitHub metadata, selects deterministic GoReleaser names, verifies `checksums.txt`, extracts exactly the Courier executable with traversal checks, writes a same-directory partial, and renames. It never executes downloaded installers.

## Risks / Trade-offs

- In-place Windows executable replacement may require a handoff process; the updater reports a stable update error when the OS denies replacement.
- `os.Root` intentionally rejects source symlinks escaping the selected root, which is safer than silently broadening authority.
- GitHub API rate limits unauthenticated checks; the command may use a user-provided `GITHUB_TOKEN` from the environment without logging it.

## Migration Plan

Existing internal registry code is removed after Cobra parity tests. No user data migration is required.

## Open Questions

None.

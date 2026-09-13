## Context

The project must remain self-contained, cross-platform, and extensible. The foundation makes no SSH, filesystem, or rsync decisions; it only maps top-level arguments to isolated handlers.

## Goals / Non-Goals

**Goals:**

- minimal Go entrypoint;
- deterministic command registry;
- stable startup error model;
- injectable build metadata;
- a plan of sequential OpenSpec changes.

**Non-Goals:**

- transfer and SSH implementation;
- release artifact publication.

## Decisions

### Small first-party command registry

The `Command` and `Registry` interfaces replace a global switch or heavyweight CLI framework. Dependencies stay explicit and the parser remains fully testable.

### Output and environment are injected

Handlers receive `io.Writer` values, while the entrypoint only handles signals, exit codes, and dependency assembly. Commands can therefore be tested without subprocesses.

### Stacked task branches

Each task uses `feat/NNN-kebab-case`, contains its own `openspec/changes/<name>` path, and produces one conventional commit. Every next branch starts from the previous one, so the final branch contains the complete auditable history.

## Risks / Trade-offs

- A first-party parser owns help and validation, but Courier's grammar is small and fully tested.
- Stacked branches simplify final assembly, but pull requests must merge in order.

## Migration Plan

No migration is required for a new repository. Reverting the initial commit rolls back the change.

## Open Questions

None.

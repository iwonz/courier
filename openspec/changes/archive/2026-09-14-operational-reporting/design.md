# Context

Courier already has renderer-neutral progress stages, transactional transfer errors, three delivery counters, private immutable history files, and command-level exit mapping. The pieces need one semantic model and one redaction boundary before the final acceptance and release work.

# Decisions

## Monotonic accounting

`progress.Event` carries read, sent, and confirmed counters. `Current` remains a compatibility projection of confirmed bytes for the progress-bar integration, and every emission normalizes it from confirmed state. Tracker methods enforce non-negative, overflow-safe increments and the invariant `confirmed <= sent <= read`. The existing `Add` operation records a locally read, sent, and confirmed chunk atomically; transport-specific code can advance individual legs.

Progress stages are a closed public vocabulary. Events expose the active stage, total payload bytes, elapsed time, and speed calculated from confirmed bytes. Renderers show all counters in deterministic line mode and use confirmed bytes for completion in interactive mode.

## Diagnostics and exit classification

A small dependency-free diagnostics package sanitizes arbitrary error text before it reaches reports, IPC, or history. It removes URL user-info and query/fragment data and masks conventional password, token, secret, authorization, and cookie assignments. Reports retain useful stage context without echoing sensitive values.

Command execution classifies context cancellation before ordinary command categories and returns 130. CLI, connection, transfer, update, and control codes retain their existing values. Success summaries report confirmed bytes; failures report all known counters when available.

## Bounded private history

History entries retain counter snapshots and optionally a controlled stage and sanitized message. Validation enforces monotonic counter relationships and the closed stage vocabulary. `AppendHistory` writes the immutable private event first and then retains only the newest configured number of events. Rotation deletes only filenames and private regular files owned by Courier, orders by timestamp plus UUID, and syncs the state directory.

The default retention is 256 events. Tests can construct a store with an explicit positive retention without changing the on-disk schema.

# Risks / Trade-offs

Read, sent, and confirmed can differ for network protocols, but local filesystem writes become confirmed at the successful write boundary and therefore advance together. A crash between appending a history event and pruning old files can temporarily exceed the retention bound; the next append repairs it without losing the new event.

# Migration Plan

Add the OpenSpec delta, implement diagnostics and accounting primitives, migrate renderers and stream writers, add bounded history behavior, exercise exact coverage and races, then archive the change and commit once.

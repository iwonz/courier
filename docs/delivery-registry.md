# Delivery registry

Courier's delivery registry is the durable discovery and recovery boundary for long-lived HTTP deliveries. It is intentionally separate from worker IPC: a registry record helps locate a worker, but only a live, authenticated control-channel exchange can prove that the worker still owns an endpoint.

## State model

The schema is versioned and contains five record groups:

- servers, identified by UUID, with bind and control endpoints, diagnostic process metadata, lifecycle state, and timestamps;
- deliveries, identified by UUID and linked to a server UUID, with route, policy, counters, lifecycle state, and timestamps;
- tombstones, which retain the terminal reason for a removed server or delivery;
- owned temporary resources, linked to an active or tombstoned owner so cleanup can continue after interruption;
- immutable history events with target UUID, event kind, timestamp, and final counter snapshot.

Every identity must be a non-zero canonical UUID. Snapshots reject duplicate identities, deliveries whose server is absent, and temporary resources whose owner is unknown. Server and delivery lifecycle transitions use the same validated state machine:

```text
starting -> active -> stopping -> stopped
    |          |          |
    +----------+----------+-> failed -> stopped
```

Read, sent, and confirmed counters are non-negative and monotonic. Their live accumulator is safe for concurrent updates; persisted snapshots contain plain immutable values.

## Policy and secrets

The durable policy contains only controlled values: authentication mode, attempt and failure behavior, delivery and size limits, aggregate transfer rates, canonical IP prefixes, and UI visibility. Policy changes use an optimistic version: a writer must supply the current version and exactly the next version.

The schema has no fields for passwords, Basic credentials, session keys, opaque resource tokens, URL credentials, or free-form error text. Those values may later cross private IPC when required, but they must never be written to registry or history files.

## Persistence protocol

The default state directory is the platform user configuration directory under `courier/state`. Courier verifies a real directory, sets private directory permissions where the operating system supports them, and opens a rooted filesystem handle so state object names cannot escape that directory.

A registry update uses this no-replace protocol:

1. Load the highest committed revision and compare it with the caller's expected revision.
2. Apply and fully validate the mutation in memory.
3. Write a uniquely named `0600` temporary file and flush it.
4. Atomically link it to the previously absent revision name.
5. Remove the temporary name and flush the directory.

Revision filenames contain a fixed-width sequence. A concurrent writer cannot replace an existing revision; it receives a revision conflict and must reload. Interrupted temporary files do not supersede the previous complete revision. History uses the same private immutable commit protocol and is ordered by timestamp and event UUID when read.

Courier flushes each immutable state file before committing it. It also synchronizes the containing directory on operating systems that expose portable directory `fsync` semantics. Windows does not support `Sync` on an `os.File` opened for a directory, so the post-commit directory flush is a documented best-effort boundary there; an unsupported directory flush never turns a successfully flushed no-replace state commit into `Access is denied`.

State objects must be private regular files on systems with POSIX permission bits. Symlinks, malformed JSON, unknown JSON fields, filename/content mismatches, invalid references, and invalid lifecycle or policy values fail closed.

## Trust and cleanup boundaries

A stored PID is diagnostic metadata only. Courier never sends a signal, reports a server as live, or transfers lease ownership based on a PID or registry state alone. Worker lifecycle code must confirm ownership through the versioned private IPC endpoint; this also prevents PID-reuse mistakes.

Owned temporary metadata records only Courier-created local or remote paths. The component that creates a temporary resource must register it, clean it on success, failure, and cancellation, and remove its registry record only after cleanup succeeds. Stale recovery may act only on validated Courier-owned entries.

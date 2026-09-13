# Change: Add the delivery registry domain

## Why

Long-lived browser and webhook deliveries need stable identities and durable state before worker processes, control commands, or HTTP policy enforcement can be implemented. Persisting ad hoc process records would make PID reuse, partial writes, leaked secrets, and ambiguous cleanup unavoidable.

## What Changes

- Add UUID identities and validated server, delivery, policy, counter, tombstone, history, and owned-temporary-resource models.
- Add explicit lifecycle transitions and optimistic policy versions.
- Add concurrency-safe monotonic live counters with immutable snapshots.
- Add a private application state directory backed by rooted filesystem access.
- Persist registry revisions and history events as immutable, atomically committed files rather than overwriting live state.
- Validate all references and reject credentials, control characters, duplicate identities, invalid transitions, and malformed policy data before persistence.

## Impact

This change introduces an internal domain and persistence boundary only. It does not start workers, open sockets, trust PIDs as liveness evidence, expose public commands, or transmit IPC messages.

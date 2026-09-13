## 1. IPC boundary

- [x] 1.1 Add strict bounded versioned frames, stable operations, and stable IPC error codes
- [x] 1.2 Add private Unix-socket and Windows named-pipe transports with context deadlines
- [x] 1.3 Add handshake, registration, listing, policy, stop, shutdown, and subscription messages

## 2. Worker ownership

- [x] 2.1 Add bind-scoped startup locks and live compatibility probing
- [x] 2.2 Add foreground connection leases, background ownership, and last-delivery shutdown
- [x] 2.3 Add deterministic stale reconciliation and owned-temporary cleanup

## 3. Quality

- [x] 3.1 Cover framing, malformed peers, races, lease loss, shared binds, incompatible listeners, and shutdown at 100%
- [x] 3.2 Document IPC trust, worker ownership, and recovery behavior
- [x] 3.3 Run full verification, strict OpenSpec validation, archive the change, and commit once

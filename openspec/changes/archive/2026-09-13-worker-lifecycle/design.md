# Context

Multiple CLI invocations may create deliveries on the same bind concurrently. The durable registry from task 019 is useful for discovery, but it is not a lock and its PIDs are not trustworthy liveness evidence. Courier needs a cross-platform control boundary that survives independent processes without exposing helper commands publicly.

# Decisions

## Private IPC is the authority

Control endpoints are Unix sockets inside Courier's private state directory or Windows named pipes restricted to the current user. Every connection begins with a protocol handshake containing the expected server UUID, protocol version, and worker compatibility fingerprint. A PID, open TCP port, or registry state never substitutes for that exchange.

Messages use a four-byte length prefix and strict JSON. Both requests and responses are capped before allocation, carry a request UUID, and use stable operation and error codes. Deadlines and context cancellation bound all reads, writes, and dials. Secrets may be added to private IPC payloads by later policy work, but are never placed in argv, URLs, registry files, or logs.

## Startup is serialized per canonical bind

A SHA-256 digest of the canonical bind names a private advisory lock. A coordinator holds that lock while it loads discovery state, probes matching server records, reconciles stale records, and either registers with one compatible live worker or invokes the injected worker launcher. The launcher contract reports readiness only after the control listener and durable server record exist.

An incompatible live worker fails without mutation. A foreign TCP listener cannot satisfy the control handshake and is never adopted. Distinct binds use distinct locks and can start in parallel.

## Foreground ownership follows a control connection

A foreground registration receives a lease UUID and keeps its lease connection open. The runtime owns that lease until an explicit release or connection loss. Losing the lease stops only its delivery. Background deliveries have no foreground lease and survive the initiating CLI process.

Stopping a delivery is idempotent at the runtime boundary. When the final delivery and lease are gone, the worker shuts down unless a keepalive owner is registered. Server stop ends all deliveries before the server tombstone and listener shutdown.

## Recovery acts only on validated owned state

If a registry record cannot complete a matching live handshake while its bind lock is held, the coordinator marks its deliveries and server stale. Cleanup is attempted only for `OwnedTemp` records that passed registry validation, through a location-aware injected callback. A cleanup record is removed only after the callback succeeds; failures remain discoverable and abort recovery.

# Risks / Trade-offs

OS-private IPC authenticates the local account rather than introducing a second persisted secret. Any process running as that account can control Courier, which matches the permissions of Courier's state and destination files. Remote administration is deliberately excluded.

Advisory locks serialize cooperating Courier processes, not hostile same-user processes. No registry mutation or signal is attempted until strict state validation and the control handshake decision complete.

# Migration Plan

Add and test the protocol codec, platform transports, runtime state machine, bind coordinator, leases, and stale reconciliation. Cross-build every supported client target, run race tests and full verification, then archive this change before adding policy middleware.

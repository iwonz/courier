# Context

The private registry is durable discovery state, not proof that a process still owns a listener or UUID. Conversely, enumerating operating-system PIDs cannot prove that a reused PID is Courier. Server control therefore needs both registry ownership and a versioned IPC handshake before displaying an entry as live or issuing a stop.

# Decisions

## Secret-free discovery metadata

New server records persist their worker compatibility identifier. New delivery records persist source and destination display strings. Hosted resource tokens, runtime definitions, credentials, HTTP Authorization values, prompt responses, and helper consent remain ephemeral and are never written to the registry. Older records without display metadata remain valid and render an explicit unavailable value.

## Authoritative inventory

A control service loads the private per-user registry, sorts identities deterministically, and probes each recorded control endpoint using the recorded server UUID and compatibility. A successful hello followed by a valid live snapshot supplies the displayed state and nested deliveries. A failed probe never becomes a live claim; the registry entry is retained and labeled unreachable for later automatic stale reconciliation.

## Scoped stop resolution

`stop <uuid>` parses one canonical UUID, resolves whether it is a server, delivery, or tombstoned target, identifies the owning server, repeats the authoritative hello/list check, and then uses the matching IPC stop operation. Delivery stop affects only that delivery. Server stop affects all deliveries owned by that worker. A known tombstone is an idempotent success; an unknown UUID is an error.

`stop --all` snapshots all registered data servers, independently verifies and stops each one, continues after individual failures, and returns a joined control error if any target could not be authoritatively stopped. It does not inspect, signal, or stop the separate administrative UI.

# Risks / Trade-offs

An unreachable registry entry may represent a crashed worker or a transient local IPC failure. Listing it as unreachable is honest and non-destructive; cleanup remains governed by stale reconciliation rather than unsafe PID signaling. A server can disappear between probe and stop, so the command reports the IPC result and a later idempotent retry can observe its tombstone.

# Migration Plan

Extend optional registry metadata compatibly, add and test the control service, register isolated Cobra providers, update the canonical contract/reference, run full verification, archive this change, and commit it once.

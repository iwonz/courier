# Worker lifecycle and private IPC

Courier uses one long-lived worker for compatible deliveries that share a canonical TCP bind. Worker execution is an internal runtime mode and is deliberately absent from the public Cobra command tree.

## Startup and bind ownership

Before inspecting or starting a worker, the coordinator takes an advisory lock whose filename is the SHA-256 digest of the canonical bind. Locks for different binds are independent. While holding that lock, the coordinator:

1. loads and validates the delivery registry;
2. probes matching server records through their private control endpoints;
3. reuses a worker only after its server UUID and compatibility fingerprint match;
4. reconciles an unreachable stale record and its validated owned temporary resources; or
5. starts one detached worker and waits for its control handshake before registering the delivery.

The worker reserves the requested TCP bind before it reports readiness. A foreign listener therefore causes startup to fail; it is never mistaken for Courier. The stored PID is diagnostic and is never probed or signalled as an ownership test.

Launch configuration is a short-lived private file containing only the state directory, bind, UUID, control endpoint, and compatibility fingerprint. The child removes it before opening listeners. The child receives a small allowlisted environment, so unrelated credentials such as npm or GitHub tokens are not inherited. Future delivery secrets must be sent over private IPC after startup and must not be placed in the launch file, environment, argv, URL, registry, or log output.

## Control transport and protocol

Unix-like systems use a `0600` Unix-domain socket inside Courier's private state directory. Windows uses a named pipe with a current-user-only security descriptor. This operating-system access control is the local authentication boundary; Courier does not expose remote administration through this transport.

Every message has a four-byte network-order length followed by strict JSON. Frames are limited to 1 MiB before payload allocation. Envelopes contain protocol version `1`, a canonical non-zero request UUID, an operation, and an optional typed payload. Unknown fields, unknown operations, unsupported versions, malformed UUIDs, trailing JSON, and oversized frames fail before dispatch.

Stable operations cover:

- handshake and compatibility probing;
- delivery registration and lease claiming/release;
- server and delivery snapshots;
- optimistic policy updates;
- delivery stop, server stop, and shutdown;
- bounded progress subscriptions.

Responses repeat the request UUID and contain either a typed payload or one stable error code. Internal failures are redacted at the protocol boundary. Calls use context-aware dials and deadlines; subscriptions and lease connections remain open only for their ownership lifetime.

## Foreground and background deliveries

A foreground delivery is associated with a random lease UUID and the control connection held by its initiating CLI. Explicit lease release or connection loss stops only that delivery. An independent background delivery has no foreground lease and survives the initiating process.

Delivery stop is idempotent. Server stop commits terminal state for every delivery before tombstoning the server. A worker shuts down after its final delivery ends unless a registered keepalive owner remains; releasing the final keepalive completes shutdown. Listener closure, lease closure, registry updates, and subscription cancellation are race-tested.

## Stale recovery and cleanup

A matching PID or open TCP port cannot make a registry record live. Only a valid private-IPC handshake with the recorded server UUID can do that. Connection refusal, a missing endpoint, or a closed local transport identifies a stale record; timeouts, malformed responses, permission failures, and other uncertain errors fail closed without deleting state.

Recovery first verifies or removes only the expected Courier control endpoint. It then processes `OwnedTemp` records whose validated owner is the stale server or one of its deliveries. Cleanup is location-aware and injected by the component that knows how to remove that local or remote resource. Each ownership record is removed only after cleanup succeeds. A failure stops reconciliation and leaves the remaining metadata available for a safe retry. Finally, deliveries and the server receive `stale` tombstones through an optimistic immutable registry update.

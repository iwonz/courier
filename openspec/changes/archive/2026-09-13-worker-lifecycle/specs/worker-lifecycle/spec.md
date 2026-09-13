## ADDED Requirements

### Requirement: Authoritative private IPC

Courier SHALL treat a versioned handshake over a private local control transport as the only authority that a discovered worker is live and owns its server UUID.

#### Scenario: Registry PID was reused

- **WHEN** a server record has a running PID but its control endpoint cannot complete the expected handshake
- **THEN** Courier does not adopt or signal that process and reconciles the record as stale

### Requirement: Bounded versioned protocol

Courier SHALL reject unsupported protocol versions, unknown operations, duplicate or malformed request identities, unknown JSON fields, and frames larger than the fixed protocol limit before dispatch.

#### Scenario: Peer advertises an oversized frame

- **WHEN** the length prefix exceeds the protocol limit
- **THEN** Courier rejects the message without allocating the advertised payload or mutating lifecycle state

### Requirement: Bind-safe worker reuse

Courier SHALL serialize startup per canonical bind and SHALL reuse a discovered worker only after a live handshake confirms the expected UUID and a compatible runtime fingerprint.

#### Scenario: Concurrent starts target one bind

- **WHEN** two Courier processes concurrently request compatible deliveries on the same bind
- **THEN** exactly one worker owns the bind and both deliveries register with that worker

### Requirement: Foreground lease ownership

Courier SHALL bind every foreground delivery to one private control connection and SHALL stop that delivery when its lease is released or the connection is lost.

#### Scenario: Foreground owner is interrupted

- **WHEN** the initiating CLI process loses its lease connection
- **THEN** the worker stops only the leased delivery and preserves independent background or foreground deliveries

### Requirement: Scoped shutdown lifecycle

Courier SHALL make delivery stop idempotent, stop all deliveries before server shutdown, and stop a worker after its last delivery ends unless an explicit keepalive owner remains.

#### Scenario: Final delivery ends

- **WHEN** a worker has no delivery, foreground lease, or keepalive owner remaining
- **THEN** it commits terminal state, closes its control listener, and exits cleanly

### Requirement: Safe stale reconciliation

Courier SHALL reconcile stale server and delivery records while holding the bind startup lock and SHALL invoke cleanup only for validated Courier-owned temporary records.

#### Scenario: Stale owned temporary cleanup fails

- **WHEN** cleanup of a validated owned temporary resource returns an error
- **THEN** recovery fails closed and retains the ownership record for a later safe retry

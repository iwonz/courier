# worker-lifecycle Specification

## Purpose
Define cross-platform worker discovery, private IPC authority, bind-safe reuse, foreground lease ownership, scoped shutdown, and fail-closed stale recovery.

## Requirements

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

### Requirement: Platform-valid Unix control endpoints

Courier SHALL create new Unix control sockets inside its private state directory using a deterministic, lossless compact representation of the full server UUID that fits the supported platform pathname limit under Courier's default state-directory layout. Courier SHALL continue using a persisted legacy control endpoint exactly as recorded, and SHALL leave Windows named-pipe generation unchanged.

#### Scenario: Administration starts from the default macOS state directory

- **WHEN** Courier generates an administration control endpoint beneath `~/Library/Application Support/courier/state`
- **THEN** the endpoint fits Darwin's Unix-socket pathname boundary and the private listener can start

#### Scenario: A web delivery starts from the default macOS state directory

- **WHEN** Courier acquires a new data worker for a web delivery
- **THEN** the worker receives the same deterministic compact endpoint that its launch configuration validates

#### Scenario: A legacy process remains live across upgrade

- **WHEN** persisted state records a live `control-<UUID>.sock` endpoint created by an older Courier version
- **THEN** the upgraded client probes, controls, and removes that exact recorded endpoint without rewriting it

### Requirement: Hosted endpoint validity before acquisition

Courier SHALL complete hosted endpoint preflight before acquiring a worker, while reusing SSH trust, authentication, helper preparation, and bounded runtime credentials without persisting new secrets.

#### Scenario: Hosted endpoint is invalid

- **WHEN** hosted preflight finds a missing unauthorized path, a file in a directory role, or a directory in a file-only role
- **THEN** no worker is launched, acquired, or registered

### Requirement: Safe registration-time path failures

Courier SHALL revalidate the hosted path during registration and return a stable, sanitized, path-specific preflight error if the path disappeared or changed type. Arbitrary filesystem and worker failures SHALL remain redacted.

#### Scenario: Path disappears during registration

- **WHEN** a successfully preflighted hosted directory disappears before worker registration opens it
- **THEN** Courier returns an actionable preflight failure rather than `Courier IPC internal: internal worker error`

#### Scenario: Unclassified filesystem failure occurs

- **WHEN** registration encounters an arbitrary filesystem failure that is neither absence nor wrong type
- **THEN** Courier does not expose the underlying filesystem detail

### Requirement: Observable foreground termination

Courier SHALL retain each progress subscriber's delivery filter, publish a final filtered snapshot or tombstone before closing subscribers, and classify foreground termination from that evidence.

#### Scenario: Delivery is stopped externally

- **WHEN** a second Courier process stops the foreground delivery or its server by UUID
- **THEN** the initiating command observes confirmed removal, closes its watcher and lease, prints a neutral stopped result, and exits `0`

#### Scenario: Worker disappears unexpectedly

- **WHEN** the subscription closes without a confirmed stop or terminal snapshot
- **THEN** the initiating command reports a control failure and exits `40`

#### Scenario: Stop races lease release

- **WHEN** external stop and foreground cleanup occur concurrently
- **THEN** watcher and lease cleanup remains idempotent and does not turn a confirmed stop into a failure

# delivery-registry Specification

## Purpose
Define Courier's validated server and delivery state model, immutable private persistence, monotonic counters, and the trust boundary between durable discovery metadata and authoritative live IPC.

## Requirements

### Requirement: Stable delivery identities

Courier SHALL identify every server, delivery, tombstone, history event, and owned temporary resource with a validated UUID and SHALL reject duplicate or dangling references.

#### Scenario: Delivery references an unknown server

- **WHEN** a registry snapshot contains a delivery whose server UUID is absent
- **THEN** validation fails before the snapshot is persisted

### Requirement: Secret-free durable state

Courier SHALL persist only controlled discovery, policy, lifecycle, counter, and cleanup metadata and SHALL provide no durable registry field for credentials, session secrets, resource tokens, URL credentials, or free-form error text.

#### Scenario: Registry is inspected at rest

- **WHEN** an operator reads a valid registry revision
- **THEN** it contains no authentication or delivery secret material

### Requirement: Monotonic counters

Courier SHALL maintain race-safe non-negative read, sent, and confirmed byte counters and expose immutable snapshots.

#### Scenario: Concurrent progress updates

- **WHEN** multiple transfer routines add counter deltas concurrently
- **THEN** the final snapshot contains every accepted delta without regression

### Requirement: Private atomic registry revisions

Courier SHALL write registry state through private temporary files, flush each completed state file, and atomically commit immutable revision files under a rooted private application directory. Courier SHALL synchronize committed directory metadata where the operating system exposes supported directory-sync semantics and SHALL not fail a flushed Windows state commit solely because portable directory `Sync` is unsupported.

#### Scenario: Registry write is interrupted

- **WHEN** a new temporary revision is incomplete or fails to commit
- **THEN** readers continue to load the previous complete revision and the owned temporary file is removed

#### Scenario: Windows commits a registry revision

- **WHEN** a complete private state file is flushed and committed without replacement on Windows
- **THEN** the update succeeds without attempting unsupported directory-handle synchronization

### Requirement: Optimistic registry updates

Courier SHALL reject an update whose expected revision differs from the latest durable revision.

#### Scenario: Two writers use one revision

- **WHEN** the first writer commits and the second writer still expects the old revision
- **THEN** the second update fails without replacing or merging the committed state

### Requirement: Registry is not liveness authority

Courier SHALL treat stored process identifiers and endpoints as discovery metadata and SHALL require live authenticated IPC before declaring a server active.

#### Scenario: Stale process identifier

- **WHEN** a registry record contains a PID that has been reused by another process
- **THEN** later lifecycle logic does not treat that PID alone as proof of Courier ownership

### Requirement: Bounded private operational history

Courier SHALL store only validated, sanitized operational metadata in private immutable history entries and SHALL retain the newest 256 events by default.

#### Scenario: Retention bound is exceeded

- **WHEN** appending a valid event causes history to exceed its configured retention
- **THEN** Courier removes only its oldest private history files, preserves deterministic chronological order, and synchronizes the state directory

### Requirement: Durable upload activity index

Courier SHALL maintain a bounded global index of validated immutable upload receipts in its private state directory. Every receipt SHALL link a canonical receipt UUID to live-or-historical server and delivery UUIDs and SHALL contain only controlled route, destination, actor class/subject, peer IP, original/stored filename, byte count, digest, extraction count, and UTC timestamp metadata.

#### Scenario: Activity is inspected at rest

- **WHEN** an operator reads valid global or destination activity receipts
- **THEN** matching receipt and delivery identities explain who supplied which original file, where it was stored, its exact size and digest, and when it arrived without revealing credentials, sessions, resource tokens, or authorization headers

#### Scenario: Activity retention is exceeded

- **WHEN** a successful append exceeds the configured global activity retention
- **THEN** Courier removes only the oldest validated receipt objects and returns recent receipts in deterministic newest-first order

#### Scenario: A receipt is malformed

- **WHEN** a receipt filename, identity, field, hash, timestamp, or stored object type is invalid
- **THEN** the global activity reader fails closed without projecting untrusted data into administration

# server-control Specification

## Purpose
Define secret-free, IPC-authoritative inventory and UUID-scoped stop behavior for Courier data servers and deliveries without treating registry process IDs as control authority.

## Requirements

### Requirement: Authoritative server inventory

Courier SHALL list registered data servers and their nested deliveries deterministically and SHALL describe a server as live only after its UUID and compatibility are verified through private IPC and its live snapshot is valid.

#### Scenario: Registry entry is unreachable

- **WHEN** a registered control endpoint cannot complete its authoritative IPC probe
- **THEN** Courier labels the entry unreachable and does not infer liveness from its recorded PID

### Requirement: Secret-free control output

Courier SHALL show UUIDs, bind/process metadata, route, source, destination, state, policy, and counters without returning credentials, opaque resource tokens, runtime definitions, or other secret material.

#### Scenario: Protected delivery is listed

- **WHEN** a Basic- or password-protected delivery appears under `courier servers`
- **THEN** output identifies the authentication mode but contains no username, password, session, Authorization value, or resource token

### Requirement: Scoped UUID stop

Courier SHALL resolve one canonical UUID as a server, delivery, or known tombstoned target and SHALL verify the owning live worker before issuing the corresponding IPC stop operation.

#### Scenario: One delivery is stopped

- **WHEN** `courier servers stop <delivery-uuid>` succeeds
- **THEN** only that delivery is stopped and sibling deliveries on the same worker remain active

#### Scenario: One server is stopped

- **WHEN** `courier servers stop <server-uuid>` succeeds
- **THEN** that server and all deliveries it owns are stopped

### Requirement: Idempotent and unknown targets

Courier SHALL treat a UUID with an existing tombstone as already stopped and SHALL reject a UUID absent from active records and tombstone targets.

#### Scenario: Completed target is stopped again

- **WHEN** the requested UUID is already a tombstone target
- **THEN** Courier returns success without contacting or signaling another process

### Requirement: Stop all data servers

Courier SHALL independently verify and request shutdown for every registered data server, continue across individual failures, and report any partial failure without affecting the administrative UI.

#### Scenario: One of several workers is unreachable

- **WHEN** `courier servers stop --all` can stop only a subset of registered data workers
- **THEN** Courier reports the stopped count and returns a control failure describing the incomplete operation

### Requirement: No PID authority

Courier SHALL NOT signal a process solely from registry PID data and SHALL NOT expose manual clean, public worker, or `servers list` commands.

#### Scenario: PID has been reused

- **WHEN** a stale registry entry's PID now belongs to another process
- **THEN** server control leaves that process untouched because IPC identity verification did not succeed

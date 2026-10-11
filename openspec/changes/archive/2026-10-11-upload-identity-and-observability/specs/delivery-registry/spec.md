## ADDED Requirements

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

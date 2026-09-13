# delivery-policy-core Specification

## Purpose
Define the shared authentication, peer admission, abuse response, transfer reservation, and aggregate rate-limit boundary for network deliveries.

## Requirements

### Requirement: Ephemeral credential verification

Courier SHALL convert supplied delivery credentials to salted Argon2id verifier material, SHALL bound concurrent verification work, and SHALL never retain plaintext credentials in public state, persistent state, URLs, arguments, or logs.

#### Scenario: Public policy is inspected

- **WHEN** a delivery with authentication is listed or its public policy is updated
- **THEN** the snapshot contains authentication mode and limits but no credential, verifier, session, or CSRF material

### Requirement: Isolated authenticated sessions

Courier SHALL scope opaque password sessions and CSRF tokens to one delivery and SHALL enforce idle and absolute expiration while keeping Basic authentication request-scoped.

#### Scenario: Session is replayed to another delivery

- **WHEN** a valid session token from one delivery is presented to another delivery
- **THEN** authentication fails without revealing protected metadata

### Requirement: Canonical peer admission

Courier SHALL derive the peer IP from the accepted connection, canonicalize it, and apply exact-IP and CIDR rules without trusting forwarding headers.

#### Scenario: Rejected peer supplies an allowed forwarded address

- **WHEN** the connection peer is outside the allowlist but a forwarding header names an allowed address
- **THEN** admission is rejected using the connection peer

### Requirement: Atomic authentication threshold

Courier SHALL account concurrent authentication failures per delivery-scoped peer fingerprint and SHALL perform exactly one configured ban or stop transition when the attempt threshold is reached.

#### Scenario: Concurrent failures reach the threshold

- **WHEN** multiple failed authentications for one fingerprint reach the configured threshold concurrently
- **THEN** later requests are rejected by the ban or the delivery receives a stop decision without a lost update

### Requirement: Atomic transfer reservations

Courier SHALL reserve delivery slots and declared or incrementally observed incoming bytes before accepting them and SHALL release reservations exactly once on every terminal path.

#### Scenario: Concurrent uploads exceed a shared limit

- **WHEN** individually acceptable uploads would collectively exceed the delivery limit
- **THEN** the reservation that would exceed the limit is rejected before its body is consumed

### Requirement: Incoming file size limits

Courier SHALL enforce the configured maximum size independently for every incoming file, including zero-byte and unknown-length files.

#### Scenario: Unknown-length file crosses the limit

- **WHEN** incremental accounting would move a file beyond its configured maximum
- **THEN** the next bytes are rejected before they are committed

### Requirement: Aggregate directional throttling

Courier SHALL apply cancellation-aware upload and download rates through shared per-delivery direction limiters, and unlimited rates SHALL bypass throttling.

#### Scenario: Two transfers share one upload rate

- **WHEN** two upload-oriented legs run concurrently for one delivery
- **THEN** their combined admitted bytes consume the same upload limiter

### Requirement: Deterministic enforcement order

Courier SHALL enforce request bounds, peer admission, bans, authentication, session CSRF, reservations, throttling, and handler execution in that order.

#### Scenario: Unauthorized oversized upload arrives

- **WHEN** a request fails authentication and also declares a file larger than policy permits
- **THEN** authentication rejects it before a transfer reservation or body read occurs

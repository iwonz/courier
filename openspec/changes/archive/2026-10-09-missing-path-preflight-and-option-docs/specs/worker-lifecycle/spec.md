## ADDED Requirements

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


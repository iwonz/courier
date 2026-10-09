## ADDED Requirements

### Requirement: Incoming webhook QR readiness

Interactive incoming webhook readiness SHALL include a QR whose payload is the exact printed POST endpoint when the terminal can render it. The endpoint SHALL remain multipart-only and SHALL NOT gain an HTML form.

#### Scenario: Incoming webhook is ready

- **WHEN** `from webhook://` reaches readiness in a capable terminal
- **THEN** Courier shows the POST URL and exact-payload QR without changing request semantics

### Requirement: Compiled webhook acceptance

Courier SHALL exercise incoming and outgoing webhook profiles through the compiled binary and real HTTP multipart requests.

#### Scenario: Incoming multipart contains multiple files

- **WHEN** compiled-runtime acceptance submits more than one file part
- **THEN** the delivery rejects the request and commits no destination entry

#### Scenario: Directory is sent to an HTTP receiver

- **WHEN** compiled-runtime acceptance invokes `from <directory> to http://... --archive`
- **THEN** the receiver observes exactly one multipart file field named `file` containing the verified archive

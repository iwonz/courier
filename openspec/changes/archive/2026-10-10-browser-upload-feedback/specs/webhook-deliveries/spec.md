## MODIFIED Requirements

### Requirement: Compiled webhook acceptance

Courier SHALL exercise incoming and outgoing webhook profiles through the compiled binary and real HTTP multipart requests, asserting both committed destination bytes and explicit guarded HTTP outcomes.

#### Scenario: Incoming multipart is committed

- **WHEN** compiled-runtime acceptance sends one valid file part to an incoming webhook
- **THEN** the endpoint returns committed success and the destination contains the exact submitted bytes

#### Scenario: Incoming multipart name conflicts

- **WHEN** compiled-runtime acceptance sends a valid file part whose final name already exists
- **THEN** the endpoint returns an explicit guarded collision response and preserves the original destination bytes

#### Scenario: Incoming multipart contains multiple files

- **WHEN** compiled-runtime acceptance submits more than one file part
- **THEN** the delivery rejects the request and commits no destination entry

#### Scenario: Directory is sent to an HTTP receiver

- **WHEN** compiled-runtime acceptance invokes `from <directory> to http://... --archive`
- **THEN** the receiver observes exactly one multipart file field named `file` containing the verified archive

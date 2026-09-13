# webhook-deliveries Specification

## Purpose
Define Courier's bounded incoming and outgoing multipart webhook profile,
including authentication, transactional commit, redirect/retry policy, and
honest classification of uncertain remote outcomes.

## Requirements

### Requirement: Incoming webhook profile

Courier SHALL expose an opaque delivery URL that accepts exactly one `multipart/form-data` file part named `file` and SHALL not expose a browser page or password-session surface for that delivery.

#### Scenario: Multipart shape is unsupported

- **WHEN** an incoming request has a non-file first part, a field other than `file`, or more than one part
- **THEN** Courier rejects it without committing any destination object

### Requirement: Incoming commit acknowledgement

Courier SHALL acknowledge an incoming webhook only after its safe basename has been staged and committed to the final local or SSH destination, or after selected archive entries have been safely committed when `--extract` is active.

#### Scenario: Final name collides

- **WHEN** the destination already contains the uploaded basename
- **THEN** Courier returns a conflict, preserves the existing object and unrelated entries, and removes staging data

### Requirement: Shared incoming policy

Courier SHALL apply peer IP admission, optional Basic authentication, attempt actions, transfer reservation, selection, maximum file size, extraction limits, and aggregate upload rate before or while consuming incoming webhook data.

#### Scenario: Request is not admitted

- **WHEN** authentication, IP, size, or delivery-capacity policy rejects an incoming webhook
- **THEN** Courier does not consume its multipart payload or reveal destination metadata

### Requirement: Outgoing webhook profile

Courier SHALL issue exactly one HTTP POST containing exactly one multipart file part named `file`; a directory source SHALL require `--archive`, and no redirect or automatic retry SHALL occur.

#### Scenario: Receiver redirects

- **WHEN** the configured endpoint returns a redirect
- **THEN** Courier reports the received non-success response and does not send the file to the redirect target

### Requirement: Outgoing source preparation

Courier SHALL send a selected regular source file directly or create and verify one `<source-name>.tar.gz` when `--archive` is active, using bounded streaming and the aggregate upload-rate limit.

#### Scenario: Directory lacks archive mode

- **WHEN** an outgoing webhook source is a directory without `--archive`
- **THEN** Courier rejects it during preflight without opening an HTTP request

### Requirement: Outgoing authentication secrecy

Courier SHALL support no authentication or interactively acquired HTTP Basic credentials and SHALL not place those credentials in URLs, arguments, registry data, output, diagnostics, or logs.

#### Scenario: Basic receiver authenticates the request

- **WHEN** an outgoing webhook uses `--auth basic`
- **THEN** Courier sends the prompted username and password only in the request Authorization header

### Requirement: HTTP result certainty

Courier SHALL treat a received 2xx response as HTTP acceptance, a received non-2xx response as a known rejection, and a missing response after payload transmission as an unknown outcome that is not retried.

#### Scenario: Response is lost after payload transmission

- **WHEN** the transport consumes payload bytes but returns no HTTP response
- **THEN** Courier reports the sent-byte count and unknown outcome without issuing another request

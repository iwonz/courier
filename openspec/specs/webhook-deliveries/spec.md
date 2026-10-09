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

### Requirement: Webhook path preflight

Courier SHALL treat incoming webhook destinations and extraction roots as directories and SHALL treat an unarchived outgoing webhook source as a file-only role. It SHALL complete local or SSH path preflight before worker acquisition or an outgoing HTTP request.

#### Scenario: Incoming destination is absent

- **WHEN** an incoming webhook targets a missing destination directory
- **THEN** Courier confirms or force-creates it before launching or registering the worker

#### Scenario: Outgoing file-only source is absent

- **WHEN** an unarchived outgoing webhook source does not exist
- **THEN** Courier fails at preflight without creating it or opening an HTTP request

### Requirement: Discoverable detached incoming webhook lifecycle

Courier SHALL document and expose `--background` as the sole detached mode for incoming webhooks. It SHALL return after printing the URL and delivery UUID, survive terminal closure, and end through delivery or server stop, configured stop behavior, fatal worker failure, or process termination.

#### Scenario: Background incoming webhook is ready

- **WHEN** an incoming webhook starts with `--background`
- **THEN** the initiating command exits after readiness and the delivery remains listed by `courier servers`

### Requirement: Interactive incoming webhook lifecycle

Interactive incoming webhook readiness SHALL show route, URL, full delivery UUID, source, destination, foreground or background mode, and the applicable stop instruction. Foreground webhook service SHALL maintain an elapsed listening status until cancellation or fatal failure.

#### Scenario: Incoming webhook runs in foreground

- **WHEN** an incoming webhook becomes ready without `--background`
- **THEN** Courier shows its complete readiness identity and a live elapsed listening state until stopped

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

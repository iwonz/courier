## MODIFIED Requirements

### Requirement: Observable foreground termination

Courier SHALL retain each progress subscriber's delivery filter, allow a healthy subscription to remain silent without an implicit request timeout, publish a final filtered snapshot or tombstone before closing subscribers, and classify foreground termination from that evidence.

#### Scenario: Foreground delivery remains idle

- **WHEN** a foreground browser or webhook delivery emits no progress for longer than the ordinary IPC request timeout
- **THEN** its progress subscription, worker, lease, and initiating command remain active

#### Scenario: Caller bounds a subscription read

- **WHEN** a progress-subscription caller supplies a deadline or cancels its context
- **THEN** the blocked read terminates promptly with the caller's context result

#### Scenario: Delivery is stopped externally

- **WHEN** a second Courier process stops the foreground delivery or its server by UUID
- **THEN** the initiating command observes confirmed removal, closes its watcher and lease, prints a neutral stopped result, and exits `0`

#### Scenario: Worker disappears unexpectedly

- **WHEN** the subscription closes without a confirmed stop or terminal snapshot
- **THEN** the initiating command reports a control failure and exits `40`

#### Scenario: Stop races lease release

- **WHEN** external stop and foreground cleanup occur concurrently
- **THEN** watcher and lease cleanup remains idempotent and does not turn a confirmed stop into a failure

## ADDED Requirements

### Requirement: Observable foreground termination

Courier SHALL retain each progress subscriber's delivery filter, publish a final filtered snapshot or tombstone before closing subscribers, and classify foreground termination from that evidence.

#### Scenario: Delivery is stopped externally

- **WHEN** a second Courier process stops the foreground delivery or its server by UUID
- **THEN** the initiating command observes confirmed removal, closes its watcher and lease, prints a neutral stopped result, and exits `0`

#### Scenario: Worker disappears unexpectedly

- **WHEN** the subscription closes without a confirmed stop or terminal snapshot
- **THEN** the initiating command reports a control failure and exits `40`

#### Scenario: Stop races lease release

- **WHEN** external stop and foreground cleanup occur concurrently
- **THEN** watcher and lease cleanup remains idempotent and does not turn a confirmed stop into a failure

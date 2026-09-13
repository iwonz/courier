## ADDED Requirements

### Requirement: Bounded private operational history

Courier SHALL store only validated, sanitized operational metadata in private immutable history entries and SHALL retain the newest 256 events by default.

#### Scenario: Retention bound is exceeded

- **WHEN** appending a valid event causes history to exceed its configured retention
- **THEN** Courier removes only its oldest private history files, preserves deterministic chronological order, and synchronizes the state directory

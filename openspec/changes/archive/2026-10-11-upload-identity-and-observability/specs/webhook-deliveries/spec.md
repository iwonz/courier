## MODIFIED Requirements

### Requirement: Incoming commit acknowledgement

Courier SHALL acknowledge an incoming webhook only after its safe upload has been staged, committed under the configured preserve-or-unique filename policy, and recorded in matching destination and global activity receipts, or after selected archive entries and their receipt have been safely committed when `--extract` is active.

#### Scenario: Final name collides

- **WHEN** the destination already contains the uploaded basename and filename policy is `preserve`
- **THEN** Courier returns a conflict, preserves the existing object and unrelated entries, and removes staging data

#### Scenario: Repeated original names use unique policy

- **WHEN** multiple valid webhook requests submit the same original basename under filename policy `unique`
- **THEN** every request receives a different committed UUID-based name and every receipt retains the shared original name

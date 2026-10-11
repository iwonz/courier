## MODIFIED Requirements

### Requirement: Transactional browser upload

Courier SHALL reserve policy capacity before reading an upload, stage accepted bytes privately, commit only to an absent safe destination name, persist matching destination and global activity receipts before acknowledgement, and preserve unrelated entries. The optional `--filename-policy` SHALL default to `preserve`; `unique` SHALL commit a UUID-based filename while retaining the original filename in both receipts.

#### Scenario: Uploaded name conflicts

- **WHEN** the sanitized submitted upload name already exists and filename policy is `preserve`
- **THEN** Courier rejects the upload without consuming or replacing that entry and removes its staging data

#### Scenario: Uploaded name uses unique policy

- **WHEN** an accepted upload uses filename policy `unique`
- **THEN** Courier commits it under a fresh UUID-based safe name and records both the submitted and committed names

#### Scenario: Uploaded archive is extracted

- **WHEN** a browser upload delivery uses `--extract` and receives a safe tar.gz
- **THEN** Courier applies the shared extraction limits and selection policy, commits non-conflicting top-level entries into the destination root without retaining the uploaded archive, and records the archive receipt and entry count

#### Scenario: Receipt persistence fails

- **WHEN** Courier cannot durably record the destination or global receipt
- **THEN** it returns a guarded server error and does not acknowledge the upload as delivered

## ADDED Requirements

### Requirement: Private upload activity namespace

Courier SHALL write secret-free immutable upload receipts under the destination `.courier/activity/` namespace for local and SSH destinations and SHALL exclude that namespace from browser listings, downloads, archives, and Courier directory transfers.

#### Scenario: A received directory is later published

- **WHEN** a destination containing Courier activity receipts is shared through a browser delivery or Courier directory transfer
- **THEN** no receipt file, actor identity, peer address, or internal metadata path is exposed in the manifest or payload

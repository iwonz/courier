## ADDED Requirements

### Requirement: Compiled upload identity and activity acceptance

Courier SHALL exercise preserved and unique incoming browser and webhook uploads through the compiled binary, real HTTP multipart requests, local destination receipts, the private global activity index, and the administration UI.

#### Scenario: Unique browser uploads are observed

- **WHEN** compiled browser acceptance uploads repeated Unicode and binary filenames with `--filename-policy unique`
- **THEN** byte-identical files receive distinct UUID-based committed names, matching destination/global receipts retain each original name, and administration renders their actors, exact quantities, and timestamps

#### Scenario: Preserved webhook upload is observed

- **WHEN** compiled webhook acceptance uploads a safe file under the default policy
- **THEN** its basename remains unchanged and matching receipt identities survive worker stop and administration restart

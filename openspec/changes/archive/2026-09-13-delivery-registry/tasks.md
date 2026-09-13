## 1. Domain model

- [x] 1.1 Add UUID server, delivery, tombstone, history, and owned-temp identities
- [x] 1.2 Add validated policies, lifecycle transitions, references, and optimistic versions
- [x] 1.3 Add race-safe monotonic counters and immutable snapshots

## 2. Private persistence

- [x] 2.1 Create a private rooted state directory and reject unsafe state objects
- [x] 2.2 Commit immutable registry revisions atomically with optimistic revision checks
- [x] 2.3 Append and load private immutable history events with deterministic ordering

## 3. Quality

- [x] 3.1 Cover validation, corruption, interruption, races, permissions, and cleanup at 100%
- [x] 3.2 Document registry schema, secret exclusions, liveness limits, and ownership metadata
- [x] 3.3 Run the full verification gate, strict validation, and archive the change

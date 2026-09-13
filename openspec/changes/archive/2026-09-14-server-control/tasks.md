## 1. Control domain

- [x] 1.1 Persist secret-free source, destination, and compatibility metadata for new workers and deliveries
- [x] 1.2 Add deterministic authoritative server/delivery inventory with unreachable-state handling
- [x] 1.3 Add scoped UUID resolution, idempotent tombstone handling, server stop, delivery stop, and stop-all aggregation

## 2. CLI and contract

- [x] 2.1 Register independent `servers` and `servers stop` Cobra handlers without exposing aliases or cleanup commands
- [x] 2.2 Render stable nested output without credentials, tokens, runtime definitions, or unsafe diagnostics
- [x] 2.3 Mark only server control commands and `--all` shipped in the canonical contract and generated reference

## 3. Quality

- [x] 3.1 Cover live/stale inventory, mismatched snapshots, UUID kinds, repeated/unknown stops, partial stop-all failures, races, and output at 100%
- [x] 3.2 Document local control semantics and the no-PID-signal trust boundary
- [x] 3.3 Run full verification, strict OpenSpec validation, archive the change, and commit once

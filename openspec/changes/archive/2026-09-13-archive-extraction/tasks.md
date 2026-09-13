## 1. Codec and validation

- [x] 1.1 Add duplicate-safe codec registration and tar.gz resolution
- [x] 1.2 Validate entry paths, relative links, duplicates, types, depth, count, size, and expansion ratio
- [x] 1.3 Keep created archives compatible with fixed extraction limits

## 2. Extraction transaction

- [x] 2.1 Inspect without writes and reopen the archive for bounded extraction
- [x] 2.2 Stage selected entries and commit into an absent or existing extraction root without overwrite
- [x] 2.3 Preserve unrelated entries and clean owned stages on success, failure, and cancellation

## 3. CLI and quality

- [x] 3.1 Ship extract and max-extracted-size through planner, command, generated reference, and reporting
- [x] 3.2 Cover local/remote-capable backends, collisions, selection, limits, malformed archives, and failure cleanup
- [x] 3.3 Run strict OpenSpec validation and the full verification gate

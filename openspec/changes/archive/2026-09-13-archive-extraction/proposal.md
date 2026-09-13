# Change: Add safe archive extraction

## Why

Courier can create and verify tar.gz files but cannot extract them. Extraction needs stronger guarantees than streaming entries directly into a destination: the entire archive and collision set must be validated before the first commit, and archive bombs and unsafe links must be rejected independently of user size settings.

## What Changes

- Add an extensible archive codec registry with tar.gz as the only shipped codec.
- Ship `--extract` and `--max-extracted-size` for path-to-path operations.
- Inspect archive structure, paths, links, types, duplicates, limits, and destination conflicts before staging output.
- Extract through bounded buffers into a private stage and commit absent top-level entries into the destination root.
- Preserve unrelated root entries and apply the shared selection policy.
- Enforce 100,000 entries, depth 64, and 100:1 expansion ratio even when the configurable size cap is unlimited.

## Impact

Extraction never auto-detects from extension without `--extract`. Unsupported codecs and every collision fail in preflight; tar.gz remains the only accepted format.

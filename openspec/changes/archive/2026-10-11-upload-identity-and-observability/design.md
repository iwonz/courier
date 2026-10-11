# Design: Portable upload receipts and quantity presentation

## Context

Incoming uploads are committed through the shared `fsx.Backend`, which supports both rooted local filesystems and SFTP. SQLite is suitable for Courier's local private state, but a SQLite database placed directly on an SFTP-backed destination cannot provide reliable locking, journaling, or atomic page updates. Courier already has a cross-platform immutable, private, no-replace persistence protocol.

## Decisions

### Filename policy

The public flag is `--filename-policy <preserve|unique>` and applies only to `web-to-path` and `webhook-to-path`. The default `preserve` behavior remains compatible and continues to reject collisions. `unique` creates a canonical UUID basename plus a bounded safe final extension from the submitted name. The original safe filename is never discarded. Because extraction commits archive entries rather than the uploaded archive name, `unique` conflicts with `--extract`.

The selected policy is carried in the existing bounded runtime definition as an optional string. An omitted field decodes as `preserve`, preserving upgrade compatibility without a runtime-definition version bump.

### Receipt model and identity

Each acknowledged upload has one immutable receipt with:

- schema version and canonical receipt, server, and delivery UUIDs;
- route and destination display value;
- original name, committed name, extraction flag/entry count, exact byte count, and SHA-256 digest;
- authentication class (`none`, `basic`, or `password`), a bounded subject (`anonymous`, verified Basic username, or `password-authenticated`), and the canonical socket peer IP;
- UTC received time.

Passwords, session/CSRF values, authorization headers, resource tokens, URL secrets, and request headers are never recorded. Receipt validation rejects controls, traversal, malformed hashes, invalid identities, and inconsistent extracted/stored-name fields.

### Destination and global journals

The destination journal is a private `.courier/activity/` directory containing one mode-`0600` JSON object per receipt. Courier verifies every metadata directory as a real non-symlink directory, writes and flushes a private staged object, commits it without replacement, and synchronizes through the backend where supported. The metadata namespace is excluded from Courier browser navigation, directory archives, and ordinary directory transfers so audit identity data is not accidentally republished.

The worker also appends the same receipt to the private global delivery store. Filenames include timestamp and receipt UUID, reads are newest-first, and retention is bounded. Matching receipt/delivery/server IDs form the link between global and destination records. This reuses Courier's portable immutable store instead of adding a second durability model or a CGO/database dependency.

Receipt persistence is part of upload acknowledgement. A response is successful only after payload commit plus destination and global receipt commits. Post-payload receipt failures return a guarded server error and never claim an unrecorded upload as acknowledged.

### Administration projection

The admin process reads a bounded recent-receipt projection directly from the global store and includes it in the guarded snapshot/SSE API. The React UI shows receipt totals and a responsive activity table for the selected delivery, with time, actor and peer, original-to-stored mapping, size, and extraction state. It never opens destination journals or accepts arbitrary paths.

### Quantities

Shared formatters distinguish bytes (`B`, `KiB`, `MiB`, `GiB`, `TiB`, `PiB`, `EiB`) from bits (`bit`, `Kibit`, and higher) and rates (`/s`). Values below the scaling threshold use space-grouped integers. Scaled values use at most one fractional digit; detailed and accessible surfaces retain the grouped exact base-unit value. Invalid or negative domain values fail validation rather than rendering misleading quantities.

Interactive terminal and browser surfaces use human-readable quantities. Redirected CLI output and JSON APIs retain raw integer fields for compatibility and machine consumption.

## Testing

- Unit-test option applicability/default/conflict, UUID name generation, safe extension handling, collision behavior, and optional runtime-definition compatibility.
- Exercise local and mocked SFTP receipt commits, permissions, symlink rejection, retention, malformed records, secret absence, and failure after payload commit.
- Verify admin API/UI activity projection, delivery filtering, localization, narrow layouts, accessibility, and exact byte tooltips.
- Cover byte/bit/rate formatting boundaries from zero through EiB in Go and TypeScript.
- Extend compiled browser/webhook E2E with preserve and unique names, byte equality, destination/global receipt parity, restart persistence, and admin visibility.
- Run focused tests, strict OpenSpec validation, and the complete `make verify` release-candidate gate.

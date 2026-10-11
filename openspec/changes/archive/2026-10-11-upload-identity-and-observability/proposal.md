# Change: Upload identity, accounting, and human-readable quantities

## Why

Incoming browser and webhook uploads currently use the submitted basename as the committed destination name. Operators cannot choose collision-resistant names while retaining the sender's original filename, and Courier keeps only aggregate live counters rather than a durable per-file audit trail. Byte quantities are also rendered inconsistently as ungrouped raw integers in terminal and web surfaces.

## What Changes

- Add `--filename-policy <preserve|unique>` to incoming browser and webhook routes. `preserve` remains the default; `unique` commits a UUID-based name while retaining the safe extension and recording the original name.
- Persist validated, secret-free upload receipts both beside the destination under `.courier/activity/` and in Courier's private global state directory. Use immutable receipt objects rather than SQLite files so local and SSH destinations have identical atomic no-replace behavior without relying on unsafe database locking over SFTP.
- Link every receipt by receipt, delivery, and server identities and record route, destination, authenticated actor class/subject, socket peer IP, original/stored names, payload bytes, digest, extraction outcome, and timestamp.
- Expose recent receipts through the loopback-only administration API and render a localized upload activity view in the administration application.
- Introduce shared Go and TypeScript quantity formatting with grouped exact values, binary byte/bit units, explicit `B` versus `bit`, and rate units; use it across interactive terminal reports, progress, server views, updater progress, delivery UI, and administration UI while retaining stable raw redirected output.
- Extend real runtime/browser acceptance for preserved and unique uploads, durable receipts, admin visibility, and formatted quantities.

## Impact

- Affected specs: `web-deliveries`, `webhook-deliveries`, `delivery-registry`, `admin-ui`, `operation-reporting`, `cli-contract`, `cross-platform-acceptance`
- Affected code: operation planning, hosted runtime definitions, upload commit path, delivery state store, worker/admin wiring, terminal reporting, shared web UI, data UI, admin UI, generated CLI documentation, and E2E tests
- CLI contract version advances to `0.17.0`; schema remains `3`
- Existing delivery registry, IPC, worker-definition, and administration-state versions remain unchanged; the new runtime-definition field is optional and defaults to `preserve`

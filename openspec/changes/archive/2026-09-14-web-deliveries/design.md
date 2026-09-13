# Context

A single worker may host several deliveries on one bind. The listener reserved by the worker lifecycle must therefore dispatch by an unguessable resource token and cannot use filesystem paths, credentials, or delivery UUIDs as public URLs. Both local and SSH-backed endpoints must enter through the same backend and transaction abstractions.

# Decisions

## Worker-owned data host

Worker registration carries a bounded private runtime definition in the existing versioned IPC frame. The durable delivery record remains secret-free. A host adapter validates and opens the endpoint before registration commits, and registration rollback closes any opened resources. Delivery and server stop commit durable tombstones before closing host resources. Owned-temporary metadata is removed only after physical cleanup succeeds, so a cleanup failure remains recoverable instead of leaving an apparently active registry entry with an already closed resource.

SSH-backed deliveries perform an initiating-process preflight using the normal SSH configuration, host-key verification, agent, key, password, and helper-consent behavior. Only prompt responses actually needed by that preflight and the fact of successful helper approval cross private IPC. The worker clears them after opening the endpoint and never persists them. A pre-approved helper may be redeployed without a second prompt.

The worker serves a chi router on its already reserved data listener. One host instance owns all delivery routes for the bind. A randomly generated 256-bit base64url resource token selects a delivery; unknown and stopped tokens return the same not-found response.

## Protected metadata boundary

Authentication and peer admission run before filesystem metadata, directory entries, archive names, or sizes are read. Browser password login establishes a host-only session and CSRF pair. Basic authentication remains request-scoped. JSON errors do not distinguish unknown credentials from protected resource state.

## Shared transfer primitives

Uploads reserve policy capacity before consuming multipart file bytes, write to private sibling staging paths, fsync, and commit only to absent final names. Conflicts are detected before final commit and unrelated destination entries are preserved. Selection rules are evaluated using sanitized relative names. With `--extract`, the staged upload is passed to the shared tar.gz codec registry and safe staged extraction transaction instead of being retained as a file.

Downloads use bounded streaming from the opened backend and the aggregate download limiter. Directories are navigable through a normalized relative path and may be streamed as a deterministic tar.gz with one top-level source entry. Archive output rejects traversal and symlink escape and does not create an unbounded intermediate archive. Explicit `--archive` prepares and verifies one private archive before public registration, then publishes only that object.

## UI artifact ownership

The data application is a separate Vite workspace that imports the shared Lit UI package. Its versioned API paths are relative to the opaque resource root. Production assets are generated deterministically into an embed package; a freshness check rebuilds to a temporary directory and compares content. `--no-ui` omits HTML and exposes only documented JSON/download endpoints.

# Risks / Trade-offs

Opaque resource tokens are access locators, not a replacement for configured authentication. They may appear in browser history and command output, so they contain no endpoint or credential data and are revoked at delivery stop. Streaming directory archives cannot know their compressed size in advance; safety is enforced while walking and reading rather than with a full temporary artifact. Explicit archive mode intentionally uses a private verified temporary artifact and removes it when the delivery stops.

# Migration Plan

Add host registration and rollback semantics, implement and test the web data plane, add the data UI and freshness gate, then expose only the two browser routes in the canonical contract. Run all verification and cross-builds, archive this change, and commit it once before adding webhook handlers.

# Context

Data workers already expose UUID-checked private IPC for snapshots, policy updates, progress, and stop actions. The private delivery registry provides discovery but is not liveness authority. The administration UI should compose those boundaries instead of creating a second transfer or policy engine.

# Decisions

## Separate singleton process

The administration server runs as a distinct internal process mode and never registers as a data server. A private per-user lock serializes startup. A `0600` atomic state file contains only its UUID, canonical loopback bind, private control endpoint, diagnostic PID, and timestamps. Startup verifies an existing instance by UUID-bound private IPC; stop does the same and never signals a PID.

Foreground start holds the process until cancellation. Background start writes a private one-use launch file, starts the current Courier executable with an allowlisted environment, waits for the authoritative ready handshake, and then returns the local URL. The child removes the launch file. Shutdown removes only state owned by the matching UUID and releases the singleton lock.

## Guarded loopback API

Only canonical loopback binds are accepted. Every request must use the listener's exact Host value. Mutating requests additionally require an exact same-origin `Origin`; absent or foreign origins fail. Courier emits no permissive CORS headers, rejects oversized or unknown JSON, sets no-store and defensive browser headers, and applies request timeouts and a bounded SSE client count.

The API exposes secret-free server/delivery views, expected-version policy updates, UUID-scoped stop actions, and SSE snapshots. It discovers data workers from the delivery registry and repeats authoritative IPC checks. Optimistic version conflicts return HTTP 409. Credentials, session material, hosted resource tokens, runtime definitions, and arbitrary worker errors never enter response bodies.

## Shared web application

`web/admin` is a Vite/Lit workspace that imports `@courier/ui` source aliases. Its typed English and Russian catalogs cover inventory, policy editing, errors, stop actions, themes, and locale selection. Deterministic Vite output is embedded under `internal/admin/assets`; freshness, build, and exact TypeScript coverage join the existing web verification gate.

# Risks / Trade-offs

The UI is intentionally local and unauthenticated at the application layer; private-machine access plus exact Host/Origin validation is its trust boundary. Users who need remote access must establish their own authenticated tunnel while preserving the expected local Host. SSE sends complete secret-free snapshots instead of a complex mutable event log, simplifying reconnection and consistency at the cost of small repeated payloads.

# Migration Plan

Add the process/control domain and API first, register isolated CLI providers, add and embed the shared-kit web consumer, update the canonical contract, run all Go/TypeScript/race/release checks, archive the OpenSpec change, and commit once.

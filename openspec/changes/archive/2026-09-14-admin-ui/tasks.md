## 1. Administration runtime

- [x] 1.1 Add private singleton lock, atomic discovery state, UUID-bound control IPC, and owned cleanup
- [x] 1.2 Add foreground/background process launch, ready verification, cancellation, and idempotent stop
- [x] 1.3 Enforce canonical loopback-only binds and independent data-worker lifecycle scope

## 2. API and CLI

- [x] 2.1 Add exact Host/Origin guards, bounded versioned JSON endpoints, and secret-free snapshots
- [x] 2.2 Add optimistic policy updates, delivery/server stop actions, and bounded SSE snapshots
- [x] 2.3 Ship independent `ui start` and `ui stop` providers and update contract/help parity

## 3. Web and quality

- [x] 3.1 Build the English/Russian Lit admin application from `@courier/ui` with themes, accessibility, responsive state, and policy conflict handling
- [x] 3.2 Embed deterministic admin assets and add exact Go/TypeScript coverage, race, security, lifecycle, and freshness tests
- [x] 3.3 Document administration trust boundaries, run full verification, archive this change, and commit once

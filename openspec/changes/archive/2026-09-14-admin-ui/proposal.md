# Change: Ship the local administration UI

## Why

Courier can host, list, and stop data deliveries, but operators need a local browser surface for live inspection and safe policy changes. The administration process must be separate from data workers so `servers stop --all` cannot terminate it, and browser access must not turn loopback control into a cross-origin or DNS-rebinding interface.

## What Changes

- Ship `courier ui start [--listen <loopback-host:port>] [--background]` and `courier ui stop`.
- Add one per-user administration singleton with a private lock, secret-free atomic discovery state, and UUID-verified private IPC stop control.
- Add a loopback-only, versioned `net/http` API for authoritative server/delivery snapshots, optimistic policy updates, scoped stop actions, and SSE snapshot progress.
- Reject untrusted Host and Origin values, disable CORS, bound requests and event clients, and never return credentials or private runtime definitions.
- Build an English/Russian Lit administration application from the shared Courier UI kit and embed deterministic production assets in the Go binary.

## Impact

The canonical contract advances only `ui start` and `ui stop` to shipped. Data-server commands retain independent lifecycle scope. No remote administration endpoint, browser-supplied credentials, public worker command, or external repository is introduced.

# Change: Ship browser upload and download deliveries

## Why

The command contract already models `web://`, but the routes are intentionally hidden until a real data plane exists. Courier needs authenticated, bounded browser transfers that reuse the policy, selection, archive, transactional storage, and worker lifecycle layers rather than introducing a second copy engine.

## What Changes

- Ship `web://` as a source and destination for `web-to-path` and `path-to-web` operations.
- Add a chi-based multi-delivery data host registered through private worker IPC.
- Use opaque per-delivery resource tokens and never reveal protected metadata before authorization.
- Add transactional multipart uploads, filtered directory navigation, file download, and whole-directory tar.gz download.
- Add a Lit data application built from the shared Courier UI kit and an explicit JSON-only `--no-ui` surface.
- Wire browser flags, policy defaults, interactive credential acquisition, foreground leases, and background worker ownership into `courier from`.
- Embed deterministic production data-UI assets in Courier binaries and verify their freshness.

## Impact

The public contract advances only the two browser routes and their applicable flags to `shipped`. Webhook, server-control, and administration commands remain planned. Runtime configuration and credentials cross private IPC only and are not persisted in the delivery registry.

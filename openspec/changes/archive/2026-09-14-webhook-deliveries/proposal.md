# Change: Ship incoming and outgoing webhook deliveries

## Why

Courier's endpoint and option matrix already reserves `webhook://` for programmable incoming uploads and HTTP(S) destinations for one-shot outgoing delivery. Those routes must reuse the browser delivery host, policy engine, endpoint backends, selection rules, archive implementation, and bounded streaming primitives without widening the public command surface or inventing provider-specific webhook profiles.

## What Changes

- Ship `webhook://` to local/SSH directory delivery through an opaque worker-hosted multipart endpoint.
- Accept exactly one `multipart/form-data` file field named `file`, authenticate and admit the peer before consuming it, and acknowledge only after the final destination commit.
- Ship local/SSH path to HTTP(S) as one immediate multipart POST with no redirect following and no automatic retry.
- Require `--archive` for an outgoing directory and reuse the verified tar.gz preparation path.
- Apply ordered selection, Basic authentication, aggregate limits, size/rate policy, bounded buffers, cancellation, and secret clearing through existing shared components.
- Report a response-less POST after request transmission as an unknown outcome rather than success.

## Impact

The contract advances only `webhook-to-path`, `path-to-http`, and their endpoint kinds to `shipped`. No JSON/raw body, custom field/header/signature, retry, redirect, or third-party provider profile is added.

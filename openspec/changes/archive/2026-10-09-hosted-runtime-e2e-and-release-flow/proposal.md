## Why

Courier's hosted routes are covered primarily through unit and mocked browser tests, so regressions can escape in the compiled binary, embedded applications, multipart boundaries, and foreground worker lifecycle. Readiness URLs also require manual copying, usage mistakes are reported as transfer failures, and the release process permits tracked changes to bypass the intended OpenSpec-to-publication lifecycle.

## What Changes

- Add terminal QR codes for validated administration, browser-delivery, and incoming-webhook readiness URLs without changing listener reachability or plain output.
- Make an externally stopped foreground delivery or administration UI exit cleanly, while preserving cancellation and unexpected-worker-loss classifications.
- Render all CLI usage failures compactly without operational byte accounting.
- Add serial real-runtime browser acceptance for incoming browser uploads, incoming webhooks, browser downloads, outgoing HTTP webhooks, and administration startup/shutdown.
- Codify and enforce a clean-main OpenSpec workflow with guarded change creation, feature-branch shipping, publication-gated Pages, and recovery-only low-level release commands.

## Capabilities

### New Capabilities

- `development-workflow`: Defines mandatory tracked-change inception, verification, merge, publication, and clean final state.

### Modified Capabilities

- `admin-ui`: Defines QR readiness and clean external foreground shutdown.
- `cross-platform-acceptance`: Defines compiled-binary hosted route and lifecycle acceptance.
- `operation-reporting`: Defines compact typed usage diagnostics and QR capability rules.
- `release-automation`: Defines feature-branch shipping, publication ordering, guards, and recovery boundaries.
- `web-deliveries`: Defines browser-route QR readiness and externally stopped foreground ownership.
- `webhook-deliveries`: Defines incoming-webhook QR readiness and compiled multipart coverage.
- `worker-lifecycle`: Defines filtered terminal progress publication and stop/loss distinction.

## Impact

This change affects terminal rendering, Cobra error classification, worker progress subscriptions, administration shutdown, browser acceptance infrastructure, repository hooks, Make targets, release scripts, workflows, and contributor documentation. It adds `rsc.io/qr v0.2.0` but does not change network exposure, route syntax, CLI contract schema/version, IPC, registry, worker-definition, or administration-state versions.

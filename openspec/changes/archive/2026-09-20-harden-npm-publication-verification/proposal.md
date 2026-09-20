## Why

npm can accept and sign a package while continuing asynchronous registry processing for several minutes. Courier's one-minute visibility check can therefore fail an otherwise successful release and prevent the automatic Formula synchronization and Pages handoff.

## What Changes

- Extend post-publication npm visibility verification to tolerate the registry's documented processing delay.
- Force every visibility probe to revalidate against the online registry instead of accepting a stale local cache result.
- Document the bounded retry behavior and preserve a hard failure when the package never becomes public.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `release-automation`: Make post-publication npm verification resilient to bounded registry processing delays.

## Impact

The release workflow and release runbook change. Package contents, credentials, publication order, tags, application behavior, and package-manager commands remain unchanged.

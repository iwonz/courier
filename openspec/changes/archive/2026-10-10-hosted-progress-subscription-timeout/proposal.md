## Why

Foreground browser and webhook deliveries close after five idle seconds because every progress-subscription read inherits the short default IPC request timeout. The command then reports an unconfirmed worker loss even though the worker is healthy, releases its lease, and shuts the delivery down.

## What Changes

- Separate long-lived IPC stream configuration from bounded request/response configuration.
- Keep foreground progress subscriptions open while idle until an event, caller deadline, cancellation, peer close, or protocol failure occurs.
- Preserve the existing five-second default timeout for ordinary IPC calls and lease release.
- Add unit and compiled-runtime regression coverage proving a foreground hosted delivery survives more than five idle seconds and still stops cleanly through control.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `worker-lifecycle`: Defines idle-safe long-lived progress subscriptions while retaining bounded ordinary IPC calls.
- `cross-platform-acceptance`: Requires real hosted foreground acceptance to cross the former five-second timeout boundary.

## Impact

This change affects IPC connection deadline configuration, worker progress-subscription reads, focused lifecycle tests, and the serial compiled-binary browser suite. It does not change IPC messages, protocol versions, persisted state, route syntax, output fields, or background behavior.

## Why

Courier names Unix control sockets `control-<UUID>.sock` inside its private state directory. The default macOS state path for the reporting user makes that endpoint exactly 104 bytes, but Darwin's Unix-socket path buffer requires room for termination, so administration and web-delivery workers fail with `bind: invalid argument`.

## What Changes

- Encode new Unix control-socket filenames as the canonical UUID's 32 lowercase hexadecimal characters without separators or decorative affixes.
- Keep sockets inside Courier's private state directory and preserve listener permissions, stale cleanup, and stored legacy endpoint behavior.
- Cover the reported macOS default path and Darwin's 103-byte usable boundary with regression tests.
- Publish the correction as patch release v0.3.6 through the established ship flow.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `worker-lifecycle`: require newly generated private Unix control endpoints to fit the supported platform limit under Courier's default state-directory layout.

## Impact

New Unix endpoint names change without changing the registry, administration-state, or IPC schemas. Existing persisted endpoint strings remain dialable and removable as recorded. Windows named pipes, public commands, and CLI options remain unchanged.

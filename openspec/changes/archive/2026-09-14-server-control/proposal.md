# Change: Ship local server and delivery control

## Why

Background and shared-bind deliveries now outlive their initiating commands, but Courier has no public way to inspect or stop them. Control must resolve UUIDs through the private registry and verify the owning worker over IPC; a recorded PID is descriptive data and must never be trusted as authority for signaling a process.

## What Changes

- Ship `courier servers`, `courier servers stop <uuid>`, and `courier servers stop --all`.
- Persist secret-free source, destination, and worker compatibility metadata for new registrations.
- Build deterministic nested server/delivery views from live IPC snapshots, clearly labeling unreachable registry entries.
- Resolve server and delivery UUIDs against registry ownership and authoritative worker identity before stop actions.
- Make a repeated stop of a tombstoned UUID an idempotent success and reject unknown UUIDs.
- Stop all data workers while remaining independent from the future administrative UI lifecycle.

## Impact

The canonical contract advances only `servers`, `servers stop`, and `--all` to `shipped`. No `servers list`, `servers clean`, PID signaling, worker command, or manual registry mutation is exposed.

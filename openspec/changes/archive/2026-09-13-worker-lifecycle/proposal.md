# Change: Add shared worker lifecycle management

## Why

Browser and incoming-webhook deliveries need long-lived local workers that can share one bind safely. Starting one process per command without coordination would race on the listener, confuse stale PIDs with live Courier instances, leak foreground deliveries after their owner exits, and make later server controls unreliable.

## What Changes

- Add a versioned, bounded, strict request/response IPC protocol with stable operation and error identifiers.
- Add private Unix-socket and Windows named-pipe control transports whose operating-system permissions form the local authentication boundary.
- Add bind-scoped startup locks and a coordinator that reuses only live compatible workers.
- Add connection-owned foreground leases and explicit background ownership.
- Add worker registration, listing, policy update, stop, shutdown, and progress-subscription protocol operations.
- Reconcile stale registry records through live handshakes, tombstones, and cleanup callbacks for validated Courier-owned temporary resources.
- Shut down a worker after its last delivery ends unless an explicit keepalive owner remains.

## Impact

This change adds internal lifecycle and IPC packages only. It does not expose worker modes in Cobra, start an HTTP data server, implement authentication policy, or ship the planned `servers` and `ui` commands.

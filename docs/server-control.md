# Server control

Courier exposes local control commands for browser and incoming-webhook deliveries hosted by background or shared-bind data workers:

```sh
courier servers
courier servers stop <server-or-delivery-uuid>
courier servers stop --all
```

`courier servers` prints data servers in UUID order and nests each server's deliveries in UUID order. The output contains the bind address, diagnostic process ID, lifecycle state, timestamps, route, source and destination display values, controlled policy values, and byte counters. It never reads or displays credentials, sessions, Authorization headers, opaque resource tokens, or private runtime definitions. Older records without endpoint display metadata use `<unavailable>`.

## Authority and stale entries

The private per-user registry is discovery metadata, not proof that a process is still Courier or still owns a listener. Courier labels a server `live` only after the recorded private control endpoint completes the versioned UUID and compatibility handshake and returns a matching valid snapshot. A record that cannot pass this check remains visible as `unreachable` for diagnosis and later owned-resource reconciliation.

The process ID is informational only. Courier never signals it. This prevents an old registry record from terminating an unrelated process after operating-system PID reuse.

## Scoped stopping

A canonical UUID selects exactly one delivery, one data server, or a known tombstoned target. Courier repeats the authoritative IPC check before sending a stop request:

- stopping a delivery leaves its sibling deliveries running;
- stopping a server stops that data server and all deliveries it owns;
- stopping an already tombstoned target succeeds without contacting another process;
- an unknown UUID fails without changing registry or process state.

`courier servers stop --all` independently verifies every registered data server and continues after individual failures. It always reports the number successfully stopped and returns exit code `40` if any server could not be verified or stopped. It does not affect the separate administrative UI process.

There are intentionally no public `servers list`, `servers clean`, worker, PID-based, or direct registry-mutation commands.

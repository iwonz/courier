# Change: Bound the Windows pipe acceptance lifecycle

## Why

The same Windows named-pipe acceptance passed once and then blocked a later `main` run until Go's ten-minute package timeout. The test reused a fixed pipe name, dialed with an unbounded context, and registered listener cleanup only after the connection completed, so a transient native connection stall produced a leaked test lifecycle and delayed diagnostics.

## What Changes

- Give each native Windows pipe acceptance run a unique valid server identity and pipe endpoint.
- Bound the client dial with an explicit timeout.
- Register listener cancellation and cleanup immediately after creation so every failure path releases the native pipe and accept goroutine.
- Keep the real owner-scoped named-pipe listener and client in the acceptance path.

## Impact

Windows acceptance still exercises the production transport and ACL boundary, but it cannot wait for the package-wide timeout or collide with another test endpoint. Production IPC behavior and public interfaces remain unchanged.

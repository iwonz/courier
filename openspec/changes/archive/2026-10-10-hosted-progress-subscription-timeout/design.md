## Context

`ipc.Configure` intentionally gives request/response calls a five-second deadline when their context has no deadline. `worker.Subscription.Next` currently reuses that helper even though progress subscriptions are streams that can legitimately be silent indefinitely. After the initial snapshot, an idle foreground route receives no event, times out after five seconds, is classified as unexplained worker loss, and releases the foreground lease.

## Goals / Non-Goals

**Goals:**

- Keep healthy foreground browser and webhook routes alive across arbitrary idle intervals.
- Preserve context deadlines and cancellation as prompt ways to interrupt a blocked subscription read.
- Preserve bounded defaults for ordinary IPC RPCs and lease operations.
- Exercise the real compiled binary beyond the former timeout boundary.

**Non-Goals:**

- Add keepalive traffic, polling, retries, reconnects, or a configurable timeout flag.
- Change worker-loss classification, stop semantics, output, or any protocol/state version.
- Change background delivery lifetime.

## Decisions

### Streaming IPC has distinct deadline semantics

Add a streaming connection configurator alongside the existing bounded configurator. It validates the connection and context, applies an explicit caller deadline when present, clears any inherited deadline otherwise, and installs the same cancellation interrupt that advances the socket deadline to unblock I/O. Ordinary `Call`, lease release, registration, and control operations continue using the five-second fallback.

### Only progress reads use streaming configuration

The subscription handshake remains a bounded IPC request. Once established, every `Subscription.Next` read uses streaming configuration. This keeps malformed or unreachable subscription setup bounded while allowing silence after readiness.

### Acceptance crosses the real boundary

Unit tests verify deadline clearing, explicit deadlines, cancellation, and configuration failures without sleeping for five seconds. The serial compiled-runtime browser test deliberately leaves a foreground delivery idle beyond five seconds before opening its URL, then completes normal control shutdown. This protects the exact production path that escaped the earlier fast E2E scenarios.

## Risks / Trade-offs

- **[A dead peer can leave a read blocked]** -> Socket close, process death, caller cancellation, and explicit deadlines still unblock the read; the foreground command intentionally owns an indefinite healthy lease.
- **[Deadline reset can race cancellation]** -> Each read installs and stops its own context callback, and no later read occurs after its context is canceled.
- **[Runtime acceptance becomes slower]** -> Add one bounded idle interval just beyond five seconds rather than delaying every hosted scenario.

## Migration Plan

Ship as the next patch release. Existing workers and clients use the unchanged protocol; only newly executed foreground clients receive corrected local deadline handling. No state migration or cleanup is required.

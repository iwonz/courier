# Design: Bounded Windows pipe acceptance

## Context

`go-winio` creates a sentinel pipe handle in `ListenPipe` and creates a connectable server instance when `Accept` starts. A client using `context.Background` can therefore remain in the native overlapped connection call indefinitely if the handshake stalls. A fixed UUID also makes the acceptance endpoint unnecessarily reusable.

## Decisions

### Unique native endpoint

Generate a fresh Courier delivery ID with the existing domain constructor and derive the pipe name through `ControlEndpoint`. This preserves validation and naming coverage while removing cross-run endpoint reuse.

### Bounded connection

Use a short explicit context deadline for the native dial. The timeout is a test lifecycle boundary, not a production transport policy.

### Cleanup ownership

Register cancellation and listener closure before starting `Accept`. The accept result channel remains buffered, allowing the goroutine to finish during cleanup even when the test has already failed.

## Verification

Cross-compile the Windows test binary, run exact local coverage and full release verification, then require repeated native Windows CI success together with green Linux, macOS, and Pages workflows.

## Why

Foreground commands can return a wrapped or joined `context.Canceled` after Ctrl+C. Courier currently recognizes only the unwrapped sentinel when choosing the terminal presentation, so an expected user interruption can be rendered as an operation failure even though the process exits with code `130`.

## What Changes

- Classify cancellation-only error trees as planned user interruptions, including wrapped and joined cancellation values.
- Keep mixed cancellation and cleanup/control failures in the operational failure path with their stable category exit code.
- Apply the same classification to ordinary commands and internal worker, administration, helper, and updater entry points.
- Add interactive and plain-output regression coverage for exact, wrapped, joined, and mixed cancellation results.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. The existing `operation-reporting` requirements already define neutral interactive interruption and require concurrent cleanup failures to remain errors; this change restores implementation conformance.

## Impact

This change affects application error classification, top-level internal-mode reporting, and focused reporting tests. It does not change command syntax, plain-output schemas, IPC/state versions, hosted routes, or lifecycle ownership.

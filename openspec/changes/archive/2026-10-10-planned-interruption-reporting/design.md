## Context

Courier's rich error boundary compares a returned cause directly with `context.Canceled`, while hosted foreground cleanup returns `errors.Join(waitErr, watchErr, releaseErr)`. Even when cancellation is the only non-nil cause, `errors.Join` produces a distinct wrapper, so the boundary renders a red failure panel and only later maps the tree to exit code `130`. The same exact-value check exists in the internal-mode boundary.

The existing `operation-reporting` specification already distinguishes a cancellation-only result from cancellation accompanied by cleanup or lease-release failure.

## Goals / Non-Goals

**Goals:**

- Recognize exact, wrapped, and joined cancellation-only error trees.
- Preserve neutral TTY reporting and exit code `130` for planned interruption.
- Preserve failure reporting and the owning stable category code when any non-cancellation error is present.
- Keep redirected output deterministic and backward compatible.

**Non-Goals:**

- Change signal handling, foreground ownership, external-stop behavior, or worker protocols.
- Treat deadlines, worker loss, or cleanup failures as user cancellation.
- Change the stable plain-output schema.

## Decisions

### Classify the complete error tree

Add one reporting helper that recursively follows both single-error `Unwrap() error` and joined `Unwrap() []error` trees. A result is a planned interruption only when the tree is non-empty and every leaf resolves to `context.Canceled`. This accepts contextual wrappers without allowing a sibling cleanup error to disappear.

Checking only `errors.Is(err, context.Canceled)` was rejected because it returns true for mixed joins and would incorrectly convert real cleanup failures into neutral stops.

### Select a non-cancellation command error for mixed joins

The application boundary will inspect command errors in the returned tree and prefer the first one whose cause is not cancellation-only. This retains its stage, counters, and category code when deferred cleanup joins an error to an interrupted transfer. A cancellation command error remains the fallback for a pure interruption so its stage and counters remain visible.

### Preserve plain interruption output

Interactive output will use the existing `Stopped by user` panel. Redirected output will retain its current stable operational fields and `result: failed`, while the process still exits `130`; changing that schema would be a separate public contract change.

## Risks / Trade-offs

- **[A custom wrapper adds semantics without a sibling error]** -> Courier treats a conventional single-cause wrapper as context only; operational errors must be represented as a separate joined cause or command error.
- **[Tree traversal encounters cycles]** -> Go's standard wrappers are acyclic; guard visited comparable error values so a malformed custom error cannot recurse indefinitely.
- **[Mixed cleanup output includes cancellation text]** -> Prefer the non-cancellation command error node, which reports the actionable failure rather than the cancelled sibling.

## Migration Plan

Ship as the next patch release. No persisted state, protocol, configuration, or user cleanup is involved; rollback is the previous binary.

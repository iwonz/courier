## Context

The release workflow intentionally requires aggregate `100.0%` statement coverage on Windows before any publication job can run. The new error-tree classifier already behaves correctly, but its empty-join rejection and traversal depth guard are not reached by the current cross-platform unit suite.

## Goals / Non-Goals

**Goals:**

- Cover both defensive branches with deterministic platform-independent tests.
- Preserve the implementation and all public behavior shipped in the release commit.

**Non-Goals:**

- Relax or round the Windows coverage gate.
- Move, delete, or recreate `v0.3.15`.
- Change cancellation classification or terminal rendering.

## Decisions

Add purpose-built test error types and wrapper chains rather than excluding defensive code from coverage. The empty joined-error type verifies that a multi-error tree with no cancellation leaf is false. A chain longer than the traversal limit verifies bounded command-error lookup without recursion failure.

## Risks / Trade-offs

- **[The tests couple to defensive limits]** -> Assert only safe termination and no selected command error; the explicit limit remains an intentional implementation boundary.
- **[Another platform changes aggregate rounding]** -> The tests run identically on every Go platform and directly execute the previously uncovered statements.

## Migration Plan

Ship as the next patch. The failed `v0.3.15` tag remains immutable and unpublished; `v0.3.16` becomes the published correction.

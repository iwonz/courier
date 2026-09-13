# Context

The target CLI uses one `from <source> to <destination>` handler for file, browser, and webhook operations. Parsing those modes independently in later handlers would duplicate ambiguity rules and allow invalid requests to acquire resources before rejection.

# Decisions

## Endpoint kind is explicit

Every parsed endpoint has a stable kind. `web://` and `webhook://` are exact sentinels, HTTP(S) endpoints retain their parsed URL, and SSH paths retain authority plus path flavor. The legacy `Remote` field remains temporarily available to existing transfer code and is derived consistently for SSH endpoints.

## Route planning is pure

The planner accepts endpoint strings and ordered option occurrences, then returns typed endpoints, a route, a path-direction subtype, and effective options. It performs no I/O and invokes no callbacks.

## Option metadata is reusable

Definitions declare repeatability, route applicability, conflicts, and defaults. Occurrences retain their original index. Ordered selection values preserve interleaving between `--exclude`, `--exclude-regex`, and `--exclude-from` without tying domain code to Cobra.

# Risks / Trade-offs

The planner initially models planned behavior not exposed by Cobra. This deliberate separation keeps public help honest while allowing later tasks to reuse one validated route matrix.

# Migration Plan

Add endpoint kinds, add the planner and option value types, route the shipped transfer through the planner, exhaustively test the matrix, then archive the specification delta.

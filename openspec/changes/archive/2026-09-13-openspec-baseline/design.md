# Context

All existing OpenSpec tasks are complete and validated, but none has been archived into the main specification set.

# Decisions

## Preserve implementation order

Archive changes in the same order as their stacked branches so later deltas apply to the requirements established by earlier work.

## Archive future changes at completion

Each future task starts as an active change and ends with its artifacts preserved in the OpenSpec archive and its delta applied to the baseline.

# Risks / Trade-offs

Archiving exposes conflicts that were previously isolated between active changes. Resolve only structural overlaps and keep the resulting baseline faithful to the shipped repository.

# Migration Plan

Archive changes 001 through 011, validate the resulting specs, update contributor guidance, then archive this workflow-only change with `--skip-specs`.

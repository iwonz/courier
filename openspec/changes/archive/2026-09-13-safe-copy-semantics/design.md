# Context

Courier already stages data under a private partial name, but its final commit renames the old target aside and deletes it after replacement. A preflight-only collision check is also insufficient because a target can appear between validation and commit.

# Decisions

## Identity is an operation disposition

Safety evaluation returns either proceed or no-op. Plain identity is a successful no-op; archive and extraction identity remain unsafe because the transformed result can collide with its own input.

## Existing final paths are immutable

The application checks the resolved final path before archive creation or transfer. The shared engine repeats that check and the commit boundary uses an optional backend no-replace transaction so a late collision cannot trigger replacement.

## Directory containers are not merge targets

An existing destination directory still acts as a container under destination semantics. Only its computed child is the final path. Other children are never scanned for deletion or changed.

# Risks / Trade-offs

Whole-directory atomic rename is not uniformly available across all transports. Backends may reserve an absent directory and move staged children into it, guaranteeing no overwrite and rollback on failure while allowing a short-lived partial view.

# Migration Plan

Change the safety disposition, add collision preflight and no-replace commit, replace synchronization tests with preservation and collision tests, then archive the delta.

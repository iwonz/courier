# Context

The operation planner already preserves ordered selection-option occurrences. The existing copy and archive walkers currently include every supported source object.

# Decisions

## One compiled selector

The selection package consumes planner rules once, expands `--exclude-from` at its occurrence position, and emits one immutable matcher. Gitignore rules retain order and negation semantics; any matching regular expression excludes the path.

## Paths are root-relative slash paths

Every consumer presents normalized slash-separated paths relative to the selected data root. The root itself is always retained. Directory exclusions prune traversal according to gitignore parent exclusion semantics.

## Compilation precedes resource acquisition

The command compiles selection rules before local/SSH endpoint opening. Rule files are local configuration input and do not travel to remote systems.

# Risks / Trade-offs

Gitignore negation cannot re-include a descendant of an excluded parent unless the parent is also re-included, matching Git behavior. Regular expressions are exclusion-only and therefore do not override one another.

# Migration Plan

Add the selector and tests, expose the three flags through the existing provider, integrate both walkers, update the canonical contract, then archive the delta.

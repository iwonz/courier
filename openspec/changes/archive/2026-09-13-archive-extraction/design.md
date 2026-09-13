# Context

Extraction accepts local or SSH-backed archives and writes to local or SSH-backed roots. A one-pass implementation cannot prove that a late archive entry is safe before earlier entries have changed destination state.

# Decisions

## Registry owns format dispatch

An immutable registry maps normalized extensions to codecs and rejects duplicate registrations. Tar.gz provides inspection and extraction; later formats can implement the same internal contract without changing command parsing.

## Inspection and extraction are separate passes

The codec reopens the source after a complete inspection. Inspection streams every payload, validates normalized unique paths and safe relative symlink targets, applies fixed and configurable limits, and checks selected destination paths without writing.

## Extraction root is not ordinary destination resolution

The destination argument is the root itself. If absent, the complete staged root is committed. If it is an existing directory, staged top-level entries are committed individually with no-replace semantics, leaving unrelated entries untouched.

## Fixed bomb limits are non-configurable

Entry count, depth, and expansion ratio remain active when max extracted size is unlimited. The default configurable expanded-size cap is 100 GiB.

# Risks / Trade-offs

Committing several top-level entries into an existing root cannot be one atomic filesystem operation. Courier completes collision preflight before commit and uses no-replace operations for every top-level entry; a late race can yield explicitly partial confirmed output without overwriting data.

# Migration Plan

Add the registry and tar.gz codec, integrate extract orchestration and flags, expand security tests, update the canonical contract and documentation, then archive the delta.

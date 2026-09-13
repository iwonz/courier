## Why

Without an unambiguous local/remote endpoint model, Courier cannot safely resolve the actual destination or block destructive self-copy operations. These rules must stay transport-independent and run before opening the destination.

## What Changes

- Parse local paths and `[user@]host:/path` remote endpoints, including SSH aliases and bracketed IPv6.
- Distinguish Windows drive paths from remote endpoints.
- Resolve the actual destination from a directory hint/stat result and source name.
- Reject identical paths, directory self-nesting, control characters, and archive traversal.
- Add unit tests for spaces, Unicode, special characters, and symlink canonicalization.

## Capabilities

### New Capabilities

- `endpoint-model`: local/remote endpoint model, parsing, and destination semantics.
- `path-safety`: preflight equality, nesting, and safe path-join checks.

### Modified Capabilities

None.

## Impact

Adds internal `endpoint` and `safety` packages. The public CLI does not change yet, and no transport or filesystem I/O occurs.

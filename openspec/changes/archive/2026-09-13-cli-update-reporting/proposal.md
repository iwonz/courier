## Why

Courier's transfer engine and native transport need a stable user-facing command, terminal reporting, and a safe way to stay current. Established CLI and progress libraries reduce custom parsing/rendering while Courier retains ownership of transfer semantics and recovery.

## What Changes

- Replace the foundation parser with a Cobra command tree and an isolated `from` handler.
- Implement all local/remote direction orchestration with independent concurrent SSH setup through `errgroup`.
- Use `os.Root` for bounded local filesystem access instead of relying on lexical path checks alone.
- Render interactive progress through mpb while retaining confirmed-byte structured events.
- Print deterministic success/error summaries and stable exit codes.
- Add `courier update` with GitHub release checks, platform asset selection, checksum verification, private staging, and executable replacement.
- Use `x/term` for no-echo password/passphrase prompts and never expose a password flag.

## Capabilities

### New Capabilities

- `transfer-command`: complete `from <source> to <destination>` orchestration and flags.
- `self-update`: verified GitHub release discovery and binary replacement.
- `operation-reporting`: interactive progress and final success/error summaries.

### Modified Capabilities

- `cli-foundation`: use Cobra's extensible command tree instead of the initial first-party registry.
- `path-safety`: bind local filesystem operations to `os.Root` in addition to canonical preflight validation.

## Impact

Adds Cobra, mpb, `x/sync/errgroup`, and `x/term`; introduces `app`, `report`, and `update` packages; rewires `cmd/courier`; and adds bounded local filesystem support in `fsx`.

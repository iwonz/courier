# Design: Native Windows quality gate hardening

## Context

Courier preserves metadata where the host filesystem supports it. Windows does not report Unix permission bits with the same semantics as POSIX systems, and local paths must use the host separator. Tests must validate the portable contract rather than incidental Unix representations.

## Decisions

### Host-native local paths

Local destination and symlink-target expectations use `filepath.Join`. Remote paths continue to use slash semantics. Expanded SSH identity files are local resources, so Courier cleans them with `filepath.Clean` after token and home expansion. Updater assertions verify the installed file instead of comparing a pre-resolution path with the canonical path returned by Windows.

### Capability-aware metadata assertions

Functional file contents, timestamps, traversal protection, staging, cleanup, and symlink behavior remain covered on Windows. Exact Unix permission-bit assertions run only where those bits have stable POSIX meaning.

### Isolated updater handoff

The runtime-default updater test verifies platform asset and executable selection. On Windows it injects the handoff boundary because its three-byte fixture is intentionally not a runnable PE executable. Dedicated handoff tests retain coverage of the real Windows control flow.

### Fail-fast reusable Windows runner

One PowerShell script runs `go test`, stops immediately on a non-zero native command exit, validates exact aggregate statement coverage, and removes its coverage profile in `finally`. Both push CI and tag release gates invoke that script. Push CI includes `fix/**` branches so corrective work is validated before it advances to `main`.

## Testing

- Run the full local verification and GoReleaser snapshot gate.
- Cross-compile Windows test binaries before pushing.
- Require a successful native Windows GitHub Actions job with exact 100% coverage and installer acceptance.

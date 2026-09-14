# Design: Deterministic workflow shell validation

## Context

Actionlint performs its native workflow checks without external dependencies, but its shell-script integration is enabled only when ShellCheck can be executed. The developer host did not have ShellCheck while GitHub's Ubuntu runner did, creating a platform-dependent gate.

## Decisions

### Repository-cached tool

`workflow-check` resolves a pinned ShellCheck under `.cache/tools`. When absent or at the wrong version, a POSIX bootstrap selects macOS/Linux and amd64/arm64, downloads the official release archive, verifies a release-asset SHA-256 digest, extracts only the executable, and atomically moves it into the ignored cache. Unsupported hosts or checksums fail closed.

### Explicit actionlint integration

Make passes the resolved executable through actionlint's `-shellcheck` option. The validation can no longer degrade silently based on ambient `PATH`. CI therefore exercises the same ShellCheck version and configuration as a clean local verification.

### Minimal workflow correction

The npm visibility loop uses its attempt counter to avoid sleeping after the final failed lookup. This preserves six checks and five ten-second retry intervals while satisfying ShellCheck without a suppression.

## Verification

Run the bootstrap twice to cover installation and cache reuse, run workflow validation with the pinned tools, run strict OpenSpec validation and the complete verification gate, then confirm a fresh Linux CI run succeeds.

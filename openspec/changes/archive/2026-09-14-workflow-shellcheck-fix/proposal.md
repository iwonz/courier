# Change: Make workflow shell validation deterministic

## Why

The Linux runner found an unused retry variable in the release workflow after the local macOS verification gate passed. Actionlint delegates embedded shell scripts to an external ShellCheck executable, so the check was silently weaker when ShellCheck was absent locally.

## What Changes

- Bootstrap a pinned checksum-verified ShellCheck binary into the ignored repository tool cache on supported developer hosts.
- Always pass that exact executable to actionlint during `make verify`.
- Use the release-publication retry counter so the existing workflow passes the same shell validation locally and in CI.

## Impact

The local verification command gains one cached development tool download and now rejects embedded workflow-shell defects before push. Runtime binaries, release artifacts, package consumers, Pages content, and credentials are unchanged.

# Change: Make native Windows CI portable and fail fast

## Why

The first native Windows run exposed tests that assumed POSIX path separators and permission bits, plus an updater test that tried to execute fixture bytes during the Windows handoff. The workflow also continued into coverage analysis after a failed `go test`, obscuring the primary failures.

## What Changes

- Express local destination expectations with host-native path semantics while retaining slash semantics for remote endpoints.
- Normalize expanded local SSH identity paths on the host platform.
- Assert POSIX permission bits only on operating systems that expose them consistently.
- Inject the already-tested Windows updater handoff boundary when a test uses non-executable fixture bytes.
- Run native Windows tests and exact coverage through one fail-fast PowerShell script shared by CI and release workflows.

## Non-goals

- Weakening the exact 100% first-party statement coverage gate.
- Changing runtime transfer, authentication, or release behavior beyond local path normalization.
- Emulating Windows execution on a non-Windows host.

## Impact

- Affected spec: `release-automation`.
- Affected packages: `internal/app`, `internal/endpoint`, `internal/sshx`, and platform-sensitive tests.
- Affected automation: CI and release Windows jobs.

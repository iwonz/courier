## 1. Runtime acceptance

- [x] 1.1 Add black-box coverage for four path directions, no-op, collisions, archive/extraction, interruption, Unicode, special paths, and cleanup
- [x] 1.2 Add large synthetic bounded-memory/backpressure tests and focused race-suite coverage
- [x] 1.3 Keep native SSH/SFTP, helper, HTTP security, worker, registry, and IPC acceptance in the unified gate

## 2. Browser and packages

- [x] 2.1 Add pinned Chromium tests for all three UIs, English/Russian, system/light/dark, keyboard, responsive, and auth leakage
- [x] 2.2 Add labeled ephemeral package installs for Ubuntu, Debian, Arch, Manjaro, Fedora, RHEL-compatible UBI, and Alpine
- [x] 2.3 Add cleanup traps and assertions scoped only to Courier-owned temporary files and Docker resources

## 3. CI and release candidate

- [x] 3.1 Run Linux, macOS, and Windows runtime suites and verify the complete primary/BSD artifact matrix
- [x] 3.2 Gate release publication on browser, package, workflow, contract, OpenSpec, installer, and GoReleaser checks
- [x] 3.3 Run the final clean-worktree verification, archive this change, and commit once

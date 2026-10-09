## 1. Specification and terminal behavior

- [x] 1.1 Add terminal-only exact-URL QR encoding and rendering with color, NO_COLOR, dumb-terminal, narrow-width, sanitation, and payload tests
- [x] 1.2 Integrate QR readiness with administration UI, browser upload/download, and incoming webhook output without changing plain output or listener reachability
- [x] 1.3 Add typed compact usage diagnostics for Cobra and route errors while preserving operational reports

## 2. Hosted foreground lifecycle

- [x] 2.1 Retain progress subscription filters and publish final delivery snapshots/tombstones before subscriber shutdown
- [x] 2.2 Attach a race-safe watcher to foreground acquisition and distinguish external stop, Ctrl+C, and unexpected worker loss
- [x] 2.3 Make foreground administration stop cleanly through verified control and retain cancellation exit 130

## 3. Real runtime acceptance

- [x] 3.1 Build a real Courier binary for a serial Playwright runtime suite with isolated homes, ports, bounded readiness, and mandatory cleanup
- [x] 3.2 Cover browser uploads and incoming webhooks for Unicode, binary, extraction, collision, and multipart rejection
- [x] 3.3 Cover browser downloads and outgoing HTTP multipart delivery for files, nested directories, direct payloads, and archives
- [x] 3.4 Cover external delivery/server stop, route removal, clean foreground exit, worker cleanup, and real administration UI start/stop

## 4. Mandatory change and release workflow

- [x] 4.1 Add root AGENTS.md and update README, CONTRIBUTING, and release documentation with the mandatory tracked-change lifecycle and recovery boundaries
- [x] 4.2 Add fail-closed `make change-start` for clean synchronized main, a single active OpenSpec change, feature-branch creation, and scaffold generation
- [x] 4.3 Rework `make ship` for default next-patch selection, completed-change validation, one commit, unchanged-main fast-forward merge, guarded push/tag, publication waits, branch deletion, and clean synchronized completion
- [x] 4.4 Add pre-push protection for main/release tags, recovery-only low-level targets, and release-gated workflow-dispatch-only Pages publication
- [x] 4.5 Test dirty/out-of-sync starts, incomplete changes, main movement, non-fast-forward state, failed gates/publication, recovery guards, and successful cleanup in temporary repositories with stubbed external commands

## 5. Verification and release

- [x] 5.1 Run focused Go, shell, and Playwright tests plus strict OpenSpec validation
- [x] 5.2 Run the complete `make verify` release-candidate gate and perform real browser visual verification through the Playwright workflow
- [x] 5.3 Ship the archived change as the next patch and verify GitHub Release, artifacts, npm, Scoop, Homebrew, Pages, deleted feature branch, and clean synchronized main

# Tasks

## 1. Correct platform assumptions

- [x] 1.1 Normalize explicit empty CLI arguments and expanded local SSH identity paths
- [x] 1.2 Use host-native local destination expectations
- [x] 1.3 Restrict exact POSIX permission assertions to supporting hosts
- [x] 1.4 Isolate non-executable updater fixtures from the real Windows handoff

## 2. Harden automation

- [x] 2.1 Add a fail-fast reusable native Windows test and coverage script
- [x] 2.2 Use the same script in CI and release workflows
- [x] 2.3 Run push CI for conventional `fix/**` branches

## 3. Verify

- [x] 3.1 Pass local formatting, vet, race, exact coverage, npm, installer, and GoReleaser snapshot checks
- [x] 3.2 Pass strict OpenSpec validation and workflow linting
- [x] 3.3 Pass the native Windows CI and PowerShell installer job

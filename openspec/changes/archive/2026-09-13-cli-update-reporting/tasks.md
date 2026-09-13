## 1. CLI and filesystem safety

- [x] 1.1 Replace the initial registry with Cobra root and isolated handlers
- [x] 1.2 Implement exact `from … to …` parsing and archive flag
- [x] 1.3 Add `os.Root`-bounded local backends
- [x] 1.4 Use `x/term` for no-echo interactive secrets

## 2. Orchestration and reporting

- [x] 2.1 Open independent remote endpoints with errgroup and deterministic cleanup
- [x] 2.2 Execute all four directions and archive staging
- [x] 2.3 Render structured progress through mpb and deterministic line fallback
- [x] 2.4 Map CLI, connection, transfer, and update failures to stable exit codes

## 3. Self-update

- [x] 3.1 Query and compare GitHub releases
- [x] 3.2 Select, download, checksum, and extract the runtime asset privately
- [x] 3.3 Replace the executable safely and clean all temporary state

## 4. Verification

- [x] 4.1 Cover commands, all directions, archive, reporting, update, and failures
- [x] 4.2 Run race, 100% coverage, cross-build, and strict OpenSpec validation
- [x] 4.3 Create the conventional commit on `feat/005-cli-update-reporting`

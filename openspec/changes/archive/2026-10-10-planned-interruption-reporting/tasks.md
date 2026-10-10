## 1. Cancellation classification

- [x] 1.1 Add cancellation-only error-tree classification and verify exact, wrapped, joined, mixed, nil, and cyclic cases with report package tests
- [x] 1.2 Apply cancellation-only classification to application and internal-mode boundaries and verify pure interruption keeps exit `130` while mixed failures keep their category code

## 2. Terminal reporting regression

- [x] 2.1 Add interactive command tests proving Ctrl+C renders `Stopped by user` without an operation-failed panel and retains stage/counters
- [x] 2.2 Add mixed cancellation and cleanup/control tests proving actionable failures remain red and are not reported as successful or neutral

## 3. Verification and release

- [x] 3.1 Run focused Go tests, strict OpenSpec validation, and the complete `make verify` gate
- [x] 3.2 Ship the archived change as the next patch and verify release surfaces, Pages, deleted feature branch, and clean synchronized main

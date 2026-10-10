## 1. Upload transport

- [x] 1.1 Add a progress-aware single-file browser upload adapter with bounded HTTP failures and verify progress, success, network, abort, timeout, and malformed-response tests

## 2. Batch interface

- [x] 2.1 Implement deterministic sequential multi-file queue state and verify independent progress/results, continuation after failure, overlap prevention, and same-selection retry
- [x] 2.2 Render a responsive accessible upload list with exact sizes, percentages, progress bars, localized states, and live batch summaries; verify English/Russian and narrow-layout behavior

## 3. Runtime acceptance

- [x] 3.1 Extend compiled browser E2E to select multiple files together, observe queue rows, and verify byte-identical committed outputs plus an independent collision result
- [x] 3.2 Rebuild embedded data assets and run focused web, browser, and strict OpenSpec checks

## 4. Verification and release

- [x] 4.1 Run the complete `make verify` release-candidate gate
- [x] 4.2 Ship the archived change as the next patch and verify release surfaces, package channels, Pages, deleted feature branch, and clean synchronized main

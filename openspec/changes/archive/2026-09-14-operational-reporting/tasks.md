## 1. Accounting and diagnostics

- [x] 1.1 Add monotonic read/sent/confirmed progress accounting and the closed stage vocabulary
- [x] 1.2 Add centralized credential-safe diagnostic sanitization and focused adversarial tests
- [x] 1.3 Update deterministic and interactive reports to use unified counters and sanitized reasons

## 2. Runtime integration

- [x] 2.1 Apply shared accounting to copy, archive, extraction, and outgoing HTTP streams
- [x] 2.2 Map context cancellation and interruption to exit code 130 without changing existing code assignments
- [x] 2.3 Add bounded private history metadata, validation, ordering, and rotation

## 3. Quality

- [x] 3.1 Add exact coverage for counter divergence, overflow, redaction, interruption, and history rotation
- [x] 3.2 Update operational documentation and baseline specifications
- [x] 3.3 Run full verification, archive this change, and commit once

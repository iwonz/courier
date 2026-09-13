## 1. Incoming webhook

- [x] 1.1 Extend private worker definitions and route validation for `webhook-to-path`
- [x] 1.2 Add the opaque no-UI `/upload` endpoint with exact multipart field/cardinality validation
- [x] 1.3 Reuse staged upload, extraction, selection, policy, rate, and cleanup behavior

## 2. Outgoing webhook

- [x] 2.1 Add a bounded one-attempt HTTP(S) multipart sender with redirects disabled
- [x] 2.2 Add Basic credential handling, response classification, and unknown-outcome errors
- [x] 2.3 Reuse local/SSH opening, verified archive preparation, ordered selection, progress, and cleanup

## 3. Contract and quality

- [x] 3.1 Ship only the two webhook routes and applicable flags in the canonical contract and generated reference
- [x] 3.2 Cover multipart validation, auth/IP/limits, collision/extraction, cancellation/rates, redirect/rejection/success, unknown outcomes, SSH sources, and secret redaction at 100%
- [x] 3.3 Document the exact webhook profile and lifecycle
- [x] 3.4 Run full verification, strict OpenSpec validation, archive the change, and commit once

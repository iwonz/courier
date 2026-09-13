## 1. Authentication and admission

- [x] 1.1 Add Argon2id credential verifiers, bounded authentication work, and request-scoped Basic verification
- [x] 1.2 Add delivery-scoped opaque sessions, expiration, CSRF validation, and secure cookie metadata
- [x] 1.3 Add canonical peer/CIDR admission, concurrent attempt counters, bans, and stop decisions

## 2. Transfer enforcement

- [x] 2.1 Add atomic delivery and byte reservations with idempotent release and streaming growth
- [x] 2.2 Add aggregate cancellation-aware upload/download throttles and live reconfiguration
- [x] 2.3 Define deterministic policy evaluation order and secret-free optimistic updates

## 3. Quality

- [x] 3.1 Cover concurrent thresholds, session isolation, forged forwarding data, reservations, and throttling at 100%
- [x] 3.2 Document the policy trust boundary, defaults, ordering, and runtime-only state
- [x] 3.3 Run full verification, strict OpenSpec validation, archive the change, and commit once

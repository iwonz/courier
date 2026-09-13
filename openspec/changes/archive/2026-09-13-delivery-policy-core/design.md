# Context

Several planned delivery types share the same controls, but their byte direction and authentication presentation differ. Implementing controls inside individual handlers would create observable inconsistencies and race-prone accounting. The policy core must remain usable without `net/http` while providing small adapters for later chi middleware.

# Decisions

## Secrets remain ephemeral

Credentials are converted immediately to Argon2id verifier records with a random salt and explicit parameters. Plaintext credentials are not retained. Verifier records, sessions, CSRF values, authentication fingerprints, and ban state live only in the worker process. Public snapshots and optimistic policy updates contain configuration but never secret state.

Basic authentication is request-scoped. Password authentication may issue an opaque random browser session that is scoped to one delivery, expires on both idle and absolute deadlines, and owns an independent CSRF token. Cookies are host-only, HttpOnly, SameSite=Strict, and Secure whenever the request transport is secure.

## Peer identity is connection-derived

Admission canonicalizes only the accepted peer `host:port` address. IPv4-mapped IPv6 addresses are unmapped before exact-IP or prefix comparison. `Forwarded`, `X-Forwarded-For`, and similar request headers are never policy inputs. Authentication failure fingerprints combine the canonical peer with the delivery identity and do not include credentials.

## Threshold decisions are atomic

Each fingerprint has one synchronized failure counter. Successful authentication clears only that fingerprint. Reaching the configured threshold either creates an in-memory ban or requests delivery stop according to policy. Concurrent failures cannot miss or exceed the transition decision. Policy updates reconcile attempt thresholds and admission lists without disclosing runtime state.

## Reservations precede reads

Incoming transfers reserve their declared file bytes and one delivery slot atomically before any body is read. Reservations reject negative or excessive sizes and prevent concurrent requests from oversubscribing the configured delivery limit. Unknown-length streams are charged incrementally and fail before a write would cross the limit. Release is idempotent.

## Rate limits are aggregate

One `x/time/rate` limiter per direction belongs to a delivery policy engine. All simultaneous legs share it, so configured rates are aggregate rather than per connection. Unlimited directions bypass waiting. Reconfiguration preserves a single shared limiter and applies the new rate and bounded burst safely.

## Middleware order is explicit

The reusable evaluation order is: recover/request bounds, canonical peer admission, ban check, authentication, CSRF for state-changing session requests, transfer reservation, directional throttling, then the operation handler. Later HTTP changes must compose handlers in this order and must release reservations during all exits.

# Risks / Trade-offs

In-memory bans and sessions intentionally disappear with the worker. Durable storage would expand the secret-bearing attack surface and is not required by the delivery lifecycle. Argon2 verification is intentionally expensive; concurrent verification is bounded to prevent memory exhaustion. Rate limiting operates on bytes actually admitted to a stream and cancellation may return partial progress.

# Migration Plan

Add policy primitives and exhaustive race tests first, then integrate them into web and webhook handlers in their own changes. Cross-build all client platforms, run full verification, archive this OpenSpec change, and commit it once before starting web delivery work.

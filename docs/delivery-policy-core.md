# Delivery policy core

Courier applies one in-memory policy engine to every browser and webhook delivery. The engine is deliberately independent from HTTP routing so data pages, webhook handlers, and the administration API cannot accidentally implement different security rules.

The enforcement order is fixed: request bounds, connection-peer admission, ban state, authentication, CSRF for state-changing session requests, transfer reservation, aggregate directional rate limiting, and finally the operation handler. A rejected request never advances to a later stage.

Credentials are converted to salted Argon2id verifier material as the worker starts. Plaintext credentials, verifier records, opaque session tokens, CSRF material, fingerprints, and ban state are runtime-only and are absent from registry files, public snapshots, URLs, command arguments, and logs. Basic authentication is request-scoped. Password authentication issues a random delivery-scoped session with idle and absolute expiration; its cookie is host-only, HttpOnly, SameSite Strict, and Secure on secure transports.

IP admission uses the actual accepted connection address. Forwarding headers have no effect. Exact IPv4/IPv6 addresses and canonical CIDR prefixes are supported, including canonical handling of IPv4-mapped IPv6 peers.

`--limit` reserves aggregate concurrent transfer slots before any request body is read. `--max-file-size` is checked against declared length at reservation time and against every increment for unknown-length streams. Releases are idempotent. Upload and download rates are independent but each direction is shared by all concurrent legs of one delivery, so limits are aggregate.

Optimistic updates require the expected public policy version. They can change allowlists, thresholds, limits, rates, and presentation policy while preserving runtime-only secret state. Changing authentication mode requires a separate credential-bearing private registration flow and cannot be performed through a secret-free policy update.

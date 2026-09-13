# Context

Incoming webhooks differ from browser uploads only at the HTTP surface: they have no page, session-password mode, or browser-origin requirement. Their filesystem transaction, endpoint ownership, selection, policy reservations, and worker lifecycle are identical. Outgoing webhooks are finite command operations rather than registered deliveries and must distinguish a known HTTP response from a response-less, potentially committed request.

# Decisions

## One hosted incoming transaction

The existing multi-delivery chi host registers `webhook-to-path` definitions through the same private IPC channel and exposes `/d/<opaque-token>/upload`. The handler accepts exactly one multipart file part whose field name is `file`. It runs peer/IP/Basic admission and reserves delivery capacity before reading the multipart body, then delegates to the same staged upload/extraction transaction used by browser upload. It never serves the Lit data page or password-session API for a webhook delivery.

## One-shot outgoing transport

A dedicated webhook sender owns multipart framing and HTTP result classification while app orchestration owns endpoint opening, selection, optional archive creation, progress reporting, and cleanup. The sender streams a single regular file through a bounded buffer and aggregate upload limiter. The default native HTTP client uses normal TLS verification, bounded connect/TLS/header timeouts, proxy environment support, and a redirect policy that returns the first 3xx response without following it.

Courier performs exactly one `RoundTrip` attempt. A 2xx response means the receiver accepted the HTTP request; it does not prove durable storage, including for 202. A received non-2xx response is a known rejection. If no response arrives after any payload bytes were consumed, the result is explicitly unknown and is never retried automatically.

## Authentication and secrets

Incoming webhook authentication is `none` or Basic through the shared policy engine. Outgoing Basic credentials are prompted before opening the source and are applied with `Request.SetBasicAuth`; password mode remains invalid in the operation planner. Credentials are held in byte slices where possible, cleared after request construction, and never enter endpoints, registry state, output, diagnostics, or logs.

# Risks / Trade-offs

Generic webhooks have no universal contract. Courier intentionally implements only its documented multipart profile, so integrations that require JSON, raw bodies, signatures, redirects, idempotency keys, or custom headers need a later versioned extension. A transport error after bytes leave the process cannot prove whether a receiver accepted them; reporting unknown is safer than an implicit retry that could duplicate side effects.

# Migration Plan

Add and validate this OpenSpec, generalize the hosted upload transaction without duplicating it, implement and test the one-shot sender, wire both routes into the existing `from` handler, update the canonical contract and generated reference, run full verification, archive this change, and commit it once.

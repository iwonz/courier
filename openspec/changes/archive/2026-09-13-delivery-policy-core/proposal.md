# Change: Add the delivery policy core

## Why

Web and webhook deliveries need one security boundary before they expose HTTP handlers. Authentication, peer admission, concurrent limits, and traffic shaping must behave identically across data and administration surfaces without persisting credentials or trusting proxy-supplied addresses.

## What Changes

- Add in-memory Argon2id credential material and constant-time Basic/password verification.
- Add opaque, expiring, delivery-scoped browser sessions with CSRF tokens and secure cookie metadata.
- Add concurrency-safe authentication attempt accounting, fingerprint bans, and stop-on-threshold decisions.
- Add canonical peer IP and CIDR admission based only on the accepted connection address.
- Add atomic transfer reservations for delivery-wide and per-file size limits.
- Add aggregate upload and download throttles using `golang.org/x/time/rate`.
- Define one deterministic policy evaluation and HTTP middleware order.
- Apply optimistic public policy updates without returning, logging, or persisting secret material.

## Impact

This change adds a reusable internal policy package and extends delivery policy validation. It does not expose HTTP routes, public commands, or UI. Those integrations remain planned for later changes.

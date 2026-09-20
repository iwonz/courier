## Context

See `proposal.md` for motivation. During `v0.3.2`, npm returned a successful signed publish response and explicitly reported asynchronous processing, but the public registry remained stale beyond the workflow's six ten-second probes.

## Goals / Non-Goals

**Goals:**

- Distinguish a normal bounded npm processing delay from a missing publication.
- Keep the release workflow deterministic and non-interactive.
- Preserve a finite failure boundary for genuine publication problems.

**Non-Goals:**

- Republish an accepted immutable npm version.
- Add another registry, credential, package, or manual approval.
- Mask authentication or package-content failures in the publish job.

## Decisions

- Probe for up to six minutes at ten-second intervals. The observed delay was under four minutes; six minutes gives headroom without making a real failure unbounded.
- Pass `--prefer-online` to `npm view` so each probe revalidates registry metadata rather than relying on a stale package cache.
- Keep visibility verification after the independent publish job succeeds. The publish job remains responsible for authentication and immutable version creation; the verification job remains responsible for public availability.

## Risks / Trade-offs

- [Risk] A genuine registry failure takes five minutes longer to report. → Keep a fixed 36-attempt limit and emit the exact package/version on exhaustion.
- [Risk] Additional online probes add registry traffic. → Retain a conservative ten-second interval and stop immediately on success.

## Migration Plan

Deploy with the next patch tag. Existing releases and npm versions remain immutable and require no migration.

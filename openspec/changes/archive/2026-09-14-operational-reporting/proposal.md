# Change: Unify operational reporting

## Why

Courier reports basic confirmed-byte progress, but path transfers, HTTP deliveries, control failures, persisted history, and interruption do not yet share one complete accounting and diagnostic contract. Operators need consistent stages and monotonic read/sent/confirmed counters without exposing credentials, while automation needs the documented interruption exit code and bounded private history.

## What Changes

- Extend structured progress with explicit read, sent, and confirmed counters while preserving deterministic terminal and non-terminal rendering.
- Centralize safe diagnostic redaction and stable exit-code classification, including exit code 130 for cancellation or interruption.
- Make success and failure summaries use the same accounting vocabulary and render only sanitized reasons.
- Add bounded chronological delivery history rotation and reject regressing counter relationships.
- Feed the shared progress/accounting model through file, archive, extraction, and outgoing HTTP streams without adding unbounded buffering.

## Impact

Public command syntax is unchanged. Existing exit codes remain stable, with the already reserved interruption code becoming effective. History remains private, atomic, and bounded; no credential, URL user-info, query secret, or raw control-plane error is persisted or printed.

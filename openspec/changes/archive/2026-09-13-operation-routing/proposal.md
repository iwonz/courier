# Change: Add typed operation routing

## Why

Courier currently treats every endpoint as either a local path or an SSH path and begins opening resources before it has a complete, reusable model of the requested route and options. Browser and webhook delivery modes need unambiguous endpoint syntax and a side-effect-free preflight boundary before their runtimes are added.

## What Changes

- Classify local, SSH, `web://`, `webhook://`, HTTP, and HTTPS endpoints without confusing Windows drives, SSH aliases, IPv6, or unknown URI schemes.
- Add a typed operation planner for every supported endpoint direction.
- Validate option applicability, conflicts, and repeated non-repeatable occurrences without filesystem, network, SSH, or worker access.
- Preserve the relative order of interleaved selection options for the later shared selection engine.

## Impact

Existing path-to-path behavior remains public and compatible. Planned service routes and options become representable and testable, but remain absent from live Cobra help until their implementation changes the canonical contract status to shipped.

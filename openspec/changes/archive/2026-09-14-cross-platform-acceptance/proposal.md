# Change: Complete cross-platform acceptance

## Why

Courier's feature branches have exact unit coverage and release dry runs, but the release candidate needs one explicit acceptance layer for complete routes, real browsers, platform builds, distribution packages, bounded streams, race-sensitive lifecycle code, and owned-resource cleanup. CI must prove these surfaces before publication instead of relying on artifact existence alone.

## What Changes

- Add black-box CLI acceptance for all path directions, archive/extraction, collision/no-op, interruption, Unicode, special paths, staging cleanup, and stable reporting.
- Add large synthetic bounded-stream and backpressure checks plus focused race suites for registries, counters, reservations, policy, IPC, and worker shutdown.
- Add real Chromium coverage for data UI, admin UI, and landing across English/Russian, system/light/dark, keyboard, responsive, and protected-metadata cases.
- Add labeled ephemeral package-install checks for Ubuntu, Debian, Arch Linux, Manjaro, Fedora, RHEL-compatible UBI, and Alpine, with immediate traps and post-run ownership assertions.
- Expand GitHub Actions to run runtime suites on Linux, macOS, and Windows, cross-build all primary/helper artifacts, validate workflows, and run browser/package acceptance before release publication.
- Remove the obsolete Winget pull-request verifier from the consolidated release baseline so acceptance reflects Courier's single-repository distribution policy.

## Impact

Public CLI behavior is unchanged. `make verify` becomes the final local release-candidate gate and requires Docker plus a locally installed Playwright Chromium browser in addition to the documented Go/Node toolchain. All new temporary resources are uniquely labeled or created under private temporary directories and are checked after cleanup.

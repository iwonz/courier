## ADDED Requirements

### Requirement: Route and transformation acceptance

Courier SHALL exercise all four filesystem directions, browser/webhook route handlers, archive creation/extraction, no-op, collisions, interruption, path safety, Unicode, special paths, partial reporting, and cleanup through isolated acceptance fixtures.

#### Scenario: An acceptance operation is interrupted

- **WHEN** a transfer is canceled after staging data
- **THEN** Courier returns exit code 130, leaves the previous final namespace unchanged, and removes only its owned partial resources

### Requirement: Bounded performance acceptance

Courier SHALL verify large synthetic streams with bounded buffers, observable backpressure, cancellation, and no payload-sized memory allocation.

#### Scenario: A slow destination blocks

- **WHEN** a source produces data faster than the destination accepts it
- **THEN** transfer reads remain bounded by the configured buffer and cancellation releases the operation and its staging resources

### Requirement: Cross-platform release acceptance

Courier SHALL run compiled runtime suites on Linux, macOS, and Windows, cross-build the declared primary and BSD helper artifacts, and install matching release packages in isolated Ubuntu, Debian, Arch Linux, Manjaro, Fedora, RHEL-compatible, and Alpine containers.

#### Scenario: A distribution package check exits

- **WHEN** installation succeeds, fails, or is interrupted
- **THEN** its cleanup trap removes every uniquely labeled Courier container and proves that no labeled container, network, volume, or temporary file remains

### Requirement: Real-browser acceptance

Courier SHALL test data UI, admin UI, and landing behavior in Chromium for English/Russian catalogs, system/light/dark preferences, keyboard navigation, responsive viewports, and protected-metadata non-disclosure.

#### Scenario: Protected delivery returns an authentication error

- **WHEN** the data page receives an unauthenticated response containing an adversarial secret marker
- **THEN** the marker is not rendered in either locale and the accessible generic authentication state is shown

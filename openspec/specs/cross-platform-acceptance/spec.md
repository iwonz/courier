# cross-platform-acceptance Specification

## Purpose
Define the release-candidate acceptance layers for runtime routes, bounded streams, browsers, platforms, distribution packages, and owned-resource cleanup.

## Requirements

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

Courier SHALL run compiled runtime suites on Linux, macOS, and Windows, cross-build the declared primary and BSD helper artifacts, and install matching release packages in isolated Ubuntu, Debian, Arch Linux, Manjaro, Fedora, RHEL-compatible, and Alpine containers. Native Windows acceptance SHALL exercise owner-scoped named pipes using the current-process identity, a unique valid endpoint per test run, bounded connection contexts, cleanup registered before the handshake, rooted filesystem operations across the native-path and `fs.FS` path dialects, portable archive symlink targets, and platform-supported metadata semantics.

#### Scenario: A distribution package check exits

- **WHEN** installation succeeds, fails, or is interrupted
- **THEN** its cleanup trap removes every uniquely labeled Courier container and proves that no labeled container, network, volume, or temporary file remains

#### Scenario: Windows uses a rooted directory path

- **WHEN** a native Windows path produced by `filepath` crosses into an `fs.FS` operation
- **THEN** the adapter supplies slash-separated syntax without weakening the rooted confinement boundary

#### Scenario: Windows creates a private control pipe

- **WHEN** Courier resolves the named-pipe security principal
- **THEN** it queries the supported current-process token and grants pipe access only to that user SID

#### Scenario: Windows pipe acceptance stalls

- **WHEN** a native named-pipe test cannot complete its client-server handshake
- **THEN** its bounded context terminates the attempt and pre-registered cleanup releases the listener and accept goroutine without waiting for the package timeout

### Requirement: Real-browser acceptance

Courier SHALL test data UI, admin UI, and landing behavior in Chromium for English/Russian catalogs, system/light/dark preferences, reduced motion, touch and fine-pointer input, keyboard navigation, responsive viewports, protected-metadata non-disclosure, natural landing scrolling, responsive panorama loading, external-request isolation, and layout stability. The acceptance gate SHALL enforce a maximum of 225 KiB per landing WebP, 1.6 MiB for all eight landing panorama assets, 45 KiB gzip for landing JavaScript, and 8 KiB gzip for landing CSS.

#### Scenario: Protected delivery returns an authentication error

- **WHEN** the data page receives an unauthenticated response containing an adversarial secret marker
- **THEN** the marker is not rendered in either locale and the accessible generic authentication state is shown

#### Scenario: The landing first renders

- **WHEN** the initial landing viewport loads at a supported wide or portrait size
- **THEN** it contains no visible tagline, Run or Replay action, transcript stage, editable command, horizontal overflow, or external runtime request and it requests only the selected hero asset plus at most the next panorama segment

#### Scenario: A visitor uses natural navigation

- **WHEN** the visitor scrolls or activates route, installation, or CLI navigation
- **THEN** the page does not snap between sections, the masthead identifies the section at its header-aware reading band, and reduced-motion preference removes smooth scrolling

#### Scenario: A visitor crosses a panorama boundary

- **WHEN** adjacent landing sections are visible together at a supported viewport
- **THEN** their matching responsive segments overlap without uncovered background, visible tiling, geometry shift, or loading the inactive responsive source

#### Scenario: A visitor uses landing command surfaces

- **WHEN** route, installation, or CLI Copy is activated in English or Russian
- **THEN** the exact immutable command is offered to the clipboard and localized status is announced without a simulation, API mutation, or layout shift

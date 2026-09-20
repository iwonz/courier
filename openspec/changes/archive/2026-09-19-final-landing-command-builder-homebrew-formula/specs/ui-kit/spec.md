## MODIFIED Requirements

### Requirement: Auditable identity assets

Courier SHALL retain exactly the declared five transparent pixel Relay WebP assets, one pinned local display font, and the declared local transparent PNG third-party marks with role, dimensions, byte count, generation prompt or upstream source, lineage or revision, license, and SHA-256 digest. Pixel treatment SHALL remain limited to Courier/Relay and functional internal icons. Third-party brands SHALL preserve official geometry and color without pixelated image treatment. The canonical neutral-v1 Relay SHALL anchor identity-preserving mark, route, delivery, and administration v2 derivatives with recorded provenance and visual QA. The mark SHALL remain at most 24 KiB, neutral and route sprites at most 64 KiB each, delivery and administration sprites at most 48 KiB each, all Relay assets at most 248 KiB combined, and the display font at most 48 KiB.

#### Scenario: Asset integrity check

- **WHEN** repository verification runs
- **THEN** undeclared assets, legacy Relay names, missing provenance or licenses, incorrect dimensions, unofficial or pixelated third-party branding, broken neutral lineage, or an exceeded budget fail the gate

#### Scenario: Identity assets are validated

- **WHEN** the asset gate inspects Relay and third-party media
- **THEN** it verifies alpha, dimensions, hashes, neutral-source linkage, consumer references, official brand provenance, and the absence of third-party pixel rendering

### Requirement: Shared command readout

Courier SHALL provide an independent immutable command strip with localized heading, an optionally disabled Copy action, reserved live feedback, optional details and footer actions, and deterministic reset when command identity changes. It SHALL not execute the represented command.

#### Scenario: Copy is unavailable

- **WHEN** a consumer supplies an incomplete or invalid generated command
- **THEN** Copy is disabled and no clipboard write occurs

#### Scenario: Copy succeeds or fails

- **WHEN** a user activates Copy and clipboard access succeeds or fails
- **THEN** the exact command is offered to the clipboard and a localized result is announced in the reserved status region without geometry movement

#### Scenario: The command identity changes

- **WHEN** a consumer replaces the displayed command selection
- **THEN** prior copy feedback resets synchronously and no timed work remains

### Requirement: Local official brand marks

Courier SHALL provide a typed first-party pixel icon registry for functional endpoint, preference, file, transfer, status, and administration icons. GitHub, installation, platform, operating-system, distribution, shell, and package-manager identities SHALL use typed local transparent PNG marks in official geometry and color without pixel sampling. Every third-party raster SHALL retain pinned official source, revision, color, license, attribution, and trademark records. Channels without independent official marks, including npx and wget, SHALL remain text-only. Locale controls SHALL use local first-party pixel flags rather than platform emoji.

#### Scenario: A browser surface renders an icon

- **WHEN** Courier renders a functional icon, channel mark, or locale flag
- **THEN** it uses the declared local first-party registry or official brand PNG with an accessible name where needed and makes no runtime request

#### Scenario: A branded channel is rendered

- **WHEN** a supported package manager, shell, operating system, or distribution mark appears
- **THEN** the UI uses its transparent official-color PNG through an `img`, retains its pinned source record, avoids pixelated rendering, and exposes no remote asset URL

### Requirement: Shared Bezier route geometry

Courier SHALL preserve measured cubic route positioning while sampling and snapping its visual points to the four-pixel grid. The landing route signal SHALL be exactly four by four pixels with a centered anchor and linear motion on the existing path; reduced motion SHALL place that center at the path midpoint without travel or repeated animation.

#### Scenario: A route is rendered

- **WHEN** finite Source and Destination positions are supplied
- **THEN** a deterministic crisp polyline connects them, its terminals remain fixed-aspect, and the signal center follows that path or rests at its midpoint under reduced motion

#### Scenario: A consumer supplies two endpoint positions

- **WHEN** a route path is requested for finite coordinates
- **THEN** the helper returns a stable cubic source path and a deterministic four-pixel-grid sample for the rendered connector

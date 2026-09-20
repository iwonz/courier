## MODIFIED Requirements

### Requirement: Auditable identity assets

Courier SHALL retain exactly five transparent pixel Relay v3 WebP assets, one pinned local Pixelify Sans display font, four pinned local Overpass Mono Latin/Cyrillic 400/600 subsets totaling at most 48 KiB, their OFL notices, and the declared local transparent PNG third-party marks. Neutral-v3 SHALL be the identity authority for mark, route, delivery, and administration v3 assets. Provenance SHALL record role, dimensions, byte count, generation prompt or upstream source, lineage or revision, license, consumers, and SHA-256 digest. Legacy Relay v1/v2 files SHALL NOT ship. Existing per-sprite and combined budgets remain unchanged.

#### Scenario: Asset integrity check

- **WHEN** repository verification runs
- **THEN** undeclared assets, legacy Relay files, missing provenance or licenses, incorrect dimensions, broken v3 neutral lineage, altered third-party branding, excess font bytes, or an exceeded sprite budget fail the gate

#### Scenario: Identity assets are validated

- **WHEN** the asset gate inspects Relay, fonts, and third-party media
- **THEN** it verifies alpha, dimensions, hashes, neutral-v3 linkage, consumers, font licenses and sizes, official brand provenance, and the absence of third-party pixel rendering

### Requirement: Coherent brand identity system

Courier SHALL use one repository-local terminal 8-bit Relay v3 family across landing, delivery, administration, README, favicon, and product chrome. Every v3 role SHALL preserve the established Relay silhouette, anatomy, proportions, eye, beak, feather structure, earpiece, satchel, pose, and working object while changing only cobalt to teal, orange/coral to yellow, and plumage to a cool black/white neutral ramp. Assets SHALL retain true alpha and hard pixel clusters.

#### Scenario: A browser surface presents Courier

- **WHEN** landing, delivery, or administration renders
- **THEN** its header uses mark-v3 beside `COURIER CLI`, its optional v3 role sprite matches the surface, and product meaning remains complete without the image

#### Scenario: Relay changes roles

- **WHEN** all five v3 assets are compared on light, dark, and checkerboard backgrounds
- **THEN** the same recognizable Relay identity, role poses, object count, dimensions, transparent background, and terminal palette remain visible

### Requirement: Intentional control appearance

Courier SHALL present all three browser surfaces through one terminal data-grid grammar. Dark mode SHALL use `#000000`, `#0D1015`, `#191C20`, white, 70% white secondary text, and 55% white rules. Light mode SHALL invert those roles with white and 4%/9%/70%/55% black. Action/success SHALL use `#71FFF6`, selection/warning `#FAD14F`, destructive/error `#C94A55`, and light-theme teal text `#006B67`. Bright teal and yellow fills SHALL use black text. Controls SHALL use square corners, one-pixel rules, two-pixel focus, and instant step states without chamfer, offset shadow, blur, or soft shadow. Pixelify Sans SHALL be limited to wordmark and H1/H2; Overpass Mono 400/600 SHALL serve interface text, labels, values, and commands.

#### Scenario: A user operates a Courier control

- **WHEN** a control is rendered, focused, selected, disabled, hovered, or pressed in any theme
- **THEN** it remains contrast-safe, keyboard-operable, square, sharply ruled, and free of soft or offset effects

#### Scenario: A product surface groups related content

- **WHEN** content is grouped without its own interactive state
- **THEN** it uses a flat surface band, typography, alignment, spacing, or a meaningful one-pixel rule instead of a detached card island

### Requirement: Shared workbench primitive

Courier SHALL compose product state from flat terminal panels, horizontal row rules, semantic controls, and truthful product data. Decorative card islands, fake shell execution, fictional monitoring data, soft shadows, blur, and rounded framing SHALL NOT be introduced.

#### Scenario: A product surface presents state

- **WHEN** landing, delivery, or administration groups related information
- **THEN** hierarchy comes from typography, spacing, flat surface bands, and meaningful one-pixel rules while all displayed values originate from the CLI contract or runtime API

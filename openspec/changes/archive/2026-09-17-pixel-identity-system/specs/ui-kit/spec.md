## MODIFIED Requirements

### Requirement: Coherent brand identity system

Courier SHALL use one repository-local modern 8-bit Relay identity across landing, delivery, administration, README, favicon, and product chrome. Relay SHALL be a compact square pigeon courier with a satchel, parcel, folded wings, compact tail, short beak, visible eye, and at most one small earpiece. The transparent ImageGen-authored mark, neutral mascot, route pose, delivery pose, and administration pose SHALL use hard pixel edges and SHALL NOT depict armor, a helmet, visor, glowing face panel, exoskeleton, metallic chest plate, police or military equipment, photorealism, smooth 3D shading, scenery, text, logos, or watermarks.

The neutral sprite SHALL be the canonical identity reference. Every role sprite SHALL be an identity-preserving derivative of that canonical asset and SHALL preserve its body size, body-to-head proportions, physiology, head, eye, beak, folded wings, compact tail, base plumage, earpiece, and satchel. Only pose, role equipment, clothing, and carried or attached objects MAY vary. Independent mascot redraws SHALL NOT be accepted as role variants.

#### Scenario: A browser surface presents Courier

- **WHEN** landing, delivery, or administration renders
- **THEN** its header uses the square pixel Relay mark beside `COURIER CLI`, its optional role sprite matches the surface, and complete product meaning remains available without either image

#### Scenario: Relay changes roles

- **WHEN** the landing, delivery, administration, README, and header assets are compared
- **THEN** they depict the same recognizable mascot proportions and physiology while only the declared pose and role equipment differ

### Requirement: Intentional control appearance

Courier SHALL present semantic browser controls through one shared modern pixel grammar using the declared twelve-color palette, four-pixel geometry unit, crisp dividers, chamfered interactive frames, hard offset state shadows, and visible two-pixel focus. Ordinary layout wrappers SHALL remain transparent and SHALL NOT add card backgrounds, outer radii, blur, soft shadows, or borders. Body content SHALL remain system sans, commands SHALL remain system monospace, and the locally bundled Cyrillic Pixelify Sans font SHALL be limited to brand and display headings.

#### Scenario: A user operates a Courier control

- **WHEN** a control is rendered, focused, selected, disabled, or activated
- **THEN** its pixel state remains contrast-safe, keyboard-operable, semantically native, and free of an unnecessary surrounding island

#### Scenario: A product surface groups related content

- **WHEN** content is grouped for layout without its own interactive state
- **THEN** it remains on the continuous document canvas and uses spacing, typography, alignment, or a meaningful crisp divider instead of a card island

### Requirement: Local official brand marks

Courier SHALL provide typed first-party and third-party pixel icon registries rendered from reviewed sixteen- or twenty-four-pixel grids with crisp edges and current color. Installation, platform, operating-system, GitHub, endpoint, preference, file, transfer, status, and administration icons SHALL be local. Third-party pixel derivatives SHALL retain pinned official source, license, attribution, and trademark records. Locale controls SHALL use local pixel flags rather than platform emoji.

#### Scenario: A browser surface renders an icon

- **WHEN** Courier renders a functional icon, channel mark, or locale flag
- **THEN** it uses the declared local pixel registry with an accessible name where needed and makes no runtime request

#### Scenario: A branded channel is rendered

- **WHEN** a supported package manager, shell, operating system, or distribution mark appears
- **THEN** the UI uses its reviewed pixel-grid derivative, inherits the requested Courier color, retains its pinned source record, and exposes no remote asset URL

### Requirement: Shared Bezier route geometry

Courier SHALL preserve measured cubic route positioning while sampling and snapping its visual points to the four-pixel grid. Decorative route packets and status signals SHALL use short stepped motion only, and reduced motion SHALL render the same state without travel or repeated animation.

#### Scenario: A route is rendered

- **WHEN** finite Source and Destination positions are supplied
- **THEN** a deterministic crisp polyline connects them, its terminals remain fixed-aspect, and any packet animation stops under reduced motion

#### Scenario: A consumer supplies two endpoint positions

- **WHEN** a route path is requested for finite coordinates
- **THEN** the helper returns a stable cubic source path and a deterministic four-pixel-grid sample for the rendered connector

### Requirement: Auditable identity assets

Courier SHALL retain exactly the declared five transparent pixel Relay WebP assets and one pinned local display font with role, dimensions, byte count, generation prompt, lineage or upstream revision, license, and SHA-256 digest. The mark SHALL remain at most 24 KiB, neutral and route sprites at most 64 KiB each, delivery and administration sprites at most 48 KiB each, all Relay assets at most 248 KiB combined, and the display font at most 48 KiB.

#### Scenario: Asset integrity check

- **WHEN** repository verification runs
- **THEN** undeclared assets, legacy armored Relay files, missing provenance, missing licenses, incorrect dimensions, or an exceeded budget fail the gate

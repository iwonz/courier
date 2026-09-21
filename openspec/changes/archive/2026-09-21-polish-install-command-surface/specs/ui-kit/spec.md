## MODIFIED Requirements

### Requirement: Local official brand marks

Courier SHALL provide a typed first-party pixel icon registry for functional endpoint, preference, file, transfer, status, and administration icons. GitHub, installation, platform, operating-system, distribution, shell, and package-manager identities SHALL use typed local transparent PNG marks in official geometry and color without pixel sampling. Every third-party raster SHALL retain pinned official source, revision, color, license, attribution, and trademark records. npx SHALL reuse the official npm mark rather than invent an independent npx brand; channels without an applicable official mark, including wget, SHALL remain text-only. Locale controls SHALL use local first-party pixel flags rather than platform emoji.

#### Scenario: A browser surface renders an icon

- **WHEN** Courier renders a functional icon, channel mark, or locale flag
- **THEN** it uses the declared local first-party registry or official brand PNG with an accessible name where needed and makes no runtime request

#### Scenario: A branded channel is rendered

- **WHEN** a supported package manager, shell, operating system, or distribution mark appears
- **THEN** the UI uses its transparent official-color PNG through an `img`, retains its pinned source record, avoids pixelated rendering, and exposes no remote asset URL

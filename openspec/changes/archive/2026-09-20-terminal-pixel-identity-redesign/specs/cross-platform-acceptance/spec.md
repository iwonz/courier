## MODIFIED Requirements

### Requirement: Pixel browser acceptance

Courier SHALL verify landing, delivery, and administration in English/Russian and system/light/dark at 320×568, 390×844, 1024-wide, and 1440×900. Acceptance SHALL assert exact terminal tokens, local Overpass Mono and Pixelify Sans scope, five local Relay v3 assets and correct per-surface resolution, unmodified local brand PNGs, square controls, one-pixel rules, two-pixel focus, absence of shadows/chamfers/card islands/external requests/overflow, centered route signal and reduced-motion midpoint, truthful contract/API data, and all existing Copy, builder, auth, upload/download, administration mutation, and preference flows. Landing JavaScript SHALL remain at most 150 KiB gzip and CSS at most 9 KiB gzip.

#### Scenario: Browser acceptance runs

- **WHEN** the production surfaces are exercised at every required locale, theme, motion setting, and viewport
- **THEN** visual structure, asset resolution, fonts, colors, overflow, interaction, API behavior, and bundle limits satisfy the terminal identity contract

#### Scenario: Browser surfaces render the pixel identity

- **WHEN** each production surface renders in a supported theme and viewport
- **THEN** the expected v3 assets, fonts, tokens, rules, square controls, truthful values, and no external requests or overflow are observed

#### Scenario: Product workflows are exercised

- **WHEN** Playwright performs route, Copy, builder, auth, upload/download, administration, and preference scenarios
- **THEN** the redesign preserves the exact established functional outcomes in English and Russian

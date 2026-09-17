## ADDED Requirements

### Requirement: Pixel browser acceptance

Courier SHALL test landing, delivery, and administration in Chromium for the complete shared modern 8-bit identity in English and Russian, system/light/dark themes, reduced motion, touch and fine-pointer input, keyboard operation, supported responsive viewports, fixed transparent landing masthead, protected-metadata isolation, API behavior, exact role sprites, local font and icon loading, contrast-safe selected controls, grid-snapped route geometry, stepped motion, no horizontal overflow, no ordinary card islands, and no external runtime requests.

#### Scenario: Browser surfaces render the pixel identity

- **WHEN** landing, delivery, and administration render at supported desktop and mobile sizes
- **THEN** each uses the expected mark, role sprite, pixel font scope, icon registry, local flag, palette, control geometry, and reduced-motion behavior without loading legacy artwork, emoji flags, Lucide icons, generated backgrounds, remote assets, or hidden protected data

#### Scenario: Product workflows are exercised

- **WHEN** acceptance selects routes, copies commands, authenticates, transfers data, receives SSE snapshots, changes policy, or stops a target
- **THEN** the visual redesign preserves all existing contract-backed and guarded behavior without layout instability or external effects

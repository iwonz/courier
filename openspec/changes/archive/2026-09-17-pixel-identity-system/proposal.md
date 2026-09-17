# Change: Introduce the 8-bit Courier identity system

## Why

Courier's current armored 3D pigeon and soft rounded interface read as a robot product rather than a compact courier mascot. The generated character, smooth gradients, Lucide icons, emoji flags, and rounded controls also form several competing visual languages across the landing, delivery, and administration applications.

## What Changes

- Replace the current Relay assets with one square pixel-art courier pigeon mark, a neutral mascot, and route, delivery, and administration role sprites.
- Apply one accessible 12-color pixel grammar, local Cyrillic display font, crisp icons, pixelized third-party marks, stepped route motion, and grid-snapped route geometry across every browser surface.
- Preserve transparent borderless composition, native semantics, theme and locale behavior, API/security boundaries, CLI behavior, and release behavior.
- Add deterministic asset, font, provenance, budget, unit, and real-browser gates for the new system.

## Impact

The change affects browser presentation and repository-owned identity assets only. Public CLI syntax, generated command data, HTTP endpoints, authentication, SSE, transfer behavior, and release artifacts remain compatible.

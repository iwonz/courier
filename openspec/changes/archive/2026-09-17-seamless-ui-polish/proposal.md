# Change: Seamless UI polish

## Why

The first React/shadcn pass preserved too many visible card borders and nested surfaces, which fragments otherwise related workflows. The header also reuses the full mascot at icon size, where its rounded body reads as a perched bird rather than a compact technical mark, and the locale control does not visually expose the active language.

## What Changes

- Introduce a dedicated square ImageGen Relay mark for compact product chrome while retaining the full Relay mascot for editorial use.
- Render `COURIER CLI` beside the header mark and position landing navigation directly after the brand.
- Show the active locale through the locale control icon while retaining localized current/next accessible names.
- Reduce non-functional borders, nested cards, shadows, and section seams across landing, delivery, and administration.
- Merge related landing and administration regions into continuous shared compositions without changing product behavior.

## Impact

This changes browser presentation and local identity assets only. CLI behavior, HTTP APIs, authentication isolation, release behavior, route selection, copying, and administrative mutations remain unchanged.

# Change: Borderless UI flow

## Why

The seamless pass reduced borders but retained large rounded tonal containers around the landing route, CLI registry, delivery manifest, and administration workspace. Those containers still read as disconnected islands instead of one continuous Courier interface.

## What Changes

- Remove large card backgrounds, outer radii, backdrop panels, and nested tonal blocks from landing, delivery, and administration.
- Make shared cards and command/route presentations structural and transparent by default.
- Preserve hierarchy through spacing, typography, responsive grids, and only meaningful row or column dividers.
- Keep visible surfaces for semantic controls, form inputs, alerts, selection, focus, status, and destructive actions.
- Compact adjacent landing spacing, keep selected endpoints contrast-safe, and use a fixed backgroundless masthead.

## Impact

This is a presentation-only change. CLI behavior, generated contract data, HTTP APIs, authentication isolation, preferences, copy actions, and release behavior remain unchanged.

# Change: Terminal 8-bit UI and identity redesign

## Why

Courier's browser surfaces use the correct product structure but still read as a conventional card-based blue interface. The landing, delivery, and administration applications need one recognizable terminal/data-grid language, a restrained black/white/teal/yellow palette, and one consistently recolored Relay family without changing runtime behavior.

## What Changes

- Replace the twelve-color presentation system with accessible light/dark terminal tokens, square one-pixel rules, step states, and local Overpass Mono UI typography.
- Recompose landing, delivery, and administration as route, transfer, and operations consoles using only contract or runtime data.
- Recolor the five Relay assets as identity-preserving v3 derivatives, promote neutral-v3 to identity authority, remove legacy shipped sprites, and update asset QA.
- Extend exact unit, asset, browser, documentation, and bundle acceptance for the redesigned surfaces.

## Impact

This changes presentation tokens, asset resolution, and browser layout only. CLI contracts, API contracts, preferences, localization behavior, command construction, authentication, upload/download, and administration mutations remain compatible. Release versioning, commit, and deployment are outside this change.

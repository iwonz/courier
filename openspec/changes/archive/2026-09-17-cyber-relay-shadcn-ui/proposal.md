# Change: Technological Relay shadcn UI

## Why

The current swift identity feels generic and the browser surfaces still rely on a bespoke Lit component system. Courier needs a recognizable technological pigeon mascot and one established, accessible component grammar across the landing, protected delivery UI, and administration UI.

## What Changes

- Replace the swift with Relay, an ImageGen-authored technological pigeon mascot with a stable transparent reference and no illustrated page background.
- Migrate all browser applications and the shared UI package from Lit custom elements to React and repository-owned shadcn components built on Radix primitives.
- Introduce shared shadcn Button, Card, Badge, Checkbox, Input, Select, Tabs, Tooltip, Progress, Separator, ScrollArea, and Alert components with one calm technological Courier theme.
- Rebuild the landing as three natural-height React sections while preserving the route, installation, command, copy, theme, locale, accessibility, and external-link behavior.
- Rebuild delivery and administration with the same components while preserving API contracts, authentication isolation, uploads, downloads, SSE updates, optimistic policy versions, and stop scoping.
- Keep every font, icon, image, script, and style local at runtime and retain exact coverage, asset, bundle, embed, and browser acceptance gates.

## Impact

This changes browser implementation, presentation, dependencies, and generated identity assets. Public CLI behavior, the machine-readable command contract, HTTP APIs, release behavior, and security boundaries remain unchanged.

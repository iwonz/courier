## Why

Courier already exposes structured parameter metadata and stable runtime output, but required inputs are not visually prioritized and the interactive terminal experience does not distinguish readiness, progress, success, interruption, and failure. Server inventory and self-update activity are especially difficult to scan during an interactive session.

## What Changes

- Sort arguments and options required-first in every generated presentation while preserving contract order inside each group.
- Replace visible Required/Optional landing labels with an accessible red required marker and split parameter metadata into distinct compact surfaces and documentation columns.
- Add decorative POSIX and PowerShell icons to the shell selector without changing its accessible text or command generation.
- Add one terminal-capability renderer for transfer, hosted delivery, administration UI, server control, update, and version output while preserving deterministic plain output.
- Render responsive nested server and delivery tables, hosted readiness and elapsed state, neutral user interruption, and staged self-update progress.
- Add canonical release URLs to update and version output without terminal hyperlink escape sequences.

## Capabilities

### Modified Capabilities

- `cli-contract`: Defines stable required-first ordering and separate generated metadata columns.
- `project-landing`: Defines required markers, separate metadata surfaces, and shell icons.
- `operation-reporting`: Defines TTY detection, semantic rich output, sanitation, and neutral interruption.
- `server-control`: Defines responsive server and delivery tables.
- `web-deliveries`: Defines browser-delivery readiness and foreground status presentation.
- `webhook-deliveries`: Defines incoming-webhook readiness and foreground status presentation.
- `self-update`: Defines typed update stages, byte progress, and canonical release links.

## Impact

The change affects generated README/reference/help ordering, landing presentation and tests, runtime reporting, hosted lifecycle output, server inventory, update events, version output, and Go/web verification. Contract schema `3`, contract version `0.16.0`, `target_release`, exit codes, plain-output field names, IPC, registry, worker-definition, and administration-state versions remain unchanged.

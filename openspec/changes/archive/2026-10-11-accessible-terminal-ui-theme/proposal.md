# Change: Accessible theme-neutral terminal UI

## Why

Courier's interactive panels force a dark surface and white foreground while nested Lip Gloss styles emit independent resets. On a light terminal those resets expose the terminal's white background inside the dark panel, producing unreadable white bands and inconsistent spacing. The same fixed palette is shared by update, version, hosted delivery, UI, server, success, interruption, and failure reporting.

## What Changes

- Replace filled, terminal-width panels with a compact utility layout that inherits the terminal's foreground and background.
- Keep semantic color as a supplemental marker selected from the terminal's ANSI palette; retain visible text and symbols when color is disabled.
- Make multiline fields, release notes, long URLs, paths, and narrow terminals wrap with stable indentation without padded color surfaces.
- Remove filled table rows and headers while preserving bold headings, borders, responsive wrapping, and complete values.
- Preserve QR contrast, TTY detection, `NO_COLOR`, `TERM=dumb`, sanitization, redirected output, exit codes, and command behavior.

## Impact

- Affected spec: `operation-reporting`
- Affected code: shared `internal/terminalui` renderer and its runtime presentation tests
- No CLI, IPC, registry, worker, administration-state, or contract version changes

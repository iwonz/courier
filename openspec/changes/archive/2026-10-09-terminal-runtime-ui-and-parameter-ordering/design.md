## Context

Courier writes progress on stderr and final results on stdout. Existing non-interactive output is used by scripts and tests, while real terminals can support richer hierarchy and status. The landing and generated documentation already receive parameter metadata from the CLI contract but present it differently.

## Goals / Non-Goals

**Goals:**

- Keep redirected output stable, cursor-free, and free of ANSI sequences.
- Give interactive runtime commands one square, dark, semantic Courier presentation.
- Preserve full UUIDs, paths, and policy values when a terminal is narrow.
- Make requiredness and option relationships consistent across landing, help, README, and reference generation.

**Non-Goals:**

- Change structured help into decorative panels.
- Add a color flag, detach alias, OSC-8 hyperlink, protocol version, or persisted-state migration.
- Replace mpb transfer progress or change transfer accounting.

## Decisions

### Output capability is derived per stream

The common renderer detects a real terminal from the output file descriptor and disables color when `NO_COLOR` is present or `TERM=dumb`. Explicit renderer modes make tests deterministic. Dynamic values are stripped of ANSI and unsafe control characters before layout.

### Plain behavior remains the compatibility surface

Every command retains its existing redirected fields, counters, ordering, and exit code. Version and update append one final `release: <URL>` line. Rich rendering is selected only for a real TTY, so pipes and snapshot automation do not receive borders, cursor motion, or color.

### Server inventory chooses layout before rendering

Wide terminals receive one server overview followed by one delivery table per server. Narrow terminals receive two-column field/value tables, which preserve complete UUID, endpoint, counter, policy, and timestamp values rather than truncating identifiers.

### Updater emits typed, optional events

`Updater.Run` remains the compatibility entry point and delegates to `RunWithProgress` with no sink. The progress-capable path emits release check, archive download, checksum download, verification, extraction, installation, and completion events. Downloads report current bytes even when content length is unknown.

### Contract ordering is stable

Generators copy parameters into two stable groups: required first and optional second. Relative contract order is retained inside each group. The landing uses an `aria-hidden` red asterisk plus a localized screen-reader label, while documentation retains explicit status and gives every relationship its own column.

## Risks / Trade-offs

- **[Rich tables can exceed small terminals]** → Select field/value layout below the safe multi-column width and allow values to wrap without dropping content.
- **[Terminal data may contain escape sequences]** → Sanitize all dynamic labels and values before styling.
- **[Progress redraw can leave a partial line]** → Clear only the live status line on completion, interruption, or failure.
- **[A new dependency affects release builds]** → Use Lip Gloss only in the shared renderer and cover cross-platform builds in the existing verification gate.

## Migration Plan

Ship the renderer, generated documentation, and landing presentation together. No persisted data or protocol migration is required. A rollback restores the prior presentation; redirected command output remains compatible except for the intentional final release URL line on `version` and `update`.

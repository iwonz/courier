## Context

The data application currently uses `fetch` to send one `FormData` request and stores one global upload phase. Fetch does not expose broadly supported request-body progress, while the server correctly enforces exactly one `file` multipart part per transaction. The shared UI kit already provides a bounded accessible progress component.

## Goals / Non-Goals

**Goals:**

- Show the complete current selection before its first request completes.
- Represent every selected file independently through queued, active, and terminal states.
- Drive the active visualization from actual browser upload progress.
- Preserve one-file transactional server semantics and guarded errors.

**Non-Goals:**

- Add parallel uploads, drag-and-drop, directory selection, pause, cancel, reordering, or automatic retry.
- Change HTTP routes, multipart shapes, server limits, collision policy, or extraction behavior.
- Persist upload history across navigation or reload.

## Decisions

### Use a sequential client queue

Convert the selected `FileList` to a stable ordered batch, render all rows, then await one upload before starting the next. This keeps the existing endpoint and policy reservation model unchanged, avoids a batch monopolizing concurrency, and gives every failure an isolated result. A failed request updates only its row and does not short-circuit the loop.

### Use XMLHttpRequest for observable request progress

Move the upload adapter from fetch to `XMLHttpRequest`, which exposes request-body progress without inventing progress from timers. Normalize the multipart progress ratio to `0..100`, keep the exact file size as separate metadata, and set 100 percent only after a successful response. Reuse the existing bounded JSON error classification for non-success responses; network, abort, and timeout events receive stable status-free errors.

### Keep one current batch in React state

Each row receives a batch-local stable identifier so duplicate file names remain distinct. A state update helper changes only the matching row. The chooser is disabled while any row is queued or uploading, its value is cleared after the files are captured, and a later selection replaces the completed list rather than growing an unbounded history.

### Make progress and results accessible

Render a localized batch summary as a polite live status, use a semantic list, label each progress bar with its file name, and expose text status and percentage alongside color/icon changes. Narrow layouts stack metadata without horizontal overflow.

## Risks / Trade-offs

- **[Multipart progress includes framing overhead]** -> Display normalized percentage rather than claiming the event's raw byte count equals payload bytes; show the exact file size separately.
- **[A large batch is slower when sequential]** -> Sequential delivery preserves policy and transactional behavior; parallelism would require a separate resource-control design.
- **[Component unmount leaves an in-flight browser request]** -> Browser navigation owns cancellation as before; no resumability is promised.

## Migration Plan

Rebuild the embedded data assets and ship as a patch-compatible UI enhancement. The unchanged endpoint remains compatible with current workers and clients, and rollback requires no state cleanup.

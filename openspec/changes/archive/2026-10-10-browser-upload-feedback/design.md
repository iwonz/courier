## Context

`DataContent.sendFile` currently awaits the multipart request, clears the input only on success, and refreshes metadata. A `web-to-path` metadata refresh always returns the same upload manifest, so successful completion has no visible effect. Its catch path only sets the page-level `failed` flag, while the corresponding alert is deliberately rendered only when metadata is absent; upload errors are therefore silent. The HTTP handler itself remains transactional and the compiled runtime already commits successful browser and webhook payloads.

## Goals / Non-Goals

**Goals:**

- Make pending, committed, and rejected browser-upload outcomes explicit and accessible.
- Allow a failed file to be selected again immediately.
- Preserve the server as the authority for commit, collision, policy, and extraction outcomes.
- Exercise the real compiled binary for both browser and incoming-webhook outcomes.

**Non-Goals:**

- Add overwrite, rename, multi-file selection, drag-and-drop, resumable uploads, or byte-level progress.
- Change webhook HTML behavior or add a webhook form.
- Change server error disclosure, multipart framing, policy, routes, protocols, or persisted state.

## Decisions

### Upload state is independent of metadata state

The data application will keep a small upload state machine (`idle`, `pending`, `success`, `failure`) alongside metadata/authentication state. Pending disables the chooser and exposes a polite live status. Success is rendered only after the upload request returns a successful JSON response. Failure uses an assertive alert and does not discard already-authorized metadata.

### Guarded API errors retain status and safe server reason

The data API will throw a typed request error for non-2xx responses. It records the status and, when the JSON body has the existing `{error: string}` shape, the server-provided guarded reason. Malformed, non-JSON, or unbounded values fall back to the status-only message. The UI renders the typed message as text, never markup.

### Every attempt releases the file input

The input value is cleared in `finally`, whether the request succeeds or fails. Pending state prevents overlapping selections, and clearing permits browsers to emit a new change event when the same path is selected again.

### Runtime acceptance verifies observation and storage

The compiled-runtime browser test will assert pending/success UI, exact destination bytes, a visible collision failure, and a same-file retry request. Incoming-webhook acceptance will continue using real multipart requests and additionally assert the guarded collision body, so both hosted ingress paths are covered without claiming that webhooks have a UI.

## Risks / Trade-offs

- **[Very fast uploads may make pending UI brief]** -> Unit tests hold the request open to verify pending semantics; runtime acceptance treats success as the durable state.
- **[Server reasons are English]** -> The surrounding copy is localized and reasons are the deliberately safe existing API strings; no filesystem detail is exposed.
- **[No numeric progress]** -> `fetch` does not expose upload progress portably, so Courier reports an honest indeterminate pending state rather than fabricated percentages.

## Migration Plan

Ship as the next patch release. The embedded asset is replaced with the binary; existing workers and routes require no migration or cleanup.

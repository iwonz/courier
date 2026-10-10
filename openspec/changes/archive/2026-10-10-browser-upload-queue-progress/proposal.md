## Why

The browser upload workspace accepts only one selected file and reduces the whole transfer to a binary pending/success/failure alert. Users cannot prepare a batch or see which file is queued, how far an active upload has progressed, or which files in a selection completed independently.

## What Changes

- Allow one or multiple files to be selected from the browser upload control.
- Present the selected batch as a responsive file list with name, exact size, individual progress, percentage, and queued/uploading/delivered/failed state.
- Send the batch as sequential single-file transactions so the existing one-file multipart contract, policy reservations, collision protection, and delivery accounting remain authoritative.
- Continue the queue after an individual failure and retain every result until the next selection, while permitting the same files to be selected again after completion.
- Use real browser upload progress, sanitized server failures, localized English/Russian labels, and accessible live/progress semantics.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `web-deliveries`: Extend the observable browser-upload lifecycle from one binary state to a multi-file queue with per-file progress and independent terminal results.

## Impact

This change affects the embedded data React application, its browser API adapter, localization, unit tests, embedded build assets, and compiled-runtime Playwright acceptance. It does not change HTTP endpoints, multipart field names, delivery policy, worker/IPC/state versions, or CLI behavior.

## Why

The embedded browser-upload page performs the real multipart transfer, but it renders no pending or successful state afterward. Upload failures are also invisible once metadata is present, and the file input is left populated after a failure, so choosing the same file again may not emit another change event. A successful transfer and a rejected collision therefore both look like “nothing happened.”

## What Changes

- Show an accessible pending state while the selected file is being uploaded.
- Show an explicit confirmation only after Courier returns a committed-upload response.
- Show a safe, actionable error containing the HTTP status and guarded server reason when an upload is rejected.
- Reset the file input after every completed attempt so the same file can be retried.
- Disable concurrent selection while one upload is active.
- Extend unit and compiled-runtime browser coverage for success, collision feedback, retryability, and committed bytes.
- Retain and strengthen real multipart webhook acceptance so the corresponding non-UI route is proven to commit bytes and return explicit success/collision responses.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `web-deliveries`: Makes the embedded browser-upload lifecycle observable and retryable without changing upload transaction semantics.
- `webhook-deliveries`: Strengthens compiled incoming-webhook acceptance for explicit committed-success and collision responses.

## Impact

This change affects the embedded data application, its guarded API client, localized copy, component/API tests, and serial compiled-runtime browser/webhook acceptance. It does not change routes, multipart formats, authentication, overwrite policy, IPC/state versions, CLI schema, or terminal output.

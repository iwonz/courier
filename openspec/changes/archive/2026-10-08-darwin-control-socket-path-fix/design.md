## Context

Both the administration process and data workers call `ipc.ControlEndpoint`. On macOS, `os.UserConfigDir` places state under `~/Library/Application Support/courier/state`; combining that directory with `control-`, a canonical 36-character UUID, and `.sock` produces the reported 104-byte path. Darwin accepts at most 103 pathname bytes for a Unix-domain socket.

## Decisions

- Keep the control socket in the existing private state directory so the local authentication and cleanup boundary does not move.
- Remove the `control-` prefix, UUID hyphens, and `.sock` suffix for newly generated Unix endpoints. The resulting 32 lowercase hexadecimal characters retain every UUID bit, are deterministic, and remain safe on case-insensitive filesystems.
- Leave Windows pipe names unchanged.
- Continue treating persisted control endpoints as authoritative. Do not rewrite registry or administration state and do not bump any schema or protocol version.
- Exercise a real 103-byte listener and dial on Darwin in addition to asserting the reported default path is shortened from 104 to 87 bytes.

## Risks / Trade-offs

- Compact names are less self-describing in directory listings. The surrounding private state directory and full lossless UUID representation retain sufficient ownership and diagnostic identity.
- Arbitrary caller-supplied state directories can still exceed an operating system's socket limit. Courier has no public custom state-directory option, and this fix intentionally avoids moving the security boundary or adding a temporary-directory fallback.

## Migration

No persisted-state migration is required. New processes generate compact endpoints; clients continue dialing stored legacy endpoints for live older processes and stale reconciliation removes the exact recorded path. Publish v0.3.6 after strict OpenSpec and the complete release gate pass.

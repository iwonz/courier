# Operational reporting

Courier uses one reporting model for filesystem copies, archives, extraction, HTTP deliveries, background workers, and control commands.

## Stages and counters

The public stage vocabulary is `preflight`, `archive`, `extract`, `transfer`, `commit`, `verify`, `cleanup`, and `complete`. Progress events always identify a stage and carry:

- `read`: payload bytes consumed from the source;
- `sent`: payload bytes submitted to the next transport or target boundary;
- `confirmed`: payload bytes accepted by the target boundary;
- `total`, elapsed time, and speed calculated from confirmed bytes.

The invariant is `0 <= confirmed <= sent <= read`. A local committed write advances all three counters together. An outgoing webhook can report read and sent bytes while confirmed remains zero until a successful response is received. Courier never labels an uncertain remote outcome as success.

Interactive terminals receive an `mpb` progress bar. Redirected stderr receives deterministic cursor-free lines with stable named fields. The final success summary includes source, actual destination, confirmed bytes, elapsed time, and result. Failure output includes the failing stage, all known counters, a sanitized reason, and result.

## Diagnostic safety

All final reports and internal-mode failures pass through the same sanitizer. URL user-info, query strings, fragments, and conventional password, token, secret, authorization, and cookie assignments are removed before output. Runtime credentials and sessions remain outside registry and history data.

## Private history

Operational history files use the same private immutable write-and-sync protocol as the delivery registry. Optional diagnostic text is sanitized before validation and persistence. Events are read in timestamp-plus-UUID order, and the store retains the newest 256 events by default. Rotation considers only Courier history filenames, verifies that every removal target is a private regular file, removes oldest entries, and synchronizes the state directory.

## Exit codes

Courier returns `0` for success, `2` for CLI errors, `10` for SSH connection/trust/authentication errors, `20` for transfer-stage errors, `30` for self-update errors, `40` for registry/IPC/control errors, and `130` when an operation or internal runtime is canceled or interrupted.

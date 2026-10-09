## Context

Hosted foreground commands currently own a lease but wait only for process context cancellation. A second Courier process can stop the delivery and close the worker while the initiating process continues waiting. The worker already supports progress subscriptions, but subscriber filters are not retained for later publications and final worker shutdown closes subscribers before the terminal state is observable. Browser acceptance currently targets Vite/mock surfaces rather than the compiled embedded runtime. Shipping also commits a dirty `main` directly, leaving no enforceable start-to-finish change boundary.

## Goals / Non-Goals

**Goals:**

- Preserve stable plain output and protocol versions while making TTY readiness scannable from another device where the printed URL is actually reachable.
- Distinguish confirmed operator stop, user interruption, and unexplained worker loss without lease races.
- Exercise every hosted direction through the real compiled binary and embedded browser applications.
- Make the documented OpenSpec, verification, merge, release, and Pages lifecycle executable and fail closed.

**Non-Goals:**

- Expose loopback services to a LAN, add a public tunnel, or add an HTML page for incoming webhooks.
- Add `to webhook://`, change outgoing webhook syntax, or add a QR/plain-output flag.
- Bump any persisted or IPC version, require pull requests, or remove explicit recovery operations.

## Decisions

### QR rendering is a terminal-only projection

The shared terminal renderer encodes only a URL already validated and sanitized by the caller. `rsc.io/qr v0.2.0` supplies the module bitmap. Courier adds a four-module quiet zone and combines two module rows with Unicode half blocks. ANSI black/white backgrounds are used only in a color TTY; `NO_COLOR` receives an ANSI-free half-block form. Redirects, `TERM=dumb`, and widths below the complete QR width omit the QR. The original text URL always remains.

### Foreground lifetime follows a delivery-filtered progress subscription

Worker subscribers retain their requested delivery filter. Registration returns a foreground watcher alongside the lease. Stop operations publish a terminal snapshot or tombstone before closing subscribers. The initiating command treats an observed removal following a confirmed stop as success, context cancellation as exit 130 after lease release, and subscription loss without terminal evidence as control failure 40. Closing a watcher and lease is idempotent and ordered so `Release` cannot race a remote stop into a false failure.

### Administration distinguishes parent cancellation from control shutdown

The administration process uses its local cancellation only for verified `ui stop`. Parent context cancellation remains observable separately. Foreground startup prints a neutral stopped result only for a previously started local process that shuts down through control; Ctrl+C remains cancellation.

### Usage failures are typed before reporting

Cobra parse, unknown-command, argument-count, unknown-flag, and unsupported-route errors are wrapped as a usage category with exit code 2. A compact renderer prints sanitized error text, an available Cobra suggestion, and the relevant help command. Operational preflight and runtime failures continue through the counter-bearing report path.

### Runtime E2E owns isolated processes and ports

The browser test launcher builds one real Courier binary into a temporary root and exports its path. A serial Playwright suite creates a separate HOME and deterministic free ports per scenario, starts foreground commands with piped readiness output, drives the embedded application or real multipart HTTP, verifies exact bytes and extracted trees, stops the target from a second Courier process, and asserts exit code zero plus empty server state. Cleanup always attempts scoped Courier shutdown before terminating leftover child processes and deleting the fixture.

### Every tracked mutation is one feature-branch release unit

`make change-start TYPE=<type> CHANGE=<name>` fetches, requires a clean `main` equal to `origin/main`, rejects another active change, creates `TYPE/CHANGE`, and scaffolds OpenSpec. `make ship` runs only from that feature branch, validates the one completed change, derives the next patch unless overridden, archives, generates, and verifies before creating one conventional commit. It records the original main commit, switches locally to unchanged `main`, performs `--ff-only`, pushes under a narrow guard token, tags, waits for all publication surfaces and Pages, then deletes the local feature branch and proves clean `main == origin/main`.

Direct main and release-tag pushes are rejected by the installed pre-push hook unless a ship or explicitly documented recovery token is present. `make release` and `make pages-publish` require the recovery token. Pages loses its `push` trigger and runs only after the release matrix or an explicit recovery dispatch.

## Risks / Trade-offs

- **[Subscription closure can be ambiguous]** → Publish the final filtered snapshot first and classify EOF without it as worker loss.
- **[QR blocks can wrap]** → Compute the rendered cell width and omit the QR when the terminal is narrower.
- **[Real E2E can leak processes]** → Use isolated HOME/ports, serial execution, scoped stop commands, bounded waits, and a final process cleanup registry.
- **[A release can fail after main/tag publication]** → Fail at the publication boundary, preserve immutable evidence, and document the narrowly authorized recovery command rather than rewriting history.
- **[The Formula job advances main after tagging]** → Shipping fast-forwards to the verified Formula commit before Pages and final clean-state assertions.

## Migration Plan

Land the renderer, lifecycle, E2E, hooks, scripts, workflows, and documentation in one feature branch. The first use of the new ship workflow releases the change as the next patch. Existing clones run `make hooks` to install the updated pre-push guard. No runtime state migration is required.

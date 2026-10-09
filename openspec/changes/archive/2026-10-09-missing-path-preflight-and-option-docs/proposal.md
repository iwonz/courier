## Why

Missing directory endpoints currently fail late, and hosted deliveries can surface a redacted worker-internal error after a worker has already been registered. At the same time, Courier's existing detached mode and option applicability are difficult to discover because terminal help, README, and the landing page do not share structured parameter documentation.

## What Changes

- Add a shared preflight that classifies local and SSH path roles, reports every missing directory before mutation, confirms creation interactively, or recursively creates it with `--force-source-creation`.
- Preserve exact-file destination behavior and reject missing paths in file-only roles without creating them as directories.
- Run hosted-path preflight before worker acquisition and return stable, path-specific errors if the path becomes invalid during registration.
- Keep `--background` as the only detached-mode flag and document its lifecycle consistently for browser transfers, incoming webhooks, and the administration UI.
- Extend the CLI contract with localized descriptions and explicit requiredness, then generate structured Cobra help, the CLI reference, landing data, and managed README parameter tables from that contract.
- Present localized, human-readable option metadata and applicability on the landing page.

## Capabilities

### New Capabilities

- `path-preflight`: Defines missing-path classification, confirmation, ordered creation, permissions, SSH handling, and safe preflight failures.

### Modified Capabilities

- `cli-contract`: Adds localized parameter descriptions, requiredness, the creation flag, structured help metadata, and managed README generation.
- `endpoint-model`: Clarifies directory-role inference while preserving exact-file destinations and file-only endpoint behavior.
- `web-deliveries`: Requires hosted path validation before acquisition and documents foreground and detached browser-delivery lifecycles.
- `webhook-deliveries`: Requires incoming and outgoing webhook paths to follow the shared preflight roles and detached lifecycle.
- `worker-lifecycle`: Requires hosted preflight before registration and safe path-specific registration failures.
- `project-landing`: Adds localized structured parameter rows and human-readable applicability labels.

## Impact

The change affects transfer planning and execution, local and SSH filesystem access, hosted worker preparation and error mapping, Cobra help rendering, the versioned CLI contract and its generators, README and reference documentation, landing-page data and presentation, and their Go, browser, and acceptance tests. IPC, registry, worker-definition, and administration-state versions remain unchanged.

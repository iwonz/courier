## Context

See `proposal.md` for motivation. Path behavior is currently distributed among endpoint parsing, rooted local and SFTP resources, finite transfer execution, and worker registration. Hosted SSH preparation already authenticates before acquisition but does not inspect the endpoint. The CLI contract generates documentation and landing data, while Cobra help and README parameter descriptions are maintained separately. Existing persisted worker definitions and registry records must remain readable without a version bump.

## Goals / Non-Goals

**Goals:**

- Make path-role classification and missing-directory authorization identical for finite and hosted local or SSH endpoints.
- Ensure confirmation is atomic with respect to user intent: inspect and authorize the full set before the first mutation.
- Preserve rooted filesystem boundaries, ephemeral SSH credentials, exact-file destinations, and redaction of unexpected failures.
- Make one validated contract drive every public parameter description.

**Non-Goals:**

- Add a detach alias, TTL, service-manager integration, overwrite mode, temporary-directory fallback, or automatic creation of source files.
- Change persisted worker definitions, IPC messages, registry formats, administration state, or existing collision semantics.
- Roll back a directory after its successful creation.

## Decisions

### Preflight uses explicit endpoint roles

The application layer will derive source and destination roles from the selected route and options, then run one shared preflight over opened local or SSH resources. Roles distinguish directory-required, file-required, and flexible source paths; the latter resolves to the existing type or to a directory when absent. This keeps route policy outside filesystem backends while allowing both backends to share inspection and creation behavior.

Inferring solely from filename extensions was rejected because extensions are neither reliable nor consistent across local and SSH paths. Treating every missing destination as a directory was rejected because it would break exact-file renaming.

### Authorization precedes mutation

Preflight first inspects and records all missing directories in source-then-destination order. Without the force flag it obtains every confirmation before calling recursive creation. It then creates and revalidates in the same order with `0700` permissions. This prevents a declined second prompt from leaving a surprising first directory, while intentionally avoiding unsafe rollback after a real creation failure.

### Hosted preparation occurs before worker acquisition

The hosted runner will prepare its endpoint before building or acquiring a worker. Remote preparation will reuse the existing SSH connection, trust, authentication, and helper workflow, capture only the current bounded runtime credentials for the worker, and close the preflight connection. The worker still reopens and revalidates the endpoint at registration, so no persistent format needs to change.

Adding creation instructions to worker definitions was rejected because it would move user consent behind acquisition, complicate compatibility, and require persisted-schema changes.

### Expected path races receive typed safe errors

Missing and wrong-type endpoint failures will use a small typed preflight classification with sanitized role and display path. The CLI maps those errors to stage `preflight` and exit code `20`. Worker registration maps only these expected conditions through IPC; arbitrary filesystem details remain internal and redacted.

### Contract schema version 3 owns public prose

Arguments and flags will carry `{en, ru}` descriptions and explicit requiredness. Validation requires both locales and optional flags. The generator will project a compact Go metadata table for Cobra, Markdown reference content, landing JSON, and README sections between stable markers. Freshness tests compare all generated outputs.

Embedding the YAML contract in the binary was rejected because it would increase runtime packaging complexity. Hand-maintained help text was rejected because it would preserve the current drift risk.

### Structured help remains native to Cobra

A custom help renderer will consume generated metadata but retain Cobra's command discovery and standard error behavior. It will group arguments and options and render status, descriptions, applicability, defaults, repeatability, dependencies, and conflicts in a terminal-readable layout.

## Risks / Trade-offs

- **[Remote preflight adds an SSH round trip before hosted startup]** → Reuse the existing authentication preparation connection and keep worker credentials bounded exactly as today.
- **[A path can change after successful preflight]** → Revalidate during creation and worker registration, and return safe expected-race errors.
- **[Recursive creation mode is affected by process umask]** → Request `0700` and verify tests never observe permissions broader than owner access on supported systems.
- **[Generated documentation expands the contract surface]** → Enforce schema validation, Cobra parity, deterministic generation, and freshness checks in `make verify`.
- **[Large localized option rows can crowd small screens]** → Use responsive stacked metadata and verify keyboard, accessibility, and overflow behavior in browser tests.

## Migration Plan

Ship the schema and generator updates together with regenerated help metadata, CLI reference, README sections, and landing JSON. Existing registry, administration, IPC, and worker-definition data require no migration. A rollback restores the previous binary and generated presentation files; directories already created by an authorized invocation remain user-owned and are not removed.

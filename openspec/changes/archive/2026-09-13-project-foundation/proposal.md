## Why

Courier starts as a greenfield project, so it needs a reproducible foundation, a single CLI contract, and an explicit delivery plan. This enables independently verifiable transfer changes without coupling the public interface to a specific transport.

## What Changes

- Create the Go module and minimal `courier` executable.
- Introduce a command registry that accepts new handlers without parser changes.
- Add `help` and `version`, build metadata, and baseline tests.
- Record the architecture plan, branch policy, and conventional commits.
- License the project under MIT.

## Capabilities

### New Capabilities

- `cli-foundation`: CLI startup, help, version output, and extensible command registration.

### Modified Capabilities

None.

## Impact

Adds the root Go module files, `cmd/courier`, internal `cli` and `buildinfo` packages, documentation, and license. This change adds no external runtime dependencies.

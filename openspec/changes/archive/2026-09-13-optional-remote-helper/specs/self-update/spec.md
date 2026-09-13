## ADDED Requirements

### Requirement: Reliable Windows self-update handoff

Courier SHALL replace a running Windows executable through a verified staged handoff and SHALL remove handoff artifacts after replacement.

#### Scenario: Windows update is available

- **WHEN** a verified newer Windows binary has been staged
- **THEN** Courier starts the staged binary in hidden handoff mode, which waits for the original process to exit before replacing the target

#### Scenario: Replacement completes

- **WHEN** the handoff process replaces the target
- **THEN** it starts the installed binary in hidden cleanup mode, and cleanup removes the handoff executable after its process exits

#### Scenario: Unsafe internal arguments

- **WHEN** an internal handoff or cleanup invocation references an invalid PID, a non-Courier staging name, or a staging file outside the target directory
- **THEN** Courier rejects the operation without replacing or deleting the target

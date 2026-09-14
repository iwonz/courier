## ADDED Requirements

### Requirement: Coherent brand identity system

Courier SHALL define and apply one mature, utilitarian identity across browser surfaces, documentation, and retained assets using shared semantic tokens and components rather than consumer-specific brand implementations.

#### Scenario: A browser surface presents Courier

- **WHEN** the landing, delivery, or administration application renders
- **THEN** it uses the same positioning, wordmark, semantic color roles, typography, focus treatment, and operational visual grammar from the shared UI package

### Requirement: Non-essential mascot guidance

Courier SHALL use the original Relay courier-pigeon character only as a restrained orientation and state-communication aid, with locally verified provenance and without making meaning depend on the image.

#### Scenario: Mascot artwork is unavailable

- **WHEN** Relay cannot be loaded or is hidden from assistive technology
- **THEN** headings, status text, controls, and progress information still communicate the complete workflow

### Requirement: Calm operational communication

Courier SHALL use concise, concrete, non-alarmist language that identifies actions and verified states without unsupported guarantees, secret disclosure, or humor during failures.

#### Scenario: A delivery changes state

- **WHEN** a user sees preparation, transfer, verification, completion, or failure feedback
- **THEN** the message names the current or stopped operation and preserves an actionable, technically accurate tone in every supported locale

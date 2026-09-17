## MODIFIED Requirements

### Requirement: UI and JSON surfaces

Courier SHALL serve the embedded React data application by default as a shared shadcn destination workspace with typed localization, cyclic theme and locale controls, protected authentication, verified route presentation, breadcrumbs, manifests, and semantic upload/download actions. The workspace SHALL use a continuous canvas and reserve borders for form inputs, focus, and meaningful manifest boundaries rather than nesting bordered cards. Unauthenticated rendering SHALL contain only the real authentication form and generic state, and protected metadata SHALL appear only after authorization. Optional local Relay artwork SHALL remain decorative and SHALL reveal no delivery state. The application SHALL use only the existing guarded delivery API and SHALL expose only versioned JSON/download surfaces when `--no-ui` is active.

#### Scenario: Authentication is required

- **WHEN** an unauthenticated visitor opens a protected delivery
- **THEN** a lightweight authentication region and real password form render without protected names, paths, sizes, routes, credentials, command input, or terminal chrome

#### Scenario: An authorized delivery opens

- **WHEN** authentication and metadata admission succeed
- **THEN** lightweight shadcn controls and one continuous manifest surface present the verified route, navigation, entries, uploads, and downloads without changing API semantics

#### Scenario: An authorized browser delivery opens

- **WHEN** authentication and metadata admission succeed
- **THEN** the destination workspace presents the verified route, navigation, entries, and permitted actions without arbitrary command execution

#### Scenario: UI-disabled root is requested

- **WHEN** a browser opens a no-UI delivery root
- **THEN** Courier returns a JSON description without embedding HTML or protected metadata

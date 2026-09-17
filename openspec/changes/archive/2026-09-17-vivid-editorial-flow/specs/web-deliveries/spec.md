## MODIFIED Requirements

### Requirement: UI and JSON surfaces

Courier SHALL serve the embedded Lit data application by default as a shared destination workbench with typed localization, cyclic theme and locale controls, protected authentication, verified route presentation, breadcrumbs, manifests, and semantic upload/download actions. Unauthenticated rendering SHALL contain only the real authentication form and generic state, and protected metadata SHALL appear only after authorization. Optional local swift artwork SHALL remain decorative and SHALL reveal no delivery state. The application SHALL use only the existing guarded delivery API and SHALL expose only versioned JSON/download surfaces when `--no-ui` is active.

#### Scenario: Authentication is required

- **WHEN** an unauthenticated visitor opens a protected delivery
- **THEN** the vivid destination workbench and generic form render without protected names, paths, sizes, routes, credentials, terminal chrome, or command input

#### Scenario: An authorized delivery opens

- **WHEN** authentication and metadata admission succeed
- **THEN** the workbench presents the verified route, navigation, entries, and permitted API-backed actions

#### Scenario: An authorized browser delivery opens

- **WHEN** authentication and metadata admission succeed
- **THEN** the destination workbench presents the verified route, navigation, entries, and permitted actions without arbitrary command execution

#### Scenario: UI-disabled root is requested

- **WHEN** a browser opens a no-UI delivery root
- **THEN** Courier returns a JSON description without embedding HTML or protected metadata

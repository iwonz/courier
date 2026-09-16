## MODIFIED Requirements

### Requirement: UI and JSON surfaces

Courier SHALL serve the embedded Lit data application by default as a shared terminal workspace with typed localization, theme controls, protected authentication, manifest transcripts, and semantic upload/download actions, and SHALL expose only versioned JSON and download surfaces when `--no-ui` is active. The application SHALL provide no command prompt and SHALL continue to use only the existing guarded delivery API.

#### Scenario: An authorized browser delivery opens

- **WHEN** authentication and metadata admission succeed
- **THEN** the terminal workspace presents the verified route, navigation, entries, and permitted actions without exposing arbitrary command execution

#### Scenario: UI-disabled root is requested

- **WHEN** a browser opens a no-UI delivery root
- **THEN** Courier returns a JSON description without embedding HTML or protected metadata

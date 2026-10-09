## ADDED Requirements

### Requirement: Interactive browser delivery lifecycle

Interactive browser download and upload readiness SHALL show route, URL, full delivery UUID, source, destination, foreground or background mode, and the applicable stop instruction. A foreground delivery SHALL maintain an elapsed status line until cancellation or fatal failure; a background delivery SHALL return after readiness.

#### Scenario: Browser delivery runs in foreground

- **WHEN** a browser download or upload becomes ready without `--background`
- **THEN** Courier shows its complete readiness identity and a live elapsed status until the user stops it

#### Scenario: Browser delivery runs in background

- **WHEN** the same delivery becomes ready with `--background`
- **THEN** Courier prints `courier servers stop <delivery-uuid>` and returns after the readiness panel

## ADDED Requirements

### Requirement: Interactive incoming webhook lifecycle

Interactive incoming webhook readiness SHALL show route, URL, full delivery UUID, source, destination, foreground or background mode, and the applicable stop instruction. Foreground webhook service SHALL maintain an elapsed listening status until cancellation or fatal failure.

#### Scenario: Incoming webhook runs in foreground

- **WHEN** an incoming webhook becomes ready without `--background`
- **THEN** Courier shows its complete readiness identity and a live elapsed listening state until stopped

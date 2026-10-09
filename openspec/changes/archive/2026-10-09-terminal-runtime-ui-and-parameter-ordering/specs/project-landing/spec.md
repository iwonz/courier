## ADDED Requirements

### Requirement: Accessible required and metadata presentation

The landing SHALL mark required parameters with a red asterisk immediately before syntax, expose a localized required label to assistive technology, omit visible Required and Optional badges, and render non-empty Default, Repeatable, Applies to, Requires, and Conflicts values as separate content-width metadata surfaces in that order.

#### Scenario: A required parameter is inspected

- **WHEN** a required argument or option renders in English or Russian
- **THEN** its red decorative asterisk precedes the syntax and its localized required state remains available to a screen reader

#### Scenario: An optional parameter is inspected

- **WHEN** an optional option renders
- **THEN** no visible Optional badge or required asterisk is shown and each relevant metadata category occupies a separate lightweight background row

### Requirement: Icon-led shell selection

The landing SHALL show a code-native terminal icon with POSIX and the local official PowerShell raster icon with PowerShell while preserving visible text, accessible names, keyboard focus, selected state, and normal raster sampling.

#### Scenario: Shell controls render

- **WHEN** the command shell selector is visible
- **THEN** both controls retain their text labels and semantic button state while their decorative icons are hidden from assistive technology

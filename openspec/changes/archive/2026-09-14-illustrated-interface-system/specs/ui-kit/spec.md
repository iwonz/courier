## ADDED Requirements

### Requirement: Intentional control appearance

Courier SHALL present theme, locale, form, file, and policy controls through one shared visual and interaction grammar without exposing unstyled browser-native chrome.

#### Scenario: A user operates a Courier control

- **WHEN** the control is rendered, focused, selected, disabled, or activated with a keyboard or pointer
- **THEN** it uses Courier tokens and visible state treatment while preserving the expected semantic role, accessible name, focus order, and change behavior

### Requirement: Contextual Relay illustration family

Courier SHALL use role-specific Relay illustrations that preserve the character identity while matching the operational context of each appearance.

#### Scenario: Relay appears on a product surface

- **WHEN** a landing, authentication, delivery, or administration view includes Relay
- **THEN** the illustration's wardrobe, tools, posture, and surrounding scene communicate that view's purpose without duplicating an unrelated pose or carrying required information

#### Scenario: Illustration assets are audited

- **WHEN** the asset verification gate runs
- **THEN** every shipped illustration has a declared role, local path, dimensions, digest, prompt, and authorship notice and no undeclared identity raster is bundled


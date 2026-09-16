## MODIFIED Requirements

### Requirement: Contextual Relay illustration family

Courier SHALL use role-specific Relay illustrations that preserve the character identity while matching the operational context and live layout of each appearance. Full-scene artwork SHALL reserve purpose-composed quiet regions for controls, concentrate narrative detail around the edges of those regions, and contain no embedded product text, third-party logos, or fake UI.

#### Scenario: Relay appears on a product surface

- **WHEN** a landing, authentication, delivery, or administration view includes Relay
- **THEN** the illustration's wardrobe, tools, posture, and surrounding scene communicate that view's purpose without duplicating an unrelated pose or carrying required information

#### Scenario: Illustration assets are audited

- **WHEN** the asset verification gate runs
- **THEN** every shipped illustration has a declared role, local path, dimensions, digest, prompt, and authorship notice and no undeclared identity raster is bundled

#### Scenario: Relay appears behind a product surface

- **WHEN** a landing section overlays live controls on a Relay scene
- **THEN** the character, props, control bay, and responsive crop frame those controls as one composition without obscuring either the interface or the operational narrative

## ADDED Requirements

### Requirement: Smooth stationary scene response

Courier SHALL provide a shared responsive scene primitive whose base illustration remains stationary and whose optional fine-pointer glow and refraction converge through elapsed-time interpolation. The primitive SHALL stop animation after convergence, cancel scheduled work when disconnected, return gently to an ambient position after pointer exit, and disable moving refraction for coarse pointers or reduced motion.

#### Scenario: Fine pointer coordinates change abruptly

- **WHEN** a pointer target moves between distant points in consecutive events
- **THEN** the rendered lens advances smoothly according to elapsed time and reaches the target without a visible single-frame jump

#### Scenario: A scene is removed during interpolation

- **WHEN** the component disconnects while an animation frame is pending
- **THEN** it cancels the frame and retains no active pointer loop or pointer-driven transform on the base illustration

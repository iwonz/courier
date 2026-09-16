## ADDED Requirements

### Requirement: Shared terminal primitives

Courier SHALL provide one Lit and TypeScript UI package for delivery pages, administration pages, and the project landing page, including shared terminal surfaces, transcript rows, and read-only command demonstrations, with no copied component implementations between consumers.

#### Scenario: A consumer presents operational output

- **WHEN** a Courier browser surface renders a command, stage, result, warning, or failure transcript
- **THEN** its structure, tokens, focus treatment, localization hooks, and accessible live behavior come from the shared package

### Requirement: Stable amorphous scene response

Courier SHALL provide a shared responsive scene primitive whose stable host tracks a fine pointer while every decorative descendant ignores hit testing. Its base illustration SHALL remain stationary and its multi-lobed glow and restrained refraction SHALL converge through one elapsed-time interpolation loop, stop after convergence, cancel scheduled work when disconnected, return gently to an ambient position after pointer exit, and disable moving effects for coarse pointers or reduced motion.

#### Scenario: The lens crosses its own decorative image

- **WHEN** a fine pointer moves through a refracted or highlighted region
- **THEN** the host retains pointer tracking without alternating leave events, snapping to ambient coordinates, or transforming the base illustration

#### Scenario: A scene is removed during interpolation

- **WHEN** the component disconnects while an animation frame is pending
- **THEN** it cancels the frame and retains no active pointer loop

### Requirement: Fixed-aspect route terminals

Courier SHALL provide deterministic cubic Bezier paths and fixed-aspect code-native terminal markers for source-to-destination presentations without coupling consumers to a stretched SVG viewport.

#### Scenario: A route is rendered in a non-square viewport

- **WHEN** the connector SVG scales to a wide or portrait layout
- **THEN** Source and Destination markers remain circular within one rendered pixel

### Requirement: Read-only terminal demonstrations

Courier SHALL provide accessible immutable command demonstrations with copy, activation, replay, bounded deterministic transcript stages, selection reset, disconnect cleanup, and reduced-motion completion without exposing editable command input or invoking a shell, installer, transfer, or network operation.

#### Scenario: A visitor runs a browser demonstration

- **WHEN** the Run demo control is activated with pointer or keyboard
- **THEN** the localized transcript identifies itself as a preview, advances through its declared stages, and states that no browser-side command or transfer occurred

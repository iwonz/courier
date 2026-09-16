## MODIFIED Requirements

### Requirement: Smooth stationary scene response

Courier SHALL provide a shared responsive scene primitive whose base illustration remains stationary and activates once either eagerly or within a declared viewport proximity. Its optional fine-pointer glow SHALL use elapsed-time interpolation, and its refracted duplicate SHALL be created only after fine-pointer interaction and removed after returning to ambient. The primitive SHALL stop animation after convergence, disconnect pending observers, cancel scheduled work when disconnected, and never create moving refraction for coarse pointers or reduced motion.

#### Scenario: Fine pointer coordinates change abruptly

- **WHEN** a pointer target moves between distant points in consecutive events
- **THEN** the rendered lens advances smoothly according to elapsed time and reaches the target without a visible single-frame jump while the base illustration remains stationary

#### Scenario: A scene is removed during interpolation

- **WHEN** the component disconnects while an observer or animation frame is pending
- **THEN** it disconnects the observer, cancels the frame, and retains no active pointer loop or pointer-driven transform on the base illustration

#### Scenario: A scene has not been approached

- **WHEN** a non-eager scene remains outside its preload boundary
- **THEN** it preserves layout without creating an image element or requesting its source

### Requirement: Shared terminal primitives

Courier SHALL provide one Lit and TypeScript UI package for delivery pages, administration pages, and the project landing page, including a terminal workspace for real API-backed product surfaces and an immutable command readout for landing commands, with no copied component implementations between consumers.

#### Scenario: A product surface presents operational output

- **WHEN** delivery or administration renders real API-backed state, output, warning, or failure content
- **THEN** its terminal workspace structure, tokens, focus treatment, localization hooks, and accessible behavior come from the shared package

#### Scenario: A consumer presents operational output

- **WHEN** a delivery or administration consumer presents live operational output
- **THEN** its terminal workspace comes from the shared package while landing-only commands use the immutable command readout

#### Scenario: The landing presents a command

- **WHEN** a generated route, installation, or CLI command is available
- **THEN** the shared readout exposes immutable text, localized context, Copy, live feedback, details, and footer actions without demo execution state

### Requirement: Stable amorphous scene response

Courier SHALL provide a shared responsive scene primitive whose stable host tracks a fine pointer while every decorative descendant ignores hit testing. Its base illustration SHALL remain stationary, while its multi-lobed glow and interaction-only restrained refraction SHALL converge through one elapsed-time interpolation loop, stop after convergence, cancel scheduled work when disconnected, return gently to an ambient position after pointer exit, and disable moving effects for coarse pointers or reduced motion.

#### Scenario: The lens crosses its own decorative image

- **WHEN** a fine pointer moves through a refracted or highlighted region
- **THEN** the host retains pointer tracking without alternating leave events, snapping to ambient coordinates, or transforming the base illustration

#### Scenario: The pointer returns to ambient

- **WHEN** pointer exit interpolation converges on the ambient position
- **THEN** the refracted duplicate is removed while the stationary base and non-essential ambient treatment remain

#### Scenario: A scene is removed during interpolation

- **WHEN** the component disconnects while an animation frame is pending
- **THEN** it cancels the frame and retains no active pointer loop or refracted duplicate

## ADDED Requirements

### Requirement: Shared command readout

Courier SHALL provide an accessible immutable command readout with localized heading, description, Copy action, reserved live copy status, details and footer-action slots, deterministic reset when the command identity changes, and no timers, transcripts, editable input, Run action, Replay action, or execution behavior.

#### Scenario: Copy succeeds or fails

- **WHEN** a user activates Copy and clipboard access succeeds or fails
- **THEN** the exact command is offered to the clipboard and a localized result is announced in the reserved status region without geometry movement

#### Scenario: The command identity changes

- **WHEN** a consumer replaces the displayed command selection
- **THEN** prior copy feedback resets and no timed work remains

### Requirement: Shared control frame

Courier SHALL provide shared control-frame tokens and a semantic icon-link component so external icon links and segmented preference controls use the same dimensions, border, padding, radius, surface, hover, focus, and theme behavior while preserving their native accessible roles.

#### Scenario: Header controls are rendered together

- **WHEN** an external icon link appears beside theme or locale radiogroups
- **THEN** their outer chrome matches across light, dark, and system themes while the link remains a link and the preferences remain radio controls

## REMOVED Requirements

### Requirement: Read-only terminal demonstrations

**Reason**: Browser command simulations add timers and implied execution without performing Courier operations.

**Migration**: Use the immutable shared command readout for copyable landing commands and retain the terminal workspace only for real API-backed delivery and administration surfaces.

# ui-kit Specification

## Purpose
Define the single reusable visual, accessibility, theme, localization, component, and identity-asset foundation shared by every Courier browser surface.

## Requirements

### Requirement: Shared UI package

Courier SHALL provide one Lit and TypeScript UI package for delivery pages, administration pages, and the project landing page, with no copied component implementations between consumers.

#### Scenario: Consumer imports a control

- **WHEN** a Courier web application imports a public component
- **THEN** its behavior, tokens, icons, and accessible states come from the shared package

### Requirement: Auditable identity assets

Courier SHALL retain only approved Courier identity assets and notices with local digests and SHALL exclude injected scripts, credentials, remote resources, obsolete mirror messaging, and reference-page executable code.

#### Scenario: Asset integrity check

- **WHEN** the repository quality gate runs
- **THEN** every retained identity asset matches its provenance digest and every manifest entry resolves locally

### Requirement: User-selectable themes

Courier SHALL provide visible `system`, `light`, and `dark` theme choices, default to `system`, persist explicit preference locally, meet accessible contrast, and respect reduced-motion preferences.

#### Scenario: System theme changes

- **WHEN** the stored preference is `system` and the operating-system color scheme changes
- **THEN** the resolved theme updates without replacing the stored preference

### Requirement: Extensible localization

Courier SHALL provide typed English and Russian catalogs, use English as the fallback, negotiate the initial locale from browser languages, and permit additional catalog modules without backend changes.

#### Scenario: Unsupported browser locale

- **WHEN** no supported locale matches the browser language list
- **THEN** English is selected and every requested message key resolves

### Requirement: Exact UI source coverage

Courier SHALL enforce 100% statements, branches, functions, and lines for first-party TypeScript sources while excluding generated output and third-party code.

#### Scenario: Untested branch is introduced

- **WHEN** a first-party TypeScript branch is not executed by the test suite
- **THEN** the repository verification gate fails

### Requirement: Coherent brand identity system

Courier SHALL define and apply one mature, utilitarian identity across browser surfaces, documentation, and retained assets using shared semantic tokens and components rather than consumer-specific brand implementations.

#### Scenario: A browser surface presents Courier

- **WHEN** the landing, delivery, or administration application renders
- **THEN** it uses the same positioning, wordmark, semantic color roles, typography, focus treatment, and operational visual grammar from the shared UI package

### Requirement: Non-essential mascot guidance

Courier SHALL use the original Relay courier-pigeon character only as a restrained orientation and state-communication aid, with locally verified provenance and without making meaning depend on the image.

#### Scenario: Mascot artwork is unavailable

- **WHEN** Relay cannot be loaded or is hidden from assistive technology
- **THEN** headings, status text, controls, and progress information still communicate the complete workflow

### Requirement: Calm operational communication

Courier SHALL use concise, concrete, non-alarmist language that identifies actions and verified states without unsupported guarantees, secret disclosure, or humor during failures.

#### Scenario: A delivery changes state

- **WHEN** a user sees preparation, transfer, verification, completion, or failure feedback
- **THEN** the message names the current or stopped operation and preserves an actionable, technically accurate tone in every supported locale

### Requirement: Intentional control appearance

Courier SHALL present theme, locale, form, file, and policy controls through one shared visual and interaction grammar without exposing unstyled browser-native chrome.

#### Scenario: A user operates a Courier control

- **WHEN** the control is rendered, focused, selected, disabled, or activated with a keyboard or pointer
- **THEN** it uses Courier tokens and visible state treatment while preserving the expected semantic role, accessible name, focus order, and change behavior

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

### Requirement: Compact icon preference controls

Courier SHALL render theme and locale selectors as shared icon-only segmented radiogroups without visible group or option labels while retaining localized accessible names, selected state, persisted preference, roving focus, and arrow/Home/End keyboard behavior. Locale options SHALL use native flag emoji as their visible symbols.

#### Scenario: A user operates a locale selector

- **WHEN** the user points to, focuses, or navigates an English or Russian locale option
- **THEN** a flag emoji identifies the option visually while its localized name and radio state remain available programmatically

#### Scenario: A user operates an icon preference selector

- **WHEN** the user points to, focuses, or navigates a theme or locale option
- **THEN** its meaning is available programmatically, its state is visibly distinguishable, and changing it has the same persisted behavior as the labeled control

### Requirement: Relay compact mark

Courier SHALL use a recognizable, repository-local, generated Relay mascot image as its compact product mark across browser components and favicons.

#### Scenario: A compact Courier identity is rendered

- **WHEN** the wordmark has limited space or a favicon is displayed
- **THEN** the locally stored transparent mark depicts Relay with sufficient light/dark contrast and no remote image dependency

### Requirement: Local official brand marks

Courier SHALL provide a shared brand-icon primitive that renders pinned, repository-bundled official geometry in Courier monochrome without runtime network requests. Third-party source, license, attribution, and trademark constraints SHALL be recorded separately from Courier-owned identity assets.

#### Scenario: A branded channel is rendered

- **WHEN** a supported package manager, shell, operating system, or distribution mark appears
- **THEN** the UI uses its registered official geometry, inherits the requested Courier color, and exposes no remote asset URL

### Requirement: Shared styled checkbox

Courier SHALL provide a shared checkbox primitive with custom visual treatment, native checked and disabled semantics, localized accessible labeling, keyboard activation, focus indication, and a composed change event.

#### Scenario: A visitor toggles a view filter

- **WHEN** the checkbox is activated by pointer or keyboard
- **THEN** its checked state and accessible state update once and consumers receive the new boolean value

### Requirement: Shared Bezier route geometry

Courier SHALL provide deterministic cubic Bezier path generation for source-to-destination route presentations without coupling consumers to a specific layout.

#### Scenario: A consumer supplies two endpoint positions

- **WHEN** a route path is requested for finite coordinates
- **THEN** the helper returns a stable horizontal cubic Bezier path spanning those positions

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

### Requirement: Fixed-aspect route terminals

Courier SHALL provide deterministic cubic Bezier paths and fixed-aspect code-native terminal markers for source-to-destination presentations without coupling consumers to a stretched SVG viewport.

#### Scenario: A route is rendered in a non-square viewport

- **WHEN** the connector SVG scales to a wide or portrait layout
- **THEN** Source and Destination markers remain circular within one rendered pixel

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

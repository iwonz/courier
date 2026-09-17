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

Courier SHALL expose theme through one semantic icon button that cycles `system`, `light`, and `dark`, defaults to browser-resolved `system`, persists an explicit valid preference when storage is available, meets accessible contrast, and respects reduced-motion preferences. The system choice SHALL continue following browser color-scheme changes until the user stores a different choice.

#### Scenario: Theme cycles from its default

- **WHEN** a visitor with no saved theme repeatedly activates the theme control
- **THEN** the preference cycles system to light to dark to system while system resolution continues following the browser media query

#### Scenario: System theme changes

- **WHEN** the stored preference is `system` and the operating-system color scheme changes
- **THEN** the resolved theme updates without replacing the stored preference

### Requirement: Extensible localization

Courier SHALL provide typed English and Russian catalogs, use English as the fallback, negotiate the initial locale from browser languages, and expose locale through one semantic icon button that cycles English and Russian. A valid explicit choice SHALL persist when storage is available, and additional catalog modules SHALL require no backend changes.

#### Scenario: Locale cycles from its default

- **WHEN** a visitor activates the locale control
- **THEN** English and Russian alternate, persist when storage is available, and dispatch the existing composed locale-change event

#### Scenario: Unsupported browser locale

- **WHEN** no supported locale matches the browser language list
- **THEN** English is selected and every requested message key resolves

### Requirement: Exact UI source coverage

Courier SHALL enforce 100% statements, branches, functions, and lines for first-party TypeScript sources while excluding generated output and third-party code.

#### Scenario: Untested branch is introduced

- **WHEN** a first-party TypeScript branch is not executed by the test suite
- **THEN** the repository verification gate fails

### Requirement: Coherent brand identity system

Courier SHALL provide one repository-owned vivid editorial identity across landing, delivery, administration, README, favicons, and wordmarks using deep ink, cobalt, warm cream, safety coral, and restrained mint. The compact mascot SHALL be an ImageGen-authored transparent raster swift with a mature aerodynamic courier silhouette, and every contextual raster SHALL be checksum-recorded with no remote dependency.

#### Scenario: A browser surface presents Courier

- **WHEN** the landing, delivery, or administration application renders
- **THEN** it uses the local vivid identity, accessible contrast, rounded kinetic geometry, system typography, and shared workbench grammar without engraved scenery, remote fonts, or remote images

### Requirement: Non-essential mascot guidance

Courier SHALL use the swift courier only as a restrained orientation and state-communication aid, with locally verified provenance and without making meaning depend on the image.

#### Scenario: Mascot artwork is unavailable

- **WHEN** swift artwork cannot be loaded or is hidden from assistive technology
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

### Requirement: Smooth stationary scene response

Courier SHALL provide a responsive scene primitive that activates eager imagery immediately and non-eager imagery once at viewport proximity, creates only the browser-selected wide or portrait base image, and disconnects its observer after activation or removal. It SHALL not register pointer movement, create refracted duplicates, run requestAnimationFrame, or encode essential information.

#### Scenario: A non-eager scene approaches the viewport

- **WHEN** its proximity observer first intersects
- **THEN** one responsive base picture is created, the observer disconnects, and no decorative interaction layer is allocated

#### Scenario: Fine pointer coordinates change abruptly

- **WHEN** pointer coordinates change over a rendered scene
- **THEN** no pointer-driven visual state, refraction, animation, or base-image transform is created

#### Scenario: A scene is removed during interpolation

- **WHEN** the component disconnects while a proximity observer is pending
- **THEN** it disconnects the observer and retains no active animation, pointer loop, or duplicate image

#### Scenario: A scene has not been approached

- **WHEN** a non-eager scene remains outside its preload boundary
- **THEN** it preserves layout without creating an image element or requesting its source

#### Scenario: A scene is removed before activation

- **WHEN** a non-eager scene disconnects before its observer intersects
- **THEN** the observer disconnects and no image request or scheduled animation remains

### Requirement: Compact icon preference controls

Courier SHALL render theme and locale selectors as separate semantic icon buttons without visible labels or radiogroups. Each button SHALL expose localized current and next values, use shared control-frame styling, preserve keyboard activation and focus indication, and delegate persistence and document updates to one shared browser preference controller.

#### Scenario: A user operates an icon preference selector

- **WHEN** the user points to, focuses, or activates the theme or locale button
- **THEN** its current and next states are available programmatically and activation advances exactly once through the documented cycle

#### Scenario: A user operates a locale selector

- **WHEN** the visitor activates the locale button by pointer or keyboard
- **THEN** its localized accessible name identifies the current and next locale while the visible icon advances once without exposing a radiogroup

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

### Requirement: Fixed-aspect route terminals

Courier SHALL provide deterministic cubic Bezier paths and fixed-aspect code-native terminal markers for source-to-destination presentations without coupling consumers to a stretched SVG viewport.

#### Scenario: A route is rendered in a non-square viewport

- **WHEN** the connector SVG scales to a wide or portrait layout
- **THEN** Source and Destination markers remain circular within one rendered pixel

### Requirement: Shared command readout

Courier SHALL provide an independent immutable command strip with localized heading, Copy action, reserved live feedback, optional details and footer actions, and deterministic reset when command identity changes. It SHALL not depend on the workbench component or render a fake shell prompt, editable field, transcript, timer, Run action, or execution behavior.

#### Scenario: Copy succeeds or fails

- **WHEN** a user activates Copy and clipboard access succeeds or fails
- **THEN** the exact command is offered to the clipboard and a localized result is announced in the reserved status region without geometry movement

#### Scenario: The command identity changes

- **WHEN** a consumer replaces the displayed command selection
- **THEN** prior copy feedback resets synchronously and no timed work remains

### Requirement: Shared control frame

Courier SHALL provide shared control-frame tokens and a semantic icon-link component so external icon links and cyclic preference buttons use the same dimensions, border, padding, radius, surface, hover, focus, and theme behavior while preserving their native accessible roles.

#### Scenario: Header controls are rendered together

- **WHEN** an external icon link appears beside theme or locale buttons
- **THEN** their outer chrome matches across light, dark, and system themes while the link remains a link and each preference remains a button

### Requirement: Contextual Vector illustration family

Courier SHALL use the locally generated swift mark and one responsive hero illustration pair as its complete identity raster set. The landing SHALL use the hero inside the first section; delivery and administration MAY reuse it only as restrained decorative atmosphere behind complete semantic workbenches. Artwork SHALL contain no required information, product text, third-party logo, credential, watermark, or fake UI.

#### Scenario: Vector appears on a product surface

- **WHEN** a landing, authentication, delivery, or administration view includes the swift
- **THEN** the semantic interface remains complete without the image and no surface requires a separate stitched or role-specific background

#### Scenario: Illustration assets are audited

- **WHEN** the asset verification gate runs
- **THEN** every shipped illustration has a declared role, local path, dimensions, digest, prompt, sequence metadata, and authorship notice and no undeclared identity raster is bundled

### Requirement: Vector compact mark

Courier SHALL use a recognizable repository-local ImageGen-authored raster swift mark across browser components, wordmarks, and favicons. The transparent mark SHALL retain a clean silhouette and controlled cobalt/coral accents while remaining readable at compact product-chrome sizes.

#### Scenario: A compact Courier identity is rendered

- **WHEN** the wordmark has limited space or a favicon is displayed
- **THEN** the local raster swift remains clear against light and dark surfaces without a remote image dependency or vector substitute

### Requirement: Shared workbench primitive

Courier SHALL provide a neutral accessible workbench with optional heading, actions, status, body, and footer regions for real product state. It SHALL use high-contrast translucent surfaces, generous rounded geometry, compact spacing, and normal sans-serif content without viewport-sized framing, engraved texture, shell prompts, fake execution, or arbitrary command input.

#### Scenario: A product surface presents state

- **WHEN** landing, delivery, or administration renders structured content
- **THEN** shared workbench structure presents it with readable contrast and lightweight spatial hierarchy without copied consumer-specific implementations

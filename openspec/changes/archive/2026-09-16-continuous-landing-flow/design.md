## Context

The landing currently renders four viewport-snapped sections, each with a separately composed responsive scene and a timer-driven command demonstration. The shared scene eagerly constructs both base and refracted image layers, while the landing builds all below-fold scenes immediately. The delivery and administration applications reuse the same scene and preference controls but rely on the terminal component for real API-backed workspaces. See `proposal.md` for motivation and the capability deltas for observable requirements.

## Goals / Non-Goals

**Goals:**

- Keep selectors and exact copyable commands while removing every simulated execution path from the landing.
- Make four variable-height sections read as one responsive illustrated journey with no visible tile boundaries.
- Bound initial image work and bundle size without layout shifts or inaccessible fallback behavior.
- Share command-readout and icon-link behavior through the UI package while preserving delivery and administration terminal behavior.

**Non-Goals:**

- Changing the CLI contract, HTTP APIs, authentication, administration operations, releases, or README banner.
- Adding editable command input, browser-side command execution, analytics, external runtime assets, or a new UI framework.
- Replacing the delivery or administration compositions in this change.

## Decisions

### Command readouts are not terminals

The timer-driven command-demo element will be replaced by a smaller immutable command-readout element with heading, description, copy status, details, and footer action slots. The real terminal structural primitive remains available to delivery and administration. This separation removes simulation state and avoids implying that copying a command is execution. Keeping a disabled Run control was rejected because it would preserve the wrong product model and dead localization surface.

### Navigation follows a header-aware observation band

The landing will use natural document flow and CSS smooth scrolling rather than scroll snap. Section observation will use a narrow root-margin band below the measured masthead, with a deterministic scroll-position fallback for initial and boundary states. This tracks the section occupying the reading line instead of comparing ratios that favor short sections. Reduced-motion navigation calls immediate scrolling. Fixed section heights were rejected because narrow localized content must expand without clipping.

### One journey is delivered as bounded responsive segments

The art direction is generated as one continuous Relay journey, then shipped as four wide and four portrait WebP segments. Matching transition bands overlap behind masked section edges so each image can load independently while the page reads as one composition. A monolithic panorama was rejected because it would force a large initial request and inefficient mobile cropping; unrelated cross-fades were rejected because they do not create spatial continuity.

### Scene activation is one-shot and proximity based

The scene component owns activation independently from image rendering. Eager scenes render immediately with high fetch priority. Non-eager scenes observe a one-viewport root margin and create their picture only once when activated; environments without IntersectionObserver render a native-lazy image fallback. Disconnect cancels the observer and animation frame. A consumer-controlled global preloader was rejected because it duplicates viewport policy and makes shared scenes inconsistent.

### Refraction is an interaction-only secondary image

Only the stationary base scene is present before interaction. A fine-pointer move activates the refracted duplicate, and returning to the ambient state removes it after interpolation settles. Coarse pointers and reduced motion never construct the duplicate. The ambient glow remains CSS-only. Keeping a permanent second image was rejected because it doubles decode and memory cost even when the effect is never used.

### Shared header actions use one control frame

The UI package will expose semantic control-frame tokens and an icon-link component. Segmented preference controls and the external GitHub action consume the same dimensions, border, radius, surface, hover, and focus tokens while retaining link versus radio semantics. Styling the GitHub anchor only in the landing was rejected because shared chrome would drift across themes and browser applications.

### Budgets are deterministic repository gates

Asset verification records dimensions, byte counts, hashes, prompts, roles, and sequence metadata and rejects a landing segment above 225 KiB or all eight above 1.6 MiB. The production Pages build additionally rejects landing JavaScript above 45 KiB gzip and CSS above 8 KiB gzip. Browser acceptance records requests to prove the initial viewport loads only the selected hero source and at most the next responsive segment.

## Risks / Trade-offs

- [Generated segments reveal seams at uncommon aspect ratios] → Compose explicit low-detail transition bands, test both picture sources at every supported viewport, and retain masked overlap independent of image content.
- [IntersectionObserver timing differs between browsers] → Make activation idempotent, use generous one-viewport margins, and provide native lazy-loading when the API is unavailable.
- [Natural section heights make active navigation ambiguous at boundaries] → Anchor observation to one measured header-aware reading band and cover exact boundaries in browser tests.
- [On-demand refraction briefly decodes a second image] → Reuse the already cached source, keep the base visible, and make the effect non-essential.
- [Strict compressed budgets become toolchain-sensitive] → Measure deterministic production artifacts with the repository's pinned Node/build configuration.

## Migration Plan

Create and validate the OpenSpec delta, replace the shared command component and control framing, implement one-shot scene activation, recompose landing layout/navigation, generate and optimize the continuous artwork, update provenance and documentation, then run unit, browser, asset, bundle, OpenSpec, release-snapshot, and clean-worktree gates. Archive the completed change before the single conventional commit. Rollback is the single feature commit because no persisted data, public API, or contract migration is introduced.

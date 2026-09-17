## Context

The current browser system uses Relay-specific asset modules, four landing scenes, segmented preferences, a pointer-driven refractive scene, and terminal-framed product layouts. The implementation already has contract-generated routes and commands, responsive asset provenance, exact web coverage, and guarded delivery/admin APIs. This change preserves those functional boundaries while replacing the visual and component grammar.

## Goals / Non-Goals

**Goals:**

- Make Courier identifiable through Vector, `From here to anywhere.`, and one restrained aerospace-editorial system.
- Present routing in the first viewport without duplicating the CLI route contract.
- Use one-button theme and locale cycling with correct browser defaults and persistence.
- Remove decorative pointer work and fake terminal affordances while preserving responsive performance and accessibility.
- Give delivery and administration coherent, task-oriented workbench layouts.

**Non-Goals:**

- Changing CLI syntax, endpoint behavior, HTTP payloads, authentication, policy semantics, or release automation.
- Adding editable commands, browser execution, external fonts/assets, analytics, or unguarded administration state.
- Encoding required product information into generated images.

## Decisions

### Vector is a generated character with a deterministic mark

The canonical adult moth and editorial scenes are generated as repository-owned raster assets. The compact mark is a deterministic SVG derived from angular wing and waypoint geometry so it remains crisp and auditable at favicon size. Generated images never contain text, interface chrome, credentials, third-party logos, or required information.

### Identity references are role based

Asset modules and exports use landing, delivery, and administration roles rather than the mascot's name. This keeps application architecture independent from future illustration changes. Provenance remains the source of truth for the character, generation lineage, prompts, hashes, dimensions, roles, and sequence.

### Preferences are cyclic buttons

Theme and locale keep their existing custom-element names and change events, but each renders one semantic button. Theme cycles system, light, and dark; locale cycles English and Russian. A shared document preference controller owns persistence, media changes, and the resolved document dataset, preventing duplicate listeners and inconsistent state.

### Routing is the hero interaction

The first section combines the headline, generated scene, endpoint banks, measured Bezier route, payload signal, command, Copy action, and applicable options. Endpoint explanations are assistive-only. Selection changes only on activation; hover and focus never mutate route state. Narrow or short viewports expand naturally rather than clipping.

### Workbenches replace terminals

`courier-workbench` supplies neutral heading, actions, status, body, and footer regions without a shell metaphor. `courier-command-readout` becomes an independent immutable code strip. Delivery uses workbenches for access and manifests. Administration uses a navigator and selected inspector, preserving selection by UUID across authoritative snapshots.

### Scenes are static and proximity loaded

`courier-scene` retains eager and proximity activation plus responsive picture selection, but removes pointer listeners, duplicate imagery, masks, requestAnimationFrame work, and ambient effects. The landing uses three wide/portrait sequence segments with masked overlap; delivery and administration each use one wide/portrait scene.

## Risks / Trade-offs

- [The combined hero is too dense on small screens] → Keep desktop first-viewport containment but allow natural mobile expansion, compact endpoint controls, and invariant command geometry.
- [Generated character consistency drifts] → Establish the canonical transparent Vector reference first and use it for every contextual generation.
- [A one-button preference is less discoverable than a segmented group] → Expose current and next values through localized accessible names, visible icons, title text, focus treatment, and persisted feedback.
- [Master-detail administration hides simultaneous detail] → Keep every server/delivery visible in the navigator, default deterministically, preserve selection across SSE, and stack navigator plus inspector on narrow screens.
- [Large identity replacement regresses bundles] → Preserve JavaScript/CSS budgets, reduce landing scenes from eight to six, and enforce per-surface image totals.

## Migration Plan

Create and validate the OpenSpec change; implement shared identity, preference, scene, workbench, and readout primitives; restructure the landing; migrate delivery and administration; generate and optimize new artwork; replace and verify provenance; update embedded assets, documentation, unit/browser acceptance, and exact coverage; archive the change; commit once; publish the branch and fast-forward main after every gate passes. Rollback is the single feature commit because no persisted application data or public protocol changes.

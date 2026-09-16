# Courier brand system

Courier is a dependable delivery layer for files and directories. Its identity should feel like a field manual connected to a live network control room: direct, organized, technically credible, and calm under pressure.

## Core idea

**Positioning:** Courier is the cross-platform CLI for moving files across local, SSH, browser, and webhook boundaries while keeping the route, destination, and result explicit.

**Promise:** Move files. Keep control.

**Supporting line:** One self-contained binary. Explicit routes. Verified arrival.

Courier is not positioned as magic infrastructure or as a replacement for every synchronization product. It is a focused transfer tool whose useful guarantees are visible before, during, and after delivery.

## Personality

- **Capable:** knows the route, checks the handoff, and reports what actually happened.
- **Composed:** treats routine work, interruption, and failure with the same measured tone.
- **Plain-spoken:** uses concrete verbs and technical facts instead of abstractions.
- **Protective:** makes unsafe collisions, trust failures, and credential boundaries visible.
- **Quietly human:** warm details and Relay add character without turning operational work into a joke.

The result should be more mature than a playful developer toy and more welcoming than military, cyberpunk, or enterprise-security theater.

## Relay, the field operator

Relay is Courier's original pigeon mascot: an adult courier pigeon wearing a compact graphite utility harness with signal-lime route tabs and a sealed data capsule. Relay is observant, prepared, and approachable. The character represents dependable movement and a verified handoff, not speed at any cost.

Use Relay:

- in the landing hero and brand storytelling;
- to orient a user at an authentication or empty state;
- at meaningful confirmation moments;
- in editorial material where a human presence improves comprehension.

Do not use Relay:

- as the only way to communicate state or meaning;
- inside dense tables, repeated list rows, or every panel;
- to make light of a failure, security decision, or lost connection;
- with speech bubbles that imitate a human support agent;
- as a child, toy, robot, superhero, soldier, or weapon-bearing character;
- with third-party uniforms, delivery-company symbols, or another product's mascot treatment.

[`relay-mascot.png`](../web/ui/assets/relay-mascot.png) is the canonical character reference. Product surfaces use context-specific scenes instead of repeating that neutral pose:

| Role | Surface | Situation |
|---|---|---|
| Dispatch navigator | Landing hero | Plans a visible local-to-remote route across a full dispatch landscape |
| Route cartographer | Landing routing | Maps local, SSH, browser, and webhook endpoints around the route explorer |
| Field installer | Landing installation | Fits the self-contained Courier capsule to available systems around the channel chooser |
| Integrity inspector | Landing command reference | Verifies the capsule beside the generated command and option registry |
| Access controller | Delivery authentication | Presents a credential at a private checkpoint |
| Operations controller | Administration | Observes routes and adjusts delivery policy from a control room |
| Courier in flight | README | Carries a confirmed handoff from source to destination |

Wardrobe and tools change only to explain the role. Relay's adult anatomy, graphite harness, cream capsule, signal-lime route tabs, orange beak, and editorial line work remain stable. Every scene is decorative in product UI unless surrounding editorial copy specifically describes Relay; the interface must remain complete when the image is unavailable.

## Visual language

### Palette

| Role | Color | Use |
|---|---|---|
| Graphite 900 | `#151714` | Dark canvas, primary ink, inverse surfaces |
| Paper 50 | `#f3f4e9` | Light canvas, text on graphite |
| Signal lime | `#d4ff45` | Primary route, active state, decisive action |
| Beak orange | `#ff8758` | Focus rings and small human accents |
| Danger coral | `#ff6b5f` | Stopped or unsafe states only |
| Field gray | `#737b6d` | Secondary operational detail |

Signal lime is a locator, not wallpaper. Large surfaces stay graphite, paper, or neutral field colors. Orange never competes with the main action; its main product role is keyboard focus. Status colors always appear with text, never as color-only meaning.

### Typography

Courier uses system-resident fonts and ships no remote font dependency. Display text is compact, heavy, and slightly tightened. Body copy is neutral and readable. Labels, routes, byte counts, versions, and commands use a monospaced stack with tabular numerals.

- Headlines use sentence case, not title case.
- Operational labels may use tracked uppercase at small sizes.
- Long paths must wrap or truncate intentionally; they must never force viewport overflow.
- Commands retain their exact capitalization and punctuation.

### Shape and composition

The layout uses a disciplined grid, visible registration lines, route nodes, compact status dots, and framed operational panels. Corners are slightly rounded rather than pill-shaped. Shadows are restrained. Route lines and Relay's directional posture provide motion. Full-bleed editorial scenes may form an environmental backdrop; restrained translucent instrument panels can align with deliberately composed equipment bays, but floating glass cards, decorative gradients, neon spectacle, and detached 3D objects do not define the system.

The Courier mark is a compact profile of Relay with the orange beak, alert eye, and signal-lime capsule harness retained at favicon size. It may appear without the wordmark in compact contexts. Do not rotate it, add speed lines, recolor individual features arbitrarily, or place it on insufficient contrast.

### Motion

Motion confirms state changes and direction. Transitions remain short and interruptible. No essential information depends on animation. Illustrated scene bases do not pan or scale with the pointer; any pointer response is confined to a local amorphous highlight and restrained refractive lens. Coarse pointers use a fixed ambient treatment, and `prefers-reduced-motion` removes non-essential travel, morphing, refraction animation, and smooth scrolling.

### Controls

Courier controls never expose unstyled browser chrome. Small fixed choices such as theme and locale use icon-only shared segmented radio groups with localized programmatic names, visible selection, arrow-key navigation, and one focus stop. Form inputs, selects, and checkboxes keep dependable platform semantics but normalize appearance, spacing, indicators, hover, focus, and disabled states through the shared UI kit. Boolean policy values use a switch treatment; view filters use a designed checkbox; file selection uses a designed action surface while retaining a real keyboard-operable file input.

Custom appearance must not recreate browser responsibilities badly. Labels remain programmatic, Enter and Space activate buttons, arrow keys move within segmented groups, validation remains available, and the orange focus ring is never removed.

## Communication system

### Voice rules

1. Start with the current state or required action.
2. Name the object: route, source, destination, archive, server, or delivery.
3. Use a concrete verb: prepare, connect, send, receive, verify, stop, retry.
4. Report only what Courier can prove. “Confirmed” is stronger and more useful than “done.”
5. State what remains safe after failure when known.
6. Keep credentials, query secrets, and private metadata out of copy, URLs, logs, and examples.

### Operational vocabulary

| Prefer | Avoid |
|---|---|
| route | pipeline, tunnel magic |
| delivery | job, payload blast |
| preparing | working on it |
| sending / receiving | syncing when no synchronization occurs |
| verifying | almost done |
| verified arrival | guaranteed delivery |
| stopped | killed, died |
| blocked by a collision | something went wrong |
| confirmed bytes | transferred bytes when confirmation is unknown |

Canonical state sequence: **Preparing → Routing → Sending or receiving → Verifying → Complete**. A failure message follows: **what stopped + why + safe next action**.

Examples:

- “Route ready. Choose a file to dispatch.”
- “Verifying the archive before commit.”
- “Delivery stopped: the final path already exists. No destination data was changed.”
- “Administration data is temporarily unavailable. Retry the local connection.”

Avoid unsupported guarantees, breathless superlatives, fake urgency, blame, jokes in error states, and claims that Courier is “military-grade,” “unbreakable,” or “the fastest.”

## Product surfaces

- **Landing:** an immersive dispatch brief. Lead with the promise and a non-interactive `Source` to `Destination` handoff, then give routing, installation, and the generated command registry one full viewport each. Every scene is stationary and composed around its full-width working interface rather than placed beside it as a decorative tile. Route and installation state changes only through activation; hover communicates affordance without silently changing content.
- **Delivery UI:** a transfer terminal. Preserve focus on authentication, destination contents, and the single next action. Protected metadata appears only after authorization.
- **Administration UI:** an operations control room. Prioritize live state, bindings, delivery identity, confirmed counters, policy changes, and explicit stop actions.
- **CLI and documentation:** use the same route, delivery, verification, collision, and confirmed-byte vocabulary. Decorative identity never interferes with copy-and-paste commands.

## Accessibility and integrity

Every control is reachable by keyboard with a visible orange focus ring. Light, dark, and system themes retain text and state contrast. English is the fallback language and Russian catalogs preserve intent rather than translating slogans literally. Layouts support 360-pixel viewports, long UUIDs, Unicode paths, zoom, and reduced motion.

Identity assets are local, checksum-recorded, and covered by the MIT license notice. No UI build may fetch a font, image, script, tracking pixel, or theme resource from a third party.

## Generated asset record

Relay and every contextual scene were produced with the built-in OpenAI image generation tool for this repository using the canonical character as the identity reference. The family keeps an adult courier pigeon, compact graphite utility harness, sealed cream data capsule, sparse signal-lime route tabs, orange beak accent, and textured flat editorial/screen-print treatment. Landing scenes have separately composed wide and portrait variants: neither mobile delivery nor desktop delivery is an automatic crop of the other. Their equipment bays, route tables, and headline space are planned around the overlaid interface. Role-specific prompts add dispatch, installation, routing, verification, access, operations, or in-flight tools and environments while prohibiting text, logos, watermarks, weaponry, robot parts, childlike proportions, security theater, and third-party mascot imitation.

The auditable prompt summaries, roles, dimensions, byte counts, and SHA-256 digests are stored in [`web/ui/assets/provenance.json`](../web/ui/assets/provenance.json). Product illustrations use optimized local WebP delivery assets; no runtime surface fetches generated art remotely.

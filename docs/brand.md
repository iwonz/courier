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

The production asset is [`relay-mascot.png`](../web/ui/assets/relay-mascot.png). It is decorative in product UI unless surrounding editorial copy specifically describes Relay. The interface must remain complete when the asset is unavailable.

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

The layout uses a disciplined grid, visible registration lines, route nodes, compact status dots, and framed operational panels. Corners are slightly rounded rather than pill-shaped. Shadows are restrained. Route arrows and the square Courier mark provide motion; decorative gradients, glass surfaces, neon glow, and floating 3D objects do not.

The Courier mark is a source bracket, a parcel node, and a forward route combined in one square. It may appear without the wordmark at favicon or compact-control sizes. Do not rotate it, add speed lines, recolor individual shapes arbitrarily, or place it on insufficient contrast.

### Motion

Motion confirms state changes and direction. Transitions remain short and interruptible. No essential information depends on animation. `prefers-reduced-motion` removes non-essential movement and smooth scrolling.

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

- **Landing:** a dispatch brief. Lead with the promise, one valid installation command, the route model, and evidence for safety and distribution.
- **Delivery UI:** a transfer terminal. Preserve focus on authentication, destination contents, and the single next action. Protected metadata appears only after authorization.
- **Administration UI:** an operations control room. Prioritize live state, bindings, delivery identity, confirmed counters, policy changes, and explicit stop actions.
- **CLI and documentation:** use the same route, delivery, verification, collision, and confirmed-byte vocabulary. Decorative identity never interferes with copy-and-paste commands.

## Accessibility and integrity

Every control is reachable by keyboard with a visible orange focus ring. Light, dark, and system themes retain text and state contrast. English is the fallback language and Russian catalogs preserve intent rather than translating slogans literally. Layouts support 360-pixel viewports, long UUIDs, Unicode paths, zoom, and reduced motion.

Identity assets are local, checksum-recorded, and covered by the MIT license notice. No UI build may fetch a font, image, script, tracking pixel, or theme resource from a third party.

## Generated asset record

Relay was produced with the built-in OpenAI image generation tool for this repository, then stored locally with transparency and deterministic digest verification. Final prompt summary: an original, adult courier pigeon in a compact graphite utility harness, carrying a sealed data capsule with sparse signal-lime route tabs and an orange beak accent; premium flat editorial and screen-print treatment; transparent background; strong small-size silhouette; no text, logo, watermark, weaponry, robot parts, childlike proportions, or third-party mascot imitation.

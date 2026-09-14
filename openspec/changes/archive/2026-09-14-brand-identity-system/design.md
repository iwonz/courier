# Context

Courier already has one Lit-based UI kit and three browser consumers. Its graphite, paper, signal-lime, and beak-orange colors are useful identity seeds, but composition and copy do not yet express the product's core promise: controlled delivery with verified arrival. The existing cute robot-pigeon artwork also competes with the mature operational tone needed for security-sensitive transfer workflows.

# Decisions

## Brand idea: a field manual for dependable delivery

The visual system combines a field manual's clarity with a network control room's live state. Dense technical information uses a disciplined grid, compact labels, route lines, registration marks, and tabular numerals. Spacious reading surfaces, plain language, and warm illustration keep the system from feeling militaristic or hostile.

## Relay is a guide, not decoration

Relay is an original courier pigeon depicted as a capable field operator with a graphite utility harness, signal-lime route markers, and a small orange beak accent. The mascot appears at high-value orientation and confirmation moments, never inside dense tables or as a substitute for status text. The silhouette remains clear at small sizes and avoids childlike, weaponized, robotic, or derivative treatments.

## One semantic token and component system

The shared UI package owns the wordmark, route mark, mascot presentation, status labels, masthead, buttons, panels, progress, typography, and semantic colors. Landing, delivery, and administration applications compose those primitives without introducing parallel brand tokens or copied controls.

## Operational language is calm and verifiable

English is the source product language and Russian is a complete typed localization. Headlines are short. Actions use concrete verbs. Progress language distinguishes preparing, routing, sending, receiving, verifying, and complete states. Errors state what stopped, why, and what remains safe. Marketing avoids superlatives, fear, jokes during failures, and guarantees that the transfer engine cannot prove.

## Generated raster art remains auditable

The selected Relay image is generated once, stored locally, optimized without losing transparency, and covered by the existing digest manifest. Code-native route marks and interface graphics remain SVG/CSS rather than generated raster assets. The final prompt and authorship context are recorded in the asset notice and brand manual.

# Risks / Trade-offs

A distinctive visual system adds CSS surface area. Keeping the changes in the shared UI kit and continuing exact TypeScript coverage limits drift. Generated mascot art is intentionally non-essential: every screen remains complete, understandable, and accessible if the image cannot load.

# Migration Plan

Document the brand, generate and approve Relay, extend shared primitives and tokens, redesign all three applications, update catalogs and tests, rebuild embedded assets, complete browser review, run the full repository gate, then archive this change.

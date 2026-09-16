# Context

Courier already shares Lit components, themes, localization, contract-generated landing data, and responsive scenes. The redesign must reuse those sources of truth, preserve exact test coverage, and distinguish a browser demonstration from real CLI execution.

# Decisions

## Terminal surfaces are presentation and control boundaries

The shared UI kit owns terminal framing, transcript rows, and deterministic read-only demonstrations. Landing commands are immutable generated text. Copy copies only that text; Run demo advances an explicitly simulated transcript and never invokes a shell, transfer, installer, or network operation. Delivery and administration keep their existing forms and versioned APIs and receive no command prompt.

## Demonstrations are deterministic state machines

Each demo owns idle, running, and complete states, one bounded timer at a time, replay, selection reset, disconnect cleanup, and reduced-motion immediate completion. Route, installation, and help transcripts are constructed from the existing route matrix, installer registry, and generated command contract.

## Pointer response uses stable hit testing

The scene host owns pointer tracking and all decorative descendants ignore pointer events. A single elapsed-time loop converges on the latest coalesced target and stops after settling. The refracted duplicate uses an irregular multi-lobed mask centered on the rendered pointer position; the base illustration never transforms. Coarse pointers and reduced motion use a fixed ambient treatment.

## Route terminals are independent of stretched SVG geometry

SVG draws only cubic Bezier connectors. Terminal markers are fixed-aspect HTML elements positioned with percentages, so viewport scaling cannot turn circles into ellipses.

## Artwork is regenerated around the terminal layouts

Twelve responsive cinematic editorial images cover hero, routing, installation, CLI, delivery, and administration roles. Wide and portrait images are separately composed rather than cropped. They contain no text, fake UI, logos, credentials, or required information. The existing README banner is retained.

# Risks / Trade-offs

Simulated output could be mistaken for execution, so every transcript identifies itself as a preview and states that no browser-side command or transfer occurred. Additional art increases repository size, so production WebP assets are optimized and non-hero scenes remain lazy. Timer-driven demos and pointer loops can leak work, so cancellation and convergence are explicitly tested.

# Migration Plan

Add and validate the OpenSpec delta; implement shared terminal and scene primitives; recompose landing, delivery, and administration surfaces; generate, inspect, optimize, and register responsive artwork; update copy, documentation, tests, and embedded bundles; run full verification; archive the change into baseline specifications; and create one conventional commit on the numbered feature branch.

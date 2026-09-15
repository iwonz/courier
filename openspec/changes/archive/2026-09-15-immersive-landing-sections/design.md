# Context

Courier already has a contract-backed route explorer and a coherent Relay illustration family. The next iteration should improve pacing and viewport use without duplicating route or installation data, introducing a JavaScript carousel, or weakening keyboard and reduced-motion behavior.

# Decisions

## The landing opts the document into snap scrolling

The landing document carries a dedicated root class that enables mandatory vertical snap alignment without affecting embedded data or administration surfaces. The hero and the three navigable content sections each occupy at least one `100svh` viewport and account for the persistent masthead through shared spacing tokens. Anchor navigation targets the section starts across the shadow boundary. Short viewports may scroll bounded content inside a section where displaying every command or option at once is impossible.

## The masthead separates navigation by alignment, not rules

The logo and section navigation form the left cluster in routing, installation, and CLI order. GitHub and preference controls form the right cluster. There is no bottom border and no separator between GitHub and the selectors. Mobile retains the same order while allowing the clusters to fit the available width.

## The hero is an interactive scene, not a call-to-action panel

The hero contains only the product promise and a full-bleed Relay dispatch scene. The illustration is masked and color-blended into the canvas instead of sitting in a framed card. Pointer movement changes a local spotlight and route position; activation sends a visible route pulse. Focus and Enter/Space provide the same activation, and reduced-motion users receive an immediate state change without animation. Installation commands and release buttons move entirely out of the hero.

## Every slide owns wide and portrait full-bleed scenes

Dispatch, cartography, field installation, and integrity inspection art fills the hero, route, installation, and CLI slides respectively. Each role has a separately composed wide and portrait image selected by a shared responsive `<picture>` component; mobile is not an automatic crop of desktop art. Quiet headline regions and equipment bays are built into each composition so bounded functional panels appear installed in the scene. Gradient and mask layers preserve contrast and connect the art to the current light or dark canvas. Pointer position may shift a non-essential scene layer and signal glow; reduced motion removes the shift. The illustrations are never framed as a separate image tile; only functional route, installation, and reference controls receive bounded surfaces.

## Dense sections expose one focused surface

The route explorer retains its source/destination model and contract projection but stretches a focused control surface across the free area of its background. Installation becomes an icon-led channel chooser with one width-safe command readout plus repository-owned native/direct download actions. Commands and options remain a single section with bounded columns over an integrity-inspection background so their full generated data remains reachable without extending the snap slide.

## Locale flags use emoji; the product mark uses generated Relay art

The generic segmented option supports a decorative symbol in addition to a shared icon. Locale options render native flag emoji while retaining localized names, radio semantics, persistence, and keyboard behavior. The shared Courier brand and favicons use a new transparent, text-free Relay head mark produced through the built-in image-generation workflow; the wordmark remains live text.

# Risks / Trade-offs

Mandatory snap can be uncomfortable when content exceeds a short viewport, so only the landing owns it, every section provides a bounded internal overflow region when necessary, and reduced-motion disables smooth movement. Emoji appearance varies by operating system, which is acceptable because accessible localized names remain authoritative. Generated raster art is less resolution-independent than SVG, so the mark is generated at a high square resolution and displayed at small sizes with local optimized derivatives.

# Migration Plan

Add and strictly validate the OpenSpec delta; generate and record the compact mark plus wide/portrait scene pairs; extend the shared symbol and responsive mascot options; recompose the landing and catalogs; update unit and browser acceptance; rebuild embedded assets; complete headed desktop/mobile visual QA; run `make verify`; archive this change into the baseline; and create one matching conventional commit on the numbered feature branch.

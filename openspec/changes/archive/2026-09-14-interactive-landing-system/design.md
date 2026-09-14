# Context

Courier's illustrated identity is established, but the landing still reads as a sequence of separate inventories. The page should answer three questions in order: how to install Courier, which endpoints can be connected, and which commands or options are available. Compact controls and a recognizable Relay mark should support that hierarchy without weakening accessibility.

# Decisions

## Preference controls are icon-only radiogroups

The generic segmented control gains an icon option and a visually-hidden group-label mode. Theme uses system, sun, and moon symbols. Locale uses compact language badges built in the shared icon system. Each button keeps its localized accessible name, selected state, roving tab stop, arrow/Home/End behavior, and visible focus ring. The group keeps its localized programmatic label even though no `Theme` or `Language` legend is visible.

## The compact mark is Relay

The brand component and the three standalone SVG lockups use one simple head-and-capsule Relay silhouette. It preserves the graphite, paper, signal-lime, and beak-orange identity at favicon size and replaces the source-node-arrow symbol. This is repository-native vector work, so it does not add generated raster provenance.

## One route explorer projects the contract

The landing computes individual source/destination pairs from `contractData.routes`. Source and destination endpoint buttons are keyboard-operable and respond to focus, hover, and activation. Selecting a source exposes only supported destinations; selecting a valid pair shows its `courier from <source> to <destination>` example and permitted flags. HTTP(S) remains a destination presentation for the contract's `http` endpoint, while the four requested conceptual families remain local, remote, web, and webhook. No route or flag is handwritten outside presentation-only example syntax for each endpoint kind.

## Commands and options share one reference surface

Commands and flags remain generated contract data but are rendered within one section. Command rows lead into a compact option registry in the same bounded surface. The full generated reference remains reachable from that section and GitHub remains reachable from the masthead; a separate documentation block is unnecessary.

## Installation is compact and width-safe

Installation channels carry shared platform/package-manager glyphs, use compact rows rather than large generic panels, and give commands the flexible grid column. Long commands scroll inside their own code region while the document stays within the viewport. Native Linux packages and direct binaries are first-class compact channels rather than prose-only cards.

# Risks / Trade-offs

Icon-only controls need stronger accessible naming and less visible onboarding than text controls, so localized tooltips and programmatic labels remain available. Hover cannot be the only route interaction, so the same state changes occur on focus and activation. The route projection must handle future endpoint and route additions without silently inventing examples; unknown endpoint kinds receive contract syntax as their fallback representation.

# Migration Plan

Add the OpenSpec delta, extend the shared icons and segmented control, update and checksum the SVG identity assets, recompose the landing and catalogs, update documentation and tests, rebuild embedded web artifacts, complete headed desktop/mobile visual QA, run the repository verification gates, archive this change into the baseline, and create one conventional commit on the numbered feature branch.

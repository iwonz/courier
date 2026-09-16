# Context

Courier already has a contract-backed route matrix, a shared Lit UI kit, responsive full-slide scenes, and a masked refraction layer. This change must improve visual integration without duplicating product capabilities in presentation-only data or introducing remote runtime assets.

# Decisions

## Route vocabulary is role-aware presentation

The generated contract remains authoritative. Landing view models map the `ssh` endpoint to `Remote`, incoming `webhook://` to source `Web Hook`, and outgoing `http(s)://` to destination `Web Hook`. `Web` retains its scheme name but receives direction-specific iconography and copy. Route examples and applicability continue to come from generated contract data.

The source selector contains Local, Remote, Web, and Web Hook. The destination selector contains Local, Remote, Web, and Web Hook. The outgoing HTTP(S) contract endpoint is folded into the visual destination Web Hook choice; it is not advertised as a new URI scheme.

## Pointer motion is filtered by elapsed time

The scene stores target and rendered pointer coordinates separately. A single animation loop converges through exponential smoothing derived from elapsed milliseconds, clamps long frame gaps, and stops beneath a small distance threshold. Pointer leave eases toward the ambient position before stopping. The base image never receives a pointer-driven transform. Coarse pointers and reduced motion bypass the loop.

## Artwork is composed around controls

Four new wide and four portrait raster scenes share one mature graphite, paper, and signal-lime visual language. Each places Relay and detailed props around a deliberately quiet control bay matching the corresponding layout. Raster assets contain no text, logos, or fake controls; live HTML remains authoritative and accessible. New versioned filenames are selected only after visual inspection and their dimensions, hashes, prompts, roles, and authorship are recorded.

## Controls use compact operational geometry

Route endpoint choices become horizontally economical chips with a strong terminal node, role-direction icon, and short label. The selected connector attaches to those nodes. Installation choices become wrapping tag buttons above a stable command surface. Fixed grid tracks and bounded overflow preserve panel geometry across choices and viewports.

## Masthead belongs to the scene

The masthead uses translucent graphite glass with a soft lower fade and light controls. It remains measured, safe-area-aware, permanently visible, and section-aware. Every absolute off-landing link uses a new browsing context with `noopener noreferrer` so external navigation does not replace the landing page or expose the opener.

## Installation commands are directly copyable

The stable command readout includes one shared-icon copy button. Clipboard access is attempted only after activation, the original command remains selectable, and a bounded live status communicates success or failure without resizing the panel. The behavior is isolated behind an injected helper so denied and unavailable clipboard paths are deterministic under test.

# Risks / Trade-offs

Additional raster assets increase repository and Pages payload size, so selected WebP files are optimized and responsive sources remain mandatory. Continuous interpolation can waste animation frames if termination is imprecise, so convergence and cleanup are unit tested. Presentation aliases can become misleading if detached from the route matrix, so the mapping is centralized and every visual choice resolves back to a shipped contract endpoint.

# Migration Plan

Add and strictly validate the OpenSpec delta; generate and inspect responsive artwork; implement smoothed scene coordinates and role-aware endpoint view models; recompose masthead, route, installation, and slide layers; add command copying and external-link protection; update localized copy, provenance, tests, and documentation; complete visual and full verification; archive the change into baseline specs; and create one matching conventional commit on the numbered feature branch.

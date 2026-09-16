# Context

Courier already ships four purpose-composed responsive landing scenes and a generated CLI registry. This iteration must improve precision and stability without adding another scene inventory, duplicating contract data, or requesting remote assets at runtime.

# Decisions

## The masthead is measured and always visible

The landing uses `viewport-fit=cover`, safe-area padding, and a `ResizeObserver` to publish the masthead's actual block size to slide spacing. A light paper-tinted glass fade replaces the canvas-colored dark overlay. An `IntersectionObserver` marks the active navigation target without hiding or collapsing the masthead.

## Scene motion is local, not positional

The base scene never translates or scales in response to input. A shared scene component renders the responsive base picture plus a cached duplicate masked into an irregular cursor-centered lens. The duplicate scales very slightly around the cursor while layered glows produce an amorphous highlight. Pointer updates are scheduled through one animation frame. Touch uses a fixed ambient treatment and reduced motion removes lens morphing and travel.

## Source and Destination are roles

`Source` and `Destination` remain literal technical role labels in both locales. `Local`, `SSH`, `Web`, `Webhook`, and `HTTP(S)` remain canonical endpoint type values. The hero has no control semantics or activation state; a decorative Bezier path continuously communicates the handoff.

## Functional selection requires activation

Route endpoints and installation channels update only from button activation. Hover and focus style the control without changing the current choice. Invalid destinations are native disabled buttons. The route connector measures the selected button edges and draws a cubic Bezier path without participating in layout.

## Dynamic data occupies invariant tracks

Every working section uses an `auto minmax(0, 1fr)` slide grid. Full-width panels use fixed internal tracks and bounded scrolling so flag counts, command lengths, channel choices, and filtered results cannot change panel bounds or move adjacent controls.

## Brand geometry is authoritative and local

A shared monochrome brand-icon component renders pinned Simple Icons geometry plus official PowerShell and Scoop geometry where the collection has no entry. Marks inherit Courier color, carry accessible names only when requested, and are bundled into the compiled UI. Notices record upstream source, license, attribution, and trademark constraints.

Consumer Vite builds preserve their existing stable asset names while assigning source-derived names to emitted SVG marks. This prevents same-extension brand assets from swapping numbered filenames between the production build and isolated embedded-asset verification.

## CLI compatibility comes from generated command flags

The command registry starts with no selection and all flags visible. Command rows toggle selection. A custom checked checkbox becomes enabled when a command is selected; while enabled it filters by that command's generated `flags` array. Commands with no flags render a localized empty state. Route-specific applicability remains confined to the route explorer.

# Risks / Trade-offs

Duplicating the current scene inside a masked lens adds compositing work, so the effect is limited to fine pointers, uses the browser image cache, and is disabled under reduced motion. Official marks carry their upstream licensing and trademark obligations even though Courier remains MIT, so the shipped notice is authoritative and no brand asset is represented as Courier-owned.

# Migration Plan

Add and strictly validate the OpenSpec delta; implement shared scene, brand-icon, checkbox, and Bezier helpers; recompose the landing state and fixed grids; update catalogs, notices, and documentation; maintain exact coverage and browser geometry assertions; run the complete verification pipeline; archive the change; and create one matching conventional commit on the numbered feature branch.

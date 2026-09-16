# Change: Integrate landing scenes and route controls

## Why

The landing uses full-viewport illustrations, but the current artwork and dark instrument panels read as separate layers rather than one product scene. Pointer refraction follows raw coordinates and visibly snaps, route nodes lack visual authority, endpoint labels do not explain their source/destination roles, and route and installation choices consume more vertical space than their information requires. The light masthead also conflicts with the graphite operational identity.

## What Changes

- Replace all four landing scenes with purpose-composed wide and portrait artwork whose quiet regions are designed around each section's real controls.
- Smooth the fine-pointer lens with time-based interpolation while leaving the base scene stationary and preserving touch and reduced-motion fallbacks.
- Replace flattened hero labels with crisp, semantic Source and Destination route points.
- Present SSH as `Remote`, incoming and outgoing webhook behavior as `Web Hook`, and use role-specific endpoint icons and explanations without changing the public CLI contract.
- Compact route choices and installation channels into precise tag-like controls with stable section geometry.
- Restyle the masthead as dark, scene-integrated glass, add an accessible command-copy action, and open every off-landing link in a safe new tab.

## Impact

This change affects responsive landing artwork, the shared scene primitive, landing composition, localized presentation copy, documentation, identity-asset provenance, and unit/browser acceptance. It does not alter endpoint parsing, route availability, CLI syntax, contract version, release channels, or runtime transfer behavior.

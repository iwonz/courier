## Why

The landing still exposes an unnecessary visible shell-syntax label, uneven nested padding, and a button-shaped Copy action that make the command surfaces feel stepped rather than aligned. The route illustration also reads as a separate heavy block instead of supporting the localized product promise.

## What Changes

- Remove the visible shell-syntax label while retaining an accessible name for the POSIX/PowerShell selector.
- Align shell selection, command text, and Copy on one inline axis with even surface padding.
- Present Copy as an icon-led text action whose hover state changes the link treatment without adding a background.
- Scale the decorative route illustration to roughly half the available hero width, soften its edges and opacity, and allow it to sit partially behind the headline.
- Extend regression coverage and document the revised composition.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-landing`: Refine the command-surface alignment and the responsive relationship between the headline and route artwork.

## Impact

This change affects the shared command readout presentation, the project landing layout, shared CSS, generated embedded browser CSS, landing documentation, and UI regression coverage. It does not change the CLI contract, copy behavior, route selection, image assets, browser APIs, or publication workflow.

# Design: 8-bit Courier identity system

## Relay becomes a courier mascot

Relay is a blocky pigeon courier with feathers, a satchel, a parcel, and one small earpiece. The character has no armor, helmet, visor, glowing face panel, exoskeleton, metallic chest plate, police equipment, or smooth 3D shading. A separately generated square portrait serves as product chrome; transparent route, delivery, and administration poses supply optional orientation without carrying required meaning.

The canonical neutral sprite is generated first with the built-in image workflow. The mark and role sprites are identity-preserving derivatives that use it as their explicit reference; independent redraws are rejected. Anatomy, proportions, plumage, eye, beak, earpiece, and satchel stay fixed, while only pose, role equipment or clothing, and carried or attached objects may vary. All accepted assets are transparent, hard-edged pixel art, encoded losslessly, rendered with pixelated sampling, recorded in provenance, and capped by role-specific budgets.

## Modern pixel grammar

The shared UI owns a twelve-color light/dark palette, a locally pinned Cyrillic Pixelify Sans display font, a four-pixel geometry unit, chamfered semantic controls, crisp dividers, hard offset interaction shadows, pixel-grid icons, and stepped decorative motion. Body text remains system sans and operational values remain system monospace.

Ordinary layout stays transparent. Pixel frames are reserved for controls, fields, alerts, selection, focus, status, and boundaries that explain separate behavior. This preserves the borderless flow while replacing smooth gradients, blur, soft shadows, pill shapes, and rounded browser chrome.

## Icons and route geometry

First-party actions and endpoint concepts use code-native sixteen- or twenty-four-pixel grids. Third-party channel marks are reviewed pixel-grid derivatives of pinned official sources and retain source, license, attribution, and trademark records. Locale buttons use local pixel flags rather than platform emoji.

The measured cubic route remains the positioning source. A pure helper samples that curve and snaps the samples to the shared four-pixel grid; consumers render the resulting crisp polyline and a step-timed packet. Reduced motion removes the packet animation without changing route state.

## Surface loading

Brand owns only the compact mark. A role-aware Relay sprite selects the neutral, route, delivery, or administration asset, and each application renders only its role. Landing requests the mark and route sprite initially; README uses neutral; delivery and administration use their dedicated sprites. No generated background, panorama, remote font, remote icon, or external runtime request is introduced.

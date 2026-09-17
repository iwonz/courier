# Design

## Identity

Relay returns as a deliberately manufactured technological pigeon mascot: a recognizable pigeon silhouette built from simplified graphite and cobalt body volumes, segmented wing plates, an integrated delivery module, restrained cyan routing light, and a single orange waypoint beacon. It is less like a living bird and more like a mature product character, without becoming cute or toy-like. One transparent ImageGen asset is the canonical mark and hero character.

The visual system uses a calm neutral canvas, deep navy type, restrained cobalt surfaces, cool border light, cyan operational accents, and orange only for active delivery state and focus. Relay is composited over interface-owned route geometry; there is no cyberpunk city, scenic background, page panorama, or generated UI. Generated art contains no text, interface, logos, credentials, or required information.

## React and shadcn architecture

`@courier/ui` becomes a React component library. Its repository-owned components follow the shadcn composition model: semantic React components, Radix primitives where behavior requires them, Tailwind utilities, `class-variance-authority`, and the shared `cn` helper. Courier owns the generated source and can update it without runtime access to shadcn infrastructure.

Landing, delivery, and administration are independent React roots composed from the same library. They do not wrap React in Lit or retain parallel custom-element implementations. Vite compiles Tailwind locally, and the Go embed pipeline continues to receive static self-contained assets.

Shared preferences use React context and one browser controller. Theme cycles `system → light → dark`; locale cycles English → Russian. Browser defaults remain authoritative until a stored explicit choice exists.

## Product surfaces

The landing remains three ordinary document blocks: route hero, installation, and CLI reference. The hero composes the transparent Relay mascot with code-native route lines, nodes, and restrained interface fields. Route controls, install channels, command selection, compatibility filtering, and copy feedback use shadcn controls and keep contract-backed behavior.

The delivery application renders only generic authentication content before authorization. Authorized manifests, uploads, downloads, and navigation use cards, buttons, inputs, badges, separators, and progress components without changing endpoints or CSRF behavior.

The administration application uses cards, tabs, selects, checkbox, badges, alerts, and scroll areas for metrics, server/delivery navigation, policy editing, refresh, and stop actions. SSE snapshots preserve the selected delivery when possible.

## Verification

Tests exercise shared shadcn variants and semantics, preferences, every product interaction, protected-data isolation, and responsive geometry. Coverage includes TSX and remains exact. Browser acceptance verifies local-only assets, no horizontal overflow, three natural-height landing sections, theme/locale persistence, route and install selection, copying, delivery authentication, upload/download behavior, administration policy changes, stop actions, and SSE updates.

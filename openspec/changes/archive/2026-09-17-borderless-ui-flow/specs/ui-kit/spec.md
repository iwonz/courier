## MODIFIED Requirements

### Requirement: Intentional control appearance

Courier SHALL present theme, locale, form, file, policy, navigation, and command controls through one shared lightweight shadcn visual grammar. Ordinary layout wrappers SHALL be transparent and SHALL NOT add card backgrounds, outer radii, shadows, blur, or borders. Explicit chrome SHALL be reserved for interactive controls, form inputs, alerts, focus, selection, destructive emphasis, status, or a divider that explains separate scrolling or interaction regions.

#### Scenario: A user operates a Courier control

- **WHEN** the control is rendered, focused, selected, disabled, or activated with a keyboard or pointer
- **THEN** it uses Courier tokens and visible semantic state while preserving the expected role, accessible name, focus order, and change behavior without an unnecessary surrounding border

#### Scenario: A product surface groups related content

- **WHEN** content is grouped for layout without its own interactive state
- **THEN** it remains on the continuous document canvas and uses spacing, typography, alignment, or a meaningful hairline divider instead of a card island

### Requirement: Shared workbench primitive

Courier SHALL use repository-owned shadcn layout and control primitives for structured product state without a separate terminal or workbench component. Card SHALL remain a semantic structural wrapper with no default surface chrome; command readouts, route displays, and tab lists SHALL likewise avoid a containing panel. The system SHALL use normal sans-serif content without viewport-sized framing, shell prompts, fake execution, or arbitrary command input.

#### Scenario: A product surface presents state

- **WHEN** landing, delivery, or administration renders structured content
- **THEN** shared shadcn composition presents it directly on the document canvas with readable hierarchy and no ordinary card islands

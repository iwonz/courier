## Why

The terminal landing still carries decorative metrics, rules, and repeated labels that compete with its primary route, installation, and command workflows. The remaining surfaces need a quieter hierarchy with stable channel selection and a completely unframed command builder.

## What Changes

- Remove the endpoint, route, and command count strip from the route console.
- Remove visual rules that divide the masthead, landing sections, and route-console regions.
- Give installation tabs immediate, stable active and hover states without background flashing.
- Stop repeating the selected installation channel below its command.
- Remove all borders and divider rules from the command-and-parameter builder while preserving its layout and behavior.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-landing`: Simplify the landing hierarchy, installation chooser, and command builder presentation without changing their functional contracts.

## Impact

The change affects the landing React composition, presentation classes, focused unit/browser assertions, and landing acceptance documentation. It does not change CLI data, routes, installation commands, copy behavior, preferences, or backend APIs.

## Why

The landing currently presents two independent command outputs: a route example beside Source/Destination and a separate command builder after installation. They can disagree and make the route controls feel disconnected from the copyable CLI command. Installation also repeats two links to the same release destination beneath its command.

## What Changes

- Combine the route selector and the contract-backed CLI builder into one workspace with one generated readout and Copy action.
- Select `courier from` initially, fill its arguments from the selected route template, and show route controls only for that command.
- Keep the command list responsive and preserve all argument, option, shell, and copy behavior.
- Remove the route-console eyebrow, put one direct-binaries link beside the installation tabs, and remove the duplicate Linux-packages link.
- Update landing documentation and acceptance for the two-section flow.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-landing`: unify route templates with command construction and simplify installation navigation.
- `cross-platform-acceptance`: verify the unified workspace, link placement, and exact copy behavior in the browser.

## Impact

Only landing presentation and client-side state change. The CLI contract, generated data, backend APIs, install commands, release URLs, Relay assets, and other web surfaces stay unchanged. The release is a patch version through the existing one-command ship flow.

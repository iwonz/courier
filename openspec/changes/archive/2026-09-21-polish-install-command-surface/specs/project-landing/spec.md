## ADDED Requirements

### Requirement: Transparent labeled-only landing workspace

The landing SHALL place the headline, Relay, and unified command workspace directly on the theme canvas without a tinted wrapper or sprite-cell fill. Its command selector and both immutable readouts SHALL omit the redundant visible `Command` label while retaining one Copy action per readout. The masthead SHALL vertically center the `COURIER CLI` wordmark on the Relay mark's axis.

#### Scenario: A visitor scans the command workspace

- **WHEN** the hero, command selector, and command readout render in either theme
- **THEN** the hero and Relay cell have transparent backgrounds, the command form remains functional, and no `Command` heading precedes the selector or appears beside Copy

#### Scenario: A visitor scans the installation readout

- **WHEN** an installation channel is selected
- **THEN** its exact command and Copy action render without a repeated visible `Command` label

#### Scenario: A visitor scans the masthead

- **WHEN** the compact brand renders
- **THEN** the mark and wordmark share a horizontal center axis within one pixel

## MODIFIED Requirements

### Requirement: Viewport-aware installation chooser

Installation SHALL keep horizontally scrollable channel tabs and one immutable command readout. A single direct-binaries link with a leading first-party download icon SHALL remain outside the tab scroller to its right, including at narrow widths. The redundant Linux-packages link SHALL not render. The npx tab SHALL show the official local npm mark beside its npx label; wget remains text-only.

#### Scenario: Installation is narrow

- **WHEN** installation renders at 320 or 390 pixels
- **THEN** channels remain scrollable while the direct-binaries link and its icon stay visible to their right without document overflow

#### Scenario: A visitor follows a direct download

- **WHEN** the direct-binaries link is activated
- **THEN** it opens the current GitHub Releases latest destination in a new protected tab

#### Scenario: Third-party branding is inspected

- **WHEN** installation and GitHub marks render
- **THEN** supported brands use local official-color PNGs without pixelated sampling, npx shares the npm mark, and wget stays text-only

#### Scenario: A user inspects an installation channel

- **WHEN** a channel tab receives pointer, hover, or keyboard activation
- **THEN** its exact command replaces the prior selection with stable active and hover states, no flash, and no repeated channel label below

#### Scenario: A user copies an installation command

- **WHEN** Copy is activated for the selected channel
- **THEN** the exact immutable command is written to the clipboard and localized status is announced

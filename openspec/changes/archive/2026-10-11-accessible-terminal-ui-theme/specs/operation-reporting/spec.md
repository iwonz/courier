## MODIFIED Requirements

### Requirement: Terminal-capability runtime presentation

Courier SHALL use the shared compact panel and table presentation for runtime commands only when their target output stream is a real TTY, SHALL inherit the terminal's configured foreground and background for textual surfaces, SHALL use color only as a redundant semantic signal, SHALL disable color for `NO_COLOR` or `TERM=dumb`, and SHALL sanitize dynamic terminal values. Redirected output SHALL retain deterministic plain fields without ANSI or cursor-control sequences.

#### Scenario: Runtime output is redirected

- **WHEN** transfer, hosted, UI, servers, update, or version output is piped or redirected
- **THEN** Courier emits its stable plain representation without borders, color, or cursor movement

#### Scenario: Dynamic text contains terminal controls

- **WHEN** a path, remote error, release note, or other dynamic value contains escape or unsafe control sequences
- **THEN** the interactive renderer removes those sequences before layout

#### Scenario: Runtime output uses a light terminal theme

- **WHEN** Courier renders an interactive panel or table over a light terminal background
- **THEN** textual content inherits the terminal foreground and background without dark filled surfaces, white reset gaps, or color-dependent meaning

#### Scenario: Runtime output is narrow or multiline

- **WHEN** a panel contains release notes, a long URL or path, or the detected terminal is narrow
- **THEN** Courier wraps complete values with stable indentation and does not paint or pad unused terminal columns

#### Scenario: Color is disabled in an interactive terminal

- **WHEN** `NO_COLOR` disables ANSI color while the output remains interactive
- **THEN** panel titles, semantic markers, labels, values, and table headings remain structurally distinguishable

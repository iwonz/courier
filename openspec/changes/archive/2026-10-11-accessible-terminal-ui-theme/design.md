# Design: Theme-neutral terminal presentation

## Context

The rich renderer currently composes an outer style with a fixed dark background and inner label/value styles. ANSI reset sequences emitted by an inner style also reset the outer background. Dark terminals mask this artifact; light terminals show default-background gaps as bright rectangles. Fixed white body text and hard-coded RGB surfaces also ignore user terminal themes and accessibility preferences.

## Decisions

### Inherit the terminal canvas

Panels and tables will not emit foreground or background colors for body text, values, padding, or empty cells. The terminal's configured foreground/background pair remains authoritative. The renderer will not attempt an OSC background query and will not add a theme flag.

### Use color only as redundant semantics

Each panel/table title receives a visible marker and bold uppercase text. The marker uses the terminal ANSI palette for info/success/warning/failure. Labels and values use the default foreground. Thus title text, markers, status words, and existing footer symbols preserve meaning with color disabled or with a custom terminal palette.

### Compact panels without painted whitespace

Panels render only occupied lines rather than a terminal-width rectangle. Fields align in sufficiently wide terminals; narrow layouts put the label above an indented value. Multiline and long values are wrapped to the detected width with continuation indentation. No line is padded merely to hold a background.

### Theme-neutral tables

Tables retain borders, bold headers, wrapping, and complete values, but rows and headers have no filled backgrounds or forced foregrounds. Titles use the same marker/title treatment as panels.

### QR exception

Color QR rendering keeps explicit black and white cells and its quiet zone. A QR code requires controlled luminance contrast to remain machine-readable; it is separate from the textual panel surface and already has an ANSI-free `NO_COLOR` form.

## Testing

- Assert colored panels/tables contain no background-color escape sequences and no fixed RGB foreground escapes.
- Snapshot compact wide, multiline, long-value, and narrow panel layouts without color.
- Verify every semantic tone remains distinguishable by a marker without color.
- Retain sanitization, plain-output, table completeness, QR, runtime command, and exact-coverage tests.

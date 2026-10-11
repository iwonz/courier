# Design: Deterministic cross-platform renderer coverage

## Context

Wide interactive panels align labels and values. When a labeled value exceeds the remaining width, the first wrapped line follows the label and each continuation is indented to the value column. Existing coverage exercised an unlabeled multiline body and a narrow stacked field, but not this aligned continuation branch. Platform-specific source sets made that single missed statement visible only in the Windows aggregate release gate.

## Decision

Extend the wide-layout renderer test with one long labeled value that always exceeds the calculated value width at the fixed 72-column fixture size. Verify the continuation text begins at the same value column, the entire sanitized value remains present after layout whitespace is removed, and the panel remains within its width budget without fixed background colors.

No production branch, coverage threshold, generated artifact, or release workflow is changed. The test closes the gap rather than weakening the exact-coverage gate.

## Testing

- Run the `internal/terminalui` package with a coverage profile and require 100.0 percent statements.
- Run focused terminal UI, report, and application tests plus strict OpenSpec validation.
- Run the complete `make verify` release-candidate gate before shipping the next patch.

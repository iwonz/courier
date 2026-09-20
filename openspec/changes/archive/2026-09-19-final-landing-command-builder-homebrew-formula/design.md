# Design

## Landing chrome and brand media

The masthead is one static `h-16` flex row in normal document flow. It contains the Courier brand on the left and GitHub, theme, and locale actions on the right, with no internal navigation. Third-party marks are transparent local PNG files rendered by `img`; channels without an independent official mark use text only.

The neutral Relay v1 remains the identity source. Mark, route, delivery, and admin are identity-preserving v2 derivatives with recorded prompts, hashes, dimensions, alpha validation, consumer links, and a light/dark/checkerboard contact sheet.

## Structured command builder

The canonical contract owns command paths, ordered arguments, parameter value kinds, choices, placeholders, dependencies, conflicts, repeatability, and defaults. The generated landing projection preserves this order. UI state stores arguments and explicit parameter values separately, clears it on command changes, and produces tokens without parsing the display-only usage string.

Validation is deterministic. Required arguments, UUIDs, unsigned numbers, enum choices, and dependencies gate Copy. Defaults are hints only. Enabling a conflicting value clears its peer; `--all` clears the stop UUID and entering a UUID clears `--all`. Repeatable values preserve row order.

POSIX quoting uses single-quoted tokens and the standard `'"'"'` sequence for embedded apostrophes. PowerShell uses single-quoted tokens with doubled apostrophes. Safe tokens remain unquoted in both modes.

## Homebrew release flow

GoReleaser publishes a deterministic `courier_<version>_source.tar.gz` plus checksums and no cask. One renderer owns both snapshot and tagged Formula output. The Formula depends on Go at build time, disables CGO, injects the release version, commit, and date, installs `./cmd/courier`, and verifies `courier version`.

After the GitHub Release exists, a macOS job reads the published source checksum, renders the Formula, runs `brew audit`, installs it with `--build-from-source`, executes the installed binary, and only then commits the exact verified Formula to `main`. Release verification requires the source archive, Formula, and Scoop manifest and rejects cask expectations.

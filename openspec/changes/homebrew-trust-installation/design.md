# Design: Trust the exact Courier cask

## Context

Homebrew 6 treats tap definitions as executable Ruby and refuses to load non-official items until the user grants trust. Courier's repository name also does not follow the `homebrew-<tap>` shortcut convention, so it must be registered with an explicit URL.

## Decision

The documented sequence first runs `brew trust --cask iwonz/courier/courier`, then registers `iwonz/courier` with its explicit GitHub URL, and finally installs the fully qualified cask. Item-scoped trust is narrower than `brew trust iwonz/courier` and satisfies Homebrew's default security policy without bypass flags.

GoReleaser's `homebrew_casks` publisher remains in use because it is the supported replacement for the deprecated `brews` formula publisher. The cask continues to verify the architecture-specific GitHub Release checksum.

## Verification

Start with no Courier tap or trust entry, grant cask-only trust, register the custom URL, and run `brew fetch --cask iwonz/courier/courier`. The fetch must resolve `0.1.1` and pass SHA-256 verification. Remove the temporary tap, trust entry, and cache after the test.

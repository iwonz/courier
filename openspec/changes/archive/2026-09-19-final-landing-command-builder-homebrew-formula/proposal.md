# Change: Final landing polish, command builder, and Homebrew Formula

## Why

The landing masthead, third-party marks, Relay variants, route signal, and CLI registry still contain visual inconsistencies. The CLI registry exposes syntax but cannot build a valid copy-ready command. The Homebrew cask installs an unsigned prebuilt macOS binary, which triggers Gatekeeper without an Apple Developer ID.

## What Changes

- Simplify the masthead, replace third-party pixel glyphs with local official raster marks, align the route signal to its path, and remove CLI registry dividers.
- Regenerate the four Relay role derivatives from the canonical neutral Relay while preserving identity and record visual QA evidence.
- Extend the canonical CLI contract with typed arguments and parameter values, then build exact POSIX and PowerShell commands from that structure.
- Replace the unsigned Homebrew cask with a source-built Formula generated and verified from the deterministic release source archive.
- Extend unit, asset, browser, release, and documentation acceptance for the new behavior.

## Impact

This changes the landing presentation and command-construction contract plus the Homebrew release channel. Runtime CLI behavior and the unsigned direct macOS binary channel remain unchanged. No Gatekeeper bypass or `xattr` instruction is introduced.

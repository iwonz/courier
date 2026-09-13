# Change: Document Homebrew 6 cask trust

## Why

Homebrew 6 requires explicit trust before evaluating a cask from a non-official tap. Because Courier intentionally keeps its cask in the main repository and uses an explicit custom tap URL, attempting to tap it before granting trust fails closed even though the generated cask and checksums are valid.

## What Changes

- Add the official item-scoped `brew trust --cask` step before registering the Courier repository.
- Explain why the trust is required and why Courier trusts only its cask rather than the complete tap.
- Verify the clean trust, tap, and checksum-fetch sequence against the published `v0.1.1` cask.

## Non-goals

- Disabling Homebrew tap trust.
- Asking users to trust every formula or cask in the Courier repository.
- Replacing GoReleaser's recommended Homebrew cask publisher with its deprecated formula publisher.

## Impact

- Affected spec: `install-channels`.
- Affected documentation: primary and detailed Homebrew installation instructions.

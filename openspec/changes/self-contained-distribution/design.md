# Design: Single-repository package distribution

## Context

GoReleaser Community generated valid Homebrew, Scoop, and Winget metadata for `v0.1.0`. The first two formats do not require a dedicated repository: Homebrew accepts an explicit custom tap URL, and Scoop accepts any Git repository containing a `bucket` directory. Winget's public catalog is different because `winget install <id>` only discovers manifests merged into Microsoft's repository.

## Decisions

### Courier is its own tap and bucket

GoReleaser writes `Casks/courier.rb` and `bucket/courier.json` to the `main` branch of `iwonz/courier`. Homebrew registers the nonstandard repository name with an explicit URL, while Scoop registers the same repository as a custom bucket. Initial `v0.1.0` manifests are migrated before the auxiliary repositories are removed so the channels remain usable.

### Releases use only short-lived GitHub credentials

The release job already has `contents: write` and receives GitHub's repository-scoped `GITHUB_TOKEN`. GoReleaser uses that token for the GitHub Release and the two manifest commits. The only user-managed release secret is the granular npm publication token.

### Winget is not advertised as a supported catalog

An official Winget package necessarily creates a branch in a fork and a pull request against `microsoft/winget-pkgs`. Keeping a manifest only inside Courier would not make `winget install iwonz.Courier` discoverable and would misrepresent support. Courier therefore documents the self-contained Windows channels and removes the Winget publisher.

### Verification checks the deployed version

Snapshot validation still requires both generated manifest formats. Post-publication verification reads `Casks/courier.rb` and `bucket/courier.json` from `main` and requires their version to match the release tag.

## Cleanup sequence

1. Publish and verify the in-repository manifests.
2. Close the obsolete upstream Winget pull request.
3. Delete `iwonz/homebrew-tap`, `iwonz/scoop-bucket`, and the `iwonz/winget-pkgs` fork.
4. Delete the no-longer-used Homebrew, Scoop, and Winget repository secrets.

## Risks

- Homebrew requires the explicit repository URL because the repository is not named `homebrew-courier`.
- A same-repository manifest commit advances `main` after the immutable release tag; it changes only generated catalog files and does not alter the released source.

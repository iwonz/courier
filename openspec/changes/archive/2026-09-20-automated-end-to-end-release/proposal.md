# Change: Automated end-to-end release lifecycle

## Why

Courier has strong verification and publication jobs, but finishing a change still requires separate manual OpenSpec archival, version synchronization, commit, push, tag, workflow monitoring, Formula synchronization, and Pages publication steps. A completed change needs one fail-closed command that owns that whole lifecycle.

## What Changes

- Add `make ship` as the complete OpenSpec-to-publication command.
- Require completed OpenSpec tasks, strict validation, archival, and no unrelated active changes before release.
- Synchronize the target release and generated contract artifacts before verification.
- Commit and push the verified worktree, publish the semantic tag, wait for GitHub/Scoop/Homebrew/npm verification, synchronize the Formula commit, and deploy final GitHub Pages.
- Keep `make release` for an already committed clean `main`, but make it wait for public release and Pages completion.

## Impact

The change affects maintainer automation and documentation only. Runtime, CLI, browser, and package contracts are unchanged.

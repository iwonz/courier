# Change: Publish the project landing

## Why

Courier has a deterministic landing application and a least-privilege Pages workflow, but the workflow has not reached `main` and GitHub Pages is not enabled for the repository. Maintainers also need one safe command that builds the exact synchronized revision, dispatches publication, waits for completion, and reports the public URL.

## What Changes

- Keep automatic GitHub Pages publication on every push to `main` so any landing, shared UI, contract, or documentation-link change is redeployed without manual maintenance.
- Add `make pages-build` for a deterministic local contract check, landing test, and production build.
- Add `make pages-publish` to require an authenticated GitHub CLI, clean synchronized `main`, run the local build, enable workflow-based Pages when necessary, dispatch the Pages workflow, wait for the new run, and report its URL.
- Document the automatic and manual publication paths without introducing a token, external repository, generated branch, or mutable credential in the site artifact.

## Impact

The public CLI contract is unchanged. Publication mutates only the GitHub Pages configuration and `github-pages` deployment for `iwonz/courier`. GitHub Actions continues to use repository permissions and short-lived OIDC credentials.

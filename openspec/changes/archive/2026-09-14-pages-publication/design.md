# Context

The landing already builds under `web/landing`, uses the generated shipped-command projection, and has an Actions-native Pages workflow. Remote `main` predates that work and the Pages REST endpoint currently reports no site.

# Decisions

## One build definition

`pages-build` performs the contract freshness check, installs the locked web workspace without lifecycle scripts, runs the landing's exact-coverage tests, and produces the Vite artifact. The Pages workflow calls this same Make target rather than maintaining a parallel command sequence.

## Automatic and explicit publication

Every push to `main` triggers `.github/workflows/pages.yml`. `workflow_dispatch` remains available for recovery and explicit republishing. `pages-publish` delegates to a shell script so the repository, branch, synchronization, workflow-run selection, and error handling remain testable and readable outside Make syntax.

The script refuses detached heads, feature branches, dirty trees, and local `main` revisions that differ from `origin/main`. It builds before changing remote state. It creates workflow-based Pages configuration only when the site is absent, dispatches against the exact main revision, ignores older workflow runs for that commit, waits with `gh run watch --exit-status`, and prints the Pages URL only after success.

## Credentials and repository ownership

Local publication uses the maintainer's existing `gh auth` session and never reads or writes a token value. GitHub Actions receives only `contents: read` during build and `pages: write` plus `id-token: write` during deployment. The artifact contains static files only and stays in `iwonz/courier`; no `gh-pages` branch is created.

# Risks / Trade-offs

Publishing every `main` push may rebuild when a backend-only file changes, but it guarantees the public site always corresponds to the repository head and avoids an incomplete path filter as new generated inputs are added. Concurrency cancels an obsolete deployment when a newer main revision arrives.

# Migration Plan

Add and validate the OpenSpec change, implement the shared Make build and safe publication script, use the Make build in the workflow, verify locally, archive the change, fast-forward `main`, enable workflow-based Pages, push once, and wait for the automatic deployment. Run the manual Make publication once to verify the recovery path.

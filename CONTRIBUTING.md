# Contributing

Courier is developed through small spec-first changes. Every task that modifies tracked files must complete the entire path from OpenSpec through a public patch release and the final Pages deployment. Read-only review, diagnosis, status reporting, and planning are exempt.

1. Start on clean `main` exactly synchronized with `origin/main`.
2. Run `make change-start TYPE=feat CHANGE=capability-name`. This creates both `feat/capability-name` and its OpenSpec scaffold, and refuses dirty, stale, or concurrent work.
3. Complete the OpenSpec artifacts, implement the change in English, and run relevant focused tests.
4. Mark every task complete and run `make verify`.
5. Run `make ship CHANGE=capability-name MESSAGE="feat: describe the change"`. It archives OpenSpec, selects the next patch by default, repeats verification, creates the one conventional commit, fast-forwards `main`, publishes every release surface, deploys Pages, deletes the feature branch, and verifies clean synchronized `main`.

Run `make hooks` once to activate the tracked pre-commit and pre-push guards. Direct pushes to `main` and semantic release tags are rejected outside `ship` or the documented recovery flow. Tests must use local fixtures or short-lived containers, register cleanup immediately, and leave no containers, volumes, networks, or temporary files after success, failure, or interruption.

`make test` is the fast race-enabled runtime gate. `make verify` is the required release-candidate gate and also needs Docker, OpenSpec 1.11.0, and the pinned Playwright Chromium runtime. A change is not complete at this gate: continue through `make ship`. See [Acceptance and release-candidate verification](docs/acceptance.md) and the [release runbook](docs/releasing.md).

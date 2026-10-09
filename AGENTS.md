# Courier change workflow

Any task that changes tracked files must complete the whole repository-owned lifecycle. This applies to code, tests, documentation, assets, configuration, and automation.

1. Start from a clean `main` exactly synchronized with `origin/main`.
2. Run `make change-start TYPE=<conventional-type> CHANGE=<kebab-case-name>`; do not create the feature branch or OpenSpec change by hand.
3. Complete the OpenSpec proposal, design, specs, and task list before or alongside implementation.
4. Implement the change and run relevant focused tests while working.
5. Mark every OpenSpec task complete and run the full `make verify` gate.
6. Finish with `make ship CHANGE=<name> MESSAGE="<conventional commit>"`. `ship` archives OpenSpec, derives the next patch release unless `VERSION` is explicitly supplied, regenerates release assets, repeats verification, creates the single change commit, fast-forwards and pushes `main`, publishes and verifies the release matrix, deploys Pages, removes the feature branch, and confirms clean synchronized `main`.

Do not stop after implementation, verification, a commit, a merge, or a tag. A tracked-file change is complete only after the release and final Pages deployment succeed and the repository is again on clean `main == origin/main`.

Read-only review, diagnosis, status reporting, and planning that do not modify tracked files are exempt. `make release` and `make pages-publish` are recovery commands, not alternative workflows; use them only as documented in `docs/releasing.md` after a partially completed ship.

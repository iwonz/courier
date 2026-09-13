# Contributing

Courier is developed through small spec-first changes.

1. Create an OpenSpec path under `openspec/changes/<task-name>`.
2. Create a conventional branch such as `feat/007-capability-name`.
3. Implement and test the change in English.
4. Run `make verify`.
5. Mark the task list complete, run `openspec validate --all --strict`, and archive the change into the baseline.
6. Create one conventional commit for the task, including the archived change history and updated baseline specs.

Run `make hooks` once to activate the tracked pre-commit gate. Tests must use local fixtures or short-lived containers, register cleanup immediately, and leave no containers, volumes, networks, or temporary files after success, failure, or interruption.

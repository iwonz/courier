## 1. Local publication commands

- [x] 1.1 Add one deterministic `pages-build` target for contract freshness, landing coverage, and the production artifact
- [x] 1.2 Add a guarded `pages-publish` target that enables workflow-based Pages, dispatches the synchronized main revision, waits, and reports the URL

## 2. Automatic workflow and documentation

- [x] 2.1 Keep automatic deployment on every main push and make the workflow consume the shared build target
- [x] 2.2 Document prerequisites, automatic behavior, manual recovery, permissions, and the public URL

## 3. Verification and publication

- [x] 3.1 Run landing coverage/build, workflow validation, strict OpenSpec validation, and the repository verification gate
- [x] 3.2 Archive the OpenSpec change and create exactly one matching conventional commit
- [x] 3.3 Fast-forward and push main, enable Pages, verify the automatic deployment, and exercise the manual Make publication command

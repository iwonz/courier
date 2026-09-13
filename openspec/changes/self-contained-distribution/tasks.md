# Tasks

## 1. Consolidate package metadata

- [x] 1.1 Move the current Homebrew and Scoop manifests into the Courier repository
- [x] 1.2 Configure GoReleaser to update both manifests with `GITHUB_TOKEN`
- [x] 1.3 Remove the Winget publisher and generated-manifest requirement

## 2. Simplify release automation

- [x] 2.1 Require only `iwonz/courier` and `NPM_TOKEN` in release preflight
- [x] 2.2 Verify deployed Homebrew and Scoop versions from the Courier repository
- [x] 2.3 Update installation and release documentation

## 3. Verify and clean external state

- [x] 3.1 Pass strict OpenSpec validation, workflow lint, and the complete local verification gate
- [ ] 3.2 Publish a release with in-repository manifests
- [ ] 3.3 Close the obsolete Winget pull request and delete auxiliary repositories and secrets

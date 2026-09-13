# Courier CLI implementation plan

Work is split into stacked branches. Each branch is based on the previous branch, owns an OpenSpec change, and ends with one conventional commit.

| # | Branch | OpenSpec change | Outcome |
|---|---|---|---|
| 1 | `feat/001-project-foundation` | `project-foundation` | Go/CLI foundation, MIT, architecture contract |
| 2 | `feat/002-endpoints-and-safety` | `endpoints-and-safety` | Local/remote parsing, destination semantics, path safety |
| 3 | `feat/003-transfer-and-archive` | `transfer-and-archive` | Filesystem engine, sync, atomic staging, archives, progress events |
| 4 | `feat/004-native-ssh-transport` | `native-ssh-transport` | SSH config/auth/known_hosts/ProxyJump, SFTP, remote OS detection |
| 5 | `feat/005-cli-update-reporting` | `cli-update-reporting` | `from … to …`, reporting, exit codes, self-update |
| 6 | `feat/006-distribution-and-quality` | `distribution-and-quality` | GoReleaser, installers, npm/Brew, CI, integration and cleanup tests |
| 7 | `feat/007-optional-remote-helper` | `optional-remote-helper` | Consent-gated helper fallback for genuine remote capability gaps and Windows runtime hardening |

## Definition of Done

- `go test ./...` and the race detector pass;
- first-party Go package statement coverage is 100%;
- `openspec validate --all --strict` passes;
- primary binaries build for macOS, Linux, and Windows on amd64 and arm64, with exact-platform BSD helper archives;
- tests use local fixtures or containers only and always register cleanup;
- local and remote temporary artifacts are removed after success, failure, and cancellation;
- installation and security documentation match the shipped CLI;
- one release command creates and pushes the tag, builds with GoReleaser Community, publishes the GitHub Release with notes, and publishes `@iwonz/courier` when credentials are available.

# План реализации Courier CLI

Работа разбита на накопительные ветки. Каждая ветка основана на предыдущей, имеет собственный OpenSpec change и завершается одним conventional commit.

| № | Ветка | OpenSpec change | Результат |
|---|---|---|---|
| 1 | `feat/001-project-foundation` | `project-foundation` | Go/CLI-каркас, MIT, архитектурный контракт |
| 2 | `feat/002-endpoints-and-safety` | `endpoints-and-safety` | parsing local/remote, destination semantics, path safety |
| 3 | `feat/003-transfer-and-archive` | `transfer-and-archive` | файловый engine, sync, atomic staging, archive, progress events |
| 4 | `feat/004-native-ssh-transport` | `native-ssh-transport` | SSH config/auth/known_hosts/ProxyJump, SFTP, remote OS detection |
| 5 | `feat/005-cli-update-reporting` | `cli-update-reporting` | `from … to …`, reporting, exit codes, self-update |
| 6 | `feat/006-distribution-and-quality` | `distribution-and-quality` | installers, npm/brew, cross-build, CI, integration/cleanup tests |

## Definition of Done

- `go test ./...` и race detector проходят;
- покрытие собственных Go-пакетов составляет 100%;
- `openspec validate --all --strict` проходит;
- binaries собираются для macOS, Linux и Windows на amd64/arm64;
- тесты используют только локальные fixtures/containers и регистрируют cleanup;
- временные local/remote artifacts удаляются при успехе, ошибке и отмене;
- документация установки и безопасного поведения соответствует CLI.

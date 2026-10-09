# Courier CLI

<p align="center">
  <img src="web/ui/assets/courier-relay-pixel-neutral-v3.webp" width="420" alt="Relay, Courier's compact 8-bit pigeon courier mascot">
</p>

Courier is an extensible, cross-platform CLI for safely transferring files and directories across local, SSH, browser, and webhook endpoints.

**From here to anywhere.** Courier makes the route, destination, transfer state, and verified result explicit.

```text
courier from <source> to <destination> [flags]
```

Courier uses native Go implementations for SSH, SFTP, archive creation, transfer, and progress reporting. The released executable does not require rsync, an OpenSSH client, or a system archiver.

## Install

POSIX installer with curl:

```sh
curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh
```

POSIX installer with wget:

```sh
wget -qO- https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh
```

Windows PowerShell installer:

```powershell
irm https://raw.githubusercontent.com/iwonz/courier/main/install.ps1 | iex
```

Package managers:

```sh
npm install --global @iwonz/courier
npx @iwonz/courier --help
yarn global add @iwonz/courier       # Yarn Classic
yarn dlx @iwonz/courier --help       # modern Yarn
pnpm add --global @iwonz/courier
pnpm dlx @iwonz/courier --help

brew tap iwonz/courier https://github.com/iwonz/courier && brew install iwonz/courier/courier

scoop bucket add courier https://github.com/iwonz/courier
scoop install courier/courier
```

Homebrew builds the Formula from Courier's checksummed source archive, so it does not depend on an unsigned quarantined cask. Existing cask users should run `brew uninstall --cask courier` once before the Formula command above. GoReleaser also publishes raw binaries, source and platform tar.gz archives, Windows zip archives, checksums, and deb/rpm/apk/Arch Linux packages on [GitHub Releases](https://github.com/iwonz/courier/releases). Direct macOS binaries remain unsigned. See [Installation](docs/installation.md) for the complete OS, architecture, distribution, direct-download, and verification matrix.

## Use

```sh
courier from ./report.pdf to server:/srv/inbox/
courier from root@203.0.113.10:/opt/node/data to ./backup/
courier from source-server:/opt/node/data to backup-server:/srv/data/
courier from ./data to ./backup/ --archive
courier from ./backup/data.tar.gz to ./restore/ --extract
courier from ./data to ./backup/ --exclude '*.tmp' --exclude-from ./courier.ignore
courier from webhook:// to ./inbox/ --auth basic
courier from ./report.pdf to https://example.com/hooks/courier
courier from ./data to https://example.com/hooks/courier --archive
courier servers
courier servers stop 550e8400-e29b-41d4-a716-446655440000
courier servers stop --all
courier ui start
courier ui start --background --listen 127.0.0.1:9090
courier ui stop
```

The destination rules are deterministic:

- an existing directory, a path ending in `/` or `\`, or a missing destination for a directory source receives the source under its source name;
- every other destination is the exact final path, so `courier from ./report.pdf to ./renamed.pdf` does not create `renamed.pdf` as a directory;
- an existing final path is a collision and is never overwritten or merged; unrelated entries in a destination directory are preserved;
- `--archive` transfers a verified `<source-name>.tar.gz` instead of the source tree.
- `--extract` treats the destination as an extraction root and accepts tar.gz only. It performs a complete read-only inspection before staging, rejects collisions, and preserves unrelated destination entries. The expanded-size default is `100GiB` and can be changed with `--max-extracted-size <size|unlimited>`; fixed limits of 100,000 entries, depth 64, and a 100:1 expansion ratio always apply.

Courier never deletes the source. An identical plain source/destination is a successful no-op; transformed identity and copying a directory into itself are rejected. Files are staged under private partial names and committed only into an absent final path after preflight and transfer complete.

Selection is shared by ordinary copy, archive creation, and extraction. `--exclude` uses ordered gitignore syntax, `--exclude-regex` uses Go regular expressions, and `--exclude-from` expands a local gitignore-style rule file at its exact command-line position. The rule file is read before transfer endpoints are opened.

Before transfer or server startup, Courier checks every local or SSH path. Missing directory endpoints are listed for confirmation in source-then-destination order and are created privately with mode `0700` only after every confirmation succeeds. Non-interactive use must pass `--force-source-creation`; this flag authorizes only missing-directory creation and never bypasses SSH trust, credentials, authentication, collisions, or overwrite protection. Missing extraction archives and unarchived outgoing-webhook files fail instead of being created.

Browser downloads, browser uploads, incoming webhooks, and `courier ui start` stay attached to the initiating terminal by default. `--background` prints the URL and UUID after readiness, returns immediately, and leaves the detached process running after terminal closure. Stop it with `courier servers stop <uuid>`, `courier servers stop --all`, or `courier ui stop` as appropriate; configured stop behavior, fatal worker failure, and process termination also end it.

## Command parameters

<!-- BEGIN GENERATED CLI PARAMETERS -->

### Arguments

| Command | Argument | Status | Description |
|---|---|---|---|
| `from` | `<source>` | Required | File, directory, browser upload, or webhook input to read from. |
| `from` | `<destination>` | Required | Path, browser download, or HTTP endpoint that receives the data. |
| `servers stop` | `<uuid>` | Required | Delivery or server UUID to stop; omitted when --all is used. |
| `help` | `<command>` | Optional | Command path whose detailed help should be shown. |

### Options

| Option | Status | Description | Default | Repeatable | Applies to | Requires | Conflicts |
|---|---|---|---|---:|---|---|---|
| `--archive` | Optional | Pack the source into a verified <source-name>.tar.gz before transfer. | `false` | false | Path to path, Path to browser download, Path to HTTP webhook | none | extract |
| `--extract` | Optional | Safely extract a tar.gz source into the destination root. | `false` | false | Path to path, Browser upload to path, Incoming webhook to path | none | archive |
| `--listen <host:port>` | Optional | Bind incoming browser, webhook, or administration traffic to this address. | `127.0.0.1:8080` | false | Browser upload to path, Path to browser download, Incoming webhook to path, Administration UI start | none | none |
| `--background` | Optional | Detach browser downloads, browser uploads, incoming webhooks, or ui start after printing the URL and UUID; the process survives terminal closure until stopped by courier servers stop, courier ui stop, configured stop behavior, a fatal worker failure, or process termination. | `false` | false | Browser upload to path, Path to browser download, Incoming webhook to path, Administration UI start | none | none |
| `--force-source-creation` | Optional | Create every missing directory endpoint recursively without asking; no other safety confirmation is bypassed. | `false` | false | Path to path, Browser upload to path, Path to browser download, Incoming webhook to path, Path to HTTP webhook | none | none |
| `--auth <none\|basic\|password>` | Optional | Choose no authentication, HTTP Basic authentication, or a browser password. | `none` | false | Browser upload to path, Path to browser download, Incoming webhook to path, Path to HTTP webhook | none | none |
| `--auth-attempts <N>` | Optional | Set how many failed authentication attempts are allowed before the configured action. | `5` | false | Browser upload to path, Path to browser download, Incoming webhook to path | none | none |
| `--auth-fail-action <ban\|stop>` | Optional | Ban the peer or stop the delivery when the authentication-attempt limit is reached. | `ban` | false | Browser upload to path, Path to browser download, Incoming webhook to path | none | none |
| `--limit <N>` | Optional | Limit concurrent transfers for a hosted delivery, or allow an unlimited count. | `unlimited` | false | Browser upload to path, Path to browser download, Incoming webhook to path | none | none |
| `--no-ui` | Optional | Expose only the versioned browser-delivery data API, without the web interface. | `false` | false | Path to browser download | none | none |
| `--allow-ip <IP/CIDR>` | Optional | Allow one peer IP address or CIDR; repeat the option to add more networks. | `none` | true | Browser upload to path, Path to browser download, Incoming webhook to path | none | none |
| `--exclude <pattern>` | Optional | Exclude paths with an ordered gitignore-style pattern. | `none` | true | Path to path, Browser upload to path, Path to browser download, Incoming webhook to path, Path to HTTP webhook | none | none |
| `--exclude-regex <regex>` | Optional | Exclude paths matching a Go regular expression. | `none` | true | Path to path, Browser upload to path, Path to browser download, Incoming webhook to path, Path to HTTP webhook | none | none |
| `--exclude-from <file>` | Optional | Read ordered gitignore-style exclusion rules from a local file. | `none` | true | Path to path, Browser upload to path, Path to browser download, Incoming webhook to path, Path to HTTP webhook | none | none |
| `--max-file-size <size\|unlimited>` | Optional | Set the largest accepted incoming file size, or remove the configurable limit. | `10GiB` | false | Browser upload to path, Incoming webhook to path | none | none |
| `--max-extracted-size <size\|unlimited>` | Optional | Set the maximum total expanded archive size; requires --extract. | `100GiB` | false | Path to path, Browser upload to path, Incoming webhook to path | extract | none |
| `--upload-rate <rate\|unlimited>` | Optional | Limit the aggregate upload rate, or allow an unlimited rate. | `unlimited` | false | Local path to SSH path, SSH path to SSH path, Browser upload to path, Incoming webhook to path, Path to HTTP webhook | none | none |
| `--download-rate <rate\|unlimited>` | Optional | Limit the aggregate download rate, or allow an unlimited rate. | `unlimited` | false | SSH path to local path, SSH path to SSH path, Path to browser download | none | none |
| `--all` | Optional | Stop every discovered Courier delivery and server instead of one UUID. | `false` | false | Server stop | none | none |

<!-- END GENERATED CLI PARAMETERS -->

## SSH

Remote endpoints use `[user@]host:/absolute/path`. `host` may be an alias from `~/.ssh/config`. Courier applies `HostName`, `User`, `Port`, `IdentityFile`, `ProxyJump`, SSH agent identities, and strict `known_hosts` verification through native Go libraries.

Unknown and changed host keys fail closed. Add and verify host keys out of band before transferring; Courier never uses `InsecureIgnoreHostKey` and never accepts passwords through flags. Password and encrypted-key prompts require an interactive terminal and do not echo input.

Native SFTP is always the normal transport. If an authenticated server genuinely has no SFTP subsystem, Courier explains the capability gap and asks before making any remote change. Explicit consent allows it to select an exact OS/architecture Courier helper, verify SHA-256 locally, upload through the verified SSH channel into a private temporary directory, verify SHA-256 remotely, and run the embedded SFTP server over that channel. The helper is never installed into `PATH`, never persists as a service, and is removed after success, error, or interruption. Non-interactive execution declines the fallback.

Windows uses the native OpenSSH agent named pipe when `SSH_AUTH_SOCK` is not set. See [Security](docs/security.md) for trust boundaries, bootstrap details, and cleanup guarantees.

## Update

```sh
courier update
```

The updater compares the installed semantic version with the latest GitHub Release, selects the target archive, downloads it privately, verifies its SHA-256 entry from `checksums.txt`, extracts only the Courier executable, and replaces the installation through same-directory staging. Windows uses a staged handoff process so the running executable exits before replacement, followed by a cleanup process that removes the handoff binary.

## Exit codes

| Code | Meaning |
|---:|---|
| `0` | Success |
| `2` | CLI syntax or command error |
| `10` | SSH connection, trust, or authentication error |
| `20` | Preflight, archive, transfer, commit, or cleanup error |
| `30` | Self-update error |
| `40` | Registry, IPC, worker, or local control error |
| `130` | Interrupted or canceled operation |

Failure output includes the stage, sanitized reason, and separate read, sent, and confirmed byte counts. Success uses only confirmed bytes. Credentials, URL user-info, and query secrets are never included. See [Operational reporting](docs/operational-reporting.md) for the accounting, progress, history, and exit-code contract.

## Develop

Go 1.25 or newer and Node.js 24 or newer are required for source, npm-package, and shared UI checks. The full release-candidate gate also requires Docker, OpenSpec 1.11.0, and the pinned Playwright Chromium runtime. GoReleaser Community is bootstrapped locally at its pinned checksum-verified version.

```sh
make hooks       # activate the repository pre-commit and protected-push guards
make change-start TYPE=feat CHANGE=my-change # clean main → feature branch + OpenSpec scaffold
make test        # formatting, vet, race detector, exact Go coverage, and compiled runtime checks
make verify      # full browser/package/platform-ready release-candidate dry run
make pages-build # verify and build the current static landing
make ship CHANGE=my-change MESSAGE="feat: finish my change" # archive → verify → commit → merge → next patch → packages → Pages
```

Development is spec-first: every tracked-file change starts through `make change-start`, owns one OpenSpec change and conventional feature branch, and ends only after `make ship` has published the release and Pages and restored clean synchronized `main`. Read-only review, diagnosis, and planning are exempt. Low-level release and Pages targets are recovery-only. See [CONTRIBUTING](CONTRIBUTING.md), the [implementation plan](docs/implementation-plan.md), and the [release runbook](docs/releasing.md).

The machine-readable [CLI contract](docs/cli-contract.yaml) is the source of truth for shipped and planned commands. Its generated [command reference](docs/cli-reference.md) is checked against the live Cobra tree during every verification run.

Browser delivery pages, the administration interface, and the project landing page share the React/shadcn-based [`@courier/ui`](web/ui) package. Its [UI architecture guide](docs/ui.md) documents local components, assets, themes, localization, and the exact TypeScript coverage gate.

The [brand system](docs/brand.md) defines Courier's positioning, compact 8-bit Relay pigeon, pixel visual language, operational vocabulary, accessibility rules, and approved communication patterns.

See [Browser deliveries](docs/web-deliveries.md) for `web://` URLs, authentication, safe navigation, transactional uploads, and foreground/background lifecycle behavior.

See [Webhook deliveries](docs/webhook-deliveries.md) for Courier's exact incoming and outgoing multipart profile, HTTP result semantics, and intentionally unsupported provider-specific behavior.

See [Server control](docs/server-control.md) for authoritative inventory, UUID-scoped stopping, partial `--all` behavior, and the no-PID-signal trust boundary.

See [Administration UI](docs/admin-ui.md) for foreground/background startup, local API guards, live state, and optimistic policy editing.

See [Project landing](docs/project-landing.md) for contract-backed static content, shared UI behavior, verification, and repository-owned GitHub Pages deployment.

See [Acceptance](docs/acceptance.md) for local prerequisites, bounded-stream and browser coverage, labeled container cleanup, and CI publication gates.

## License

[MIT](LICENSE)

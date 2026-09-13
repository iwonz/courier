# Courier CLI

Courier is an extensible, cross-platform CLI for safely transferring files and directories between local and SSH endpoints.

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

brew tap iwonz/tap
brew install --cask iwonz/tap/courier

scoop bucket add iwonz https://github.com/iwonz/scoop-bucket
scoop install iwonz/courier

winget install --exact --id iwonz.Courier
```

GoReleaser also publishes raw binaries, tar.gz archives, Windows zip archives, checksums, and deb/rpm/apk/Arch Linux packages on [GitHub Releases](https://github.com/iwonz/courier/releases). See [Installation](docs/installation.md) for the complete OS, architecture, distribution, direct-download, and verification matrix.

## Use

```sh
courier from ./report.pdf to server:/srv/inbox/
courier from root@203.0.113.10:/opt/node/data to ./backup/
courier from source-server:/opt/node/data to backup-server:/srv/data/
courier from ./data to ./backup/ --archive
```

The destination rules are deterministic:

- an existing directory or a path ending in `/` receives the source under its source name;
- every other destination is the exact final path;
- an existing target tree is synchronized by staging the complete source, then replacing the target;
- `--archive` transfers a verified `<source-name>.tar.gz` instead of the source tree.

Courier never deletes the source. It rejects identical source/destination paths and copying a directory into itself. Files are staged under private partial names and committed only after preflight and transfer complete.

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

Failure output includes the stage, safe reason, and confirmed transferred bytes. Credentials are never included.

## Develop

Go 1.25 or newer and Node.js 18 or newer are required for source and npm-package checks. GoReleaser Community is bootstrapped locally at its pinned checksum-verified version.

```sh
make hooks       # opt into the repository pre-commit quality gate
make test        # formatting, vet, race detector, and exact 100% Go statement coverage
make verify      # tests, npm checks, and goreleaser release --snapshot --clean
```

Development is spec-first: every task owns a path under `openspec/changes`, a conventional branch name, and one conventional commit. See the [implementation plan](docs/implementation-plan.md) and [release runbook](docs/releasing.md).

## License

[MIT](LICENSE)

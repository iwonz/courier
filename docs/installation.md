# Installation

Courier release binaries are self-contained. Package managers and installers only place the correct binary; they do not add rsync, OpenSSH, tar, or another runtime dependency.

## Supported targets

| Operating system | Architectures | Primary channels |
|---|---|---|
| macOS | amd64, arm64 | Homebrew, npm ecosystem, POSIX installer, direct download |
| Linux | amd64, arm64 | deb, rpm, apk, Arch package, npm ecosystem, Homebrew Formula, POSIX installer, direct download |
| Windows | amd64, arm64 | Scoop, npm ecosystem, PowerShell installer, direct download |

Release archives are also built for FreeBSD, OpenBSD, and NetBSD on amd64/arm64 and DragonFly BSD on amd64. Courier uses these exact-platform archives for the temporary helper fallback when a BSD SSH endpoint lacks SFTP. They are published and checksummed with every release but are not currently distributed through the package-manager channels above.

Native Linux package mapping:

| Distribution | Artifact | Install command after download |
|---|---|---|
| Ubuntu, Debian | `.deb` | `sudo apt install ./courier_<version>_linux_<arch>.deb` |
| Fedora, RHEL | `.rpm` | `sudo dnf install ./courier_<version>_linux_<arch>.rpm` |
| Alpine Linux | `.apk` | `sudo apk add --allow-untrusted ./courier_<version>_linux_<arch>.apk` |
| Arch Linux, Manjaro | `.pkg.tar.zst` | `sudo pacman -U ./courier_<version>_linux_<arch>.pkg.tar.zst` |

The native package artifacts are published on each GitHub Release. They do not imply that Courier operates an apt, rpm, Alpine, or Arch repository.

## Verified bootstrap scripts

The POSIX installer uses curl when present and otherwise wget. It installs to `$HOME/.local/bin` by default and never invokes sudo. Set `COURIER_INSTALL_DIR` to choose another directory and `COURIER_VERSION` to pin a version.

```sh
curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh | sh

COURIER_VERSION=1.2.3 COURIER_INSTALL_DIR="$HOME/bin" \
  sh -c "$(curl -fsSL https://raw.githubusercontent.com/iwonz/courier/main/install.sh)"
```

The PowerShell installer uses the current user's `%LOCALAPPDATA%\Programs\Courier` directory and adds it to the user PATH. It never requests administrator access.

```powershell
irm https://raw.githubusercontent.com/iwonz/courier/main/install.ps1 | iex

$env:COURIER_VERSION = "1.2.3"
$env:COURIER_INSTALL_DIR = "$HOME\bin"
irm https://raw.githubusercontent.com/iwonz/courier/main/install.ps1 | iex
```

Both installers download the selected raw binary plus `checksums.txt`, verify SHA-256 before replacement, stage in the destination directory, and remove temporary files on success or failure.

## npm, npx, Yarn, and pnpm

All JavaScript package managers use the same public scoped package:

```sh
npm install --global @iwonz/courier
npx @iwonz/courier --help

yarn global add @iwonz/courier
yarn dlx @iwonz/courier --help

pnpm add --global @iwonz/courier
pnpm dlx @iwonz/courier --help
```

The dependency-free postinstall script selects the host tar.gz archive, verifies `checksums.txt`, and safely extracts only `courier` or `courier.exe`. Package-manager options that disable lifecycle scripts, such as `--ignore-scripts`, are not supported because they prevent native binary installation.

## Homebrew

```sh
brew tap iwonz/courier https://github.com/iwonz/courier && brew install iwonz/courier/courier
```

The explicit URL lets the main Courier repository act as its own tap. `Formula/courier.rb` verifies the deterministic `courier_<version>_source.tar.gz`, installs Go as a build-only dependency, sets `CGO_ENABLED=0`, and builds `./cmd/courier` with the release version, commit, and date. This takes longer than installing a prebuilt binary but avoids cask quarantine without an Apple Developer ID and without disabling Gatekeeper.

If Courier was previously installed through the retired cask, migrate once before installing the Formula:

```sh
brew uninstall --cask courier
brew tap iwonz/courier https://github.com/iwonz/courier
brew install iwonz/courier/courier
```

Do not use `xattr` or another quarantine bypass. Direct macOS binaries and non-Homebrew installer channels remain unsigned and are not Apple-notarized.

## Scoop

```powershell
scoop bucket add courier https://github.com/iwonz/courier
scoop install courier/courier
```

The main Courier repository acts as the custom Scoop bucket. Its `bucket/courier.json` manifest selects the amd64 or arm64 Windows zip and verifies the GoReleaser-generated hash.

To install the manifest directly without retaining a named bucket:

```powershell
scoop install https://raw.githubusercontent.com/iwonz/courier/main/bucket/courier.json
```

Winget is intentionally not a Courier catalog channel. Public Winget packages require versioned pull requests to the external `microsoft/winget-pkgs` repository, which conflicts with Courier's single-repository distribution model. Windows users can use Scoop, npm, the PowerShell installer, or a direct verified binary without that external review dependency.

## Direct download and verification

Release names use this stable pattern:

```text
courier_<version>_<os>_<arch>
courier_<version>_<os>_<arch>.tar.gz
courier_<version>_windows_<arch>.exe
courier_<version>_windows_<arch>.zip
courier_<version>_source.tar.gz
checksums.txt
```

For primary installations, `os` is `darwin`, `linux`, or `windows`; `arch` is `amd64` or `arm64`. BSD helper archives additionally use `freebsd`, `openbsd`, `netbsd`, or `dragonfly`. Download the matching artifact and checksum manifest from:

```text
https://github.com/iwonz/courier/releases/download/v<version>/
```

Verify on Linux:

```sh
sha256sum --check checksums.txt --ignore-missing
```

Verify on macOS:

```sh
shasum -a 256 -c checksums.txt
```

Verify on Windows PowerShell:

```powershell
(Get-FileHash -Algorithm SHA256 .\courier_<version>_windows_<arch>.exe).Hash
```

Compare the PowerShell output with the exact filename entry in `checksums.txt` before execution.

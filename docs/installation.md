# Installation

Courier release binaries are self-contained. Package managers and installers only place the correct binary; they do not add rsync, OpenSSH, tar, or another runtime dependency.

## Supported targets

| Operating system | Architectures | Primary channels |
|---|---|---|
| macOS | amd64, arm64 | Homebrew, npm ecosystem, POSIX installer, direct download |
| Linux | amd64, arm64 | deb, rpm, apk, Arch package, npm ecosystem, Homebrew where casks are supported, POSIX installer, direct download |
| Windows | amd64, arm64 | Winget, Scoop, npm ecosystem, PowerShell installer, direct download |

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
brew tap iwonz/tap
brew install --cask iwonz/tap/courier
```

GoReleaser updates `iwonz/homebrew-tap` after each release. The generated cask contains checksummed macOS amd64/arm64 assets and Linux assets for Homebrew environments that support binary casks.

## Scoop

```powershell
scoop bucket add iwonz https://github.com/iwonz/scoop-bucket
scoop install iwonz/courier
```

The manifest selects the amd64 or arm64 Windows zip and verifies the GoReleaser-generated hash.

## Winget

```powershell
winget install --exact --id iwonz.Courier
```

Each Courier release opens a manifest pull request against `microsoft/winget-pkgs`. A new version becomes available through Winget after the upstream review is merged.

## Direct download and verification

Release names use this stable pattern:

```text
courier_<version>_<os>_<arch>
courier_<version>_<os>_<arch>.tar.gz
courier_<version>_windows_<arch>.exe
courier_<version>_windows_<arch>.zip
checksums.txt
```

`os` is `darwin`, `linux`, or `windows`; `arch` is `amd64` or `arm64`. Download the matching raw binary and checksum manifest from:

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

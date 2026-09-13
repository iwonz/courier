## 1. Release artifacts

- [x] 1.1 Configure pinned GoReleaser Community builds, archives, checksums, and changelog
- [x] 1.2 Generate deb, rpm, apk, and Arch Linux packages with no runtime dependencies
- [x] 1.3 Configure Homebrew, Scoop, and Winget publishers against dedicated GitHub repositories

## 2. Installation channels

- [x] 2.1 Add checksum-verifying POSIX curl/wget installer
- [x] 2.2 Add checksum-verifying PowerShell installer and Windows zip path
- [x] 2.3 Add and test `@iwonz/courier` for npm, npx, yarn, and pnpm
- [x] 2.4 Document package-manager and direct-download commands for every supported ecosystem

## 3. Quality and release automation

- [x] 3.1 Add local formatting, vet, race, exact coverage, installer, npm, and GoReleaser snapshot gates
- [x] 3.2 Add GitHub Actions CI and ordered release/publication workflows
- [x] 3.3 Add one-command semantic tagging and release notes flow with credential preflight
- [x] 3.4 Verify installers and packages use only local fixtures and clean temporary state

## 4. Verification

- [x] 4.1 Run full local quality and GoReleaser snapshot dry run
- [x] 4.2 Validate release archive names against `courier update` and npm selectors
- [x] 4.3 Run strict OpenSpec validation and create the conventional commit

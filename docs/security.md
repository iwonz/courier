# Security model

Courier treats endpoint identity, paths, release artifacts, and temporary state as separate trust boundaries.

## SSH identity and authentication

Every target and ProxyJump hop is connected with `golang.org/x/crypto/ssh` and checked through `knownhosts.New`. Courier has no flag that disables host-key verification. OpenSSH configuration parsing supplies `HostName`, `User`, `Port`, `IdentityFile`, and `ProxyJump`; the parser does not perform connections. Courier independently constructs, verifies, authenticates, and closes every hop and final endpoint.

Private keys and the native SSH agent are attempted before an interactive password prompt. Passwords and key passphrases are read without terminal echo, are never accepted as command-line flags, and are zeroed after use where the Go representation permits it. On Windows, agent communication uses the standard OpenSSH named pipe directly; Courier does not invoke an SSH executable.

## Paths and transfer commits

Endpoint parsing rejects control characters and ambiguous forms. Local operations use `os.Root` so validated relative operations cannot escape through symlink traversal or a check/open race. Tar.gz creation and extraction share one two-pass inspector that rejects absolute or ambiguous paths, normalized duplicates, structural conflicts, unsupported types, and symlinks that escape their top-level archive entry. Extraction also enforces fixed entry-count, depth, and expansion-ratio limits even when the configurable expanded-size limit is unlimited.

The transfer engine completes preflight before changing the final target. Data is written under a private partial name and committed only after completion. Existing final paths are rejected rather than overwritten or merged, and unrelated entries in a destination container are preserved. Local commits use exclusive link, symlink, or directory reservation operations; SFTP uses the no-replace protocol rename. Whole-directory atomicity remains subject to the destination operating system, filesystem, and transport guarantees.

## Optional remote helper

SFTP supplied by the authenticated server is always preferred. A helper is eligible only for a typed failure to start that subsystem. Platform detection is read-only, then Courier requires an explicit interactive `yes`; EOF, an unavailable terminal, errors, and every other answer decline.

For a matching local OS and architecture, Courier hashes and uploads its running executable. For a different platform, it downloads the archive for its exact immutable version from the corresponding GitHub Release, verifies the archive against `checksums.txt`, safely extracts the single executable, and hashes it. Development builds cannot fetch a cross-platform helper because they have no immutable release version.

The already authenticated SSH channel uploads the helper into a randomly named private OS temporary directory. POSIX bootstrap uses the first available `sha256sum`, `shasum`, or `sha256`; Windows uses PowerShell `Get-FileHash`. A mismatch removes staging and prevents execution. The verified binary starts only in the hidden SFTP-helper mode and communicates over SSH standard streams. It is not added to `PATH`, does not listen on a network port, and is not installed as a service.

Closing the endpoint closes the helper SFTP stream, waits for the process, removes its remote directory, then closes SSH. Setup failures and cancellation run the same cleanup path. Cleanup failure is reported as an operation failure rather than hidden.

## Updates and releases

Self-update accepts only stable semantic release tags, bounded downloads, an exact platform archive, and its SHA-256 entry in `checksums.txt`. Extraction accepts only one correctly named regular executable. POSIX systems replace through a synchronized same-directory partial. Windows launches the verified partial in handoff mode, waits for the original PID to exit, replaces the target, and launches the installed target to remove the handoff file; unsafe internal paths and PIDs are rejected.

GitHub Actions publishes only after formatting, vet, race testing, exact 100% first-party statement coverage, installer acceptance, npm tests, GoReleaser validation, and a complete snapshot build succeed. See the [release runbook](releasing.md) for credentials and publication ordering.

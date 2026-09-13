## Why

Remote transfers must work without a local OpenSSH or rsync executable while retaining standard SSH trust and authentication behavior. Courier needs a native SSH/SFTP transport that plugs into the existing filesystem engine and identifies remote capabilities safely.

## What Changes

- Resolve SSH aliases, `HostName`, `User`, `Port`, `IdentityFile`, agent, `known_hosts`, and `ProxyJump`.
- Authenticate with keys or agent first and use an interactive password prompt only as a fallback; never accept passwords as CLI arguments.
- Require host-key verification and redact authentication material from errors and logs.
- Expose SFTP as an `fsx.Backend` with POSIX path semantics.
- Detect Linux, macOS, BSD, Windows, CPU architecture, and available archivers.
- Prefer a remote archiver when useful and retain the built-in archive implementation as the zero-install fallback.
- Detect unavailable SFTP/helper requirements and require explicit user confirmation before any optional remote helper bootstrap.

## Capabilities

### New Capabilities

- `ssh-connectivity`: native secure SSH authentication, configuration, jump hosts, and SFTP filesystem access.
- `remote-platform`: remote OS/architecture and archiver capability detection with a built-in fallback.

### Modified Capabilities

None.

## Impact

Adds `golang.org/x/crypto`, `github.com/kevinburke/ssh_config`, and `github.com/pkg/sftp`, plus the internal `sshx` package. The final executable still has no runtime executable dependency.

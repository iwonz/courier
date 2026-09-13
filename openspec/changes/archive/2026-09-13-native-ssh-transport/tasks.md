## 1. SSH configuration and authentication

- [x] 1.1 Parse required SSH config fields, includes, and host patterns
- [x] 1.2 Load identities and agent signers with interactive secret fallback
- [x] 1.3 Enforce known_hosts and implement ProxyJump chains

## 2. Remote filesystem and capabilities

- [x] 2.1 Implement and verify the SFTP `fsx.Backend`
- [x] 2.2 Detect remote OS, architecture, and archivers using constant probes
- [x] 2.3 Return a typed capability error for optional helper orchestration

## 3. Verification

- [x] 3.1 Cover SSH config, trust, auth, SFTP operations, and platform selection
- [x] 3.2 Run local in-process SSH integration tests with deterministic cleanup
- [x] 3.3 Run race, coverage, cross-build, and strict OpenSpec validation
- [x] 3.4 Create the conventional commit on `feat/004-native-ssh-transport`

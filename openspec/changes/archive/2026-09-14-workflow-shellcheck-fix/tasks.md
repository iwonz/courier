## 1. Deterministic local workflow lint

- [x] 1.1 Add the pinned checksum-verified ShellCheck bootstrap for supported developer hosts
- [x] 1.2 Make actionlint always use the resolved pinned ShellCheck executable

## 2. Release workflow correction

- [x] 2.1 Use the npm visibility retry counter without changing the number of publication checks
- [x] 2.2 Document tool parity and record the corrective branch in the implementation plan

## 3. Verification and delivery

- [x] 3.1 Verify bootstrap installation and cache reuse, workflow lint, strict OpenSpec, and the full repository gate
- [x] 3.2 Archive the OpenSpec change and create exactly one matching conventional commit
- [x] 3.3 Fast-forward main, push, and verify CI plus automatic Pages deployment

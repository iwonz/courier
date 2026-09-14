## 1. Platform durability boundary

- [x] 1.1 Split state-directory synchronization into covered non-Windows and Windows implementations
- [x] 1.2 Preserve file flush and no-replace commit semantics while documenting the Windows directory-flush limitation

## 2. Native regression behavior

- [x] 2.1 Add platform-specific directory-sync expectations
- [x] 2.2 Make optimistic-update test failures release coordination gates before cleanup

## 3. Verification and delivery

- [x] 3.1 Run formatting, exact coverage, workflow lint, strict OpenSpec, full release verification, and Windows cross-builds
- [x] 3.2 Archive the OpenSpec change and create exactly one matching conventional commit
- [x] 3.3 Fast-forward main, push, and verify native Windows plus complete CI and automatic Pages deployment

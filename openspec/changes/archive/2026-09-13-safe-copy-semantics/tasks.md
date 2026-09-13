## 1. Safety policy

- [x] 1.1 Return a successful disposition for canonically identical plain endpoints
- [x] 1.2 Reject identical archive/extract targets and directory descendants
- [x] 1.3 Reject existing resolved final paths before archive or transfer work

## 2. Transaction boundary

- [x] 2.1 Remove destination replacement and mirror semantics
- [x] 2.2 Add no-replace commit behavior with safe late-collision handling
- [x] 2.3 Preserve destination-only container entries

## 3. Quality

- [x] 3.1 Cover local, remote-alias, symlink, collision, rollback, and no-op cases
- [x] 3.2 Run strict OpenSpec validation and the full verification gate

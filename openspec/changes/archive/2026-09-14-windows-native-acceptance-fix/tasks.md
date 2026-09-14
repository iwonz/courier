## 1. Windows control and filesystem boundaries

- [x] 1.1 Resolve named-pipe ACL identity from the supported current-process token
- [x] 1.2 Adapt native rooted directory names to the slash-separated `fs.FS` contract

## 2. Archive and permission portability

- [x] 2.1 Convert local relative symlink targets to portable tar separators before safety validation
- [x] 2.2 Align private-mode assertions with the operating systems that expose POSIX permission bits

## 3. Verification and delivery

- [x] 3.1 Run formatting, exact coverage, workflow lint, strict OpenSpec, full release verification, and Windows cross-builds
- [x] 3.2 Archive the OpenSpec change and create exactly one matching conventional commit
- [x] 3.3 Fast-forward main, push, and require green native Windows, complete CI, and automatic Pages deployment

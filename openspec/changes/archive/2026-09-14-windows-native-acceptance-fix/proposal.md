# Change: Correct native Windows runtime boundaries

## Why

After registry directory synchronization became platform-correct, the native Windows suite exposed four independent portability assumptions: a zero token handle was used instead of the current-process pseudo token for named-pipe ACLs, rooted `fs.FS` directory reads received backslash paths, local relative symlink targets were not converted to POSIX tar form, and two tests required POSIX mode bits from Windows.

## What Changes

- Resolve the named-pipe owner SID from the supported current-process token API.
- Convert OS-native rooted directory paths to the slash-separated `fs.FS` contract.
- Normalize local symlink targets to portable tar separators before applying unchanged archive safety checks.
- Assert private POSIX permission bits only where the operating system exposes them, matching the existing runtime policy.

## Impact

Windows admin and worker control pipes become usable by the owning process, rooted directory commits work with native separators, and safe relative symlinks can be archived. Unix behavior, pipe ACL isolation, archive traversal rejection, and runtime interfaces remain unchanged.

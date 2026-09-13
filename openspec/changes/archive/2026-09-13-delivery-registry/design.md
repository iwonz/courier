# Context

Courier will run multiple deliveries on shared listeners and will later control them over private IPC. The durable registry is discovery and recovery metadata, not proof that a process is alive. It must remain safe across interruption, concurrent readers, Windows rename behavior, and stale process records.

# Decisions

## Domain records contain no secrets

Registry records contain UUIDs, route kinds, bind and control endpoints, controlled policy values, aggregate counters, lifecycle timestamps, tombstones, and owned temporary paths. Passwords, Basic credentials, session keys, resource tokens, URL credentials, and free-form error text have no fields in the persisted schema.

## Registry revisions are immutable

Each commit writes and flushes a private temporary file, then atomically links it to a previously absent revision filename before removing the temporary name. Readers select the highest valid revision. The no-replace link avoids cross-platform replacement semantics and ensures an interrupted write cannot corrupt the previous revision. A revision mismatch fails optimistically rather than silently merging writers.

## Filesystem access is rooted and private

The store creates or verifies an application-owned `0700` directory, opens it with `os.Root`, rejects symlinked state objects, and creates registry and history files with `0600` permissions. Temporary files are removed on success and every error path.

## Lifecycle and policy changes are validated in memory

An in-memory registry owns additions, transitions, tombstones, policy version changes, and temporary-resource ownership. A full snapshot validation runs before encoding so no partial mutation reaches disk. Counter accumulators use atomic integers and expose plain immutable snapshots to persistence and IPC.

# Risks / Trade-offs

Immutable revisions may leave superseded files after a crash. Readers deterministically choose the highest revision; bounded compaction and history rotation belong to operational reporting. The registry intentionally stores a PID only as diagnostic metadata and never treats it as authoritative liveness.

# Migration Plan

Add and test the domain, add the private immutable store and history log, document the schema and trust boundary, run race and cross-build checks, then archive the OpenSpec change.

# Browser deliveries

Courier can expose a local or SSH-backed path for browser download, or accept a browser upload into an existing local or SSH-backed directory. Both routes use the same worker, policy, selection, and staged-storage layers as other Courier operations.

## Start a delivery

Share a path for download:

```console
courier from ./report.pdf to web://
courier from server:/srv/releases to web:// --background --auth password
courier from ./release-tree to web:// --archive
```

Accept one-file-at-a-time uploads into a directory:

```console
courier from web:// to ./incoming/
courier from web:// to server:/srv/incoming/ --max-file-size 2GiB
courier from web:// to ./expanded/ --extract --max-extracted-size 20GiB
```

The command prints the browser URL and the delivery UUID. The URL contains a random 256-bit resource token; it does not contain a file path, credential, or policy value. Treat an unauthenticated URL as a bearer capability and stop it when it is no longer needed.

The default bind is `127.0.0.1:8080`. Use `--listen <host:port>` when other devices must connect. Network exposure does not disable peer admission or authentication.

## Authentication and access

`--auth password` prompts in the initiating terminal and creates a delivery-scoped browser session with CSRF protection. `--auth basic` prompts for a username and password and uses request-scoped HTTP Basic authentication. Passwords cannot be supplied as command arguments, are transferred only over private worker IPC, and are not written to the registry or logs. For SSH-backed paths, Courier performs an authenticated preflight and passes only credentials actually requested during that preflight to the worker over the same private IPC. A remote helper is made available to the worker only after the user approved it during preflight.

`--allow-ip` accepts an exact peer IP or CIDR and may be repeated. Courier evaluates the accepted connection address and ignores forwarding headers. `--auth-attempts` and `--auth-fail-action` control whether repeated failures ban a peer fingerprint or stop the delivery.

Protected names, paths, sizes, and directory entries are read and returned only after admission and authentication succeed.

## Download and navigation

A shared file is streamed with a bounded buffer. A shared directory can be navigated in the embedded data UI; an individual subdirectory can be downloaded as a deterministic tar.gz rooted at that directory name. Symlinks, traversal paths, special files, and objects rejected by the shared selection rules are not downloadable.

With `--archive`, Courier creates and verifies one `<source-name>.tar.gz` before registering the public delivery; that archive is the only published object. Without the flag, whole-directory downloads are generated from a preflighted manifest after authorization and reservation.

`--exclude`, `--exclude-regex`, and `--exclude-from` use their original command-line order. `--download-rate` is an aggregate delivery limit shared by concurrent client sends and, for SSH-backed sources, the remote-read leg.

## Upload and commit behavior

Each multipart request must contain exactly one safe base filename. Courier reserves transfer capacity before consuming file bytes, writes a private staging file, syncs it, and commits only when the final name is absent. A collision never replaces the existing entry, and staging data is removed on success, failure, or cancellation.

With `--extract`, the uploaded object must be tar.gz. Courier inspects every entry, enforces traversal, symlink, entry-count, depth, expanded-size, and expansion-ratio limits, stages the expanded tree, and commits only non-conflicting top-level entries into the destination root. The uploaded archive itself is not retained.

`--max-file-size` defaults to `10GiB`; `unlimited` removes the product cap. `--upload-rate` is aggregate across the delivery. Selection rules apply to uploaded names.

## Lifecycle and API-only mode

Without `--background`, the initiating process owns a foreground lease. Interrupting it stops only that delivery and closes its endpoint resources. With `--background`, the worker retains the delivery after the command exits; the printed UUID is the stable control identity used by the server-control commands once they ship.

`--no-ui` disables the HTML entry point and returns a small versioned JSON description after authorization. Versioned metadata, session, upload, and download endpoints remain available under the opaque delivery URL. Static UI assets contain no delivery metadata or credentials.

The embedded delivery application is a React root composed from the repository-owned shadcn [`@courier/ui`](../web/ui) package. It reuses the compact square Relay pigeon only as decoration; authentication, metadata, upload, and download behavior remain native semantic controls backed exclusively by the versioned API.

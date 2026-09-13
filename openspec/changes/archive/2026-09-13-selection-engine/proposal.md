# Change: Add the shared selection engine

## Why

Copy, archive, extraction, listing, and HTTP delivery must expose identical object selection. Implementing exclusion independently in each traversal would create security and product inconsistencies, especially for ordered gitignore negation and interleaved rule files.

## What Changes

- Add a reusable selector backed by go-git gitignore matching and Go regular expressions.
- Compile direct rules and local rule files in exact command-line occurrence order.
- Ship `--exclude`, `--exclude-regex`, and `--exclude-from` for the current path transfer.
- Apply the same selector to transfer scanning/copying and archive creation.
- Keep selector APIs reusable by later extraction, listing, upload, and download tasks.

## Impact

Users can omit selected relative paths without changing source data. Invalid expressions or unreadable rule files fail before endpoint resources are opened.

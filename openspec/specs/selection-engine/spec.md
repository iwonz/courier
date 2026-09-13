# selection-engine Specification

## Purpose
Define one ordered, preflight-compiled object selection policy shared by every data traversal.

## Requirements

### Requirement: Shared ordered selection

Courier SHALL use one compiled selection policy for copy, archive, extraction, listing, upload, and download traversal, with paths expressed as slash-separated values relative to the data root.

#### Scenario: Interleaved rule sources

- **WHEN** direct gitignore rules and rule files are interleaved on the command line
- **THEN** file contents are expanded at their occurrence position and later gitignore negation has the expected effect

#### Scenario: Regular expression exclusion

- **WHEN** any configured Go regular expression matches a relative path
- **THEN** that object is excluded regardless of other regular expressions

### Requirement: Selection preflight

Courier SHALL compile expressions and read local rule files before opening transfer endpoint resources.

#### Scenario: Unreadable rule file

- **WHEN** an exclude-from file cannot be read
- **THEN** preflight fails without opening a local root, SSH connection, server, or worker

### Requirement: Selection option surface

Courier SHALL expose repeatable `--exclude`, `--exclude-regex`, and `--exclude-from` options on the shipped `from` command and retain their global occurrence order.

#### Scenario: Interleaved CLI options

- **WHEN** `--exclude '*.tmp' --exclude-from rules --exclude '!keep.tmp'` is parsed
- **THEN** the planner receives those three selection occurrences in the same order

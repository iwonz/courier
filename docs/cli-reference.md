# Courier CLI reference

Contract version `0.15.0`; target release `0.3.5`. This file is generated from `docs/cli-contract.yaml`.

## Commands

| Status | Command | Kind |
|---|---|---|
| shipped | `courier from <source> to <destination> [options]` | product |
| shipped | `courier servers` | product |
| shipped | `courier servers stop <uuid>|--all` | product |
| shipped | `courier ui start [options]` | product |
| shipped | `courier ui stop` | product |
| system | `courier help [command]` | system |
| system | `courier version` | system |
| system | `courier update` | system |

## Endpoint kinds

| Status | Kind | Syntax |
|---|---|---|
| shipped | `local` | ./path, /absolute/path, or an explicit Windows path |
| shipped | `ssh` | [user@]host:/path or [user@]host:C:/path |
| shipped | `web` | web:// |
| shipped | `webhook` | webhook:// |
| shipped | `http` | http://host/path or https://host/path |

## Routes and options

| Status | Route | Source | Destination | Allowed options |
|---|---|---|---|---|
| shipped | `path-to-path` | local, ssh | local, ssh | `--archive`, `--extract`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--max-extracted-size`, `--upload-rate`, `--download-rate` |
| shipped | `web-to-path` | web | local, ssh | `--extract`, `--listen`, `--background`, `--auth`, `--auth-attempts`, `--auth-fail-action`, `--limit`, `--allow-ip`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--max-file-size`, `--max-extracted-size`, `--upload-rate` |
| shipped | `path-to-web` | local, ssh | web | `--archive`, `--listen`, `--background`, `--auth`, `--auth-attempts`, `--auth-fail-action`, `--limit`, `--no-ui`, `--allow-ip`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--download-rate` |
| shipped | `webhook-to-path` | webhook | local, ssh | `--extract`, `--listen`, `--background`, `--auth`, `--auth-attempts`, `--auth-fail-action`, `--limit`, `--allow-ip`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--max-file-size`, `--max-extracted-size`, `--upload-rate` |
| shipped | `path-to-http` | local, ssh | http | `--archive`, `--auth`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--upload-rate` |

## Options

| Status | Option | Repeatable | Default | Applies to | Conflicts |
|---|---|---:|---|---|---|
| shipped | `--archive` | false | `false` | path-to-path, path-to-web, path-to-http | extract |
| shipped | `--extract` | false | `false` | path-to-path, web-to-path, webhook-to-path | archive |
| shipped | `--listen <host:port>` | false | `127.0.0.1:8080` | web-to-path, path-to-web, webhook-to-path, ui-start | none |
| shipped | `--background` | false | `false` | web-to-path, path-to-web, webhook-to-path, ui-start | none |
| shipped | `--auth <none|basic|password>` | false | `none` | web-to-path, path-to-web, webhook-to-path, path-to-http | none |
| shipped | `--auth-attempts <N>` | false | `5` | web-to-path, path-to-web, webhook-to-path | none |
| shipped | `--auth-fail-action <ban|stop>` | false | `ban` | web-to-path, path-to-web, webhook-to-path | none |
| shipped | `--limit <N>` | false | `unlimited` | web-to-path, path-to-web, webhook-to-path | none |
| shipped | `--no-ui` | false | `false` | path-to-web | none |
| shipped | `--allow-ip <IP/CIDR>` | true | `none` | web-to-path, path-to-web, webhook-to-path | none |
| shipped | `--exclude <pattern>` | true | `none` | path-to-path, web-to-path, path-to-web, webhook-to-path, path-to-http | none |
| shipped | `--exclude-regex <regex>` | true | `none` | path-to-path, web-to-path, path-to-web, webhook-to-path, path-to-http | none |
| shipped | `--exclude-from <file>` | true | `none` | path-to-path, web-to-path, path-to-web, webhook-to-path, path-to-http | none |
| shipped | `--max-file-size <size|unlimited>` | false | `10GiB` | web-to-path, webhook-to-path | none |
| shipped | `--max-extracted-size <size|unlimited>` | false | `100GiB` | path-to-path, web-to-path, webhook-to-path | none |
| shipped | `--upload-rate <rate|unlimited>` | false | `unlimited` | local-to-ssh, ssh-to-ssh, web-to-path, webhook-to-path, path-to-http | none |
| shipped | `--download-rate <rate|unlimited>` | false | `unlimited` | ssh-to-local, ssh-to-ssh, path-to-web | none |
| shipped | `--all` | false | `false` | servers-stop | none |

## Explicitly unsupported

- `courier servers list`
- `courier servers clean`
- `public worker commands`
- `--mirror`
- `--zip`
- `--unzip`
- `--archive --extract`
- `HTTP URL as a source`
- `web:// to web://`
- `webhook:// to web://`

## Examples

```sh
courier from ./report.pdf to server:/srv/inbox/
courier from root@203.0.113.10:/opt/node/data to ./backup/
courier from source-server:/opt/node/data to backup-server:/srv/data/
courier from ./data to ./backup/ --archive
courier from webhook:// to ./inbox/ --auth basic
courier from ./report.pdf to https://example.com/hooks/courier
courier servers
courier servers stop 550e8400-e29b-41d4-a716-446655440000
courier servers stop --all
courier ui start --background
courier ui stop
```

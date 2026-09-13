# Courier CLI reference

Contract version `0.8.0`; target release `0.2.0`. This file is generated from `docs/cli-contract.yaml`.

## Commands

| Status | Command | Kind |
|---|---|---|
| shipped | `courier from <source> to <destination> [options]` | product |
| planned | `courier servers` | product |
| planned | `courier servers stop <uuid>|--all` | product |
| planned | `courier ui start [options]` | product |
| planned | `courier ui stop` | product |
| system | `courier help [command]` | system |
| system | `courier version` | system |
| system | `courier update` | system |

## Endpoint kinds

| Status | Kind | Syntax |
|---|---|---|
| shipped | `local` | ./path, /absolute/path, or an explicit Windows path |
| shipped | `ssh` | [user@]host:/path or [user@]host:C:/path |
| planned | `web` | web:// |
| planned | `webhook` | webhook:// |
| planned | `http` | http://host/path or https://host/path |

## Routes and options

| Status | Route | Source | Destination | Allowed options |
|---|---|---|---|---|
| shipped | `path-to-path` | local, ssh | local, ssh | `--archive` |
| planned | `web-to-path` | web | local, ssh | `--extract`, `--listen`, `--background`, `--auth`, `--auth-attempts`, `--auth-fail-action`, `--limit`, `--allow-ip`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--max-file-size`, `--max-extracted-size`, `--upload-rate` |
| planned | `path-to-web` | local, ssh | web | `--archive`, `--listen`, `--background`, `--auth`, `--auth-attempts`, `--auth-fail-action`, `--limit`, `--no-ui`, `--allow-ip`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--download-rate` |
| planned | `webhook-to-path` | webhook | local, ssh | `--extract`, `--listen`, `--background`, `--auth`, `--auth-attempts`, `--auth-fail-action`, `--limit`, `--allow-ip`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--max-file-size`, `--max-extracted-size`, `--upload-rate` |
| planned | `path-to-http` | local, ssh | http | `--archive`, `--auth`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--upload-rate`, `--download-rate` |

## Options

| Status | Option | Repeatable | Default | Applies to | Conflicts |
|---|---|---:|---|---|---|
| shipped | `--archive` | false | `false` | path-to-path, path-to-web, path-to-http | extract |
| planned | `--extract` | false | `false` | path-to-path, web-to-path, webhook-to-path | archive |
| planned | `--listen <host:port>` | false | `127.0.0.1:8080` | web-to-path, path-to-web, webhook-to-path, ui-start | none |
| planned | `--background` | false | `false` | web-to-path, path-to-web, webhook-to-path, ui-start | none |
| planned | `--auth <none|basic|password>` | false | `none` | web-to-path, path-to-web, webhook-to-path, path-to-http | none |
| planned | `--auth-attempts <N>` | false | `5` | web-to-path, path-to-web, webhook-to-path | none |
| planned | `--auth-fail-action <ban|stop>` | false | `ban` | web-to-path, path-to-web, webhook-to-path | none |
| planned | `--limit <N>` | false | `unlimited` | web-to-path, path-to-web, webhook-to-path | none |
| planned | `--no-ui` | false | `false` | path-to-web | none |
| planned | `--allow-ip <IP/CIDR>` | true | `none` | web-to-path, path-to-web, webhook-to-path | none |
| planned | `--exclude <pattern>` | true | `none` | path-to-path, web-to-path, path-to-web, webhook-to-path, path-to-http | none |
| planned | `--exclude-regex <regex>` | true | `none` | path-to-path, web-to-path, path-to-web, webhook-to-path, path-to-http | none |
| planned | `--exclude-from <file>` | true | `none` | path-to-path, web-to-path, path-to-web, webhook-to-path, path-to-http | none |
| planned | `--max-file-size <size|unlimited>` | false | `10GiB` | web-to-path, webhook-to-path | none |
| planned | `--max-extracted-size <size|unlimited>` | false | `100GiB` | path-to-path, web-to-path, webhook-to-path | none |
| planned | `--upload-rate <rate|unlimited>` | false | `unlimited` | local-to-ssh, ssh-to-ssh, web-to-path, webhook-to-path, path-to-http | none |
| planned | `--download-rate <rate|unlimited>` | false | `unlimited` | ssh-to-local, ssh-to-ssh, path-to-web | none |
| planned | `--all` | false | `false` | servers-stop | none |

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
```

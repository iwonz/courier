# Courier CLI reference

Contract version `0.16.0`; target release `0.3.12`. This file is generated from `docs/cli-contract.yaml`.

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

## Arguments

| Command | Argument | Requirement | Description |
|---|---|---|---|
| `from` | `<source>` | Required | File, directory, browser upload, or webhook input to read from. |
| `from` | `<destination>` | Required | Path, browser download, or HTTP endpoint that receives the data. |
| `servers stop` | `<uuid>` | Required | Delivery or server UUID to stop; omitted when --all is used. |
| `help` | `<command>` | Optional | Command path whose detailed help should be shown. |

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
| shipped | `path-to-path` | local, ssh | local, ssh | `--archive`, `--extract`, `--force-source-creation`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--max-extracted-size`, `--upload-rate`, `--download-rate` |
| shipped | `web-to-path` | web | local, ssh | `--extract`, `--listen`, `--background`, `--force-source-creation`, `--auth`, `--auth-attempts`, `--auth-fail-action`, `--limit`, `--allow-ip`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--max-file-size`, `--max-extracted-size`, `--upload-rate` |
| shipped | `path-to-web` | local, ssh | web | `--archive`, `--listen`, `--background`, `--force-source-creation`, `--auth`, `--auth-attempts`, `--auth-fail-action`, `--limit`, `--no-ui`, `--allow-ip`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--download-rate` |
| shipped | `webhook-to-path` | webhook | local, ssh | `--extract`, `--listen`, `--background`, `--force-source-creation`, `--auth`, `--auth-attempts`, `--auth-fail-action`, `--limit`, `--allow-ip`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--max-file-size`, `--max-extracted-size`, `--upload-rate` |
| shipped | `path-to-http` | local, ssh | http | `--archive`, `--force-source-creation`, `--auth`, `--exclude`, `--exclude-regex`, `--exclude-from`, `--upload-rate` |

## Options

| Status | Option | Requirement | Description | Default | Repeatable | Applies to | Requires | Conflicts |
|---|---|---|---|---|---:|---|---|---|
| shipped | `--archive` | Optional | Pack the source into a verified <source-name>.tar.gz before transfer. | `false` | false | Path to path, Path to browser download, Path to HTTP webhook | none | extract |
| shipped | `--extract` | Optional | Safely extract a tar.gz source into the destination root. | `false` | false | Path to path, Browser upload to path, Incoming webhook to path | none | archive |
| shipped | `--listen <host:port>` | Optional | Bind incoming browser, webhook, or administration traffic to this address. | `127.0.0.1:8080` | false | Browser upload to path, Path to browser download, Incoming webhook to path, Administration UI start | none | none |
| shipped | `--background` | Optional | Detach browser downloads, browser uploads, incoming webhooks, or ui start after printing the URL and UUID; the process survives terminal closure until stopped by courier servers stop, courier ui stop, configured stop behavior, a fatal worker failure, or process termination. | `false` | false | Browser upload to path, Path to browser download, Incoming webhook to path, Administration UI start | none | none |
| shipped | `--force-source-creation` | Optional | Create every missing directory endpoint recursively without asking; no other safety confirmation is bypassed. | `false` | false | Path to path, Browser upload to path, Path to browser download, Incoming webhook to path, Path to HTTP webhook | none | none |
| shipped | `--auth <none\|basic\|password>` | Optional | Choose no authentication, HTTP Basic authentication, or a browser password. | `none` | false | Browser upload to path, Path to browser download, Incoming webhook to path, Path to HTTP webhook | none | none |
| shipped | `--auth-attempts <N>` | Optional | Set how many failed authentication attempts are allowed before the configured action. | `5` | false | Browser upload to path, Path to browser download, Incoming webhook to path | none | none |
| shipped | `--auth-fail-action <ban\|stop>` | Optional | Ban the peer or stop the delivery when the authentication-attempt limit is reached. | `ban` | false | Browser upload to path, Path to browser download, Incoming webhook to path | none | none |
| shipped | `--limit <N>` | Optional | Limit concurrent transfers for a hosted delivery, or allow an unlimited count. | `unlimited` | false | Browser upload to path, Path to browser download, Incoming webhook to path | none | none |
| shipped | `--no-ui` | Optional | Expose only the versioned browser-delivery data API, without the web interface. | `false` | false | Path to browser download | none | none |
| shipped | `--allow-ip <IP/CIDR>` | Optional | Allow one peer IP address or CIDR; repeat the option to add more networks. | `none` | true | Browser upload to path, Path to browser download, Incoming webhook to path | none | none |
| shipped | `--exclude <pattern>` | Optional | Exclude paths with an ordered gitignore-style pattern. | `none` | true | Path to path, Browser upload to path, Path to browser download, Incoming webhook to path, Path to HTTP webhook | none | none |
| shipped | `--exclude-regex <regex>` | Optional | Exclude paths matching a Go regular expression. | `none` | true | Path to path, Browser upload to path, Path to browser download, Incoming webhook to path, Path to HTTP webhook | none | none |
| shipped | `--exclude-from <file>` | Optional | Read ordered gitignore-style exclusion rules from a local file. | `none` | true | Path to path, Browser upload to path, Path to browser download, Incoming webhook to path, Path to HTTP webhook | none | none |
| shipped | `--max-file-size <size\|unlimited>` | Optional | Set the largest accepted incoming file size, or remove the configurable limit. | `10GiB` | false | Browser upload to path, Incoming webhook to path | none | none |
| shipped | `--max-extracted-size <size\|unlimited>` | Optional | Set the maximum total expanded archive size; requires --extract. | `100GiB` | false | Path to path, Browser upload to path, Incoming webhook to path | extract | none |
| shipped | `--upload-rate <rate\|unlimited>` | Optional | Limit the aggregate upload rate, or allow an unlimited rate. | `unlimited` | false | Local path to SSH path, SSH path to SSH path, Browser upload to path, Incoming webhook to path, Path to HTTP webhook | none | none |
| shipped | `--download-rate <rate\|unlimited>` | Optional | Limit the aggregate download rate, or allow an unlimited rate. | `unlimited` | false | SSH path to local path, SSH path to SSH path, Path to browser download | none | none |
| shipped | `--all` | Optional | Stop every discovered Courier delivery and server instead of one UUID. | `false` | false | Server stop | none | none |

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

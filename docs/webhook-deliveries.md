# Webhook deliveries

Courier implements one deliberately small webhook profile. It is suitable for software that can send or receive one `multipart/form-data` file, and it does not claim compatibility with every service that uses the word “webhook”.

## Incoming webhook

Start an incoming endpoint whose final destination is an existing local or SSH directory:

```sh
courier from webhook:// to ./inbox/
courier from webhook:// to server:/srv/inbox/ --auth basic --background
```

Courier prints an opaque URL shaped like:

```text
http://127.0.0.1:8080/d/<random-resource-token>/upload
```

POST exactly one `multipart/form-data` file part named `file` to that URL. For example:

```sh
curl --fail --form 'file=@./report.pdf' 'http://127.0.0.1:8080/d/<token>/upload'
```

With `--auth basic`, Courier prompts the owner for the username and password before registration. A client supplies standard HTTP Basic authentication; credentials never appear in the Courier command, URL, registry, or logs. Password-session mode is browser-only and is rejected for webhooks.

The endpoint applies `--allow-ip`, `--auth-attempts`, `--auth-fail-action`, `--limit`, ordered selection rules, `--max-file-size`, and aggregate `--upload-rate`. It validates the untrusted filename as one safe basename, rejects a second multipart part, stages bytes privately, and returns success only after an absent final path is committed. An existing name is a conflict and is never overwritten. With `--extract`, the uploaded tar.gz is inspected and extracted through Courier's shared codec and safety limits instead of being retained.

There is no HTML page, metadata API, or password-session endpoint for an incoming webhook. The default listener is loopback. Use a non-loopback `--listen` value only with an appropriate network and authentication policy.

## Outgoing webhook

Send one local or SSH-backed file:

```sh
courier from ./report.pdf to https://example.com/hooks/courier
courier from server:/srv/report.pdf to https://example.com/hooks/courier --auth basic
```

Courier issues exactly one POST with exactly one multipart file part named `file`. A directory is rejected unless `--archive` is explicit:

```sh
courier from ./results to https://example.com/hooks/courier --archive
```

Archive mode creates and verifies `<source-name>.tar.gz`, applies the shared ordered selection rules before packing, streams that one artifact, and removes it afterward. Direct and archived requests use bounded memory and honor the aggregate `--upload-rate` value.

The native HTTP client performs normal TLS certificate verification for HTTPS and may use the standard proxy environment. Courier does not follow redirects, retry failed requests, add idempotency keys, or support JSON/raw bodies, custom field names, arbitrary headers, signatures, and provider-specific profiles.

## Result semantics

- Any received 2xx status means the remote HTTP endpoint accepted the request. A `202 Accepted` response does not prove durable storage or completed downstream processing.
- A received non-2xx status is a known rejection. A redirect is reported in this category and its target is not contacted.
- If payload bytes were consumed by the HTTP transport but no response arrived, Courier reports an unknown outcome and the byte count. It does not retry because doing so could duplicate a receiver-side effect.

Incoming foreground ownership and `--background` follow the same worker and lease rules as browser deliveries. Outgoing webhook delivery is a finite operation and does not create a server or durable delivery record.

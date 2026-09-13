# Courier CLI

Courier is an extensible, cross-platform CLI for safely transferring files and directories between local and SSH endpoints.

```text
courier from <source> to <destination> [flags]
```

The project is developed spec-first: every task has its own path under `openspec/changes`, task branch, and conventional commit. See [docs/implementation-plan.md](docs/implementation-plan.md) for the complete roadmap.

## Development

Go 1.25 or newer is required.

```sh
go test ./...
go build -o bin/courier ./cmd/courier
./bin/courier help
```

## License

[MIT](LICENSE)

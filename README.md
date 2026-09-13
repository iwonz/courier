# Courier CLI

Courier — расширяемая кроссплатформенная утилита для безопасного переноса файлов и каталогов между локальными и SSH endpoints.

```text
courier from <source> to <destination> [flags]
```

Проект разрабатывается spec-first: каждая задача имеет отдельный путь в `openspec/changes`, task branch и conventional commit. Полный roadmap находится в [docs/implementation-plan.md](docs/implementation-plan.md).

## Разработка

Требуется Go 1.24 или новее.

```sh
go test ./...
go build -o bin/courier ./cmd/courier
./bin/courier help
```

## Лицензия

[MIT](LICENSE)

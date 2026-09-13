GORELEASER_VERSION := v2.18.1
GORELEASER ?= .cache/tools/goreleaser

.PHONY: all fmt-check vet contract-check test npm-test ui-test goreleaser-check snapshot verify precommit hooks release clean

all: verify

fmt-check:
	@files="$$(gofmt -l $$(find cmd internal -name '*.go' -type f))"; \
	if [ -n "$$files" ]; then printf '%s\n' "Go files need formatting:" "$$files" >&2; exit 1; fi

vet:
	go vet ./...

contract-check:
	go run ./cmd/contractdoc --check

test: fmt-check vet contract-check
	go test -race ./... -covermode=atomic -coverprofile=coverage.out
	@total="$$(go tool cover -func=coverage.out | awk '/^total:/ { print $$3 }')"; \
	if [ "$$total" != "100.0%" ]; then printf '%s\n' "statement coverage is $$total, expected 100.0%" >&2; exit 1; fi

npm-test:
	npm test --prefix npm --loglevel=error
	cd npm && npm pack --dry-run --loglevel=error >/dev/null

ui-test:
	npm ci --prefix web --ignore-scripts --no-audit --no-fund --loglevel=error
	npm run verify --prefix web --loglevel=error

goreleaser-check:
	@if [ ! -x "$(GORELEASER)" ] && ! command -v "$(GORELEASER)" >/dev/null 2>&1; then \
		GORELEASER_VERSION=$(GORELEASER_VERSION) ./scripts/install-goreleaser.sh; \
	fi
	$(GORELEASER) check

snapshot: goreleaser-check
	$(GORELEASER) release --snapshot --clean
	./scripts/verify-dist.sh
	./scripts/test-release-installer.sh
	node ./scripts/test-npm-dist.js

verify: test npm-test ui-test snapshot

precommit: verify

hooks:
	git config core.hooksPath .githooks

release:
	./scripts/release.sh $(VERSION)

clean:
	rm -rf dist .cache coverage.out coverage.html web/ui/dist web/ui/coverage

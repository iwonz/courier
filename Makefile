GORELEASER_VERSION := v2.18.1
GORELEASER ?= .cache/tools/goreleaser
ACTIONLINT_VERSION := v1.7.12
ACTIONLINT ?= go run github.com/rhysd/actionlint/cmd/actionlint@$(ACTIONLINT_VERSION)
SHELLCHECK_VERSION := v0.11.0
SHELLCHECK ?= .cache/tools/shellcheck
OPENSPEC ?= openspec

.PHONY: all fmt-check vet contract-check workflow-check openspec-check test npm-test ui-test browser-test pages-build pages-publish goreleaser-check snapshot package-test verify precommit hooks release clean

all: verify

fmt-check:
	@files="$$(gofmt -l $$(find cmd internal -name '*.go' -type f))"; \
	if [ -n "$$files" ]; then printf '%s\n' "Go files need formatting:" "$$files" >&2; exit 1; fi

vet:
	go vet ./...

contract-check:
	go run ./cmd/contractdoc --check

test: fmt-check vet contract-check
	./scripts/test-runtime.sh --race

workflow-check:
	@if [ ! -x "$(SHELLCHECK)" ] && ! command -v "$(SHELLCHECK)" >/dev/null 2>&1; then \
		SHELLCHECK_VERSION=$(SHELLCHECK_VERSION) ./scripts/install-shellcheck.sh; \
	fi
	$(ACTIONLINT) -shellcheck "$(SHELLCHECK)" -color

openspec-check:
	$(OPENSPEC) validate --all --strict --no-interactive

npm-test:
	npm test --prefix npm --loglevel=error
	cd npm && npm pack --dry-run --loglevel=error >/dev/null

ui-test:
	npm ci --prefix web --ignore-scripts --no-audit --no-fund --loglevel=error
	npm run verify --prefix web --loglevel=error

browser-test: ui-test
	npm run browser:test --prefix web --loglevel=error

pages-build: contract-check
	npm ci --prefix web --ignore-scripts --no-audit --no-fund --loglevel=error
	npm run test --prefix web --workspace @courier/landing --loglevel=error
	npm run build --prefix web --workspace @courier/landing --loglevel=error

pages-publish:
	./scripts/publish-pages.sh

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

package-test: snapshot
	./scripts/test-linux-packages.sh

verify: test npm-test browser-test workflow-check openspec-check package-test

precommit: verify

hooks:
	git config core.hooksPath .githooks

release:
	./scripts/release.sh $(VERSION)

clean:
	rm -rf dist .cache coverage.out coverage.html web/ui/dist web/ui/coverage web/data/coverage web/admin/coverage web/landing/dist web/landing/coverage web/test-results web/playwright-report

#!/bin/sh

set -eu

repository_directory=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)

temporary_directory=$(mktemp -d "${TMPDIR:-/tmp}/courier-browser.XXXXXX")
runtime_directory=$(mktemp -d "/tmp/cr.XXXXXX")
port_base=${COURIER_BROWSER_PORT_BASE:-$((20000 + ($$ % 20000)))}
cleanup() {
  rm -rf "$temporary_directory"
  rm -rf "$runtime_directory"
}
finish() {
  status=$?
  trap - EXIT HUP INT TERM
  cleanup
  exit "$status"
}
trap finish EXIT
trap 'exit 130' HUP INT TERM

export COURIER_LANDING_PORT=${COURIER_LANDING_PORT:-$port_base}
export COURIER_DATA_PORT=${COURIER_DATA_PORT:-$((port_base + 1))}
export COURIER_ADMIN_PORT=${COURIER_ADMIN_PORT:-$((port_base + 2))}
export COURIER_RUNTIME_BINARY="$temporary_directory/courier"
export COURIER_RUNTIME_ROOT="$runtime_directory"
export COURIER_RUNTIME_PORT_BASE=${COURIER_RUNTIME_PORT_BASE:-$((port_base + 20))}

mkdir -p "$COURIER_RUNTIME_ROOT"
(cd "$repository_directory" && go build -o "$COURIER_RUNTIME_BINARY" ./cmd/courier)

cd "$repository_directory/web"

COURIER_PLAYWRIGHT_OUTPUT_DIR="$temporary_directory/test-results" \
  ./node_modules/.bin/playwright test --config playwright.config.ts

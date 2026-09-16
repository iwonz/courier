#!/bin/sh

set -eu

cd "$(dirname "$0")/../web"

temporary_directory=$(mktemp -d "${TMPDIR:-/tmp}/courier-browser.XXXXXX")
port_base=${COURIER_BROWSER_PORT_BASE:-$((20000 + ($$ % 20000)))}
cleanup() {
  rm -rf "$temporary_directory"
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

COURIER_PLAYWRIGHT_OUTPUT_DIR="$temporary_directory/test-results" \
  ./node_modules/.bin/playwright test --config playwright.config.ts

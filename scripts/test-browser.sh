#!/bin/sh

set -eu

cd "$(dirname "$0")/../web"

temporary_directory=$(mktemp -d "${TMPDIR:-/tmp}/courier-browser.XXXXXX")
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

COURIER_PLAYWRIGHT_OUTPUT_DIR="$temporary_directory/test-results" \
  ./node_modules/.bin/playwright test --config playwright.config.ts

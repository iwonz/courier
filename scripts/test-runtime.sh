#!/bin/sh

set -eu

cd "$(dirname "$0")/.."

case $(uname -s) in
  Darwin|Linux) ;;
  *) printf '%s\n' "Skipping POSIX runtime acceptance on this operating system"; exit 0 ;;
esac

temporary_directory=$(mktemp -d "${TMPDIR:-/tmp}/courier-runtime.XXXXXX")
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

coverage="$temporary_directory/coverage.out"
case ${1:-} in
  "") go test ./... -covermode=atomic -coverprofile="$coverage" ;;
  --race) go test -race ./... -covermode=atomic -coverprofile="$coverage" ;;
  *) printf '%s\n' "usage: scripts/test-runtime.sh [--race]" >&2; exit 2 ;;
esac

total=$(go tool cover -func="$coverage" | awk '/^total:/ { print $3 }')
[ "$total" = "100.0%" ] || {
  printf '%s\n' "POSIX statement coverage is $total, expected 100.0%" >&2
  exit 1
}

binary="$temporary_directory/courier"
go build -trimpath -o "$binary" ./cmd/courier
"$binary" help | grep -Fq "Safely transfer files and directories"
"$binary" version | grep -Fq "courier dev"

printf '%s\n' "Verified compiled $(uname -s) runtime with exact first-party coverage"

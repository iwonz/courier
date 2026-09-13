#!/bin/sh

set -eu

version=${GORELEASER_VERSION:-v2.18.1}
destination=${GORELEASER_INSTALL_PATH:-.cache/tools/goreleaser}

if [ -x "$destination" ] && "$destination" --version 2>/dev/null | grep -q "GitVersion:[[:space:]]*${version#v}"; then
  exit 0
fi

case $(uname -s) in
  Darwin) target_os=Darwin ;;
  Linux) target_os=Linux ;;
  *) printf '%s\n' "GoReleaser bootstrap supports macOS and Linux" >&2; exit 1 ;;
esac
case $(uname -m) in
  x86_64|amd64) target_arch=x86_64 ;;
  arm64|aarch64) target_arch=arm64 ;;
  *) printf '%s\n' "Unsupported GoReleaser host architecture: $(uname -m)" >&2; exit 1 ;;
esac

temporary_dir=$(mktemp -d "${TMPDIR:-/tmp}/courier-goreleaser.XXXXXX")
trap 'rm -rf "$temporary_dir"' EXIT HUP INT TERM
archive="goreleaser_${target_os}_${target_arch}.tar.gz"
base="https://github.com/goreleaser/goreleaser/releases/download/${version}"

download() {
  if command -v curl >/dev/null 2>&1; then
    curl --fail --silent --show-error --location "$1" --output "$2"
  elif command -v wget >/dev/null 2>&1; then
    wget --quiet --output-document="$2" "$1"
  else
    printf '%s\n' "curl or wget is required to bootstrap GoReleaser" >&2
    exit 1
  fi
}

download "$base/$archive" "$temporary_dir/$archive"
download "$base/checksums.txt" "$temporary_dir/checksums.txt"
expected=$(awk -v name="$archive" '$2 == name || $2 == "*" name { print tolower($1); exit }' "$temporary_dir/checksums.txt")
[ -n "$expected" ] || { printf '%s\n' "GoReleaser checksum is missing" >&2; exit 1; }
if command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$temporary_dir/$archive" | awk '{ print tolower($1) }')
else
  actual=$(shasum -a 256 "$temporary_dir/$archive" | awk '{ print tolower($1) }')
fi
[ "$actual" = "$expected" ] || { printf '%s\n' "GoReleaser checksum mismatch" >&2; exit 1; }

tar -xzf "$temporary_dir/$archive" -C "$temporary_dir" goreleaser
mkdir -p "$(dirname "$destination")"
chmod 755 "$temporary_dir/goreleaser"
mv "$temporary_dir/goreleaser" "$destination"

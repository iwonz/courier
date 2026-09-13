#!/bin/sh

set -eu

version=$(sed -n 's/.*"version":"\([^"]*\)".*/\1/p' dist/metadata.json)
case $(uname -s) in
  Darwin) target_os=darwin ;;
  Linux) target_os=linux ;;
  *) printf '%s\n' "Skipping POSIX release installer on this host"; exit 0 ;;
esac
case $(uname -m) in
  x86_64|amd64) target_arch=amd64 ;;
  arm64|aarch64) target_arch=arm64 ;;
  *) printf '%s\n' "Skipping POSIX release installer on this architecture"; exit 0 ;;
esac

temporary_dir=$(mktemp -d "${TMPDIR:-/tmp}/courier-release-install.XXXXXX")
trap 'rm -rf "$temporary_dir"' EXIT HUP INT TERM
release_dir="$temporary_dir/releases/v$version"
scratch="$temporary_dir/scratch"
install_dir="$temporary_dir/install"
mkdir -p "$release_dir" "$scratch"
asset="courier_${version}_${target_os}_${target_arch}"
binary=$(find dist -path "*/courier-unix_${target_os}_${target_arch}_*/courier" -type f | sed -n '1p')
[ -n "$binary" ] || { printf '%s\n' "cannot locate host snapshot binary" >&2; exit 1; }
cp "$binary" "$release_dir/$asset"
awk -v name="$asset" '$2 == name || $2 == "*" name' dist/checksums.txt >"$release_dir/checksums.txt"

COURIER_VERSION="$version" \
COURIER_RELEASE_BASE_URL="file://$temporary_dir/releases" \
COURIER_INSTALL_DIR="$install_dir" \
TMPDIR="$scratch" \
  ./install.sh >/dev/null

output=$("$install_dir/courier" version)
printf '%s\n' "$output" | grep -Fq "courier $version"
[ -z "$(find "$scratch" -mindepth 1 -print -quit)" ]
printf '%s\n' "Verified POSIX installer against the $version snapshot"

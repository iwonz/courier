#!/bin/sh

set -eu

[ -f dist/metadata.json ] || { printf '%s\n' "dist/metadata.json is missing" >&2; exit 1; }
version=$(sed -n 's/.*"version":"\([^"]*\)".*/\1/p' dist/metadata.json)
[ -n "$version" ] || { printf '%s\n' "cannot read release version" >&2; exit 1; }

require_file() {
  [ -f "$1" ] || { printf '%s\n' "missing release artifact: $1" >&2; exit 1; }
  name=$(basename "$1")
  expected=$(checksum_for "$name")
  [ -n "$expected" ] || { printf '%s\n' "checksums.txt has no entry for $name" >&2; exit 1; }
  if command -v sha256sum >/dev/null 2>&1; then
    actual=$(sha256sum "$1" | awk '{ print tolower($1) }')
  else
    actual=$(shasum -a 256 "$1" | awk '{ print tolower($1) }')
  fi
  [ "$actual" = "$expected" ] || { printf '%s\n' "checksum mismatch for $name" >&2; exit 1; }
}

checksum_for() {
  awk -v name="$1" '$2 == name || $2 == "*" name { print tolower($1); exit }' dist/checksums.txt
}

require_checksum() {
  [ -n "$(checksum_for "$1")" ] || { printf '%s\n' "checksums.txt has no entry for $1" >&2; exit 1; }
}

require_named_artifact() {
  require_checksum "$1"
  grep -Fq "\"name\":\"$1\"" dist/artifacts.json || {
    printf '%s\n' "artifacts.json has no release artifact: $1" >&2
    exit 1
  }
}

for target_os in darwin linux windows; do
  for target_arch in amd64 arm64; do
    require_file "dist/courier_${version}_${target_os}_${target_arch}.tar.gz"
    raw="dist/courier_${version}_${target_os}_${target_arch}"
    if [ "$target_os" = windows ]; then
      raw="${raw}.exe"
      require_file "dist/courier_${version}_${target_os}_${target_arch}.zip"
    fi
    require_named_artifact "$(basename "$raw")"
  done
done

for helper_os in freebsd openbsd netbsd; do
  for helper_arch in amd64 arm64; do
    require_file "dist/courier_${version}_${helper_os}_${helper_arch}.tar.gz"
  done
done
require_file "dist/courier_${version}_dragonfly_amd64.tar.gz"
require_file "dist/courier_${version}_source.tar.gz"

for target_arch in amd64 arm64; do
  require_file "dist/courier_${version}_linux_${target_arch}.deb"
  require_file "dist/courier_${version}_linux_${target_arch}.rpm"
  require_file "dist/courier_${version}_linux_${target_arch}.apk"
  require_file "dist/courier_${version}_linux_${target_arch}.pkg.tar.zst"
done

[ -f dist/homebrew/Formula/courier.rb ]
[ -f dist/scoop/bucket/courier.json ]

printf '%s\n' "Verified Courier release artifact matrix for $version"

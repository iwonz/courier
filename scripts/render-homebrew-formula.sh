#!/bin/sh

set -eu

if [ "$#" -ne 6 ]; then
  printf '%s\n' "usage: $0 VERSION SHA256 COMMIT DATE URL OUTPUT" >&2
  exit 2
fi

version=$1
sha256=$2
commit=$3
date=$4
url=$5
output=$6

case "$version" in *[!0-9A-Za-z.-]*|'') printf '%s\n' "invalid Formula version" >&2; exit 2;; esac
case "$sha256" in *[!0-9a-f]*|'') printf '%s\n' "invalid source SHA-256" >&2; exit 2;; esac
[ "${#sha256}" -eq 64 ] || { printf '%s\n' "invalid source SHA-256 length" >&2; exit 2; }
case "$commit" in *[!0-9a-f]*|'') printf '%s\n' "invalid commit" >&2; exit 2;; esac
case "$date" in *[!0-9T:+.Z-]*|'') printf '%s\n' "invalid build date" >&2; exit 2;; esac
case "$url" in https://*|file://*) ;; *) printf '%s\n' "invalid source URL" >&2; exit 2;; esac

template=$(CDPATH='' cd -- "$(dirname "$0")" && pwd)/homebrew/courier.rb.tmpl
directory=$(dirname "$output")
mkdir -p "$directory"
temporary="$output.tmp.$$"
trap 'rm -f "$temporary"' EXIT HUP INT TERM
sed \
  -e "s|@VERSION@|$version|g" \
  -e "s|@SHA256@|$sha256|g" \
  -e "s|@COMMIT@|$commit|g" \
  -e "s|@DATE@|$date|g" \
  -e "s|@URL@|$url|g" \
  "$template" >"$temporary"
mv "$temporary" "$output"
trap - EXIT HUP INT TERM

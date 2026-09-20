#!/bin/sh

set -eu

[ -f dist/metadata.json ] || { printf '%s\n' "dist/metadata.json is missing" >&2; exit 1; }
version=$(sed -n 's/.*"version":"\([^"]*\)".*/\1/p' dist/metadata.json)
commit=$(sed -n 's/.*"commit":"\([^"]*\)".*/\1/p' dist/metadata.json)
date=$(sed -n 's/.*"date":"\([^"]*\)".*/\1/p' dist/metadata.json)
source="dist/courier_${version}_source.tar.gz"
[ -f "$source" ] || { printf '%s\n' "snapshot source archive is missing: $source" >&2; exit 1; }
sha256=$(awk -v name="$(basename "$source")" '$2 == name || $2 == "*" name { print tolower($1); exit }' dist/checksums.txt)
[ -n "$sha256" ] || { printf '%s\n' "snapshot source checksum is missing" >&2; exit 1; }
url="https://github.com/iwonz/courier/releases/download/v${version}/$(basename "$source")"
./scripts/render-homebrew-formula.sh "$version" "$sha256" "$commit" "$date" "$url" dist/homebrew/Formula/courier.rb
ruby -c dist/homebrew/Formula/courier.rb >/dev/null

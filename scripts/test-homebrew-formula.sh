#!/bin/sh

set -eu

formula=dist/homebrew/Formula/courier.rb
[ -f "$formula" ] || { printf '%s\n' "snapshot Formula is missing" >&2; exit 1; }
ruby -c "$formula" >/dev/null
grep -Fq 'depends_on "go" => :build' "$formula"
grep -Fq 'ENV["CGO_ENABLED"] = "0"' "$formula"
grep -Fq '"./cmd/courier"' "$formula"
grep -Fq 'courier version' "$formula"
if grep -Eq '^[[:space:]]+version "' "$formula"; then
  printf '%s\n' "Formula version must be inferred from the source archive URL" >&2
  exit 1
fi

version=$(sed -n 's/.*"version":"\([^"]*\)".*/\1/p' dist/metadata.json)
commit=$(sed -n 's/.*"commit":"\([^"]*\)".*/\1/p' dist/metadata.json)
date=$(sed -n 's/.*"date":"\([^"]*\)".*/\1/p' dist/metadata.json)
source=$(pwd)/dist/courier_${version}_source.tar.gz
sha256=$(shasum -a 256 "$source" | awk '{ print tolower($1) }')
temporary=$(mktemp -d)
trap 'rm -rf "$temporary"' EXIT HUP INT TERM
archive_root="$temporary/courier-$version"
tar -xzf "$source" -C "$temporary"
[ -d "$archive_root" ] || { printf '%s\n' "snapshot source prefix is invalid" >&2; exit 1; }
ldflags="-s -w -X github.com/iwonz/courier/internal/buildinfo.Version=$version -X github.com/iwonz/courier/internal/buildinfo.Commit=$commit -X github.com/iwonz/courier/internal/buildinfo.Date=$date"
(
  cd "$archive_root"
  CGO_ENABLED=0 go build -trimpath -ldflags "$ldflags" -o "$temporary/courier" ./cmd/courier
)
"$temporary/courier" version | grep -Fq "courier $version"

if ! command -v brew >/dev/null 2>&1; then
  printf '%s\n' "Homebrew is unavailable; Formula syntax and equivalent source build verified"
  exit 0
fi
doctor=$(brew doctor 2>&1 || :)
case "$doctor" in
  *"Your Xcode "*" is too outdated."*)
    printf '%s\n' "Homebrew audit/install skipped because the host Xcode is below Homebrew's minimum; Formula syntax and equivalent source build verified"
    exit 0
    ;;
esac
if brew list --formula courier >/dev/null 2>&1; then
  printf '%s\n' "refusing to replace an existing courier Formula during snapshot verification" >&2
  exit 1
fi

tap="courier-local/snapshot-$$"
installed=false
tapped=false
developer_was_on=false
if brew config | grep -q '^HOMEBREW_DEVELOPER:'; then developer_was_on=true; fi
cleanup() {
  if [ "$installed" = true ]; then brew uninstall --force courier >/dev/null; fi
  if [ "$tapped" = true ]; then brew untap --force "$tap" >/dev/null; fi
  if [ "$developer_was_on" = false ]; then brew developer off >/dev/null; fi
  rm -rf "$temporary"
}
trap cleanup EXIT HUP INT TERM
HOMEBREW_NO_AUTO_UPDATE=1 brew tap-new --no-git "$tap" >/dev/null
tapped=true
if [ "$developer_was_on" = false ]; then brew developer off >/dev/null; fi
tap_root=$(brew --repository "$tap")
local_formula="$tap_root/Formula/courier.rb"
./scripts/render-homebrew-formula.sh "$version" "$sha256" "$commit" "$date" "file://$source" "$local_formula"
brew audit --strict --formula "$tap/courier"
HOMEBREW_NO_AUTO_UPDATE=1 brew install --build-from-source --formula "$tap/courier"
installed=true
brew test "$tap/courier"
"$(brew --prefix courier)/bin/courier" version
brew uninstall --force courier
installed=false
brew untap --force "$tap" >/dev/null
tapped=false
rm -rf "$temporary"
trap - EXIT HUP INT TERM

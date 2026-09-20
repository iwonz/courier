#!/bin/sh

set -eu

cd "$(dirname "$0")/.."

version=${1:-}
change=${2:-}
message=${3:-}

case "$version" in
  v[0-9]*.[0-9]*.[0-9]*) version=${version#v} ;;
  [0-9]*.[0-9]*.[0-9]*) ;;
  *) printf '%s\n' "Usage: make ship VERSION=MAJOR.MINOR.PATCH CHANGE=<openspec-change> MESSAGE='<conventional commit>'" >&2; exit 2 ;;
esac
printf '%s\n' "$version" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$' || {
  printf '%s\n' "Version must use semantic MAJOR.MINOR.PATCH syntax" >&2
  exit 2
}
[ -n "$message" ] || message="chore: release v$version"

branch=$(git branch --show-current)
[ "$branch" = main ] || { printf '%s\n' "Shipping requires main, not ${branch:-a detached HEAD}" >&2; exit 1; }
[ -n "$(git status --porcelain)" ] || { printf '%s\n' "There are no completed changes to ship" >&2; exit 1; }

./scripts/check-release-config.sh

if [ "${COURIER_SHIP_YES:-0}" != 1 ]; then
  git status --short
  printf '%s' "Verify, commit, push main, release v$version, and deploy Pages? [y/N] "
  read -r answer
  case "$answer" in
    y|Y|yes|YES) ;;
    *) printf '%s\n' "Ship cancelled"; exit 1 ;;
  esac
fi

if [ -n "$change" ] && [ -d "openspec/changes/$change" ]; then
  if grep -Eq '^- \[ \]' "openspec/changes/$change/tasks.md"; then
    printf '%s\n' "OpenSpec change still has incomplete tasks: $change" >&2
    exit 1
  fi
  openspec validate "$change" --strict --no-interactive
  openspec archive "$change" --yes
elif [ -n "$change" ]; then
  found=false
  for archived in openspec/changes/archive/*-"$change"; do
    if [ -d "$archived" ]; then found=true; break; fi
  done
  [ "$found" = true ] || { printf '%s\n' "OpenSpec change does not exist: $change" >&2; exit 1; }
fi

for active in openspec/changes/*; do
  [ -d "$active" ] || continue
  [ "$(basename "$active")" = archive ] && continue
  printf '%s\n' "Active OpenSpec change remains unarchived: $(basename "$active")" >&2
  exit 1
done

version_file=$(mktemp "${TMPDIR:-/tmp}/courier-version.XXXXXX")
trap 'rm -f "$version_file"' EXIT HUP INT TERM
awk -v version="$version" '
  /^target_release:/ { print "target_release: " version; found=1; next }
  { print }
  END { if (!found) exit 1 }
' docs/cli-contract.yaml >"$version_file"
mv "$version_file" docs/cli-contract.yaml
trap - EXIT HUP INT TERM
go run ./cmd/contractdoc --write

make verify
[ -n "$(git status --porcelain)" ] || { printf '%s\n' "Verification left nothing to commit" >&2; exit 1; }

git add -A
git diff --cached --quiet && { printf '%s\n' "There are no staged changes to commit" >&2; exit 1; }
# The exact pre-commit target already completed immediately above; avoid running
# the same multi-platform release-candidate gate twice for the identical tree.
git commit --no-verify -m "$message"
commit=$(git rev-parse HEAD)
git push origin main

COURIER_RELEASE_YES=1 COURIER_RELEASE_VERIFIED_COMMIT="$commit" ./scripts/release.sh "$version"

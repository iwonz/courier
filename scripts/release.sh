#!/bin/sh

set -eu

cd "$(dirname "$0")/.."

./scripts/check-release-config.sh

branch=$(git branch --show-current)
[ "$branch" = main ] || { printf '%s\n' "Releases must be created from main, not $branch" >&2; exit 1; }
[ -z "$(git status --porcelain)" ] || { printf '%s\n' "The working tree must be clean before release" >&2; exit 1; }

git fetch origin --tags --quiet
git diff --quiet HEAD origin/main || { printf '%s\n' "Local main must exactly match origin/main" >&2; exit 1; }

version=${1:-}
if [ -z "$version" ]; then
  latest=$(git tag --list 'v[0-9]*' --sort=-v:refname | sed -n '1p')
  printf '%s\n' "Latest release: ${latest:-none}"
  printf '%s' "New version (MAJOR.MINOR.PATCH): "
  read -r version
fi
case "$version" in
  v*) tag=$version ;;
  *) tag="v$version" ;;
esac
printf '%s\n' "$tag" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+$' || {
  printf '%s\n' "Version must use semantic vMAJOR.MINOR.PATCH syntax" >&2
  exit 1
}
if git rev-parse --verify --quiet "refs/tags/$tag" >/dev/null; then
  printf '%s\n' "Tag already exists locally: $tag" >&2
  exit 1
fi
if git ls-remote --exit-code --tags origin "refs/tags/$tag" >/dev/null 2>&1; then
  printf '%s\n' "Tag already exists on origin: $tag" >&2
  exit 1
fi

make verify
[ -z "$(git status --porcelain)" ] || { printf '%s\n' "Verification changed tracked files; review them before release" >&2; exit 1; }

if [ "${COURIER_RELEASE_YES:-0}" != 1 ]; then
  printf '%s' "Create and push $tag from $(git rev-parse --short HEAD)? [y/N] "
  read -r answer
  case "$answer" in
    y|Y|yes|YES) ;;
    *) printf '%s\n' "Release cancelled"; exit 1 ;;
  esac
fi

git tag -a "$tag" -m "Courier $tag"
if ! git push origin "$tag"; then
  printf '%s\n' "Push failed. The local annotated tag remains at $tag; fix connectivity and run: git push origin $tag" >&2
  exit 1
fi
printf '%s\n' "Release workflow started for $tag: https://github.com/iwonz/courier/actions/workflows/release.yml"

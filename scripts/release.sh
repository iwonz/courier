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

verified_commit=${COURIER_RELEASE_VERIFIED_COMMIT:-}
current_commit=$(git rev-parse HEAD)
if [ "$verified_commit" = "$current_commit" ]; then
  printf '%s\n' "Using the full verification already completed for $current_commit"
else
  make verify
fi
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

run_id=
attempt=0
while [ "$attempt" -lt 60 ]; do
  run_id=$(gh run list --repo iwonz/courier --workflow release.yml --event push --limit 30 --json databaseId,headBranch,headSha --jq ".[] | select(.headBranch == \"$tag\" and .headSha == \"$current_commit\") | .databaseId" | sed -n '1p')
  [ -n "$run_id" ] && break
  attempt=$((attempt + 1))
  sleep 2
done
[ -n "$run_id" ] || {
  printf '%s\n' "GitHub did not register the release workflow for $tag" >&2
  exit 1
}

printf '%s\n' "Watching release workflow $run_id for $tag"
gh run watch "$run_id" --repo iwonz/courier --exit-status --interval 10

git fetch origin main --quiet
git merge --ff-only origin/main
./scripts/publish-pages.sh

release_url=$(gh release view "$tag" --repo iwonz/courier --json url --jq .url)
pages_url=$(gh api repos/iwonz/courier/pages --jq .html_url)
printf '%s\n' "Published Courier $tag: $release_url"
printf '%s\n' "Published Courier landing: $pages_url"

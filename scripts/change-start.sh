#!/bin/sh

set -eu

cd "$(dirname "$0")/.."

type=${1:-}
change=${2:-}

case "$type" in
  feat|fix|chore|docs|refactor|test|build|ci|perf|style) ;;
  *)
    printf '%s\n' "TYPE must be one of: feat, fix, chore, docs, refactor, test, build, ci, perf, style" >&2
    exit 2
    ;;
esac
printf '%s\n' "$change" | grep -Eq '^[a-z0-9]+(-[a-z0-9]+)*$' || {
  printf '%s\n' "CHANGE must use lowercase kebab-case" >&2
  exit 2
}

branch=$(git branch --show-current)
[ "$branch" = main ] || {
  printf '%s\n' "New changes must start from main, not ${branch:-a detached HEAD}" >&2
  exit 1
}
[ -z "$(git status --porcelain)" ] || {
  printf '%s\n' "The working tree must be clean before starting a change" >&2
  exit 1
}

git fetch origin main --quiet
local_revision=$(git rev-parse HEAD)
remote_revision=$(git rev-parse origin/main)
[ "$local_revision" = "$remote_revision" ] || {
  printf '%s\n' "Local main must exactly match origin/main" >&2
  exit 1
}

for active in openspec/changes/*; do
  [ -d "$active" ] || continue
  [ "$(basename "$active")" = archive ] && continue
  printf '%s\n' "Finish and ship the active OpenSpec change first: $(basename "$active")" >&2
  exit 1
done

feature_branch="$type/$change"
git show-ref --verify --quiet "refs/heads/$feature_branch" && {
  printf '%s\n' "Local branch already exists: $feature_branch" >&2
  exit 1
}
git ls-remote --exit-code --heads origin "refs/heads/$feature_branch" >/dev/null 2>&1 && {
  printf '%s\n' "Remote branch already exists: $feature_branch" >&2
  exit 1
}

git switch -c "$feature_branch"
if ! openspec new change "$change"; then
  git switch main
  git branch -d "$feature_branch"
  exit 1
fi

printf '%s\n' "Started $feature_branch with OpenSpec change $change"

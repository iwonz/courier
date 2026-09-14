#!/bin/sh

set -eu

cd "$(dirname "$0")/.."

repository=iwonz/courier
workflow=pages.yml

remote=$(git remote get-url origin 2>/dev/null || true)
case "$remote" in
  *github.com/iwonz/courier.git|*github.com:iwonz/courier.git) ;;
  *) printf '%s\n' "origin must point to https://github.com/iwonz/courier" >&2; exit 1 ;;
esac

branch=$(git branch --show-current)
[ "$branch" = main ] || {
  printf '%s\n' "Pages must be published from main, not ${branch:-a detached HEAD}" >&2
  exit 1
}
[ -z "$(git status --porcelain)" ] || {
  printf '%s\n' "The working tree must be clean before Pages publication" >&2
  exit 1
}

command -v gh >/dev/null 2>&1 || {
  printf '%s\n' "Install GitHub CLI and run gh auth login." >&2
  exit 1
}
gh auth status --hostname github.com >/dev/null

git fetch origin main --quiet
local_revision=$(git rev-parse HEAD)
remote_revision=$(git rev-parse origin/main)
[ "$local_revision" = "$remote_revision" ] || {
  printf '%s\n' "Local main must exactly match origin/main" >&2
  exit 1
}

make pages-build
[ -z "$(git status --porcelain)" ] || {
  printf '%s\n' "The Pages build changed tracked files; review and commit them first" >&2
  exit 1
}

if current_build_type=$(gh api "repos/$repository/pages" --jq .build_type 2>/dev/null); then
  if [ "$current_build_type" != workflow ]; then
    gh api --method PUT "repos/$repository/pages" -f build_type=workflow >/dev/null
  fi
else
  gh api --method POST "repos/$repository/pages" -f build_type=workflow >/dev/null
fi

previous_run=$(gh run list --repo "$repository" --workflow "$workflow" --branch main --event workflow_dispatch --limit 1 --json databaseId --jq '.[0].databaseId // 0')
gh workflow run "$workflow" --repo "$repository" --ref main

run_id=
attempt=0
while [ "$attempt" -lt 30 ]; do
  run_id=$(gh run list --repo "$repository" --workflow "$workflow" --branch main --event workflow_dispatch --limit 20 --json databaseId,headSha --jq ".[] | select(.headSha == \"$local_revision\" and .databaseId != $previous_run) | .databaseId" | sed -n '1p')
  [ -n "$run_id" ] && break
  attempt=$((attempt + 1))
  sleep 2
done
[ -n "$run_id" ] || {
  printf '%s\n' "GitHub did not register the Pages workflow dispatch" >&2
  exit 1
}

gh run watch "$run_id" --repo "$repository" --exit-status --interval 5
pages_url=$(gh api "repos/$repository/pages" --jq .html_url)
printf '%s\n' "Published Courier landing: $pages_url"

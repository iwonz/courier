#!/bin/sh

set -eu

cd "$(dirname "$0")/.."

version=${1:-}
change=${2:-}
message=${3:-}
openspec_command=${COURIER_OPENSPEC:-openspec}
make_command=${COURIER_MAKE:-make}
contract_writer=${COURIER_CONTRACT_WRITER:-go}
release_config=${COURIER_RELEASE_CONFIG:-./scripts/check-release-config.sh}
release_script=${COURIER_RELEASE_SCRIPT:-./scripts/release.sh}

[ -n "$change" ] || {
  printf '%s\n' "Usage: make ship [VERSION=MAJOR.MINOR.PATCH] CHANGE=<openspec-change> MESSAGE='<conventional commit>'" >&2
  exit 2
}
printf '%s\n' "$change" | grep -Eq '^[a-z0-9]+(-[a-z0-9]+)*$' || {
  printf '%s\n' "CHANGE must use lowercase kebab-case" >&2
  exit 2
}
printf '%s\n' "$message" | grep -Eq '^(feat|fix|chore|docs|refactor|test|build|ci|perf|style)(\([a-z0-9._/-]+\))?!?: .+' || {
  printf '%s\n' "MESSAGE must be a non-empty conventional commit" >&2
  exit 2
}

feature_branch=$(git branch --show-current)
case "$feature_branch" in
  */"$change") ;;
  main|"")
    printf '%s\n' "Shipping requires the feature branch created for $change, not ${feature_branch:-a detached HEAD}" >&2
    exit 1
    ;;
  *)
    printf '%s\n' "Branch $feature_branch does not match OpenSpec change $change" >&2
    exit 1
    ;;
esac

case "$feature_branch" in
  feat/*|fix/*|chore/*|docs/*|refactor/*|test/*|build/*|ci/*|perf/*|style/*) ;;
  *) printf '%s\n' "Unsupported feature branch type: $feature_branch" >&2; exit 1 ;;
esac

git fetch origin main --tags --quiet
base_revision=$(git rev-parse origin/main)
head_revision=$(git rev-parse HEAD)
local_main=$(git rev-parse main)
[ "$head_revision" = "$base_revision" ] || {
  printf '%s\n' "The feature branch must contain no commits before ship; ship creates the single change commit" >&2
  exit 1
}
[ "$local_main" = "$base_revision" ] || {
  printf '%s\n' "Local main must exactly match origin/main" >&2
  exit 1
}
[ -n "$(git status --porcelain)" ] || {
  printf '%s\n' "There are no completed changes to ship" >&2
  exit 1
}

active_count=0
active_name=
for active in openspec/changes/*; do
  [ -d "$active" ] || continue
  [ "$(basename "$active")" = archive ] && continue
  active_count=$((active_count + 1))
  active_name=$(basename "$active")
done
archive_count=0
archived_change=
for archived in openspec/changes/archive/*-"$change"; do
  [ -d "$archived" ] || continue
  [ -n "$(git status --porcelain -- "$archived")" ] || continue
  archive_count=$((archive_count + 1))
  archived_change=$archived
done

change_state=active
change_path="openspec/changes/$change"
if [ "$active_count" -eq 1 ] && [ "$active_name" = "$change" ]; then
  "$openspec_command" validate "$change" --strict --no-interactive
elif [ "$active_count" -eq 0 ] && [ "$archive_count" -eq 1 ]; then
  change_state=archived
  change_path=$archived_change
  "$openspec_command" validate --all --strict --no-interactive
  printf '%s\n' "Resuming ship from archived OpenSpec change $change"
else
  printf '%s\n' "Exactly one active OpenSpec change named $change, or its uncommitted ship archive, is required" >&2
  exit 1
fi
[ -f "$change_path/tasks.md" ] || {
  printf '%s\n' "OpenSpec task list is missing for $change" >&2
  exit 1
}
if grep -Eq '^- \[ \]' "$change_path/tasks.md"; then
  printf '%s\n' "OpenSpec change still has incomplete tasks: $change" >&2
  exit 1
fi
"$release_config"

if [ -z "$version" ]; then
  latest=$(git tag --list 'v[0-9]*.[0-9]*.[0-9]*' --sort=-v:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sed -n '1p')
  [ -n "$latest" ] || {
    printf '%s\n' "Cannot derive the next patch without an existing semantic release tag" >&2
    exit 1
  }
  current=${latest#v}
  major=${current%%.*}
  remainder=${current#*.}
  minor=${remainder%%.*}
  patch=${remainder#*.}
  version="$major.$minor.$((patch + 1))"
fi
case "$version" in
  v[0-9]*.[0-9]*.[0-9]*) version=${version#v} ;;
  [0-9]*.[0-9]*.[0-9]*) ;;
  *) printf '%s\n' "VERSION must use semantic MAJOR.MINOR.PATCH syntax" >&2; exit 2 ;;
esac
printf '%s\n' "$version" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$' || {
  printf '%s\n' "VERSION must use semantic MAJOR.MINOR.PATCH syntax" >&2
  exit 2
}
tag="v$version"
git rev-parse --verify --quiet "refs/tags/$tag" >/dev/null && {
  printf '%s\n' "Tag already exists locally: $tag" >&2
  exit 1
}
git ls-remote --exit-code --tags origin "refs/tags/$tag" >/dev/null 2>&1 && {
  printf '%s\n' "Tag already exists on origin: $tag" >&2
  exit 1
}

if [ "${COURIER_SHIP_YES:-0}" != 1 ]; then
  git status --short
  printf '%s' "Archive $change, verify, merge, publish $tag, and deploy Pages? [y/N] "
  read -r answer
  case "$answer" in
    y|Y|yes|YES) ;;
    *) printf '%s\n' "Ship cancelled"; exit 1 ;;
  esac
fi

if [ "$change_state" = active ]; then
  "$openspec_command" archive "$change" --yes
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
"$contract_writer" run ./cmd/contractdoc --write

"$make_command" verify
[ -n "$(git status --porcelain)" ] || {
  printf '%s\n' "Verification left nothing to commit" >&2
  exit 1
}

git add -A
git diff --cached --quiet && {
  printf '%s\n' "There are no staged changes to commit" >&2
  exit 1
}
git commit --no-verify -m "$message"
feature_commit=$(git rev-parse HEAD)

git fetch origin main --tags --quiet
current_remote_main=$(git rev-parse origin/main)
current_local_main=$(git rev-parse main)
[ "$current_remote_main" = "$base_revision" ] || {
  printf '%s\n' "origin/main changed during ship; the feature commit is preserved on $feature_branch" >&2
  exit 1
}
[ "$current_local_main" = "$base_revision" ] || {
  printf '%s\n' "Local main changed during ship; the feature commit is preserved on $feature_branch" >&2
  exit 1
}

git switch main
git merge --ff-only "$feature_branch"
[ "$(git rev-parse HEAD)" = "$feature_commit" ] || {
  printf '%s\n' "The feature branch did not fast-forward main to the verified commit" >&2
  exit 1
}
COURIER_PUSH_FLOW=ship git push origin main

COURIER_RELEASE_FLOW=ship \
COURIER_PUSH_FLOW=ship \
COURIER_PAGES_FLOW=ship \
COURIER_RELEASE_YES=1 \
COURIER_RELEASE_VERIFIED_COMMIT="$feature_commit" \
"$release_script" "$version"

git fetch origin main --quiet
git merge --ff-only origin/main
[ -z "$(git status --porcelain)" ] || {
  printf '%s\n' "Ship completed publication but the final main worktree is not clean" >&2
  exit 1
}
[ "$(git rev-parse HEAD)" = "$(git rev-parse origin/main)" ] || {
  printf '%s\n' "Ship completed publication but local main does not match origin/main" >&2
  exit 1
}
git branch -d "$feature_branch"

printf '%s\n' "Shipped $tag; main is clean and synchronized with origin/main"

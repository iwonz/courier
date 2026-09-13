#!/bin/sh

set -eu

command -v gh >/dev/null 2>&1 || { printf '%s\n' "Install GitHub CLI and run gh auth login." >&2; exit 1; }
gh auth status --hostname github.com >/dev/null

expected_repository=github.com/iwonz/courier
remote=$(git remote get-url origin 2>/dev/null || true)
case "$remote" in
  *"$expected_repository"*|*github.com:iwonz/courier*) ;;
  *) printf '%s\n' "origin must point to https://github.com/iwonz/courier" >&2; exit 1 ;;
esac

for repository in iwonz/courier; do
  gh repo view "$repository" --json nameWithOwner >/dev/null || {
    printf '%s\n' "Required GitHub repository is unavailable: $repository" >&2
    exit 1
  }
done

configured=$(gh secret list --repo iwonz/courier --json name --jq '.[].name')
for secret in NPM_TOKEN; do
  printf '%s\n' "$configured" | grep -qx "$secret" || {
    printf '%s\n' "Missing GitHub Actions secret: $secret" >&2
    exit 1
  }
done

printf '%s\n' "Courier repository and GitHub secret names are configured."

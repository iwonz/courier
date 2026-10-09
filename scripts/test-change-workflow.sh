#!/bin/sh
# shellcheck disable=SC2016

set -eu

source_root=$(CDPATH='' cd -- "$(dirname "$0")/.." && pwd)
suite_root=$(mktemp -d "${TMPDIR:-/tmp}/courier-workflow.XXXXXX")
trap 'rm -rf "$suite_root"' EXIT HUP INT TERM
fixture_number=0

fail() {
  printf '%s\n' "workflow test failed: $*" >&2
  exit 1
}

new_fixture() {
  fixture_number=$((fixture_number + 1))
  fixture_root="$suite_root/fixture-$fixture_number"
  fixture_work="$fixture_root/work"
  fixture_remote="$fixture_root/origin.git"
  fixture_bin="$fixture_root/bin"
  mkdir -p "$fixture_work/scripts" "$fixture_work/.githooks" "$fixture_work/docs" \
    "$fixture_work/openspec/changes/archive" "$fixture_bin"

  cp "$source_root/scripts/change-start.sh" "$fixture_work/scripts/change-start.sh"
  cp "$source_root/scripts/ship.sh" "$fixture_work/scripts/ship.sh"
  cp "$source_root/scripts/release.sh" "$fixture_work/scripts/release.sh"
  cp "$source_root/scripts/publish-pages.sh" "$fixture_work/scripts/publish-pages.sh"
  cp "$source_root/.githooks/pre-push" "$fixture_work/.githooks/pre-push"
  chmod +x "$fixture_work/scripts/change-start.sh" "$fixture_work/scripts/ship.sh" \
    "$fixture_work/scripts/release.sh" "$fixture_work/scripts/publish-pages.sh" \
    "$fixture_work/.githooks/pre-push"

  printf '%s\n' '# fixture' >"$fixture_work/README.md"
  printf '%s\n' 'target_release: 1.2.3' >"$fixture_work/docs/cli-contract.yaml"
  printf '%s\n' keep >"$fixture_work/openspec/changes/archive/.gitkeep"
  printf '%s\n' '#!/bin/sh' 'exit 0' >"$fixture_work/scripts/check-release-config.sh"
  chmod +x "$fixture_work/scripts/check-release-config.sh"

  printf '%s\n' \
    '#!/bin/sh' \
    'set -eu' \
    'action=${1:-}' \
    'case "$action" in' \
    '  new)' \
    '    [ "${2:-}" = change ]' \
    '    change=${3:-}' \
    '    mkdir -p "openspec/changes/$change/specs/example"' \
    '    printf "%s\n" "# Proposal" >"openspec/changes/$change/proposal.md"' \
    '    printf "%s\n" "- [ ] 1. Implement" >"openspec/changes/$change/tasks.md"' \
    '    ;;' \
    '  validate)' \
    '    [ "${TEST_FAIL_STAGE:-}" != openspec ]' \
    '    ;;' \
    '  archive)' \
    '    change=${2:-}' \
    '    destination="openspec/changes/archive/2026-01-01-$change"' \
    '    [ ! -e "$destination" ]' \
    '    mv "openspec/changes/$change" "$destination"' \
    '    ;;' \
    '  *) exit 2 ;;' \
    'esac' >"$fixture_bin/openspec"

  printf '%s\n' \
    '#!/bin/sh' \
    'set -eu' \
    '[ "${1:-}" = run ]' \
    'exit 0' >"$fixture_bin/go"

  printf '%s\n' \
    '#!/bin/sh' \
    'set -eu' \
    '[ "${1:-}" = verify ]' \
    'case "${TEST_FAIL_STAGE:-}" in' \
    '  verify) exit 1 ;;' \
    '  remote-shift)' \
    '    shift_work=$(mktemp -d "$TEST_ROOT/shift.XXXXXX")' \
    '    git clone --quiet "$TEST_REMOTE" "$shift_work"' \
    '    git -C "$shift_work" config user.name Fixture' \
    '    git -C "$shift_work" config user.email fixture@example.test' \
    '    printf "%s\n" shifted >>"$shift_work/README.md"' \
    '    git -C "$shift_work" add README.md' \
    '    git -C "$shift_work" commit --quiet -m "chore: shift main"' \
    '    git -C "$shift_work" push --quiet origin main' \
    '    rm -rf "$shift_work"' \
    '    ;;' \
    '  local-shift)' \
    '    tree=$(git rev-parse main^{tree})' \
    '    parent=$(git rev-parse main)' \
    '    commit=$(printf "%s\n" "local shift" | git commit-tree "$tree" -p "$parent")' \
    '    git branch -f main "$commit" >/dev/null' \
    '    ;;' \
    'esac' >"$fixture_bin/make"

  printf '%s\n' \
    '#!/bin/sh' \
    'set -eu' \
    '[ "${COURIER_RELEASE_FLOW:-}" = ship ]' \
    '[ "${COURIER_PUSH_FLOW:-}" = ship ]' \
    '[ "${COURIER_PAGES_FLOW:-}" = ship ]' \
    '[ "${COURIER_RELEASE_VERIFIED_COMMIT:-}" = "$(git rev-parse HEAD)" ]' \
    '[ "${TEST_FAIL_STAGE:-}" != publication ] || exit 1' \
    'tag="v${1:-}"' \
    'git tag -a "$tag" -m "Fixture $tag"' \
    'git push --quiet origin "$tag"' \
    'printf "%s\n" "$tag" >"$TEST_ROOT/published"' >"$fixture_work/scripts/release-fixture.sh"
  chmod +x "$fixture_bin/openspec" "$fixture_bin/go" "$fixture_bin/make" \
    "$fixture_work/scripts/release-fixture.sh"

  git init --bare --initial-branch=main "$fixture_remote" >/dev/null
  git -C "$fixture_work" init --initial-branch=main >/dev/null
  git -C "$fixture_work" config user.name Fixture
  git -C "$fixture_work" config user.email fixture@example.test
  git -C "$fixture_work" remote add origin "$fixture_remote"
  git -C "$fixture_work" add -A
  git -C "$fixture_work" commit --quiet -m 'chore: initialize fixture'
  git -C "$fixture_work" tag -a v1.2.3 -m 'Fixture v1.2.3'
  git -C "$fixture_work" push --quiet --set-upstream origin main
  git -C "$fixture_work" push --quiet origin v1.2.3
  git -C "$fixture_work" config core.hooksPath .githooks
}

prepare_change() {
  git -C "$fixture_work" switch --quiet -c feat/sample-change
  mkdir -p "$fixture_work/openspec/changes/sample-change/specs/example"
  printf '%s\n' '# Proposal' >"$fixture_work/openspec/changes/sample-change/proposal.md"
  printf '%s\n' '- [x] 1. Implement' >"$fixture_work/openspec/changes/sample-change/tasks.md"
  printf '%s\n' changed >>"$fixture_work/README.md"
}

run_ship() {
  log=$1
  fail_stage=${2:-}
  (
    cd "$fixture_work"
    PATH="$fixture_bin:$PATH" \
      TEST_ROOT="$fixture_root" \
      TEST_REMOTE="$fixture_remote" \
      TEST_FAIL_STAGE="$fail_stage" \
      COURIER_SHIP_YES=1 \
      COURIER_RELEASE_SCRIPT=./scripts/release-fixture.sh \
      ./scripts/ship.sh '' sample-change 'feat: fixture change'
  ) >"$log" 2>&1
}

expect_failure() {
  expected=$1
  log=$2
  shift 2
  if "$@" >"$log" 2>&1; then
    fail "command unexpectedly succeeded; expected $expected"
  fi
  grep -F "$expected" "$log" >/dev/null || {
    sed -n '1,160p' "$log" >&2
    fail "missing failure text: $expected"
  }
}

# change-start rejects dirty, out-of-sync, and concurrent active work.
new_fixture
printf '%s\n' dirty >>"$fixture_work/README.md"
expect_failure 'working tree must be clean' "$fixture_root/dirty.log" \
  env PATH="$fixture_bin:$PATH" sh -c "cd '$fixture_work' && ./scripts/change-start.sh feat next-change"

new_fixture
shift_work=$(mktemp -d "$fixture_root/advance.XXXXXX")
git clone --quiet "$fixture_remote" "$shift_work"
git -C "$shift_work" config user.name Fixture
git -C "$shift_work" config user.email fixture@example.test
printf '%s\n' advanced >>"$shift_work/README.md"
git -C "$shift_work" add README.md
git -C "$shift_work" commit --quiet -m 'chore: advance remote'
git -C "$shift_work" push --quiet origin main
rm -rf "$shift_work"
expect_failure 'must exactly match origin/main' "$fixture_root/sync.log" \
  env PATH="$fixture_bin:$PATH" sh -c "cd '$fixture_work' && ./scripts/change-start.sh feat next-change"

new_fixture
mkdir -p "$fixture_work/openspec/changes/existing"
printf '%s\n' active >"$fixture_work/openspec/changes/existing/tasks.md"
git -C "$fixture_work" add openspec/changes/existing/tasks.md
git -C "$fixture_work" commit --quiet -m 'chore: add active change'
COURIER_RECOVERY=1 git -C "$fixture_work" push --quiet origin main
expect_failure 'Finish and ship the active OpenSpec change first' "$fixture_root/active.log" \
  env PATH="$fixture_bin:$PATH" sh -c "cd '$fixture_work' && ./scripts/change-start.sh feat next-change"

new_fixture
(
  cd "$fixture_work"
  PATH="$fixture_bin:$PATH" ./scripts/change-start.sh feat next-change >/dev/null
)
[ "$(git -C "$fixture_work" branch --show-current)" = feat/next-change ] || fail 'change-start did not create the feature branch'
[ -f "$fixture_work/openspec/changes/next-change/tasks.md" ] || fail 'change-start did not scaffold OpenSpec'

# ship rejects incomplete work and preserves commits at every failure boundary.
new_fixture
prepare_change
printf '%s\n' '- [ ] 1. Implement' >"$fixture_work/openspec/changes/sample-change/tasks.md"
if run_ship "$fixture_root/incomplete.log" ''; then fail 'incomplete OpenSpec shipped'; fi
grep -F 'still has incomplete tasks' "$fixture_root/incomplete.log" >/dev/null || fail 'incomplete task diagnostic missing'

new_fixture
prepare_change
if run_ship "$fixture_root/verify.log" verify; then fail 'verification failure shipped'; fi
[ "$(git -C "$fixture_work" branch --show-current)" = feat/sample-change ] || fail 'verification failure lost the feature branch'
run_ship "$fixture_root/verify-resume.log" '' || {
  sed -n '1,200p' "$fixture_root/verify-resume.log" >&2
  fail 'ship did not resume after an archived verification failure'
}
[ "$(git -C "$fixture_work" branch --show-current)" = main ] || fail 'resumed ship did not finish on main'

new_fixture
prepare_change
if run_ship "$fixture_root/remote-shift.log" remote-shift; then fail 'origin/main shift shipped'; fi
grep -F 'origin/main changed during ship' "$fixture_root/remote-shift.log" >/dev/null || fail 'origin shift diagnostic missing'
git -C "$fixture_work" show-ref --verify --quiet refs/heads/feat/sample-change || fail 'origin shift lost the feature commit'

new_fixture
prepare_change
if run_ship "$fixture_root/local-shift.log" local-shift; then fail 'non-FF local main shipped'; fi
grep -F 'Local main changed during ship' "$fixture_root/local-shift.log" >/dev/null || fail 'non-FF diagnostic missing'

new_fixture
prepare_change
if run_ship "$fixture_root/publication.log" publication; then fail 'publication failure reported success'; fi
git -C "$fixture_work" show-ref --verify --quiet refs/heads/feat/sample-change || fail 'publication failure lost the feature branch'

# Manual protected pushes fail, while a successful ship derives the next patch,
# publishes it, removes the feature branch, and leaves synchronized clean main.
new_fixture
if printf '%s\n' 'refs/heads/main a refs/heads/main b' | "$fixture_work/.githooks/pre-push" origin "$fixture_remote" >/dev/null 2>&1; then
  fail 'pre-push allowed a manual main push'
fi
printf '%s\n' 'refs/heads/main a refs/heads/main b' | COURIER_PUSH_FLOW=ship "$fixture_work/.githooks/pre-push" origin "$fixture_remote"
expect_failure 'Direct release is disabled' "$fixture_root/release-guard.log" \
  sh -c "cd '$fixture_work' && ./scripts/release.sh 1.2.4"
expect_failure 'Direct Pages publication is disabled' "$fixture_root/pages-guard.log" \
  sh -c "cd '$fixture_work' && ./scripts/publish-pages.sh"

new_fixture
prepare_change
run_ship "$fixture_root/success.log" '' || {
  sed -n '1,200p' "$fixture_root/success.log" >&2
  fail 'successful ship failed'
}
IFS= read -r published <"$fixture_root/published"
[ "$published" = v1.2.4 ] || fail 'ship did not derive the next patch release'
[ "$(git -C "$fixture_work" branch --show-current)" = main ] || fail 'ship did not finish on main'
if git -C "$fixture_work" show-ref --verify --quiet refs/heads/feat/sample-change; then fail 'ship did not remove the feature branch'; fi
[ -z "$(git -C "$fixture_work" status --porcelain)" ] || fail 'ship did not leave a clean worktree'
[ "$(git -C "$fixture_work" rev-parse HEAD)" = "$(git -C "$fixture_work" rev-parse origin/main)" ] || fail 'ship did not synchronize main'
git -C "$fixture_work" show-ref --verify --quiet refs/tags/v1.2.4 || fail 'ship did not create the release tag'

printf '%s\n' 'Change workflow acceptance passed'

# Design

## Lifecycle

`make ship` is the outer transaction. It validates repository and credential configuration, optionally archives the named completed OpenSpec change, rejects any remaining active change, updates the contract target release, regenerates derived artifacts, and runs the full local release-candidate gate before staging anything.

After verification it creates one conventional commit and pushes `main`. It delegates tagged publication to `release.sh`, passing the exact verified commit so the inner command can avoid an identical duplicate gate without weakening the clean/synchronized/tag checks. The release command watches the tag workflow to completion, fast-forwards to the verified Formula commit written by that workflow, then uses the existing Pages publisher to build, dispatch, watch, and report the final site URL.

## Failure behavior

Every stage is fail-closed. No commit occurs before local verification; no tag occurs before the committed `main` is synchronized; no final success is reported before release and Pages workflows succeed. Published tags remain immutable if a downstream service fails, following the existing rerun policy.

## Scope safety

The command prints the complete worktree and requests confirmation unless `COURIER_SHIP_YES=1` is deliberately set. It stages the reviewed worktree only after all OpenSpec and verification checks pass. The named change must be complete or already archived, and no other active OpenSpec directory may remain.

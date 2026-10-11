## 1. Contract and planning

- [x] 1.1 Add the `preserve|unique` filename policy to operation planning, Cobra, CLI contract `0.17.0`, generated help/reference/README/landing data, and route validation
- [x] 1.2 Carry the optional policy through hosted runtime definitions without bumping IPC, registry, worker-definition, or administration-state versions

## 2. Durable upload activity

- [x] 2.1 Define and validate secret-free linked upload receipts plus bounded newest-first global persistence in the private delivery store
- [x] 2.2 Add atomic private destination receipt persistence for local and SSH backends and protect the metadata namespace from Courier transfer/publication
- [x] 2.3 Apply preserve/unique naming, hashing, actor/peer attribution, extraction accounting, receipt acknowledgement, and guarded failures to browser/webhook uploads

## 3. Quantity formatting and administration

- [x] 3.1 Add exact tested Go byte/bit/rate formatting and use it across interactive reports, progress, servers, and update
- [x] 3.2 Add shared localized TypeScript quantity formatting and use it in delivery and administration surfaces
- [x] 3.3 Project recent receipts through the guarded admin API/SSE and render an accessible localized responsive activity view

## 4. Acceptance and documentation

- [x] 4.1 Update browser/webhook, registry, administration, operational, CLI, and security documentation
- [x] 4.2 Add Go and web unit tests for option behavior, names, receipt durability/security/retention, API projection, formatting boundaries, localization, and layout
- [x] 4.3 Extend compiled serial E2E for preserve/unique browser and webhook uploads, destination/global receipt parity, restart persistence, and administration visibility
- [x] 4.4 Run focused Go/web/browser tests and strict OpenSpec validation
- [x] 4.5 Run the complete `make verify` release-candidate gate

## 5. Release

- [x] 5.1 Ship the archived change as the next patch and verify GitHub Release artifacts, npm, Scoop, Homebrew, Pages, deleted feature branch, and clean synchronized main

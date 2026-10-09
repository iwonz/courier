## 1. Path preflight

- [x] 1.1 Add the `--force-source-creation` operation option and path-role model, and verify route applicability and parser tests pass
- [x] 1.2 Implement shared local and SSH inspection, confirmation, ordered `0700` creation, and revalidation, and verify unit tests cover decline, noninteractive, forced, race, wrong-type, and partial-failure cases
- [x] 1.3 Apply preflight to finite transfers, extraction, and outgoing webhooks while preserving exact-file destinations, and verify focused transfer and webhook tests pass

## 2. Hosted lifecycle

- [x] 2.1 Move local and SSH hosted-path preparation before worker acquisition while retaining bounded credentials, and verify failures launch no worker and successful foreground/background starts pass
- [x] 2.2 Add safe registration-time missing/wrong-type classification and CLI exit mapping, and verify a disappearing path reports stage `preflight` with code `20` while unexpected failures stay redacted
- [x] 2.3 Extend lifecycle acceptance for detached readiness, listing, survival, and stop by delivery/server UUID, and verify the acceptance suite passes

## 3. Contract and generated documentation

- [x] 3.1 Upgrade the CLI contract to schema `3` and version `0.16.0`, add localized argument/flag descriptions and requiredness, and verify strict contract validation and Cobra parity pass
- [x] 3.2 Generate compact structured Cobra help metadata and install the help renderer, and verify command help contains status, descriptions, applicability, defaults, repeatability, dependencies, and conflicts
- [x] 3.3 Generate CLI reference and managed README argument/option tables from the contract, and verify deterministic freshness checks pass
- [x] 3.4 Update background-mode lifecycle prose across generated and authored documentation, and verify browser, webhook, and UI applicability is consistent

## 4. Landing presentation

- [x] 4.1 Extend generated landing JSON and TypeScript types with localized descriptions, requiredness, and readable applicability, and verify catalog and command-builder tests pass
- [x] 4.2 Render responsive accessible structured parameter rows in English and Russian without changing controls or command output, and verify unit, browser, and accessibility tests pass

## 5. Verification

- [x] 5.1 Validate the OpenSpec change strictly and verify every implemented scenario has automated coverage or an explicit acceptance check
- [x] 5.2 Run focused Go and web tests, format and generated-file checks, and verify they pass
- [x] 5.3 Run the full `make verify` release-candidate gate and verify race, exact coverage, cross-platform builds, browser tests, packaging, and snapshots pass
- [x] 5.4 Perform Playwright visual QA in English and Russian at desktop and narrow viewports and verify no overflow, clipping, or inaccessible controls

# Project handoff

Updated 2026-09-11. The reproducible ordinary-Stellar-payment CLI milestone is implemented and verified. Baseline was clean at 7ad01c9 with no CLI/config/output tests. Code/evidence commit dce4769 and core-pin/candidate-output/CI commit 1861a5f were pushed normally. Existing MIT licensing and read-only purpose remain unchanged.

Implemented: explicit embedded offline demo, strict JSON config, working output overrides/JSON stdout, canonical JSON/CSV input, closed account/network/window scope, coverage-aware Horizon ingestion, input-overwrite protection and nonzero error/discrepancy/unknown exits. Removed obsolete product-specific example config. Named adapters are not selected. Runtime pins are core 89aa1ec and connectors 902fa78, with no sibling replacements.

Local go test ./..., go vet ./... and binary build passed on Windows Go 1.24.4. An exported 1861a5f archive independently passed with GOWORK=off and a fresh GOMODCACHE: 5 test cases, vet, build and embedded demo. All three modules together had 50 passing test-runner cases including core subtests (not assertion count). go mod verify passed. Final checksum cleanup removes only unused old-core entries, without changing selected dependencies. A local-link check passed for repository Markdown.

Both demo and config runs produce **1 match, 4 discrepancies, 3 UNKNOWN**, exit **3**. Tests exercise mock Horizon success and HTTP failure through the CLI, config/flag/output errors, malformed config and deterministic machine-readable output. Testing fixed '-' being resolved as a filename and preserved the conservative orphan result when duplicate input IDs exist. Terminal and JSON ambiguity findings expose candidate operation IDs.

Remote CI passed at 1861a5f: https://github.com/LedgerParity/ledger-parity-cli/actions/runs/34653811330 (Go 1.22.2 and stable, two jobs, configured Linux race tests/vet/build and offline demo). Core dda2427 and connectors 83fd7ec also passed; see docs/VERIFICATION.md for all links. This final documentation/checksum cleanup does not alter runtime/CI. Initial API limits and missing core push runs were resolved using public run pages and an authenticated manual core dispatch; no repository policy changed. Local Windows -race lacked a C compiler; do not label Linux CI as local execution.

Read-only live testnet verification passed at 2026-09-11T22:14:14Z: operation 19878182387720193, one match and zero findings. See docs/LIVE_TESTNET_RESULT.json and scripts/Test-LiveRead.ps1. Expected application data was derived from public provider data; this proves API interoperability, not independently sourced business reconciliation, adoption or production readiness. No transactions were submitted. Later report-only candidate-ID additions were not represented as a fresh live run.

Start review with docs/REVIEWER_WALKTHROUGH.md. docs/DRIPS_SUBMISSION_PREP.md is an unsent appeal draft based on current official sources and the user's recollection of relevance/impact concerns; exact rejection/application unavailable. Next product gate: a consenting operator's sanitized export and known discrepancy. No acceptance, partnership or measured impact claim.

Remaining scope: trusted-provider history, no Soroban/path payments/account creation/merge settlement, no durable checkpoints, incomplete muxed account equivalence, no StrKey checksums, unverified named adapters and input-size/duplicate-JSON-key hardening. Contributor tasks are in docs/backlog.md. These are explicit limits, not completed production features.


## Implementation update 2026-09-12

Follow-up implementation: version 0.3.0-preview adds --bundle/--replay and target_app.format=sdp-csv. examples/sdp produces 1 match, 1 amount discrepancy and 1 unknown batch, exit 3; replay must reproduce the report. Inputs/evidence are bounded to 32 MiB with duplicate-key/depth validation for CLI JSON. Original raw files/config paths are not embedded. See docs/SDP_EVIDENCE.md. No real operator case, deployed SDP integration, production-readiness or Drips acceptance claim.

Verification completed: standalone published-dependency checks and byte-identical SDP evidence replay passed. Core 954f091, connectors b99a351 and CLI 1c6a32b remote CI succeeded; exact links and evidence distinctions are recorded in docs/VERIFICATION.md. The final verification-note commit is documentation only. Real operator validation and runtime RPC ingestion remain outstanding.

## Brand and documentation handover 2026-09-12

Org brand assets (mark, tile, lockup, banner) were generated from original SVG sources in `brand/` with `@resvg/resvg-js`; palette tokens and usage are documented in `brand/palette.md` and `brand/README.md`. A static documentation site (`ledger-parity-docs/`, build-free) is live at https://ledgerparity.github.io/ (commit df9aa0f; HTTPS 200 confirmed during publish). GitHub profile content for the org (`LedgerParity/.github`) is staged locally in `org-github/`; the repository has not yet been created on GitHub, so the org profile page is not yet rendered. Repo READMEs now include the brand banner and Go/license/CI badges plus a documentation link. `SECURITY.md` existed in all three repos; governance and code of conduct live in the org-github content. Connectors/docs push required a Git Credential Manager sign-in (completed). Core CI historically required manual dispatch; connectors/cli pushes trigger runs automatically.

## Resumption update 2026-09-13

Recent commits added the HTML dashboard, local checksum flags and a separate Soroban contract repository. The contract checkout with Git history is ../-ledger-parity-contract; ../ledger-parity-contract also exists without Git metadata. Neither was changed in this follow-up.

Fixed local verification: accept the positional report for --verify-check, hash only after successful report output, preserve discrepancy/UNKNOWN exits, validate proof hashes and bounded strict JSON, stream hashing, record UTC creation time, and reject incompatible modes and proof collisions. Proof writes exclusively create new files. Regression tests cover fresh/stale reports, tampering, malformed proofs, input/output collisions and failed report writes.

Local validation: go test ./..., go vet ./..., go build ./... on Windows. No new remote CI or deployment is claimed. Local checksums are unauthenticated byte comparisons; the CLI has no contract submission integration. Next product gate remains a consenting operator's sanitized export and known discrepancy. Dashboard review and contract testing/integration remain separate follow-up work. Older CI links below apply only to their recorded revisions.

## Dashboard continuation 2026-09-13

The offline single-file viewer now validates report shape and result/count consistency, enforces a 32 MiB file limit, escapes report-controlled HTML, shows coverage assertions and candidate operation IDs, and searches both internal and on-chain identities. Internal records and on-chain operations have separate totals. Amounts and IDs remain text. The picker supports keyboard use and repeated loads; failed loads hide stale results and out-of-order reads cannot replace newer reports. Large tables show the first 500 search matches with an explicit count.

Added five Node regression tests using the actual CLI demo, plus a CI step using the built CLI. All five passed locally. A headless Chrome DOM smoke check passed for real report rendering, inert injected HTML and search events; the desktop screenshot was visually inspected. No remote CI run or publish occurred. Browser artifacts are in the workspace's .local-checks directory. Runtime contract submission and operator validation remain outstanding.

## Final local verification 2026-09-13

Final review also fixed --version bypassing incompatible --verify-check flags and rejected non-string dashboard statuses. Go tests, vet and binary build passed; all five dashboard regression tests passed against that binary. A combined SDP report/evidence/proof run returned the expected exit 3, checksum verification returned 0, and offline replay returned 3 with byte-identical report output. Final executable: ../.local-checks/ledger-parity-final.exe. Evidence artifacts: ../.local-checks/final-dd9051effbc64e5bb940e3012a8a48d3/. This completes the local CLI verification and dashboard milestone; no contract deployment, remote CI result or real operator validation is claimed.

## Remote CI and contract review 2026-09-13

CLI 2633f88 passed both Go matrix jobs: https://github.com/LedgerParity/ledger-parity-cli/actions/runs/34782675975. The separate contract review fixed missing owner authorization, broken client registration in tests and missing hash/metadata bounds; added seven contract tests, a lockfile and Linux CI; and corrected misleading CLI/on-chain claims. Contract bb7c010 passed native SDK tests and release Wasm compilation: https://github.com/LedgerParity/-ledger-parity-contract/actions/runs/34784668593. See ../-ledger-parity-contract/REVIEW.md for findings and deployment gates. Local Rust execution lacks the MSVC linker; Linux CI supplies the test evidence. No deployed-network test, signing, operator validation or production-readiness claim.

## Contract storage lifecycle 2026-09-13

The next deployment gate is implemented in contract version 0.2.0: persistent per-report records, permissionless immutable renewal, a network-capped TTL target, and contract instance/code lifetime maintenance. This is a fresh-deployment storage layout, not migration of old instance records. See ../-ledger-parity-contract/LIFECYCLE.md and ROADMAP.md for recovery and next steps. Contract 29c8a24 passed all 13 native SDK tests and the release Wasm build: https://github.com/LedgerParity/-ledger-parity-contract/actions/runs/34785885294. The next gate is signed synthetic testnet registration/read-back, renewal/restoration and fee measurements on the current protocol. No deployment or network restoration has been performed; the CLI remains read-only with local checksums.

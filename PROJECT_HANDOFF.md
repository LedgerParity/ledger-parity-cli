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

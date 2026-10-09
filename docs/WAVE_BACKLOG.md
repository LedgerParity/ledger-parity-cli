# Focused engineering backlog

Six issue drafts reviewed against current source and existing open issues for the
October 9 submission. These are engineering tasks, not approval or earned points.
Proposed complexity is subject to maintainer review and Drips app configuration.

## 1. Add a portable read-only Stellar live-check command

## Description & Context

Windows-only live-read tooling makes the documented check difficult to reproduce elsewhere.

## Proposed Complexity

Medium (150 points proposed); planning label `complexity: medium`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Implement a Go command using read-only HTTP requests and explicit network/account/window scope.
- [ ] Label any derived expectation synthetic rather than operator validation.
- [ ] Report missing history and provider failure without asserting reconciliation success.
- [ ] Keep normal CI offline and test the transport with fixtures.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

scripts/Test-LiveRead.ps1; cmd/ledger-parity/; docs/VERIFICATION.md

## Verification

`go test ./...; go vet ./...; go build ./...; node --test dashboard/viewer.test.cjs`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## 2. Record a consenting operator reconciliation case

## Description & Context

Released tooling lacks independently supplied operator evidence.

## Proposed Complexity

Medium (150 points proposed); planning label `complexity: medium`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Use a permitted sanitized export and independently known discrepancy.
- [ ] Record exact field/network/time provenance and observed false positives.
- [ ] Document reproducible commands and operator confirmation without publishing private records.
- [ ] Keep absent evidence explicitly unresolved.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

docs/OPERATOR_VALIDATION.md; docs/VERIFICATION.md; examples/

## Verification

`go test ./...; go vet ./...; go build ./...; node --test dashboard/viewer.test.cjs`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## 3. Specify and test cross-version evidence replay compatibility

## Description & Context

Evidence bundles need an explicit version compatibility contract as releases evolve.

## Proposed Complexity

Medium (150 points proposed); planning label `complexity: medium`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Inventory the supported evidence version and unknown-version rejection.
- [ ] Add fixtures for old/current/unsupported versions.
- [ ] Reject changed digests and identity fields.
- [ ] Document supported migration without silently rewriting evidence.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

pkg/evidence/; docs/SDP_EVIDENCE.md; test fixtures

## Verification

`go test ./...; go vet ./...; go build ./...; node --test dashboard/viewer.test.cjs`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## 4. Measure large-export memory and latency before streaming

## Description & Context

The CLI bounds input but has no documented scale measurements for operator-sized exports.

## Proposed Complexity

Medium (150 points proposed); planning label `complexity: medium`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Add reproducible synthetic 1k/10k-record benchmarks.
- [ ] Record parsing/matching/report memory and latency.
- [ ] Propose a streaming boundary preserving duplicate-key checks and exact amount text.
- [ ] Avoid changing semantics without measured evidence.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

cmd/ledger-parity/; pkg/evidence/; benchmarks/

## Verification

`go test ./...; go vet ./...; go build ./...; node --test dashboard/viewer.test.cjs`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## 5. Expose concise UNKNOWN coverage reasons in operator output

## Description & Context

Operators need actionable explanations for incomplete coverage and ambiguity.

## Proposed Complexity

Medium (150 points proposed); planning label `complexity: medium`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Review which UNKNOWN reasons are already visible.
- [ ] Improve only missing CLI/dashboard explanations.
- [ ] Preserve candidate IDs and report counts.
- [ ] Test malformed and HTML-controlled text.
- [ ] Include screenshots for affected dashboard behavior.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

pkg/output/; dashboard/index.html; dashboard/viewer.test.cjs

## Verification

`go test ./...; go vet ./...; go build ./...; node --test dashboard/viewer.test.cjs`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## 6. Add a release archive acceptance check on Linux

## Description & Context

Archives are cross-built while recorded independent execution is Windows-specific.

## Proposed Complexity

Medium (150 points proposed); planning label `complexity: medium`.
Actual enrollment and points must be set in the Drips app after approval.

## Requirements & Acceptance Criteria

- [ ] Download a pinned Linux release archive in an opt-in workflow.
- [ ] Verify SHA-256 before execution.
- [ ] Run version/demo/proof/replay acceptance.
- [ ] Retain expected exit codes and assert dashboard/examples are present.
- [ ] Record the exact tested release.
- [ ] Include relevant positive/negative regression evidence; required CI passes.

## Relevant Files & Architecture

.github/workflows/; scripts/; RELEASE_NOTES.md

## Verification

`go test ./...; go vet ./...; go build ./...; node --test dashboard/viewer.test.cjs`. Network checks remain opt-in.

## Contribution Guidelines

Agree bounded scope with EthTobi through GitHub. Use a focused PR
with `Closes #<issue_id>`, actual check results, and remaining limitations.

## Published issue links

- [Add a portable read-only Stellar live-check command](https://github.com/LedgerParity/ledger-parity-cli/issues/4)
- [Record a consenting operator reconciliation case](https://github.com/LedgerParity/ledger-parity-cli/issues/5)
- [Specify and test cross-version evidence replay compatibility](https://github.com/LedgerParity/ledger-parity-cli/issues/6)
- [Measure large-export memory and latency before streaming](https://github.com/LedgerParity/ledger-parity-cli/issues/7)
- [Expose concise UNKNOWN coverage reasons in operator output](https://github.com/LedgerParity/ledger-parity-cli/issues/8)
- [Add a release archive acceptance check on Linux](https://github.com/LedgerParity/ledger-parity-cli/issues/9)

These issues are published but have not been enrolled into Drips Wave or assigned.

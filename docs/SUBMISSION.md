# Stellar payment reconciliation for operators

Prepared for the October 9, 2026 Stellar Wave submission.

## Purpose and implemented utility

A read-only command-line tool compares application payment exports against ordinary classic Stellar payment operations. It produces discrepancy reports, captures bounded evidence for offline replay, and includes an offline dashboard.

## Reproduce the implementation

Go 1.22.2+; `go test ./...`, `go vet ./...`, `go build ./...`, and `node --test dashboard/viewer.test.cjs`. Run `go run ./cmd/ledger-parity --demo --format json --out -`; application exit 3 is expected for the synthetic 1-match/4-discrepancy/3-unknown demo.

## Evidence and supported scope

Released as v0.3.2-preview with cross-platform archives and SHA-256 checksums. Local checks passed October 7. No consenting operator export or adoption evidence has been supplied. Ordinary classic payments only; no signing or automatic funds movement.

Evidence reference: [https://github.com/LedgerParity/ledger-parity-cli/releases/tag/v0.3.2-preview](https://github.com/LedgerParity/ledger-parity-cli/releases/tag/v0.3.2-preview).
Baseline source revision: `82d42d6b74117e03a7bab07bfa549b7edb8cd0f7`. Final reviewed preparation revision and CI
results belong in [VERIFICATION_OCT09.md](VERIFICATION_OCT09.md).

## Maintainers and contributor work

Maintainer: EthTobi; contact via GitHub, available anytime.
See [MAINTAINERS.md](../MAINTAINERS.md), [CONTRIBUTING.md](../CONTRIBUTING.md),
[SECURITY.md](../SECURITY.md), and [CODE_OF_CONDUCT.md](../CODE_OF_CONDUCT.md).
The [focused engineering backlog](WAVE_BACKLOG.md) describes real work, relevant
files, tests, and acceptance criteria. Draft complexity values require maintainer
review and app enrollment; they do not establish approval or earned points.

## Before applying

- Confirm the Drips Wave App covers this repository and check application slots.
- Publish the reviewed backlog issues and preserve links to their acceptance checks.
- Publish these preparation changes through a reviewed PR with passing CI.
- Apply under the implemented scope above; no production/adoption claims are implied.

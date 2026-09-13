# 0.3.1-preview

This developer preview packages the ordinary Stellar-payment reconciliation workflow, scoped SDP CSV import, evidence capture/replay, local SHA-256 checks and offline dashboard.

Verification now hashes the report after it is written, accepts the documented check command, rejects malformed proofs and prevents input/proof collisions. The dashboard escapes report text, validates totals and coverage, displays candidate operation IDs, searches both sides' identities and keeps exact amounts as text. Files stay in the browser.

## Run the downloaded package

Unzip the package for your platform. On macOS/Linux, if the unzip tool drops executable permissions, run `chmod +x ledger-parity`. On Windows use `ledger-parity.exe` in place of `./ledger-parity`.

```sh
./ledger-parity --version
./ledger-parity --demo --format json --out report.json --verify auto
./ledger-parity --verify-check report.json.proof.json report.json
```

The demo returns exit 3 intentionally: 1 match, 4 discrepancies and 3 unknown results. Open `dashboard/index.html` and choose the report. SDP input, replay and operator setup are documented in the included README and docs. `SHA256SUMS` covers the download archives; checksums compare bytes and are not signatures.

## Boundaries

The CLI is read-only. It does not sign, deploy contracts or register report hashes on-chain. The optional Soroban contract is a separate repository and validation track. A local checksum does not authenticate business data or establish ledger completeness. Real operator validation remains outstanding; no production-readiness or adoption claim. Path payments, Soroban event reconciliation and durable checkpoints are deferred, not included in this preview.

Packages: Windows amd64, Linux amd64/arm64, and macOS amd64/arm64. Cross-compilation does not replace execution on each target. Go tests, Linux CI race tests, offline demo/replay and dashboard tests are the automated release checks; target-specific installation remains a preview validation task.

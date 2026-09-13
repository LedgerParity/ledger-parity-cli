# Developer-preview delivery

## Download and run

- [CLI v0.3.2-preview](https://github.com/LedgerParity/ledger-parity-cli/releases/tag/v0.3.2-preview): Windows amd64, Linux amd64/arm64, macOS amd64/arm64 archives and SHA256SUMS. Each archive includes examples, documentation and the offline dashboard.
- [Contract v0.2.0-preview](https://github.com/LedgerParity/-ledger-parity-contract/releases/tag/v0.2.0-preview): Wasm and SHA256SUMS, independently deployable from the read-only CLI.
- [Documentation](https://ledgerparity.github.io/).

Run `ledger-parity --demo --format json --out report.json --verify auto` (use `./ledger-parity` on Unix or `.\ledger-parity.exe` in Windows PowerShell). Exit 3 is the expected synthetic result: one match, four discrepancies and three unknowns. Check with `--verify-check report.json.proof.json report.json`. Open `dashboard/index.html` and choose the report. See RELEASE_NOTES.md for install details and supported scope.

## Delivered software

Ordinary classic Stellar-payment reconciliation with exact amounts and identity checks; conservative coverage/UNKNOWN outcomes; canonical JSON/CSV and release-scoped SDP export input; deterministic demo; bounded evidence capture and byte-identical offline replay; local checksum verification; an offline report dashboard; and release packaging.

The separate report-hash contract requires owner authorization, stores immutable persistent records, and supports network-capped lifetime renewal. Its read-only monitor checks report, instance and code lifetimes from an external registry. Synthetic Protocol 28 testnet deployment, signed registration, read-back, early renewal and provider fees are documented in [TESTNET.md](https://github.com/LedgerParity/-ledger-parity-contract/blob/main/TESTNET.md). Hash registration is an authorized assertion, not proof that reconciliation occurred or that its source data is correct.

## Verification

CLI [CI](https://github.com/LedgerParity/ledger-parity-cli/actions/runs/34786747678) and [release workflow](https://github.com/LedgerParity/ledger-parity-cli/actions/runs/34786775577) passed. The published Windows archive was independently downloaded, checksum-checked and executed: version, deterministic demo, proof creation and verification passed. Other targets are cross-built; no execution on those operating systems is claimed.

Contract [CI](https://github.com/LedgerParity/-ledger-parity-contract/actions/runs/34787497018) and [release workflow](https://github.com/LedgerParity/-ledger-parity-contract/actions/runs/34787502177) passed 13 Rust tests, five monitor/evidence tests and the release Wasm build. The registered report bytes are preserved by Git attributes and checked against public evidence. The deployed artifact was obtained from verified CI. Local Windows Rust execution remains unavailable because the MSVC linker is missing; Rust execution evidence comes from Linux CI.

Final local core and connector test/vet/build checks passed; all seven offline RPC corpus tests passed. Documentation deployment passed in [Pages CI](https://github.com/LedgerParity/ledgerparity.github.io/actions/runs/34787532055). These are engineering checks, not operator adoption or production validation.

## External gates and deferred scope

- A consenting operator must supply a sanitized export and known discrepancy, reproduce the finding, and report false positives/integration effort. See docs/OPERATOR_VALIDATION.md. No such operator evidence has been supplied.
- A deployment operator must own monitoring and renewal. The monitor and registry are supplied; no background service or maintenance commitment is installed or implied.
- Actual archived restoration/read-back remains pending. The synthetic records are live through ledger 5181347; the contract cannot shorten TTL. The early signed renewal was a no-op for TTL and is not a restoration test. Native tests cover threshold extensions and archival-access failures.
- Runtime RPC ingestion, path payments, Soroban event reconciliation, durable checkpoints and CLI signing/contract submission remain separately designed future scope. They are not hidden features of this preview.

The developer preview is released. Production readiness cannot be established until the external gates are met. Preserve the distinction when sharing the project.

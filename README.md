# ledger-parity-cli

A read-only tool for Stellar payment operators to compare an application export with ordinary classic Stellar payment operations. It helps investigate settlement notifications that were missed, incorrect amounts and duplicate application records. Developer preview: no demonstrated operator adoption or production-readiness claim.

From this repository alone, install Go 1.22.2+ (prefer a currently supported Go release), then:

```sh
go test ./...
go vet ./...
go build -o ledger-parity ./cmd/ledger-parity
./ledger-parity --demo --format json --out -
```

On Windows use `go build -o ledger-parity.exe ./cmd/ledger-parity`, then `.\ledger-parity.exe --demo --format json --out -`. If Go is installed but unavailable in your terminal, add its bin directory to the current shell PATH or use its full executable path.

The demo is deterministic and offline: **1 match, 4 discrepancies, 3 unknowns**, exit **3**. It includes an exact settlement, a one-unit amount difference, a missing expected operation, duplicate internal IDs, a pending record and an ambiguous batch. Duplicate input IDs prevent a definitive orphan conclusion. These synthetic addresses and transaction identifiers are deliberately not live ledger evidence.

Use the same file-input pipeline directly:

```sh
./ledger-parity --config examples/config.json --format both --out report.json
```

Configuration is strict JSON, not YAML. Source paths and config output paths resolve relative to the config file; `--out` resolves relative to the current directory. CLI flags override config values. `--out - --format json` emits only JSON, suitable for scripts. `--demo` is explicit; running without config does not silently switch to demo.

Exit codes: **0** no findings, **1** config/input/network/output error, **2** discrepancies, **3** unresolved evidence (takes precedence when both exist). Read report totals for mixed results. Ingestion errors produce no new report and a nonzero exit; an older report file may still exist, so always check the exit code. `go run` can wrap nonzero exit statuses; use the built binary in automation.

For live read-only ingestion, copy examples/config.json, replace `stellar.on_chain_path` with `stellar.horizon_url`, set the exact network passphrase and monitored accounts, and provide a canonical application export. Always use an explicit RFC3339 start/end window. The CLI expands the Horizon scan by the time tolerance. `target_app.complete` asserts that the export contains every relevant ordinary payment for those accounts/window; leave false unless you can establish that. No secrets or wallets are needed. An optional public-testnet smoke check is available with `powershell -File scripts/Test-LiveRead.ps1` after a Windows build; it synthesizes an expected row from public testnet data and documents that evidence limit.

The canonical export contract lives in [connectors](https://github.com/LedgerParity/ledger-parity-connectors), matching/coverage semantics in [core](https://github.com/LedgerParity/ledger-parity-core). Required fields include network passphrase, operation_type `payment`, exact sender/recipient, amount as a decimal string, asset/type/issuer, timestamp and status. Operation ID is recommended; reference ID means transaction hash only. Exact amounts, case-sensitive asset codes and issuers, network and direction are never bypassed by a reference. Non-completed statuses without observed settlement remain UNKNOWN. Examples are bundled into the binary for standalone use.

Supported: canonical JSON/CSV input, ordinary Horizon payment operations, bounded pagination/retries and explicit coverage. Unsupported: Soroban/RPC events, path payments, account creation/merge settlement, verified named-product integrations, memo matching, durable automatic checkpoints, signing or fund movement. Named Stellopay/Facil-Pay/Trustless Work packages are experimental examples, not CLI adapters. Provider completeness is trusted evidence, not a cryptographic proof. StrKey checksums and general muxed/base-account equivalence are not validated. Review reports before operational decisions.

[Reviewer walkthrough](docs/REVIEWER_WALKTHROUGH.md) · [verification](PROJECT_HANDOFF.md) · [roadmap](ROADMAP.md) · [appeal draft](docs/DRIPS_SUBMISSION_PREP.md) · [contributing](CONTRIBUTING.md) · [security](SECURITY.md) · [MIT](LICENSE).

# LedgerParity Reconciliation Dashboard

An offline, single-file HTML viewer for CLI reconciliation reports. Open `index.html` in a modern browser and choose or drop a report produced by `ledger-parity --format json`. Files remain in the browser; no upload or build step is required. The file picker supports keyboard use and stays available for loading another report.

The viewer shows internal record and on-chain operation counts separately, matches, discrepancies and unknown results. Coverage assertions, source, reason, monitored accounts, network and the reconciliation window are visible alongside the results. Viewing a report does not authenticate it or re-run reconciliation.

Search covers payment IDs, operation IDs (including ambiguity candidates), business references, transaction hashes, both sides' accounts, assets, issuers, amounts, statuses and notes. Amounts and IDs remain strings. On-chain-only records show their sender and destination. The table displays up to 500 matching rows; search checks the full report and the visible count explains the limit.

Reports must be at most 32 MiB and have valid report fields, result statuses, coverage metadata and consistent result/count summaries. Proof files, evidence bundles and invalid reports produce an inline error. A failed load hides the previous report. Report text is escaped before HTML rendering.

## Verification

With Go and a modern Node.js on PATH, from the repository root:

```sh
node --test dashboard/viewer.test.cjs
```

Tests generate the actual CLI demo and check rendering, HTML escaping, precision, coverage, candidate IDs, on-chain search, malformed reports, repeated loads, read failures and result limits. Set `LEDGER_PARITY_BIN` to an absolute built CLI path to skip `go run`; CI uses its built binary. These Node tests use a small DOM stub, so browser layout needs a separate visual check.

# LedgerParity Reconciliation Dashboard

A single-file, no-build-step HTML viewer for LedgerParity reconciliation reports.

## Usage

1. Run a reconciliation:
   ```sh
   ledger-parity --config config.json --format json --out report.json
   ```

2. Open `dashboard/index.html` in any browser

3. Drag & drop `report.json` onto the page (or click to browse)

## Features

- **Summary cards** — total records, matches, discrepancies, unknowns
- **Bar chart** — visual breakdown by discrepancy type
- **Searchable table** — filter by ID, sender, recipient, asset, amount
- **Brand palette** — ink, teal, violet, amber
- **Offline** — no external dependencies, works without internet

## Report format

Accepts any `DiscrepancyReport` JSON produced by `ledger-parity --format json`. The viewer reads:

- `total_internal`, `total_matched`, `total_discrepancies`, `total_unknown`
- `discrepancy_counts` (map of type → count)
- `results[]` with `status`, `discrepancy`, `internal_payment`, `on_chain_payment`, `notes`
- `coverage` for network and time window metadata

## Files

```
dashboard/
  index.html    The dashboard (single file, no build)
  README.md     This file
```

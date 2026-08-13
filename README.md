# ledger-parity-cli

[![CI](https://github.com/LedgerParity/ledger-parity-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/LedgerParity/ledger-parity-cli/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

`ledger-parity-cli` is the command-line application for **LedgerParity** — bundling `ledger-parity-core` and `ledger-parity-connectors` into a single standalone binary.

---

## ⚡ Quickstart

### Build Binary
```bash
go build -o ledger-parity ./cmd/ledger-parity/main.go
```

### Run Demonstration (Seeded Discrepancy Vector)
```bash
./ledger-parity --demo --format both
```

### Run Custom Config
```bash
./ledger-parity --config sample_config.json --format table --out report.json
```

---

## ⚙️ Configuration Example (`config.json`)

```json
{
  "target_app": {
    "name": "stellopay",
    "format": "json",
    "source_path": "./data/stellopay_export.json"
  },
  "stellar": {
    "horizon_url": "https://horizon-testnet.stellar.org",
    "accounts": [
      "GAXYZ1234567890ACCOUNT"
    ]
  },
  "reconciliation": {
    "timeframe_tolerance_sec": 600,
    "ignore_failed_on_chain": true
  },
  "output": {
    "format": "table",
    "file_path": "discrepancy_report.json"
  }
}
```

---

## 🖥️ Terminal Output Preview

```
🔍 LedgerParity Payment Reconciliation Engine v1.0.0
════════════════════════════════════════════════════════════════════════════════════
  LEDGERPARITY RECONCILIATION REPORT — STELLOPAY_DEMO
════════════════════════════════════════════════════════════════════════════════════
  Generated At:   2026-08-13 14:58:13 UTC
  Time Window:    2026-08-12 14:58 to 2026-08-13 14:58
────────────────────────────────────────────────────────────────────────────────────
  TOTAL INTERNAL: 4      | TOTAL ON-CHAIN: 3      | MATCHED: 0      | DISCREPANCIES: 5     
────────────────────────────────────────────────────────────────────────────────────
  DISCREPANCY BREAKDOWN:
    • DUPLICATE_INTERNAL     : 2
    • AMOUNT_MISMATCH        : 1
    • MISSING_ON_CHAIN       : 1
    • ORPHANED_ON_CHAIN      : 1
════════════════════════════════════════════════════════════════════════════════════
```

---

## 📄 License

[MIT License](LICENSE) © LedgerParity Maintainers.
# Contributor tasks

1. Operator validation brief: document one consenting operator's export fields and known discrepancy using sanitized data. Acceptance: record what the operator could independently reproduce, integration effort and false positives. No customer data, manufactured endorsement or impact metrics.
2. Report provenance: add hashes of local input bytes and tool/source revisions to a versioned report envelope. Acceptance: mutating an input changes its hash; no filenames/secrets need be exposed; existing report consumers have a migration fixture. Hashes prove byte identity, not input truth.
3. Portable live-read demonstration: add a Go version of scripts/Test-LiveRead.ps1 using only GETs and explicit testnet scope. Acceptance: normal CI remains offline, lack of recent ordinary payments is an honest non-success, and generated expected data is labeled synthetic. No account funding or signing.
4. Large exports and duplicate JSON keys: coordinate a bounded shared decoder with connectors. Acceptance: duplicate field names and oversized inputs fail with actionable errors while exact decimal strings survive. Avoid silently changing configured scope.

Maintainer must confirm availability before assignment. A PR needs reproduction, source evidence where relevant, tests and scope limits. No Wave issue labels, points or response-time promises are implied.

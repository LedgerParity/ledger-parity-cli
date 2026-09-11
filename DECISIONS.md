# Decisions

2026-09-11: The supported CLI selects canonical files explicitly and does not silently route a product name to an unverified adapter. Only strict JSON config is supported. Paths resolve relative to config except CLI overrides; '-' is JSON stdout. Start/end, network and monitored account scope are explicit. Internal completeness is an operator assertion, not inferred from a successful file read. Horizon failures produce exit 1 and no new report. Exit 3 for UNKNOWN takes precedence over discrepancy exit 2. Empty/unproven data is not a clean reconciliation.

The bundled demo is deterministic and entirely synthetic. The optional live-read script uses only GET requests and derives a synthetic expectation from an observation; it is interoperability evidence, not an independent business-record check. Preserve ordinary classic-payment and MIT scope; no Soroban, contract, wallet or automatic checkpoint functionality is claimed. The recalled Drips concern is Stellar relevance/impact; actual adoption remains unverified.

Ambiguity reports now expose candidate operation IDs in JSON and terminal output. Go module pins reference published source revisions. CI covers the compatibility floor and stable Go; current remote execution is unverified. Archived exports are used for isolated build verification, without new project repositories.

Verification amendment: remote CI at 1861a5f passed both Go matrix jobs, including the configured Linux race tests and offline demo. Core and connectors also have successful verified runs (docs/VERIFICATION.md). Prior unverified-CI statements are superseded; local Windows race execution remains distinct.

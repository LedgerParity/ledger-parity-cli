# Roadmap

Current bounded milestone: strict config/file workflow, exact operation-aware comparison, explicit coverage/UNKNOWN, reproducible offline demonstration and CLI failure tests. Acceptance: demo 1 match/4 discrepancies/3 unknowns, explicit exit 3; mock Horizon success and error cases; standalone test/vet/build; normal commit/push.

Next: operator-validated export/discrepancy and report provenance. See docs/backlog.md. Broader path-payment/Soroban and durable checkpoints belong to explicit core design gates, not marketing claims. Final local/live/remote verification status is in PROJECT_HANDOFF.md. Drips appeal draft is not a submitted application.

- [x] Canonical workflow, CLI negative tests, deterministic demo and read-only testnet evidence.
- [x] Local test/vet/build and documented demo/config runs.
- [x] Remote CI verified at the revisions linked in PROJECT_HANDOFF.md.
- [ ] Optional local Windows race execution (no C compiler); Linux CI race tests passed.
- [ ] Operator-validated practical impact; broader semantics remain deferred.

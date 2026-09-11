# Verification evidence

Executed 2026-09-11 on Windows with Go 1.24.4. All three modules passed `go test ./...`, `go vet ./...` and `go build ./...`. Exported commit archives were also checked with GOWORK=off and separate fresh GOMODCACHE directories; no sibling source replacements were available. Core: 37 passing test cases including named subtests; connectors: 8; CLI: 5. These are test-runner cases, not assertion counts. Isolated CLI build and embedded demo passed (exit 3; 1 match, 4 discrepancies, 3 unknowns).

Revisions tested in isolated archives: core 89aa1ec, connectors 902fa78, CLI 1861a5f. Subsequent completion records/CI updates do not alter those runtime algorithms. Removing unused older go.sum entries does not change selected versions.

The [read-only live testnet check](LIVE_TESTNET_RESULT.json) passed at 2026-09-11T22:14:14Z. It verifies provider/API interoperability by deriving one synthetic expectation from an observed ordinary payment and independently refetching the account through CLI. It is not an independently sourced business record, real customer integration, transaction submission or production proof.

Remote CI verified through public GitHub run summaries:

| Repository | Revision | Run | Result |
| --- | --- | --- | --- |
| core | dda2427 | https://github.com/LedgerParity/ledger-parity-core/actions/runs/34654056907 | Success, two Go matrix jobs |
| connectors | 83fd7ec | https://github.com/LedgerParity/ledger-parity-connectors/actions/runs/34653797659 | Success, two Go matrix jobs |
| CLI | 1861a5f | https://github.com/LedgerParity/ledger-parity-cli/actions/runs/34653811330 | Success, two Go matrix jobs, including offline demo |

Core's initial public run list was empty. Actions permission inspection confirmed enabled/all; an authenticated manual dispatch was accepted without changing repository policy. It completed successfully. Initial unauthenticated API requests were rate-limited; public run pages and normal configured authentication resolved most visibility issues.

Local race execution could not start because CGO was disabled and no C compiler was found. A portable compiler download was cancelled after very slow progress; the partial archive was removed. No compiler or global configuration was installed. Linux CI is configured to execute race tests; its results are separate from locally executed tests.

Remaining limits: trusted-provider history/continuity, opaque address strings without StrKey validation, partial muxed account scope, no Soroban/path-payment/account-creation/merge semantics, no durable atomic checkpoint, unverified named-product export contracts, no real operator adoption/impact evidence. Contributor tasks and the roadmap describe bounded next steps. No Drips application or appeal was submitted.

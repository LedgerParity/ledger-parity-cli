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


## 0.3.0-preview follow-up (2026-09-12)

Local Go tests, vet and builds passed across all three modules. Connectors and CLI were also checked with GOWORK=off against published dependencies; CLI go mod verify passed. Core's separate Node corpus passed seven offline tests using locked SDK 17.0.1. The captured public-testnet response pair is dated 2026-09-11T23:02:34.348Z, transaction 6c84fa16503524baece6dbfea456834d5612db98593ecb896fe71488f2906cb7, operation 19880664878817281. This is provider-to-provider native-payment agreement, not independent application evidence or complete history. Batch/failure/event-gap cases are synthetic mutations.

SDP mapping is tested against synthetic release-shaped CSV, not a deployed SDP export. Runtime RPC ingestion and operator adoption remain unverified. Remote CI for these revisions is recorded separately below when observed; earlier run links do not validate new changes.

The standalone SDP binary demonstration and offline replay both exited 3 with (1 match, 1 discrepancy, 1 unknown); output reports were byte-identical. Core CI passed at 954f091 ([run](https://github.com/LedgerParity/ledger-parity-core/actions/runs/34657415249)), including the separate RPC corpus and both Linux Go race jobs. Connectors CI passed at b99a351 ([run](https://github.com/LedgerParity/ledger-parity-connectors/actions/runs/34657369946)). CLI CI passed at 1c6a32b ([run](https://github.com/LedgerParity/ledger-parity-cli/actions/runs/34657582486)), including both Linux race jobs and the SDP capture/replay comparison. This verification-note update changes documentation only; CI claims refer to those exact implementation revisions.

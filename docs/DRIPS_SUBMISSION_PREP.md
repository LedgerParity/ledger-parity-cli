# Drips appeal / resubmission preparation

Draft reviewed 2026-09-11. Not submitted. The original application and rejection text are unavailable. The maintainer recalls a concern that the project should be more clearly Stellar-based and create practical ecosystem impact; this is a recollection, not a verbatim review or confirmed full rejection rationale.

## Proposed appeal text

LedgerParity is a read-only Go reconciliation tool for developers operating Stellar classic-payment applications. It compares canonical application exports with Stellar Horizon payment operations, so operators can investigate missed settlement notifications, exact amount differences and duplicated application records.

Since the earlier version, we replaced floating-point comparisons with exact stroop arithmetic, required network/asset/issuer/direction identity, removed reference matching that could bypass economic checks, added operation-level ambiguity and duplicate handling, and made unproven coverage explicit. Horizon reads now paginate with bounded retries and return errors instead of allowing network failure to become missing-payment findings. The CLI has a reproducible offline workflow and failure-path tests, strict JSON/CSV input and usable exit codes. READMEs distinguish implemented ordinary-payment support from unsupported Soroban and unverified named adapters. Contributor tasks derive from remaining data-contract and verification gaps.

This addresses the remembered relevance concern by making the Stellar payment-operator use case runnable and reviewable. We do not yet claim real operator adoption or measured ecosystem impact. The next validation step is an operator-provided sanitized export specification and a known discrepancy that they can reproduce independently. There is no signing, fund movement, token or custom contract.

Evidence: the reviewer walkthrough, committed regression tests, pinned standalone modules, local verification handoffs and any explicitly recorded read-only testnet check. Offline fixtures and testnet-derived synthetic expectations are not production integration evidence. Remote CI must be checked independently; it is not inferred from a workflow file or local tests.

## Official rules versus project recommendations

The [official maintainer guide](https://docs.drips.network/wave/maintainers/participating-in-a-wave/) requires relevant public repositories to be applied and approved by organizers. Rejected repositories use the in-app Appeal action; substantive improvements are required. It describes a two-week first-appeal wait, one-month cooldown after a declined appeal, and a maximum of three appeals. Check the repository's actual dashboard status and dates before acting. An appeal, not a fresh application of the same rejected repository, is the documented reconsideration route.

The [public Stellar Wave page](https://www.drips.network/wave/stellar) displayed no active/upcoming Wave on review. This does not establish whether appeals/applications are unavailable. [Application allowances](https://docs.drips.network/wave/maintainers/repo-application-limits/) depend on program/account configuration. [Terms](https://docs.drips.network/wave/terms-and-rules/) prohibit point manipulation and low-quality untested work. No eligibility, next deadline, budget or acceptance outcome is assumed.

Engineering recommendations: lead with the runnable CLI workflow and core correctness changes; explain the three existing repository boundaries without treating them as a funding strategy; select real bounded tasks only after maintainer review capacity is confirmed. The full rejection text, original submission, rejection date, appeal history, signed-in allowances and contributor-review availability remain unavailable. Do not invent them in the final appeal.

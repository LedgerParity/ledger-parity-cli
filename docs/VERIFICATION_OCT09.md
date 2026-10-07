# October 9 submission verification

Executed October 7, 2026 in the shared Linux workspace. Baseline source revision:
`82d42d6b74117e03a7bab07bfa549b7edb8cd0f7`. Preparation changes affect documentation/governance and, where noted,
CI checks; no runtime implementation or public API was changed.

## Local results

Runtime: go version go1.22.2 linux/amd64 and Node 22.23.2 for dashboard checks.

- `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./...`: passed.
- `gofmt -l` on tracked Go files: no unformatted files; CI now enforces this.
- `node --test dashboard/viewer.test.cjs`: 5 passed, including a real CLI demo.
- `go run ./cmd/ledger-parity --demo --format json --out -`: observed 1 match,
  4 discrepancies, 3 unknowns. Application exit 3 is expected; go run reports it
  as shell status 1. This is synthetic offline evidence, not operator adoption.

## Review and publication

- `git diff --check`: passed after preparation edits.
- CI YAML parsed locally. A syntax parse does not replace GitHub execution.
- Six engineering issues were published with bounded acceptance criteria and
  proposed complexity; links are in WAVE_BACKLOG.md. No Wave labels/enrollment
  or contributor assignments were performed.
- Maintainers xteesamz and EthTobi were owner-confirmed; GitHub contact and
  anytime availability apply. GitHub App coverage/application slots still
  require dashboard confirmation.
- Changes will be proposed through a fork PR because the available account
  cannot push directly to the organization's protected branch. Merge decisions
  remain with maintainers. Recheck the final PR checks before applying.

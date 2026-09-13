# ledger-parity-verify

Soroban smart contract for tamper-proof report verification. Stores report hashes
on-chain so operators can prove a reconciliation was performed at a specific time.

## Functions

### `store(hash, owner, metadata)`
Store a report hash with the owner's address and optional metadata.
Panics if the hash is already stored.

### `verify(hash)`
Check if a report hash exists on-chain. Returns `true` if stored.

### `get_info(hash)`
Retrieve metadata about a stored report: owner, timestamp, metadata.
Panics if not found.

### `is_owner(hash, owner)`
Check if a hash was stored by a specific owner.

## Deployment

Requires the [Stellar CLI](https://developers.stellar.org/docs/smart-contracts/getting-started/setup):

```sh
# Build
soroban contract build

# Deploy to testnet
soroban contract deploy \
  --wasm target/wasm32-unknown-unknown/release/ledger_parity_verify.wasm \
  --network testnet

# Initialize (no init function needed — contract is stateless until first store)
```

## Usage from CLI

After generating a report, the CLI can optionally store its hash:

```sh
ledger-parity --config config.json \
  --format json --out report.json \
  --verify --verify-network testnet
```

This computes the SHA-256 of the report and calls `store()` on the deployed contract.
Later, anyone can verify the report's integrity:

```sh
ledger-parity --verify-check report.json --verify-network testnet
```

## Why this exists

- **Tamper-proof audit trail** — prove a report existed at a specific time
- **Operator confidence** — on-chain verification without trusting the tool
- **Drips maintenance** — ongoing work (batch storage, access control, time queries)

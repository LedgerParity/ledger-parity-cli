#![no_std]
use soroban_sdk::{contract, contractimpl, contracttype, symbol_short, Address, Bytes, Env, Symbol};

const STORED: Symbol = symbol_short!("STORED");
const VERIFIED: Symbol = symbol_short!("VERIFIED");

#[contracttype]
pub struct ReportInfo {
    pub owner: Address,
    pub timestamp: u64,
    pub metadata: Bytes,
}

#[contract]
pub struct VerifyContract;

#[contractimpl]
impl VerifyContract {
    /// Store a report hash with owner and optional metadata.
    pub fn store(env: Env, hash: Bytes, owner: Address, metadata: Bytes) {
        if env.storage().instance().has(&hash) {
            panic!("report already stored");
        }
        let info = ReportInfo {
            owner: owner.clone(),
            timestamp: env.ledger().timestamp(),
            metadata,
        };
        env.storage().instance().set(&hash, &info);
        env.events().publish((STORED, owner), hash);
    }

    /// Verify a report hash exists. Returns true if stored.
    pub fn verify(env: Env, hash: Bytes) -> bool {
        env.storage().instance().has(&hash)
    }

    /// Get info about a stored report. Panics if not found.
    pub fn get_info(env: Env, hash: Bytes) -> ReportInfo {
        env.storage()
            .instance()
            .get(&hash)
            .expect("report not found")
    }

    /// Check if a hash was stored by a specific owner.
    pub fn is_owner(env: Env, hash: Bytes, owner: Address) -> bool {
        match env.storage().instance().get::<_, ReportInfo>(&hash) {
            Some(info) => info.owner == owner,
            None => false,
        }
    }
}

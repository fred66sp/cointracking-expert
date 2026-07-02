# Holdings Engine Specification

**Current and Historical Holdings Reconstruction**

The Holdings Engine reconstructs current and historical holdings from transaction history and compares them with holdings reported by CoinTracking, detecting discrepancies.

## Purpose

Rebuild current holdings from transaction history and compare with expected holdings, detecting unreported transactions or data inconsistencies.

## Inputs

- Complete transaction ledger
- Current holdings as reported by CoinTracking
- Historical holdings snapshots (optional)

## Outputs

- Reconstructed holdings
- Discrepancy detection and reporting
- Holdings timeline

## Responsibilities

1. Compute current holdings from complete transaction history
2. Compare with reported holdings
3. Detect and explain discrepancies
4. Build historical holdings timeline
5. Report any assets with zero holdings or transfer-only transactions

## Key Algorithms

- Holdings computation from ledger
- Discrepancy detection and attribution
- Historical timeline building

## Edge Cases

- Dust quantities (very small holdings)
- Staking rewards and locked tokens
- Wrapped tokens and bridges
- Assets with multiple instances (e.g., different DEX versions)
- Recent transactions affecting current holdings

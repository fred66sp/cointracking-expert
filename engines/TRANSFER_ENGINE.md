# Transfer Engine Specification

**Deposit and Withdrawal Matching**

The Transfer Engine matches withdrawal transactions from one account with deposit transactions to another account, identifying orphaned transfers and mismatched pairs.

## Purpose

Match transfers between accounts, detecting cases where one side of a transfer is missing or orphaned, and validating transfer consistency.

## Inputs

- All withdrawal and deposit transactions
- Multi-account transaction history
- Exchange and wallet address mappings

## Outputs

- Matched transfer pairs
- Orphaned transfers (unmatched)
- Transfer timeline and flow visualization

## Responsibilities

1. Identify potential transfer matches (same asset, similar amount, nearby dates)
2. Match withdrawals to deposits
3. Detect unmatched or orphaned transfers
4. Handle bridge transfers and wrapped tokens
5. Report on transfer timing and fee impact

## Key Algorithms

- Transfer matching algorithm (amount, asset, timestamp)
- Ambiguity resolution for multiple candidates
- Bridge and wrapped token detection

## Edge Cases

- Fees reducing transfer amount
- Partial fills or split transfers
- Delays between withdrawal and deposit
- Address format variations
- Bridge transfers (converting between chains)
- Assets with multiple versions (wrapped, bridged, etc.)

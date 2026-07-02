# FIFO Engine Specification

**Acquisition Lot Reconstruction and Cost Basis Calculation**

The FIFO Engine reconstructs acquisition lots using First-In-First-Out accounting, assigning purchase transactions to sale transactions based on chronological order. It identifies missing purchase history and calculates cost basis for tax purposes.

## Purpose

Assign acquisition lots to holdings using FIFO method, calculate cost basis, detect missing purchase history, and validate tax calculations based on lot cost.

## Inputs

- Complete ledger with all buy and sell transactions
- Historical price data (optional, for validation)
- Tax configuration (cost basis method)

## Outputs

- Acquisition lot assignment for all holdings
- Cost basis calculations
- Missing purchase history detection
- Tax gain/loss calculations by lot

## Responsibilities

1. Reconstruct acquisition lots chronologically
2. Match holdings to specific purchase transactions
3. Detect situations where more assets were sold than purchased
4. Calculate cost basis for each holding
5. Calculate realized and unrealized gains/losses

## Key Algorithms

- FIFO lot matching
- Cost basis tracking
- Gain/loss calculation
- Missing purchase history detection

## Edge Cases

- Fractional quantities
- Asset consolidations and splits
- Multiple purchase prices
- Zero-cost acquisition (airdrops, staking rewards)
- Partial holdings sales

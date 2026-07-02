# Tax Engine Specification

**Tax Liability Calculation and Validation**

The Tax Engine computes tax liabilities based on transaction history and validates CoinTracking tax calculations. Specialized support for Spanish taxation rules.

## Purpose

Calculate realized gains/losses, compute tax liabilities, validate CoinTracking tax calculations, and generate tax reports for multiple jurisdictions with emphasis on Spanish compliance.

## Inputs

- Complete transaction ledger with costs
- FIFO acquisition lot assignments
- Tax configuration (jurisdiction, method, exemptions)
- Price data for gain/loss calculation

## Outputs

- Realized gains and losses by transaction
- Tax liability by year and asset class
- Tax report in jurisdiction-specific format
- Validation results comparing to CoinTracking calculations

## Responsibilities

1. Compute realized gains/losses for all dispositions
2. Identify taxable events
3. Calculate tax liability by year
4. Support multiple accounting methods (FIFO, etc.)
5. Generate jurisdiction-specific tax reports
6. Validate against CoinTracking calculations

## Key Algorithms

- Realized gain/loss calculation
- Tax liability aggregation
- Jurisdiction-specific rule application
- Tax report generation

## Spanish Tax Rules

- Capital gains taxation
- Cryptocurrency classification
- Wash sale rules (if applicable)
- Reporting requirements
- Deduction rules

## Edge Cases

- Gifts and personal transfers
- Staking rewards and airdrops
- Hard forks and token splits
- Wrapped and bridged tokens
- Cross-jurisdiction transactions
- Non-arm's length transactions

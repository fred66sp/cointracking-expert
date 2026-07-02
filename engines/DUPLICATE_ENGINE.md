# Duplicate Engine Specification

**Transaction Duplicate Detection**

The Duplicate Engine identifies duplicate transactions resulting from multiple import sources, API failures, manual re-entry, or data corruption. It supports exact matching and probabilistic detection.

## Purpose

Identify and report duplicate transactions with varying confidence levels, supporting manual review and remediation.

## Inputs

- Complete transaction dataset
- Duplicate detection rules and thresholds
- Transaction source metadata (CSV, API, manual, etc.)

## Outputs

- Exact duplicates list
- Probable duplicates list with confidence scores
- Duplicate impact analysis (which transactions cause issues)

## Responsibilities

1. Detect exact duplicates (identical all fields)
2. Detect probable duplicates (similar but not identical)
3. Detect cross-source duplicates (CSV + API + manual)
4. Quantify impact of duplicates on balances and holdings
5. Suggest remediation (which to remove)

## Key Algorithms

- Exact matching on transaction characteristics
- Fuzzy matching for probabilistic detection
- Cross-source duplicate detection
- Impact analysis

## Edge Cases

- Partial duplicates (some fields differ)
- Fees recorded separately
- Rounding differences
- Timestamp differences
- Split transactions vs. single transactions
- Deposits vs. "buy" transactions (same action, different labeling)

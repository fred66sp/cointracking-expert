# Glossary

**Terminology and Definitions for CoinTracking Expert**

This glossary defines key terms and concepts used throughout the CoinTracking Expert project. Understanding these terms is essential for working with the framework and its documentation.

## Audit

The complete process of validating a CoinTracking database for completeness, consistency, and compliance. An audit examines transactions, reconstructs balances, detects issues, and produces a detailed report.

## Balance

The amount of a specific asset held in an account, wallet, or exchange at a given point in time. Reconstructed from transaction history.

## Duplicate Transaction

A transaction that appears more than once in the dataset, either as an exact duplicate or a probabilistic match.

## FIFO (First-In-First-Out)

Accounting method that assigns acquisition lots to holdings based on chronological order. First assets purchased are first assets sold.

## Holding

The quantity of a specific cryptocurrency asset held at a specific point in time. Reconstructed from transaction history.

## Ledger

Complete record of all transactions for an account, organized chronologically. Used to reconstruct balances and validate consistency.

## Missing Purchase History

Situation where an asset shows a negative balance at some point, indicating transactions were missing from the dataset.

## Normalization

Process of converting transaction data from various sources (CSV, API, manual) into a canonical representation.

## Reconciliation

Process of matching transactions between two sources (e.g., exchange records vs. CoinTracking database) or verifying consistency within a single source.

## Transfer

Movement of assets between two accounts, wallets, or exchanges. Includes deposits and withdrawals.

## Validation

Process of checking data for consistency, completeness, and compliance with defined rules.

## CoinTracking

Third-party cryptocurrency accounting and portfolio tracking platform. This framework validates CoinTracking databases.

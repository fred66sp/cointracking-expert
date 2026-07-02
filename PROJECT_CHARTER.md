# Project Charter

**Project:** CoinTracking Expert Framework

**Status:** Draft

**Version:** 0.1.0

**Owner:** CoinTracking Expert Project

**Last Updated:** 2026-07-02

---

# 1. Vision

Build the most reliable open-source framework for auditing, reconciling and validating CoinTracking databases, with first-class support for cryptocurrency accounting and Spanish tax reporting.

The project aims to become the reference implementation for cryptocurrency reconciliation by combining deterministic algorithms, domain knowledge and AI-assisted diagnostics.

---

# 2. Mission

Enable users to verify that their CoinTracking data is complete, internally consistent and ready to produce reliable tax reports.

The framework must explain every detected issue, identify its root cause and recommend the minimum corrective action supported by evidence.

---

# 3. Objectives

The project shall provide:

- Transaction reconciliation.
- Ledger reconstruction.
- Holdings validation.
- Duplicate detection.
- Transfer matching.
- Missing Purchase History analysis.
- Tax consistency validation.
- Professional audit reports.
- AI-assisted diagnostics.
- Extensible knowledge base.

---

# 4. Scope

The project includes:

## CoinTracking

- CSV imports
- API imports
- Manual transactions
- Reports
- Warnings
- Holdings
- Tax reports

## Exchanges

Initially:

- Binance
- Coinbase
- Kraken
- Bybit
- OKX
- KuCoin
- BingX

Additional exchanges may be supported later.

## Wallets

- Ledger Live
- MetaMask
- Trust Wallet
- Rabby
- Hardware wallets

## Blockchains

Initially:

- Bitcoin
- Ethereum
- BNB Chain
- Solana
- Polygon
- Arbitrum
- Base

---

# 5. Out of Scope

The following features are intentionally excluded from the first versions:

- Portfolio management
- Trading automation
- Exchange execution
- Investment advice
- Tax filing
- Custody services

The framework validates information.

It does not execute financial operations.

---

# 6. Guiding Principles

## Evidence First

Every conclusion must be supported by observable data.

No assumptions.

---

## Reproducibility

The same dataset must always produce the same audit result.

---

## Explainability

Every detected issue must include:

- Cause
- Evidence
- Impact
- Recommended action

---

## Minimal Intervention

Never recommend deleting or modifying transactions without sufficient evidence.

---

## Documentation Driven Development

No feature shall be implemented before its functional specification has been approved.

---

## Modular Architecture

Every component should be replaceable without affecting the rest of the system.

---

# 7. Quality Goals

The framework should always aim to produce:

- Zero false positives whenever reasonably achievable.
- Deterministic results.
- Complete traceability.
- Reproducible calculations.
- Human-readable reports.
- Machine-readable reports.

---

# 8. Success Criteria

An audit is considered complete only when:

- No unexplained negative balances exist.
- No unresolved Missing Purchase History remains.
- All duplicated transactions are identified.
- Transfers are matched or explicitly justified.
- Holdings reconstructed from history match expected balances.
- Generated reports are internally consistent.

---

# 9. Architecture Principles

The framework is composed of independent engines.

Example:

```
Import Layer
        │
Normalization Layer
        │
Audit Engine
        │
├── Duplicate Engine
├── Transfer Engine
├── Ledger Engine
├── Holdings Engine
├── FIFO Engine
└── Report Engine
```

Each engine must expose well-defined interfaces.

Business rules should remain independent from AI models.

---

# 10. Knowledge Strategy

Knowledge is a first-class component of the project.

Documentation shall be organised into:

- CoinTracking
- Exchanges
- Wallets
- Blockchains
- Taxation
- Reconciliation
- Accounting
- Audit methodology
- Real-world cases

Knowledge must be versioned alongside the source code.

---

# 11. AI Strategy

Artificial Intelligence is an interface, not the source of truth.

The framework must separate:

- Business rules
- Deterministic calculations
- AI explanations

AI models must never replace deterministic calculations.

Their role is to:

- explain;
- guide;
- diagnose;
- summarize;
- assist.

---

# 12. Project Philosophy

The framework does not attempt to "guess" accounting results.

Instead, it reconstructs them from transaction history.

Whenever uncertainty exists, it must be explicitly reported.

Silence is never preferred over uncertainty.

---

# 13. Long-Term Vision

The project should evolve into a complete ecosystem composed of:

- Python library
- Command Line Interface (CLI)
- REST API
- MCP Server
- AI agents
- Knowledge base
- Audit engine
- Tax engine
- Reporting engine

All components should share the same business rules and knowledge base.

---

# 14. Governance

Major architectural decisions must be documented as Architecture Decision Records (ADR).

Every significant change shall be:

1. Proposed.
2. Reviewed.
3. Documented.
4. Implemented.
5. Tested.

---

# 15. Core Rule

**A reconciliation is never considered complete because it "looks correct".**

A reconciliation is complete only when every balance, holding and tax calculation can be reproduced from the complete transaction history using deterministic rules.

This principle takes precedence over every other design decision in the project.
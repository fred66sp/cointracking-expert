# System Architecture

**CoinTracking Expert Framework Design**

This document describes the overall architecture of the CoinTracking Expert framework, including component responsibilities, data flow, interfaces, and design principles. The architecture emphasizes modularity, reproducibility, and clean separation of concerns.

## High-Level Overview

The framework is organized as a pipeline of independent engines that process transaction data through progressive stages of validation and analysis. Each engine has well-defined inputs, outputs, and responsibility boundaries.

```
CoinTracking Export → Normalization → Audit Engine → Validation Engines → Report Generation
                                              ↓
                                    ├─ Duplicate Engine
                                    ├─ Transfer Engine
                                    ├─ Ledger Engine
                                    ├─ Holdings Engine
                                    ├─ FIFO Engine
                                    └─ Tax Engine
```

## Core Principles

1. **Modularity**: Each engine is independent and testable
2. **Reproducibility**: Same input always produces same output
3. **Evidence-First**: All conclusions backed by transaction data
4. **Transparency**: Every issue includes cause, impact, and evidence
5. **Minimal Intervention**: Never modify without justification

## Component Structure

### Import & Normalization Layer

Responsible for reading data from various sources (CoinTracking CSV, API, manual input) and normalizing to canonical representation. Handles format conversion, data cleaning, and schema validation.

### Audit Engine

Orchestrates the audit process, manages workflow, and coordinates between specialized engines. Detects inconsistencies and produces audit reports.

### Specialized Engines

- **Duplicate Engine**: Identifies exact and probabilistic transaction duplicates
- **Transfer Engine**: Matches deposits and withdrawals between accounts
- **Ledger Engine**: Reconstructs balances chronologically
- **Holdings Engine**: Rebuilds expected holdings from transaction history
- **FIFO Engine**: Computes acquisition lots and missing purchase history
- **Tax Engine**: Validates tax calculations and generates reports

### Report Generation

Produces audit reports in multiple formats (Markdown, HTML, Excel, JSON) with configurable detail levels.

## Data Structures

All major components use standardized data structures defined in schemas/. These ensure consistency across imports, engines, and exports.

## Quality Goals

- Zero false positives whenever reasonably achievable
- Deterministic and reproducible results
- Complete traceability and audit trails
- Human-readable and machine-readable output

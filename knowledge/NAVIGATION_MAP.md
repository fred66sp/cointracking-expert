---
id: KB-B1-XXX
title: "Mapa de Navegación — Encuentra lo que Necesitas"
level: B
domain: cointracking
source: "Internal documentation"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: medium
version: 1.0

tags:
  - todo
  - needs-review

notes: "Metadatos agregados automáticamente. Verificar y actualizar conforme ADR-032."
---

# Mapa de Navegación — Encuentra lo que Necesitas

**Para:** Cuando sabes qué necesitas pero no dónde buscarlo  
**Uso:** Ctrl+F y busca tu pregunta

---

## 🔍 Por Tipo de Necesidad

### 📋 "Necesito hacer una auditoría"

| Necesito | Documento | Tiempo |
|----------|-----------|--------|
| Entender los pasos | [PROCEDURE_AUDIT_ACCOUNT.md](procedures/PROCEDURE_AUDIT_ACCOUNT.md) | 20 min |
| Diagrama de flujo | [FLOW_COMPLETE_AUDIT.md](decision-trees/FLOW_COMPLETE_AUDIT.md) | 5 min |
| Checklists (paso a paso) | [CHECKLIST_AUDIT_COMPLETE.md](checklists/CHECKLIST_AUDIT_COMPLETE.md) | 10 min |
| Ver un caso real | [CT-001](cases/ct-001-transferencia-entre-exchanges-importada-.md) a [CT-020](cases/ct-020-advertencia-tecnica-no-es-error-fiscal.md) | 30 min |

---

### 🔴 "Tengo un problema / hallazgo raro"

**Primero:** ¿Cuál es el síntoma? (ve a [TROUBLESHOOTING_INDEX.md](TROUBLESHOOTING_INDEX.md))

| Problema | Solución Rápida |
|----------|-----------------|
| **Saldo negativo** | [FLOW_NEGATIVE_BALANCE.md](decision-trees/FLOW_NEGATIVE_BALANCE.md) |
| **Duplicados** | [CHECKLIST_DUPLICATES.md](checklists/CHECKLIST_DUPLICATES.md) + [DUPLICATE_DETECTION_HEURISTICS.md](cointracking/behavioral/DUPLICATE_DETECTION_HEURISTICS.md) |
| **Falta cost basis** | [CT-002](cases/ct-002-venta-sin-historial-de-compra-previo.md) + [MISSING_PURCHASE_HISTORY_CAUSES.md](cointracking/behavioral/MISSING_PURCHASE_HISTORY_CAUSES.md) |
| **Balance no cuadra** | [BALANCE_CALCULATION_ALGORITHM.md](cointracking/behavioral/BALANCE_CALCULATION_ALGORITHM.md) |
| **API y CSV overlap** | [API_VS_CSV_OVERLAP.md](cointracking/behavioral/API_VS_CSV_OVERLAP.md) |
| **Transferencia perdida** | [PROCEDURE_RECONCILE_TRANSFERS.md](procedures/PROCEDURE_RECONCILE_TRANSFERS.md) |

---

### 💰 "Necesito preparar mi declaración de la renta"

| Necesito | Documento | Orden |
|----------|-----------|-------|
| 1. Entender el flujo fiscal | [CAPITAL_GAINS.md](../taxation/spain/CAPITAL_GAINS.md) | Primero |
| 2. Saber qué es el Modelo 721 | [INFORMATIVE_OBLIGATIONS.md](../taxation/spain/INFORMATIVE_OBLIGATIONS.md) | Segundo |
| 3. Validar mi cost basis | [PURCHASE_POOL_MECHANICS.md](cointracking/behavioral/PURCHASE_POOL_MECHANICS.md) | Tercero |
| 4. Resolver incoherencias | [PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md](procedures/PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md) | Cuarto |
| 5. Preparar el informe | `/spanish-tax-return` (skill) | Quinto |

---

### 🛠️ "No Entiendo Cómo Funciona X"

| Concepto | Dónde Aprender |
|----------|----------------|
| **CoinTracking en general** | [CSV_FORMAT.md](cointracking/official/CSV_FORMAT.md) → [BINANCE_IMPORT_WORKFLOW.md](cointracking/behavioral/BINANCE_IMPORT_WORKFLOW.md) |
| **FIFO (método de coste)** | [PURCHASE_POOL_MECHANICS.md](cointracking/behavioral/PURCHASE_POOL_MECHANICS.md) |
| **Staking & Rewards** | [STAKING_MECHANICS.md](cointracking/behavioral/STAKING_MECHANICS.md) + [STAKING_MECHANICS_BLOCKCHAIN.md](blockchains/STAKING_MECHANICS_BLOCKCHAIN.md) |
| **Airdrops** | [AIRDROPS_MECHANICS.md](cointracking/behavioral/AIRDROPS_MECHANICS.md) |
| **DeFi Swaps** | [DEFI_SWAPS_MECHANICS.md](cointracking/behavioral/DEFI_SWAPS_MECHANICS.md) |
| **Lending** | [LENDING_MECHANICS.md](cointracking/behavioral/LENDING_MECHANICS.md) |
| **Comisiones & Fees** | [FEE_HANDLING.md](cointracking/behavioral/FEE_HANDLING.md) |
| **Transacciones blockchain** | [ETHEREUM_TRANSACTION_TYPES.md](blockchains/ETHEREUM_TRANSACTION_TYPES.md) + [BITCOIN_TRANSACTION_TYPES.md](blockchains/BITCOIN_TRANSACTION_TYPES.md) |
| **Transferencias entre wallets** | [PROCEDURE_RECONCILE_TRANSFERS.md](procedures/PROCEDURE_RECONCILE_TRANSFERS.md) |

---

### 🏦 "Usó un Exchange Específico"

| Exchange | Mecánicas | Importación |
|----------|-----------|-------------|
| **Binance (Spot)** | [BINANCE_SPOT_MECHANICS.md](cointracking/behavioral/BINANCE_SPOT_MECHANICS.md) | [BINANCE_IMPORT_WORKFLOW.md](cointracking/behavioral/BINANCE_IMPORT_WORKFLOW.md) |
| **Binance (Margin)** | [BINANCE_MARGIN_MECHANICS.md](cointracking/behavioral/BINANCE_MARGIN_MECHANICS.md) | ↑ (mismo) |
| **Binance (Futures)** | [BINANCE_FUTURES_MECHANICS.md](cointracking/behavioral/BINANCE_FUTURES_MECHANICS.md) | ↑ (mismo) |
| **Binance (Earn)** | [BINANCE_EARN_MECHANICS.md](cointracking/behavioral/BINANCE_EARN_MECHANICS.md) | ↑ (mismo) |
| **Binance Convert** | [BINANCE_CONVERT_MECHANICS.md](cointracking/behavioral/BINANCE_CONVERT_MECHANICS.md) | ↑ (mismo) |
| **Kraken** | [KRAKEN_STAKING_MECHANICS.md](cointracking/behavioral/KRAKEN_STAKING_MECHANICS.md) | API directa en CT |
| **Coinbase** | [COINBASE_ADVANCED_TRADE.md](cointracking/behavioral/COINBASE_ADVANCED_TRADE.md) | API directa en CT |
| **Otro exchange** | [GENERIC_EXCHANGE_MECHANICS.md](cointracking/behavioral/GENERIC_EXCHANGE_MECHANICS.md) | 3 opciones: API / CSV / Manual |

---

### 📚 "Necesito Profundidad Técnica"

**Niveles de Detalle:**

- **Nivel A (Oficial):** [cointracking/official/](cointracking/official/) + [../taxation/spain/](../taxation/spain/)
- **Nivel B (Operativo):** [cointracking/behavioral/](cointracking/behavioral/) + [exchanges/behavioral/](exchanges/behavioral/) + [blockchains/](blockchains/)
- **Nivel C (Verificado):** [cases/](cases/) + [patterns/](patterns/) + [procedures/](procedures/)
- **Nivel D (Auxiliar):** [checklists/](checklists/) + [decision-trees/](decision-trees/)
- **Nivel E (Referencia):** [reference/](reference/) (glosario, contexto, historiadores)

---

### 🎓 "Quiero Aprender el Sistema Entero"

**Orden recomendado (2 horas):**

1. [QUICK_START.md](QUICK_START.md) (5 min)
2. [INDEX_MASTER.md](INDEX_MASTER.md) (10 min)
3. [CSV_FORMAT.md](cointracking/official/CSV_FORMAT.md) (10 min)
4. [PURCHASE_POOL_MECHANICS.md](cointracking/behavioral/PURCHASE_POOL_MECHANICS.md) (15 min)
5. [PROCEDURE_AUDIT_ACCOUNT.md](procedures/PROCEDURE_AUDIT_ACCOUNT.md) (20 min)
6. [CAPITAL_GAINS.md](../taxation/spain/CAPITAL_GAINS.md) (15 min)
7. Un caso real ([CT-001](cases/ct-001-transferencia-entre-exchanges-importada-.md) a [CT-020](cases/ct-020-advertencia-tecnica-no-es-error-fiscal.md)) (15 min)

---

### 🔗 "¿Cómo Se Conecta Todo?"

**Flujo de Conocimiento:**

```
┌─ NIVEL A (Oficial) ──────────────┐
│ AEAT/BOE/CoinTracking            │
└────────────────┬──────────────────┘
                 ↓
┌─ NIVEL B (Operativo) ────────────┐
│ Cómo funciona CoinTracking        │
└────────────────┬──────────────────┘
                 ↓
┌─ NIVEL C (Verificado) ───────────┐
│ Casos reales, patrones, procedimientos
└────────────────┬──────────────────┘
                 ↓
┌─ NIVEL D (Auxiliar) ─────────────┐
│ Checklists, árboles de decisión   │
└────────────────┬──────────────────┘
                 ↓
┌─ NIVEL E (Referencia) ───────────┐
│ Glosario, contexto               │
└─────────────────────────────────────┘
```

---

## 🚨 Si Aún No Encuentras

1. Busca en [INDEX_MASTER.md](INDEX_MASTER.md) (estructura completa)
2. Busca en [GLOSSARY.md](reference/GLOSSARY.md) (término específico)
3. Busca en [TROUBLESHOOTING_INDEX.md](TROUBLESHOOTING_INDEX.md) (síntoma)
4. Busca en [adr/INDEX.md](../adr/INDEX.md) (decisión/contexto)

Si aún no está: abre issue en el repositorio.

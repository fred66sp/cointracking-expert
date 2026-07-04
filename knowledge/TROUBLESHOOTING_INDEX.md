# Índice de Troubleshooting — Diagnostica tu Problema

**Para:** Cuando algo sale mal en tu auditoría  
**Uso:** Busca el síntoma en la tabla y sigue el árbol de decisión

---

## 🎯 Tabla Rápida de Síntomas

| Síntoma | ¿Es un Error? | Árbol | Documentos |
|---------|---------------|-------|-----------|
| **Saldo negativo** | ⚠️ RARO | [FLOW_NEGATIVE_BALANCE.md](decision-trees/FLOW_NEGATIVE_BALANCE.md) | [CT-004](cases/ct-004-balance-negativo-por-orden-cronologico-i.md), [CT-012](cases/ct-012-balance-negativo-por-importacion-parcial.md) |
| **Dos operaciones idénticas** | ⚠️ PROBABLEMENTE | [FLOW_DUPLICATE_DETECTION.md](decision-trees/FLOW_DUPLICATE_DETECTION.md) | [DUPLICATE_DETECTION_HEURISTICS.md](cointracking/behavioral/DUPLICATE_DETECTION_HEURISTICS.md), [CT-003](cases/ct-003-api-y-csv-importados-simultaneamente-dup.md) |
| **Venta sin cost basis** | ❌ SÍ (error) | [MISSING_PURCHASE_HISTORY_CAUSES.md](cointracking/behavioral/MISSING_PURCHASE_HISTORY_CAUSES.md) | [CT-002](cases/ct-002-venta-sin-historial-de-compra-previo.md), [CT-017](cases/ct-017-coste-cero-por-compra-omitida-de-ej.md) |
| **Depósito sin retiro coincidente** | ⚠️ VERIFICAR | [PROCEDURE_RECONCILE_TRANSFERS.md](procedures/PROCEDURE_RECONCILE_TRANSFERS.md) | [CT-001](cases/ct-001-transferencia-entre-exchanges-importada-.md), [CT-013](cases/ct-013-wallet-externa-no-importada-fondos-desap.md) |
| **Transacción no aparece** | ❌ SÍ (faltante) | [MISSING_PURCHASE_HISTORY_CAUSES.md](cointracking/behavioral/MISSING_PURCHASE_HISTORY_CAUSES.md) | [CT-012](cases/ct-012-balance-negativo-por-importacion-parcial.md) |
| **Staking como depósito normal** | ⚠️ CLASIFICACIÓN | [STAKING_MECHANICS.md](cointracking/behavioral/STAKING_MECHANICS.md) | [CT-005](cases/ct-005-recompensas-de-staking-clasificadas-como.md) |
| **Airdrop con coste artificial** | ⚠️ CLASIFICACIÓN | [AIRDROPS_MECHANICS.md](cointracking/behavioral/AIRDROPS_MECHANICS.md) | [CT-010](cases/ct-010-airdrop-registrado-como-compra-con-coste.md) |
| **Binance Convert duplicado** | ⚠️ IMPORTACIÓN | [BINANCE_CONVERT_MECHANICS.md](cointracking/behavioral/BINANCE_CONVERT_MECHANICS.md) | [CT-006](cases/ct-006-binance-convert-importado-como-venta-y-c.md) |
| **Comisión en moneda ajena** | ⚠️ COST BASIS | [FEE_HANDLING.md](cointracking/behavioral/FEE_HANDLING.md) | [CT-009](cases/ct-009-comision-fee-omitida-en-la-importacion.md) |
| **Swap DeFi con muchas piezas** | ⚠️ ON-CHAIN | [DEFI_SWAPS_MECHANICS.md](cointracking/behavioral/DEFI_SWAPS_MECHANICS.md) | [CT-015](cases/ct-015-swap-defi-fragmentado-en-varias-operacio.md) |
| **Balance < tenencia esperada** | ⚠️ VERIFICAR | [BALANCE_CALCULATION_ALGORITHM.md](cointracking/behavioral/BALANCE_CALCULATION_ALGORITHM.md) | [CT-004](cases/ct-004-balance-negativo-por-orden-cronologico-i.md) |
| **Dos tickers del mismo activo** | ⚠️ RENAMING | [CSV_FORMAT.md](cointracking/official/CSV_FORMAT.md) §8 | [CT-018](cases/ct-018-token-renombrado-interpretado-como-un-ac.md) |
| **Lending como transferencia** | ⚠️ CLASIFICACIÓN | [LENDING_MECHANICS.md](cointracking/behavioral/LENDING_MECHANICS.md) | [CT-011](cases/ct-011-lending-tratado-como-transferencia-gener.md) |
| **Minería como depósito** | ⚠️ CLASIFICACIÓN | Blockchain > [STAKING_MECHANICS_BLOCKCHAIN.md](blockchains/STAKING_MECHANICS_BLOCKCHAIN.md) | [CT-014](cases/ct-014-recompensas-de-mineria-mining-registrada.md) |
| **Reimportación massiva duplica** | ❌ SÍ (error) | [API_VS_CSV_OVERLAP.md](cointracking/behavioral/API_VS_CSV_OVERLAP.md) | [CT-016](cases/ct-016-duplicados-por-reimportacion-completa-del.md) |
| **Transferencia interna es venta** | ⚠️ CLASIFICACIÓN | [PROCEDURE_RECONCILE_TRANSFERS.md](procedures/PROCEDURE_RECONCILE_TRANSFERS.md) | [CT-007](cases/ct-007-transferencia-interna-confundida-con-ven.md) |
| **Ejecución parcial duplica** | ⚠️ BATCHING | [DUPLICATE_DETECTION_HEURISTICS.md](cointracking/behavioral/DUPLICATE_DETECTION_HEURISTICS.md) | [CT-008](cases/ct-008-duplicados-aparentes-por-ejecucion-parci.md) |
| **Advertencia técnica misteriosa** | ⚠️ CONTEXT | [CT-020](cases/ct-020-advertencia-tecnica-no-es-error-fiscal.md) | — |

---

## 🚨 Árbol de Decisión: "¿Tengo un Error?"

```
┌─ ¿El balance de CoinTracking ≠ exchange real?
│
├─ SÍ
│  ├─ ¿Falta una transacción completa? → CT-012 (importación parcial)
│  ├─ ¿Falta el balance de un activo? → CT-013 (wallet no importada)
│  ├─ ¿El balance está negativo? → CT-004, CT-012 (orden cronológico)
│  └─ ¿El balance es menor pero positivo? → BALANCE_CALCULATION_ALGORITHM
│
└─ NO
   ├─ ¿Hay dos operaciones idénticas?
   │  ├─ SÍ, mismo Trade ID → Legítimo (batching), NO eliminar (CT-008)
   │  ├─ SÍ, Trade ID distinto → Legítimo, NO eliminar (CT-006, CT-003)
   │  └─ NO → Continúa
   │
   ├─ ¿Hay una venta sin cost basis?
   │  ├─ SÍ → Error crítico, buscar compra (CT-002, CT-017)
   │  └─ NO → Continúa
   │
   ├─ ¿Hay un depósito sin retiro? (o viceversa)
   │  ├─ SÍ, mismo amount ± fees → Posible (CT-001 + timestamp)
   │  ├─ SÍ, diferente amount → Falta parte (CT-013)
   │  └─ NO → Continúa
   │
   └─ ✅ Probablemente OK (verifica fiscalidad en CAPITAL_GAINS)
```

---

## 📊 Por Severidad

### 🔴 CRÍTICO — Impide declarar

| Problema | Acción |
|----------|--------|
| Venta sin cost basis | Buscar compra en [CT-002](cases/ct-002-venta-sin-historial-de-compra-previo.md) |
| Saldo negativo | Diagnosticar con [FLOW_NEGATIVE_BALANCE.md](decision-trees/FLOW_NEGATIVE_BALANCE.md) |
| Balance ≠ exchange | Reconciliar con [BALANCE_CALCULATION_ALGORITHM.md](cointracking/behavioral/BALANCE_CALCULATION_ALGORITHM.md) |

### 🟡 IMPORTANTE — Afecta cifras fiscales

| Problema | Acción |
|----------|--------|
| Duplicados legítimos marcados como error | Verificar Trade ID en [DUPLICATE_DETECTION_HEURISTICS.md](cointracking/behavioral/DUPLICATE_DETECTION_HEURISTICS.md) |
| Staking/Airdrops mal clasificados | Reclasificar según [STAKING_MECHANICS.md](cointracking/behavioral/STAKING_MECHANICS.md) / [AIRDROPS_MECHANICS.md](cointracking/behavioral/AIRDROPS_MECHANICS.md) |
| Comisiones omitidas | Verificar con [FEE_HANDLING.md](cointracking/behavioral/FEE_HANDLING.md) |

### 🟢 MENOR — Cosmético o incierto

| Problema | Acción |
|----------|--------|
| Advertencias técnicas | Ver [CT-020](cases/ct-020-advertencia-tecnica-no-es-error-fiscal.md) (contexto) |
| Dos tickers del mismo activo | Fusionar si es seguro ([CT-018](cases/ct-018-token-renombrado-interpretado-como-un-ac.md)) |
| Transacción en timestamp raro | Verificar zona horaria (ADR-005) |

---

## 🔍 Búsqueda por Exchange

### Binance (Spot, Margin, Futures, Earn)
- Importación: [BINANCE_IMPORT_WORKFLOW.md](cointracking/behavioral/BINANCE_IMPORT_WORKFLOW.md)
- Spot: [BINANCE_SPOT_MECHANICS.md](cointracking/behavioral/BINANCE_SPOT_MECHANICS.md)
- Margin: [BINANCE_MARGIN_MECHANICS.md](cointracking/behavioral/BINANCE_MARGIN_MECHANICS.md)
- Futures: [BINANCE_FUTURES_MECHANICS.md](cointracking/behavioral/BINANCE_FUTURES_MECHANICS.md)
- Earn: [BINANCE_EARN_MECHANICS.md](cointracking/behavioral/BINANCE_EARN_MECHANICS.md)
- Convert: [BINANCE_CONVERT_MECHANICS.md](cointracking/behavioral/BINANCE_CONVERT_MECHANICS.md) + [CT-006](cases/ct-006-binance-convert-importado-como-venta-y-c.md)

### Kraken
- [KRAKEN_STAKING_MECHANICS.md](cointracking/behavioral/KRAKEN_STAKING_MECHANICS.md)

### Coinbase
- [COINBASE_ADVANCED_TRADE.md](cointracking/behavioral/COINBASE_ADVANCED_TRADE.md)

### Otro Exchange
- [GENERIC_EXCHANGE_MECHANICS.md](cointracking/behavioral/GENERIC_EXCHANGE_MECHANICS.md)

---

## 💡 Procedimientos Paso a Paso

| Necesito | Documento |
|----------|-----------|
| Auditar mi cuenta entera | [PROCEDURE_AUDIT_ACCOUNT.md](procedures/PROCEDURE_AUDIT_ACCOUNT.md) |
| Reconciliar transferencias | [PROCEDURE_RECONCILE_TRANSFERS.md](procedures/PROCEDURE_RECONCILE_TRANSFERS.md) |
| Resolver "missing purchase history" | [PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md](procedures/PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md) |

---

## 📞 Si Aún No Está Cubierto

1. Busca en [NAVIGATION_MAP.md](NAVIGATION_MAP.md) (por función)
2. Busca en [QUICK_START.md](QUICK_START.md) (guía de inicio)
3. Busca en [GLOSSARY.md](reference/GLOSSARY.md) (definición)
4. Consulta [adr/INDEX.md](../adr/INDEX.md) (contexto de decisión)

Si es un caso nuevo: documéntalo como nuevo caso (nivel C1) en `knowledge/cases/ct-NXX.md`.

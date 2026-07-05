# Nivel C: Casos, Patrones y Procedimientos

**Ubicación:** `knowledge/cases/` (C1), `knowledge/patterns/` (C2), `knowledge/procedures/` (C3)

**Característica:** Documentación de **casos reales auditados**, los patrones que se derivan de ellos, y los procedimientos operativos para resolverlos.

**Autoridad:** `verified` — casos reales del usuario = máxima confianza, no es teoría

**Convención de IDs:** ver ADR-036. `C1` = casos individuales, `C2` = patrones de reconciliación, `C3` = procedimientos.

---

## C1: Casos Reales Auditados (20 casos)

Cada caso documenta un hallazgo real, por qué ocurre, cómo se confunde con otra cosa, y cómo resolverlo.

1. [ct-001-transferencia-entre-exchanges-importada-.md](ct-001-transferencia-entre-exchanges-importada-.md) (KB-C1-001) — Transferencia entre exchanges importada solo en origen
2. [ct-002-venta-sin-historial-de-compra-previo-mis.md](ct-002-venta-sin-historial-de-compra-previo-mis.md) (KB-C1-002) — Venta sin historial de compra previo (Missing Purchase History)
3. [ct-003-api-y-csv-importados-simultaneamente-dup.md](ct-003-api-y-csv-importados-simultaneamente-dup.md) (KB-C1-003) — API y CSV importados simultáneamente (duplicado por doble fuente)
4. [ct-004-balance-negativo-por-orden-cronologico-i.md](ct-004-balance-negativo-por-orden-cronologico-i.md) (KB-C1-004) — Balance negativo por orden cronológico incorrecto (zona horaria)
5. [ct-005-recompensas-de-staking-clasificadas-como.md](ct-005-recompensas-de-staking-clasificadas-como.md) (KB-C1-005) — Recompensas de staking clasificadas como depósito genérico
6. [ct-006-binance-convert-importado-como-venta-y-c.md](ct-006-binance-convert-importado-como-venta-y-c.md) (KB-C1-006) — Binance Convert importado como venta y compra independientes
7. [ct-007-transferencia-interna-confundida-con-ven.md](ct-007-transferencia-interna-confundida-con-ven.md) (KB-C1-007) — Transferencia interna confundida con venta
8. [ct-008-duplicados-aparentes-por-ejecucion-parci.md](ct-008-duplicados-aparentes-por-ejecucion-parci.md) (KB-C1-008) — Duplicados aparentes por ejecución parcial de una orden
9. [ct-009-comision-fee-omitida-en-la-importacion.md](ct-009-comision-fee-omitida-en-la-importacion.md) (KB-C1-009) — Comisión (fee) omitida en la importación
10. [ct-010-airdrop-registrado-como-compra-con-coste.md](ct-010-airdrop-registrado-como-compra-con-coste.md) (KB-C1-010) — Airdrop registrado como compra con coste artificial
11. [ct-011-lending-tratado-como-transferencia-gener.md](ct-011-lending-tratado-como-transferencia-gener.md) (KB-C1-011) — Lending tratado como transferencia genérica
12. [ct-012-balance-negativo-por-importacion-parcial.md](ct-012-balance-negativo-por-importacion-parcial.md) (KB-C1-012) — Balance negativo por importación parcial vía API
13. [ct-013-wallet-externa-no-importada-fondos-desap.md](ct-013-wallet-externa-no-importada-fondos-desap.md) (KB-C1-013) — Wallet externa no importada, fondos desaparecen
14. [ct-014-recompensas-de-mineria-mining-registrada.md](ct-014-recompensas-de-mineria-mining-registrada.md) (KB-C1-014) — Recompensas de minería (mining) registradas como depósito
15. [ct-015-swap-defi-fragmentado-en-varias-operacio.md](ct-015-swap-defi-fragmentado-en-varias-operacio.md) (KB-C1-015) — Swap DeFi fragmentado en varias operaciones on-chain
16. [ct-016-duplicados-por-reimportacion-completa-de.md](ct-016-duplicados-por-reimportacion-completa-de.md) (KB-C1-016) — Duplicados por reimportación completa del mismo periodo
17. [ct-017-coste-cero-por-compra-omitida-de-ejercic.md](ct-017-coste-cero-por-compra-omitida-de-ejercic.md) (KB-C1-017) — Coste cero por compra omitida de ejercicios anteriores
18. [ct-018-token-renombrado-interpretado-como-un-ac.md](ct-018-token-renombrado-interpretado-como-un-ac.md) (KB-C1-018) — Token renombrado interpretado como un activo distinto
19. [ct-019-balance-negativo-tras-eliminar-una-compr.md](ct-019-balance-negativo-tras-eliminar-una-compr.md) (KB-C1-019) — Balance negativo tras eliminar una compra confundida con duplicado
20. [ct-020-advertencia-tecnica-interpretada-como-er.md](ct-020-advertencia-tecnica-interpretada-como-er.md) (KB-C1-020) — Advertencia técnica interpretada como error fiscal definitivo

---

## C2: Patrones de Reconciliación (derivados de los casos C1)

Generalizaciones extraídas de varios casos C1 — qué buscar, no un caso puntual.

- [`knowledge/patterns/PATTERN_DUPLICATE_DETECTION.md`](../patterns/PATTERN_DUPLICATE_DETECTION.md) (KB-C2-001) — Matriz: qué hace/no hace duplicado (deriva de CT-003, CT-008, CT-016)
- [`knowledge/patterns/PATTERN_BALANCE_RECONCILIATION.md`](../patterns/PATTERN_BALANCE_RECONCILIATION.md) (KB-C2-002) — Cómo reconocer saldos inconsistentes
- [`knowledge/patterns/PATTERN_TRANSFER_MATCHING.md`](../patterns/PATTERN_TRANSFER_MATCHING.md) (KB-C2-003) — Heurísticas para emparejar withdrawal/deposit
- [`knowledge/patterns/PATTERN_PURCHASE_POOL_EXHAUSTION.md`](../patterns/PATTERN_PURCHASE_POOL_EXHAUSTION.md) (KB-C2-004) — Síntomas de "purchase pool consumed" (deriva de CT-002, CT-017)

---

## C3: Procedimientos Operativos (cómo resolver, paso a paso)

- [`knowledge/procedures/PROCEDURE_AUDIT_ACCOUNT.md`](../procedures/PROCEDURE_AUDIT_ACCOUNT.md) (KB-C3-001) — 6 fases de auditoría completa
- [`knowledge/procedures/PROCEDURE_RECONCILE_TRANSFERS.md`](../procedures/PROCEDURE_RECONCILE_TRANSFERS.md) (KB-C3-002) — Emparejar withdrawal/deposit, incluye 4 causas raíz de transferencias huérfanas (blockchain delay, importación parcial, no acreditado, migración incompleta)
- [`knowledge/procedures/PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md`](../procedures/PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md) (KB-C3-003) — Resolver origen de coste ausente

---

## Peculiaridades de Exchange Específicas (viven en Nivel B, no C)

Las mecánicas específicas de cada exchange (qué genera cada tipo de operación, casos límite) están en `knowledge/cointracking/behavioral/`, no en Nivel C — Nivel C es para casos/patrones/procedimientos, Nivel B es para "cómo funciona la plataforma". Documentos relevantes, con hallazgos de agp2025 ya incorporados:

- [`BINANCE_SPOT_MECHANICS.md`](../cointracking/behavioral/BINANCE_SPOT_MECHANICS.md) (KB-B2-001) — incluye dust→BNB, Binance Convert, swaps DeFi, Binance Earn
- [`BINGX_MECHANICS.md`](../cointracking/behavioral/BINGX_MECHANICS.md) (KB-B2-010) — incluye el caso verificado de Copy Trading no exportado (agp2025, ~694,67 USDT)
- [`STAKING_MECHANICS.md`](../cointracking/behavioral/STAKING_MECHANICS.md) (KB-B1-001) — incluye tipos de staking (bloqueado, liquid/DeFi), clasificación fiscal RCM, regla de valor a fecha de recepción
- [`BYBIT_MECHANICS.md`](../cointracking/behavioral/BYBIT_MECHANICS.md) (KB-B2-011) — Trading Spot y Derivados en Bybit
- [`OKX_MECHANICS.md`](../cointracking/behavioral/OKX_MECHANICS.md) (KB-B2-012) — Trading completo + Web3 en OKX
- [`KRAKEN_STAKING_MECHANICS.md`](../cointracking/behavioral/KRAKEN_STAKING_MECHANICS.md) (KB-B2-006) — Cómo CoinTracking maneja staking y rewards de Kraken
- [`AIRDROPS_MECHANICS.md`](../cointracking/behavioral/AIRDROPS_MECHANICS.md) (KB-B1-002) — Cómo CoinTracking maneja airdrops (regalos de tokens)
- [`BRIDGES_AND_WRAPPING.md`](../blockchains/BRIDGES_AND_WRAPPING.md) (KB-B3-004) — Bridges y wrapped tokens

---

## Cómo Usar Nivel C

### Para Auditor/Usuario

1. **¿Hallazgo puntual?** Busca en C1 (¿ya hay un caso parecido documentado?)
2. **¿Patrón recurrente?** Busca en C2 (¿qué categoría de problema es?)
3. **¿Necesitas resolverlo paso a paso?** Busca en C3 (procedimiento operativo)
4. **¿Es específico de un exchange?** Busca en Nivel B `behavioral/`

### Para Documentación

**Antes de crear un documento nuevo en `knowledge/cases/`:**
1. Consulta ADR-036 (convención de IDs) — no reutilices IDs existentes
2. Verifica que no exista ya contenido similar en Nivel B (`cointracking/behavioral/`) o en los C1-C3 existentes
3. Si el contenido es "cómo funciona X exchange/mecánica", va en Nivel B, no en Nivel C
4. Si es un caso real nuevo, continúa la numeración: siguiente caso C1 sería `KB-C1-021`

---

## Relación con Otros Niveles

| Relación | Conexión |
|----------|----------|
| **Nivel B** | Documenta mecánicas de exchanges/protocolos; C1 son casos que ilustran fallos sobre esas mecánicas |
| **Nivel A** | Nivel C verifica y valida los principios de Nivel A (fiscal, CSV oficial) |
| **ADRs** | Nivel C operacionaliza lo decidido en ADRs (validación, vigencia, convención de IDs) |
| **Nivel D** | Checklists y árboles de decisión se construyen a partir de los casos C1 |
| **Skills** | Nivel C guía el diagnóstico en `/audit-cointracking` y `/spanish-tax-return` |

---

## Política de Vigencia (ADR-008/ADR-022)

Cada documento declara `last_verified`, `valid_until` y `source`. Antes de citar un caso en auditoría: comprobar que no esté expirado; si expira pronto, reverificar contra dato real.

---

**Última actualización:** 2026-07-05 (verificado: Bybit, OKX, Kraken staking, Airdrops y Bridges/wrapped tokens ya estaban documentados en Nivel B desde una sesión anterior — no hacía falta crearlos de nuevo, solo enlazarlos aquí)

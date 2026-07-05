---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-019: Cierre y corrección de ADR-018 — `get_gains` confirmado fiable, la reconstrucción FIFO manual era la que fallaba

**Status:** Accepted

**Date:** 2026-07-03

## Context

ADR-018 dejó la brecha BTC/USDC/OM como hipótesis `[VERIFICAR]` ("asimetría de valoración en permutas"), pendiente de contrastar contra el Tax Report oficial de CoinTracking. El usuario descargó los Tax Reports oficiales (España, FIFO) de **2024 y 2025** en Excel y se hizo el contraste real: las 39 operaciones de BTC (y todas las de OM) resultaron ser del ejercicio **2024**, no 2025 — el primer intento de mirar solo el informe de 2025 no encontraba nada porque era el año equivocado, no porque la cifra fuera cero.

## Decision

**Decisión — la corrección se basa en el contraste real:**

| Activo | Tax Report oficial (2024+2025) | `get_gains(price:"oldest")` | Reconstrucción FIFO manual |
|---|---|---|---|
| BTC | 503,50 € | 492,87 € | 94,71 € |
| USDC | 554,61 € | 553,93 € | 635,61 € |
| OM | 1.027,49 € | 1.027,49 € | 1.114,89 € |

El Tax Report oficial coincide casi al céntimo con `get_gains`; la reconstrucción manual estaba mal en los tres activos. **Se corrige la conclusión de ADR-018:** la hipótesis de "asimetría de valoración por lado de permuta" queda **descartada como causa raíz** (coincidía en magnitud por casualidad). La causa más probable real: la reconstrucción manual, operación por operación, no arrastraba bien la base de coste a través de cadenas de permutas cripto-cripto; `get_gains` sí lo hace.

Se actualiza `knowledge/cointracking/COST_BASIS_AND_VALIDATION.md` §4.4 (de "hipótesis abierta" a "resuelto"), el sub-paso de la fase 6 de `audit-cointracking/SKILL.md`, `reports/output/agp2025/REGISTRO-CAMBIOS.md` y la memoria de proyecto (`audit_state`).

## Consequences

- ✅ Regla operativa nueva y más simple que la de ADR-018: ante una discrepancia `get_gains` vs. reconstrucción propia, **confiar por defecto en `get_gains`/Tax Report oficial**, no en el cálculo manual — salvo que un contraste real diga lo contrario en ese caso.
- ✅ Dato adicional relevante para la declaración: BTC (503,50 €) y OM (1.027,49 €) son ganancias del **ejercicio 2024**, no 2025 — queda pendiente (fuera de este ADR) confirmar con el usuario/asesor si ya se declararon en su ejercicio.
- ⚠️ Sigue sin verificarse el esquema exacto de `get_trades(trade_prices=1)` (ADR-018 punto pendiente) — ya no es urgente porque la recomendación operativa ya no depende de reconstruir manualmente ese cálculo, pero queda abierto si algún día se quiere automatizar un chequeo de coherencia distinto.
- ✅ Ejemplo documentado de por qué ADR-009 (cero invención) importa incluso para el propio agente: una hipótesis con evidencia aparentemente fuerte (coincidencia de magnitud) resultó ser una pista falsa; solo el contraste contra la fuente autorizada (Tax Report oficial) lo confirmó.

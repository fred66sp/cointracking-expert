# ADR-018: Discrepancia `get_gains` vs FIFO manual — documentar como hipótesis, no automatizar en `ct_audit.py` (aún)

**Status:** Accepted

**Date:** 2026-07-03

## Context

Petición de Copilot (explotación, ADR-012) en `AGENT_CHANGE_REQUESTS.md`, tras investigar en el caso real `agp2025` por qué `cointracking_get_gains(price:"oldest")` (FIFO) da una ganancia de BTC (+492,87 €) muy distinta de una reconstrucción FIFO manual sobre `get_trades(trade_prices=1)` (+94,71 €; brecha ~398 €). Descartó comisiones, duplicados/mal tipado y FIFO-vs-pool como causa; la variable que explica una magnitud del mismo orden es la asimetría de qué lado de la permuta (compra o venta) se usa para valorar en EUR cada operación. Propuso automatizar un "delta bridge" en `tools/ct_audit.py` que descomponga la brecha por activo y avise si supera un umbral.

**Decisión:**

Se documenta el hallazgo, pero **no se automatiza todavía** en `tools/ct_audit.py`:

1. Nueva sección `knowledge/cointracking/COST_BASIS_AND_VALIDATION.md` §4.4, etiquetada explícitamente como **hipótesis empírica no confirmada por CoinTracking** (no hay artículo oficial que la respalde), con el caso de referencia, lo que se descartó como causa, y el recipe manual para reproducir el diagnóstico en otro caso.
2. Nuevo sub-paso en la fase 6 ("Cierre") del Paso 1 de `.claude/skills/audit-cointracking/SKILL.md`: si `get_gains` diverge de forma material de una reconstrucción FIFO manual, aplicar el recipe de COST_BASIS §4.4 antes de concluir, y marcar el resultado `[VERIFICAR]`.
3. **No se automatiza en `tools/ct_audit.py`** porque ese tool opera de forma determinista sobre el CSV export (esquema fijo y ya verificado, `CSV_FORMAT.md`), mientras que este diagnóstico depende de la respuesta JSON de `get_trades(trade_prices=1)` del MCP, cuyo esquema exacto de campos de valoración por lado (nombres de campo reales) **no está verificado** en `knowledge/cointracking/MCP_API.md`. Fijar nombres de campo sin verificarlos habría sido inventar contra ADR-009. Queda como trabajo pendiente real: verificar el esquema de `trade_prices=1` contra una llamada real y, solo entonces, decidir si se automatiza (¿en `ct_audit.py` extendido a JSON, o en un script nuevo?) — no se decide la forma aquí para no anticipar sin datos.

## Decision

[Decision not found]

## Consequences

- ✅ El hallazgo del caso real queda capturado como conocimiento reutilizable, sin presentarlo como regla cerrada ni como comportamiento documentado por CoinTracking.
- ✅ La skill guía explícitamente a no dejar sin explicar una brecha material entre `get_gains` y una reconstrucción manual, ni a declararla "correcta" sin pasar por el recipe.
- ✅ No se infla el alcance de `tools/ct_audit.py` con un chequeo basado en un esquema de API sin verificar (ADR-009).
- ⚠️ Pendiente real, no cerrado por este ADR: verificar el esquema exacto de `get_trades(trade_prices=1)` y decidir la forma de automatización cuando haya otro caso real o tiempo para esa verificación.
- ✅ Entrada `AGENT_CHANGE_REQUESTS.md` del 2026-07-03 ("Chequeo automático de discrepancia FIFO...") marcada como hecha, con el alcance real aplicado (documentación + playbook, no automatización de código).

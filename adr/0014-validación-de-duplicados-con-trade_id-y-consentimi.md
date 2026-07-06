---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-014: Validación de duplicados con trade_id y consentimiento explícito

**Status:** Accepted

**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)

**Date:** 2026-07-03

## Context

El 2026-07-03, al auditar CoinTracking, se detectaron como "duplicados exactos" operaciones FLOKI del 17.03.2024 que eran legítimas. El algoritmo comparaba solo campos del CSV (`(tipo, buy_amount, buy_currency, sell_amount, sell_currency, fee, exchange, fecha)`); como múltiples operaciones reales de Binance ocurrieron en el mismo segundo (batching), todas parecían idénticas. Basándose en ello, el usuario eliminó 29 copias que en realidad eran transacciones distintas (con trade_ids distintos en Binance API). El resultado: ~1,6 millones de FLOKI se perdieron del saldo hasta restaurar de backup.

**Raíz:** ct_audit.py no disponía de trade_id para distinguir operaciones; el CSV de CoinTracking tampoco lo incluye en todas las filas. La lógica de "duplicado = campos 100% idénticos" falló en presencia de transacciones legítimas separadas pero aparentemente idénticas.

## Decision

**Decisión:**

1. **Usar trade_id como identificador único cuando esté disponible (fortaleza):**
   - Si el CSV incluye trade_id, dos operaciones con trade_ids distintos **nunca son duplicados**, aunque todos los demás campos sean idénticos.
   - Si trade_id está vacío, caer a la heurística siguiente.

2. **Heurística cautelosa para duplicados sin trade_id:**
   - Si hay **exactamente 2 copias idénticas** → probable duplicado de reimportación; marcar para revisión.
   - Si hay **3-10 copias** → **ADVERTENCIA** (posibles operaciones legítimas en el mismo segundo); reportar pero **no recomendar eliminar**.
   - Si hay **más de 10 copias idénticas** → **muy probablemente legítimas** (batching de Binance, transacciones FIAT repetidas, recompensas recurrentes); reportar como "INFORMACIÓN" e **indicar que requiere confirmación manual en Binance API** antes de eliminar.

3. **Implementar consentimiento informado (refuerza ADR-009):**
   - Antes de eliminar duplicados, el agente **lista exactamente cuáles se eliminarán** con ejemplos concretos (monto, tipo, fecha, cantidad de copias).
   - **Advierte:** "Si estas operaciones son legítimas (según Binance API), eliminarlas causará saldo negativo del activo."
   - **Pide confirmación explícita** del usuario antes de proceder.

4. **Usar el MCP como árbitro (cuando esté disponible):**
   - Si el MCP de CoinTracking está conectado, consulta el número de operaciones de ese tipo con trade_ids distintos.
   - Si son más de 1, son legítimas; no eliminar.
   - Si es solo 1, es un duplicado real; OK eliminar.

## Consequences

- ✅ Evita falsos positivos como el del 2026-07-03
- ✅ Refuerza ADR-009 (consentimiento informado antes de actuar)
- ✅ El usuario toma la decisión final, no el agente
- ⚠️ Requiere que el usuario verifique en Binance API si no confía
- ⚠️ ct_audit.py debe ser más conservador; menos automatización

## Notes

**Cambios en CLAUDE.md:**

- Agregar sección ⚠️ sobre falsos positivos en duplicados.
- Instruir: "Antes de eliminar duplicados, verifica en Binance que tengan el MISMO `Tx ID`. Si tienen IDs distintos, son legítimas."

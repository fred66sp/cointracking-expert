---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-002: Fuente de verdad y resolución de conflictos

**Status:** Accepted

**Date:** 2026-07-04

## Context

Cuando una operación en CoinTracking **contradice** la realidad del exchange, la billetera, el banco o la blockchain, ¿cuál es correcta? El agente debe tener una regla clara para resolver conflictos, porque:

1. CoinTracking puede tener datos desactualizados (reimportación incompleta, API con delay, usuario editó algo manualmente)
2. Un CSV exportado puede tener timestamp truncado o comisión arreglada a mano
3. El exchange API puede devolver datos parciales (histórico limitado, operaciones canceladas no mostradas)
4. La blockchain es la verdad absoluta, pero caro verificarla operativamente

Sin una jerarquía clara, el auditor no sabe si confiar en CoinTracking o si debe buscar en el exchange real.

### Caso real

Usuario importa CSV de Binance. Luego conecta API de Binance. En CoinTracking:
- CSV dice "Compra BTC a 30.000€ el 2023-03-15"
- API dice "Compra BTC a 31.000€ el 2023-03-15"

¿Cuál es correcta? Sin regla, es ambiguo.

## Decision

Se establece una **jerarquía de fuentes de verdad**, aplicada en este orden de confiabilidad (mayor a menor):

### Nivel 1: Blockchain (verdad absoluta)

Verificable directamente en la red (Etherscan, Blockchain.com, etc.). **Nunca miente**, pero requiere esfuerzo manual.

**Cuándo usarla:**
- Verificar un depósito crítico (origen de fondos para declaración)
- Resolver un conflicto irreconciliable entre CoinTracking y exchange
- Validar una operación de gran importe (~1 BTC+)

**Límite:** No escalable para auditar 1000s de operaciones; usarla selectivamente.

### Nivel 2: Exchange API oficial (fuente primaria)

La API del exchange (Binance, Kraken, Coinbase, etc.) es la "verdad del exchange". Es el datos directo que generó la operación.

**Prevalece sobre:** CoinTracking, CSV exportado.

**Limitaciones conocidas:**
- Algunas APIs tienen rango temporal limitado (p. ej. Binance API solo devuelve 500 trades anteriores en caché)
- Operaciones canceladas pueden no aparecer
- Comisiones a veces se muestran como '0' si son cero-fee
- Timestamp puede diferir en segundos (timezone del servidor)

**Cuándo confiar:**
- Para operaciones actuales (últimos 90 días)
- Para reconciliar depósitos/retiradas (las APIs de withdrawals son muy confiables)
- Para verificar saldos actuales (holdings)

### Nivel 3: CSV exportado por el usuario desde el exchange

Snapshot de datos que el usuario exportó manualmente del exchange en algún momento.

**Prevalece sobre:** CoinTracking importado después del CSV (si CoinTracking es reimportación del mismo CSV).

**Limitaciones:**
- Es estático (refleja el estado en una fecha)
- Puede tener cortes/inconsistencias si la exportación falló
- Timestamp sin zona horaria (ambiguo)
- Depende de cómo exporte el usuario (algunos exchanges ofrecen CSV en múltiples formatos)

**Regla:** Si importaste CSV en CoinTracking hace 6 meses, y luego reimportaste la API, la API es más nueva → prevalece.

### Nivel 4: CoinTracking (verdad operativa del agente)

Los datos en CoinTracking son la **verdad del agente**: es lo que el auditor usará para calcular impuestos. Pero **no es la verdad del exchange**.

**Prevalece sobre:** Inferencias del agente, suposiciones, patrones.

**Limitaciones:**
- Puede estar corrompido si el usuario editó manualmente algo
- Puede tener duplicados, operaciones fantasma, tipos mal clasificados
- Es el resultado de N importaciones desde N fuentes — acumula errores

**Regla crítica:** Antes de confiar en CoinTracking para una decisión irreversible (p. ej. eliminar un duplicado), **verifica primero en el exchange real** (Nivel 2 o 1).

### Nivel 5: Inferencia del agente

El agente puede inferir que:
- "Esto parece un duplicado porque coinciden 5 campos"
- "Esta operación no tiene origen conocido (Missing Purchase History)"
- "Este saldo es imposible (negativo)"

**Pero la inferencia no es evidencia.**

**Regla:** Nunca actúes basándote en una inferencia sin confirmación del usuario o verificación en el exchange real.

---

## Casos de resolución de conflictos

### Caso 1: CoinTracking vs Exchange API

**Conflicto:** CoinTracking muestra "Compra 1 BTC a 50.000€" pero la API de Binance muestra "Compra 1 BTC a 50.100€".

**Resolución:**
1. ¿Cuándo se importó a CoinTracking? Si hace 6 meses, la API es más nueva → prevalece API.
2. ¿Se editó manualmente en CoinTracking? Si sí → CoinTracking podría estar corrupto → prevalece API.
3. ¿Cuál está en la blockchain? Si la blockchain muestra 50.100€ → prevalece blockchain.

**Acción:** Marcar en auditoría como `⚠️ Discrepancia: CoinTracking 50.000€ vs API 50.100€. Usar API (50.100€) para cálculos fiscales.`

### Caso 2: CSV de años atrás vs API actual

**Conflicto:** El usuario exportó CSV de Binance el 2022-01-01 (operaciones 2020-2021). Ahora conecta API de Binance (2023-presente). Aparecen 50 operaciones duplicadas con precios ligeramente distintos.

**Resolución:**
1. El CSV es más antiguo, pero es el "original".
2. La API es más nueva, pero podría tener retrasos o cambios.
3. Si el número de operaciones coincide pero los precios no → problema de fuente de datos (p. ej. CoinTracking recalculó precios).

**Acción:** Importar CSV primero (Fase 2 del ADR-027), validar. Luego importar API, comprobar solapamiento. Si hay duplicados, marcar como `[PENDIENTE DE VERIFICAR]` — no borrar automáticamente.

### Caso 3: Depósito bancario vs CoinTracking

**Conflicto:** El usuario dice "deposité 10.000€ el 2023-06-15" (confirmado en banco). CoinTracking muestra "Depósito 9.950€ el 2023-06-14".

**Resolución:**
1. Blockchain no aplica (es fiat, no on-chain).
2. Exchange no tiene constancia (los depósitos SEPA son off-exchange hasta llegar a la billetera del exchange).
3. El banco es la verdad del fiat.

**Acción:** El saldo en el banco es incuestionable. Si CoinTracking difiere, editarlo en CoinTracking (o dejar marcado como `⚠️` en auditoría) para que la discrepancia sea transparente.

---

## Consequences

**Positive:**

- **Claridad:** Cuando hay conflicto, hay regla, no ambigüedad
- **Trazabilidad:** El auditor documenta "preferí X sobre Y porque el nivel de confianza de X es mayor"
- **Protección:** Evita confiar ciegamente en CoinTracking (que puede estar corrupto)
- **Escalabilidad:** El agente sabe cuándo necesita validación manual (blockchain) vs cuándo puede confiar en la API

**Negative:**

- **Requiere trabajo manual:** Para cada conflicto crítico, verificar en el exchange real (5-10 min por operación)
- **Información incompleta:** A veces el exchange API también está limitada (histórico antiguo perdido)
- **Blockchain no es práctico:** Verificar 100s de operaciones en blockchain es inviable operativamente
- **Responsabilidad compartida:** El usuario debe confirmar que la jerarquía es correcta para su caso

## Notes

### Relación con ADRs existentes

- **ADR-004:** Reconciliación basada en datos reales antes de cerrar specs — esta es la aplicación concreta
- **ADR-009:** Protocolo crítico — nunca confiar ciegamente en una fuente
- **ADR-014:** Validación de duplicados — aplica aquí también (antes de borrar, verifica en el exchange real)
- **ADR-027:** Integración de nuevos exchanges — usa esta jerarquía al validar post-importación

### Aplicación en auditoría

Esta jerarquía se aplica en ADR-017 (orden de diagnóstico):
1. **Diagnóstico de balances:** CoinTracking vs Holdings reales del exchange (Nivel 2) vs Blockchain (Nivel 1)
2. **Diagnóstico de transferencias:** Exchange A vs Exchange B (ambos Nivel 2)
3. **Missing Purchase History:** Si no hay origen en CoinTracking, buscar en Nivel 2/1

### Pendientes

- **[PENDIENTE]** Definir un protocolo automático para detectar cuándo CoinTracking difiere de API en >5% (alerta automática)
- **[PENDIENTE]** Implementar caché de consultas Blockchain para operaciones críticas (evitar re-consultarlas)
- **[PENDIENTE]** Documentar las limitaciones de cada exchange API (rango temporal, campos disponibles)

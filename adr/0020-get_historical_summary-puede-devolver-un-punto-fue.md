# ADR-020: `get_historical_summary` puede devolver un punto fuera del rango `end` pedido — filtrar por fecha en el consumidor

**Status:** Accepted

**Date:** 2026-07-03

## Context

Petición de Copilot (explotación, ADR-012) durante la preparación de la renta 2025 (`agp2025`): `cointracking_get_historical_summary(start=1735686000, end=1767221999)` (año natural 2025) devolvió, además de la serie diaria de 2025 esperada, un punto fechado el día de la consulta (2026-07-03) — fuera del rango pedido. Relevante porque el chequeo del Modelo 721 depende de un valor exacto a 31/12.

Se revisó el código del servidor MCP propio (`cointracking-mcp/internal/tools/historical_summary.go` y `cached.go`): reenvía `start`/`end` sin modificar a la API de CoinTracking, y la clave de caché incluye ambos parámetros — se descarta que sea un bug de nuestro servidor o de la caché. Lo más probable es que sea la propia API de CoinTracking la que añade un punto "actual" adicional a la serie histórica, pero **no se ha podido confirmar contra documentación oficial** (no hay artículo público que documente la semántica exacta de `start`/`end` para este método).

## Decision

**Decisión:**

Documentar el hallazgo como advertencia empírica (no como bug confirmado) en `knowledge/cointracking/MCP_API.md`, con la mitigación práctica: **cualquier consumo que dependa de un corte exacto de fecha (Modelo 721, valoración a 31/12) debe filtrar la serie devuelta por fecha él mismo**, sin confiar en que `end` corte de forma estricta. Se añade la misma advertencia al Paso 5 de `.claude/skills/spanish-tax-return/SKILL.md`.

No se modifica el servidor Go (`cointracking-mcp`) para filtrar server-side: haría falta confirmar primero si el comportamiento es realmente de la API de CoinTracking (no de nuestro código) y si es deseable filtrar ahí o dejarlo explícito para el consumidor — se prefiere la mitigación simple (filtrar en el playbook) antes que tocar código del servidor sobre una causa no confirmada.

## Consequences

- ✅ El caso real de Copilot (Modelo 721 de 2025) no se vio afectado en la práctica porque ya filtró manualmente los puntos más cercanos a 31/12 al preparar la declaración.
- ✅ Queda documentada la mitigación para futuros casos (`spanish-tax-return`, cualquier corte a 31/12).
- ⚠️ Pendiente real, no cerrado por este ADR: confirmar si el comportamiento es de la API de CoinTracking o de alguna capa intermedia, y decidir si merece la pena filtrar server-side en `cointracking-mcp` una vez confirmado.
- ✅ Entrada `AGENT_CHANGE_REQUESTS.md` del 2026-07-03 ("Verificar/normalizar `get_historical_summary`...") marcada como hecha, con el alcance real aplicado (documentación + mitigación en el playbook, no cambio de servidor).

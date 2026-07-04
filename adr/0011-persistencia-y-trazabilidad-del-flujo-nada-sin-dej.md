# ADR-011: Persistencia y trazabilidad del flujo (nada sin dejar rastro)

**Status:** Accepted

**Date:** 2026-07-02

## Context

Se detectó un fallo grave: tras auditar Coinbase y **aplicar cambios reales** (borrado de 2 filas que afectan a datos fiscales), **no quedó ningún registro persistente** — ni informe, ni bitácora de cambios, ni se guardó la ruta de la fuente de verdad que aportó el usuario. Si se cierra el chat, se pierde todo el contexto. En un agente crítico (ADR-009) esto es inaceptable: los cambios que van a Hacienda deben ser trazables y el trabajo debe sobrevivir entre sesiones.

**Decisión — todo lo importante del flujo se persiste:**

1. **Informe por auditoría.** Toda auditoría o preparación fiscal genera un **informe persistente** en `reports/output/` (nombre con fecha), con hallazgos, acciones, verificación y pendientes. No se deja solo en el chat.
2. **Registro de cambios (append-only).** Todo cambio aplicado en CoinTracking se anota en `reports/output/REGISTRO-CAMBIOS.md`: qué se cambió, por qué, **evidencia**, estado **antes → después** y **verificación** en vivo. Nunca se borran entradas.
3. **Contexto durable en memoria.** Las rutas de fuentes de datos del usuario, el **estado de la auditoría** (cuentas hechas/pendientes) y las decisiones tomadas por chat se guardan en la **memoria** del proyecto (sobrevive entre sesiones). `CLAUDE.md` indica dónde vive todo.
4. **Ningún cambio consecuente sin rastro.** Si por chat se pide algo que altera datos o decisiones: código/decisiones → git (commit/ADR); datos de la cuenta → `reports/output/`; contexto → memoria. Al **retomar** en una sesión nueva, leer primero la memoria y `reports/output/`.

## Decision

[Decision not found]

## Consequences

- ✅ Trazabilidad completa de los cambios que acaban en Hacienda (refuerza ADR-009)
- ✅ Continuidad entre sesiones: un chat nuevo recupera el estado
- ✅ El asesor puede reconstruir qué se tocó y por qué
- ⚠️ Informes y registro contienen datos reales → viven en `reports/output/` (gitignored); la memoria es privada (`~/.claude`)

---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-006: El producto es un agente de IA auditor (Claude Code) sobre la base de conocimiento

**Status:** Accepted

**Date:** 2026-07-02

## Context

El `ROADMAP.md` y el `PROJECT_CHARTER.md` describen un framework/SDK grande: motores deterministas en Python (ledger, FIFO, fiscal, tenencias, reportes), luego CLI, API, y la IA relegada a la Fase 7. Es un plan de meses. Al revisar el objetivo real, se constata que lo que aporta valor **ahora** —y que reaprovecha todo el conocimiento ya documentado— es un **agente de IA que audita los datos de CoinTracking del usuario**, no el SDK completo.

## Decision

**Decisión:**

El **producto principal a corto plazo** es un **agente auditor de IA** que:
- Vive en **Claude Code** como **subagente + skill** (sin infraestructura de código propia).
- Usa como "cerebro" la base de conocimiento del repo (`knowledge/cointracking/*`, `knowledge/taxation/spain/*`).
- Accede a los datos por **dos vías**: el **MCP de la API de CoinTracking** (datos en vivo, cuando esté conectado) y el **CSV export** (Trade Table) como fuente/validación cruzada.
- Detecta y **explica** problemas de auditoría citando las reglas documentadas, con el formato evidencia → causa → impacto → recomendación.

**Límite de determinismo (reconciliación con FOUNDATION):**

FOUNDATION establece "la IA explica; los motores calculan" y exige reproducibilidad. Un agente LLM no es determinista. Por tanto:
- El agente **encuentra y explica** problemas (análisis cualitativo): transferencias huérfanas, ventas sin base de coste, duplicados, saldos imposibles, incoherencias fiscales. Esto es justo lo que FOUNDATION autoriza para la IA ("explicar, guiar, diagnosticar, resumir, asistir").
- El agente **no** produce cifras fiscales vinculantes por sí mismo. Las cantidades exactas (FIFO, base imponible) se marcan como **estimación no vinculante** o se delegan a un **cálculo determinista** (helper/función), nunca al criterio libre del LLM.

## Consequences

- ✅ Resultado utilizable de inmediato, reaprovechando todo el conocimiento y los principios ya escritos
- ✅ Coherente con FOUNDATION si se respeta el límite de determinismo
- ✅ Los "motores" del charter pasan a ser un **playbook de auditoría** (procedimientos del agente); pueden materializarse como helpers deterministas si se necesita rigor numérico
- ⚠️ **Supera al ROADMAP/charter en el corto plazo:** el SDK completo queda como visión futura/opcional, no como camino inmediato. Ver nota en `ROADMAP.md`.
- ⚠️ La calidad del agente depende de la cobertura del conocimiento; huecos conocidos (p. ej. fiscalidad de staking) limitan su precisión hasta cerrarse

**Próximos pasos:**

1. Definir el subagente auditor (`.claude/agents/`) con rol, principios y límite de determinismo.
2. Escribir el playbook de auditoría como skill invocable (`.claude/skills/`).
3. Conectar el MCP de CoinTracking; usar el CSV como alternativa.

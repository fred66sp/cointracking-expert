---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-017: Protocolo de diagnóstico en orden fijo para la auditoría (endurecer falsos positivos)

**Status:** Accepted

**Date:** 2026-07-03

## Context

Petición de Copilot (explotación, ADR-012) en `AGENT_CHANGE_REQUESTS.md` (2026-07-03): el playbook de reconciliación (`Paso 1` de `audit-cointracking/SKILL.md`) listaba los ocho chequeos sin un orden vinculante, lo que abre la puerta a falsos positivos — p. ej. marcar "duplicados" o "huérfanas" cuando la causa real es cobertura incompleta de exchanges, o recomendar un borrado antes de haber verificado el `Trade ID`. La propuesta se apoyó en prospección del centro de ayuda oficial de CoinTracking (`READ FIRST: General account imbalances`, `Duplicate Transactions`, `Missing Transactions Report`, `Validate Transactions`, `Roll Forward / Audit Report`, avisos de "purchase pool agotado", `Binance Import Restrictions`), cuyo contenido relevante ya estaba destilado en `knowledge/cointracking/COST_BASIS_AND_VALIDATION.md` y `CSV_FORMAT.md` — no hizo falta destilar conocimiento nuevo, solo **reordenar y explicitar** el playbook existente.

## Decision

**Decisión:**

Reescribir el Paso 1 de `.claude/skills/audit-cointracking/SKILL.md` con un **orden fijo de 6 fases** (cada una reduce falsos positivos de la siguiente):

1. Cobertura de fuentes/periodo y saldos (incluye saldos negativos).
2. Duplicados, con verificación de Trade ID/Tx ID obligatoria (ADR-014) **antes** de recomendar cualquier eliminación.
3. Transferencias huérfanas y orden temporal, por niveles (Tx Hash fuerte / heurístico con tolerancias — los umbrales exactos siguen abiertos, `CSV_FORMAT.md` §11.2, no se inventan).
4. Tipos, comisiones en tercera moneda, ventas sin base de coste y colisión de tickers.
5. Interpretación de avisos del "purchase pool" agotado.
6. Cierre: coherencia fiscal (FIFO) y riesgos residuales.

Se añade además una regla explícita, generalizando ADR-014: **nunca recomendar un borrado masivo sin evidencia por fila y confirmación explícita del usuario**, aplicable a cualquier hallazgo (no solo duplicados).

No se ha tocado `spanish-tax-return/SKILL.md` porque ya delega la reconciliación en `audit-cointracking` (Paso 1 de esa skill) sin duplicar el playbook. El subagente `.claude/agents/cointracking-auditor.md` (usado para análisis profundo dentro de la misma skill) se alinea con el mismo orden de 6 fases y con la regla de no borrado sin confirmación, para que no diverja del playbook cuando se delega en él.

## Consequences

- ✅ Reduce el riesgo de diagnósticos apresurados: el agente no declara "duplicado" o "huérfana" antes de confirmar cobertura completa de fuentes.
- ✅ No inventa umbrales numéricos no fundamentados (ventana temporal, tolerancia de importe) — se documentan como heurística abierta, coherente con ADR-009 (cero invención).
- ✅ Entrada `AGENT_CHANGE_REQUESTS.md` del 2026-07-03 marcada como hecha.
- ⚠️ Pendiente real (no cerrado por este ADR): definir umbrales de emparejamiento de transferencias con más datos (`CSV_FORMAT.md` §11.2) y el destilado del "purchase pool" ya existía, así que este ADR es de **reordenación operativa**, no de conocimiento nuevo.

# ADR-015: Integración de la base de casos ChatGPT como v2 curada (patrones de reconciliación)

**Status:** Accepted

**Date:** 2026-07-03

## Context

Copilot (explotación) propuso, vía `AGENT_CHANGE_REQUESTS.md` (petición 2026-07-02), integrar `cointracking_casos_extended.yaml` (20 casos generados con un prompt curado a un agente ChatGPT auxiliar, ver handoff `reports/output/2026-07-02_handoff_integracion_casos_chatgpt.md`) como ampliación de `cointracking_casos_base.yaml` (10 casos, esquema mínimo, ya en el repo). El candidato aportaba más cobertura y anti-patrones, pero con heterogeneidad de estilo (listas inline vs bloque), campos vacíos como `""` en vez de `null`, y profundidad desigual en evidencia/diagnóstico.

## Decision

**Decisión:**

Se ejecuta la migración por fases definida en el handoff:

- **Fase A (esquema):** se fija un esquema canónico de 16 campos (ver `knowledge/patterns/INDEX.md` §Esquema). Todos los `""` pasan a `null`; todas las listas se homogeneizan en formato bloque.
- **Fase B (curación):** los casos más resumidos del candidato (antiguos CT-004/05/06/09/11-20) se amplían con evidencia mínima accionable y pasos de diagnóstico concretos, y se enlazan con conocimiento ya existente del repo (`COST_BASIS_AND_VALIDATION.md`, `CSV_FORMAT.md`, `WEB_APP_GUIDE.md`) en vez de inventar detalle nuevo sin respaldo. Los casos de duplicados (CT-003, CT-008, CT-016, CT-019) se alinean explícitamente con **ADR-014** (validación por `trade_id` y consentimiento antes de eliminar). El caso de airdrops (antiguo CT-010) mantiene `nivel_confianza: pendiente_verificar` porque el tratamiento fiscal exacto no está cerrado en `knowledge/taxation/spain/PENDIENTES.md`.
- **Fase C (versionado):** se crea `knowledge/patterns/cointracking_casos_v2.yaml` como base **vigente**. `cointracking_casos_base.yaml` (raíz del repo) pasa a **legacy/deprecado**: no se usa en auditorías nuevas, se conserva como respaldo histórico. `cointracking_casos_extended.yaml` (raíz) queda documentado como material de origen ya superado por v2.
- **Fase D (validación):** se verifica sintaxis YAML, 100% de campos del esquema presentes en los 20 casos, y cobertura de las 5 categorías críticas de regresión (transferencias huérfanas, ventas sin base de coste, duplicados, saldos negativos, rendimientos mal clasificados) — todas presentes.

## Consequences

- ✅ El agente dispone de 20 casos con esquema homogéneo, evidencia mínima explícita y trazabilidad de fuente (`fuente_recomendada_para_revalidar`)
- ✅ Los casos de duplicados quedan coherentes con el incidente y la corrección de ADR-014
- ✅ El estado legacy/deprecado de la base anterior queda documentado (cierra la petición de `AGENT_CHANGE_REQUESTS.md` 2026-07-02)
- ⚠️ El contenido sigue siendo conocimiento de patrón (cualitativo); ningún caso constituye una cifra fiscal vinculante (ADR-006/009)
- ⚠️ Los casos `pendiente_verificar`/`hipotesis` requieren reverificación antes de usarse en un informe

## Notes

**Materiales auxiliares:** `LEEME.md` y `PROMPT_CHATGPT_AGENTE.md` (raíz del repo) eran documentación de apoyo para preparar el candidato con ChatGPT; su contenido queda absorbido por este ADR y por `knowledge/patterns/INDEX.md`, por lo que se eliminan tras la integración.

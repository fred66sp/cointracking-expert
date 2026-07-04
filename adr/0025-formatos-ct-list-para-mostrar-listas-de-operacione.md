# ADR-025: Formatos "CT-List" para mostrar listas de operaciones/hallazgos en la conversación

**Status:** Accepted

**Date:** 2026-07-04

## Context

Tras fijar el bloque-resumen CT-Task (ADR-024) para altas/correcciones manuales, el usuario propuso (adaptación de una idea generada externamente) una familia de formatos compactos para mostrar **listas** de operaciones, hallazgos de auditoría, balances y recorridos de fondos en la conversación, evitando párrafos largos o tablas Markdown pesadas de leer cuando hay muchas filas. Se revisó antes de adoptarla: el ejemplo original usaba "Recepción" para un depósito (no es un tipo real de CoinTracking) y citaba el aviso de coste faltante en inglés ("Missing Purchase History") en vez de su forma en español ya documentada. También hacía falta acotar su alcance: no puede sustituir el formato de los informes formales (`templates/AUDIT_REPORT.md`), que necesitan tablas y trazabilidad completa para el asesor (ADR-009/ADR-011).

**Decisión:**

1. Se documenta la familia de formatos (`CT-Timeline`, `CT-Audit`, `CT-Balance`, `CT-Exchange`, `CT-Asset`, `CT-Flow`; el `CT-Task` ya existente de ADR-024 completa la familia) en `knowledge/cointracking/CT_LIST_FORMATS.md`, corregidos los dos errores señalados.
2. **Ámbito exclusivo: la conversación interactiva con el usuario.** Los informes formales de `reports/output/<proyecto>/` siguen el formato de `templates/AUDIT_REPORT.md` (tablas, evidencia/causa/impacto/recomendación) sin cambios.
3. Todo hallazgo marcado `⚠`/`✗` debe ir seguido de la traducción a qué significa / por qué importa / qué hacer (regla ya existente en `CLAUDE.md`, no se relaja por usar un formato compacto).
4. Usar solo los tipos de operación ya verificados contra datos reales del proyecto; no inventar sinónimos.
5. Aplica a cualquiera que use el agente (Claude Code y Copilot, ADR-012).

## Decision

[Decision not found]

## Consequences

- ✅ Listas largas de operaciones (auditorías extensas, historiales) se vuelven más legibles de un vistazo para el usuario, sin sacrificar la trazabilidad de los informes formales.
- ✅ Corrige de raíz las imprecisiones de la propuesta original antes de fijarla como norma (mismo estándar de verificación aplicado a toda fuente externa, ADR-009).
- ⚠️ Riesgo a vigilar: no dejar un hallazgo compacto (`⚠`/`✗`) sin su traducción a lenguaje llano — el formato ahorra espacio, no explicación.

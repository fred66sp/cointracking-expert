# ADR-010: Eficiencia de tokens y caché de datos de CoinTracking

**Status:** Accepted

**Date:** 2026-07-02

## Context

Las respuestas de CoinTracking (MCP/API) pueden ser muy grandes (historial de miles de operaciones, balances con decenas de activos). Volcar ese JSON al contexto del LLM cada vez consume muchos tokens y, si se repite, malgasta también llamadas a la API (límite 60/h). Hay que trabajar de forma económica sin perder rigor (ADR-009).

## Decision

**Decisión — protocolo de eficiencia:**

1. **Caché a disco.** Al obtener datos de CoinTracking, guárdalos en `.cache/cointracking/` (ignorado por git: son datos reales) con **marca de tiempo**. Antes de llamar, comprueba si hay un snapshot reutilizable.
2. **Reutilización.** Dentro de una misma conversación, reutiliza siempre lo ya obtenido (no recalcules ni recargues). Entre sesiones, reutiliza el snapshot si está **fresco**; si es antiguo o el usuario cambió datos, **refresca** y avísalo.
3. **Consultas mínimas y dirigidas.** Pide solo lo necesario: acota por **rango de fechas** y `limit`, y usa **agregados** (`get_grouped_balance`, `get_gains`) antes que el detalle completo. No traigas todo el historial si solo hace falta un ejercicio.
4. **Procesa lo grande con código, no en el contexto.** Para volúmenes grandes (p. ej. historial de operaciones), vuelca a un fichero y usa **scripts** (python/bash) para filtrar/agregar; sube al contexto **solo el resultado compacto** (conteos, totales, filas relevantes), nunca el JSON crudo completo. Cuando sea posible, obtén los datos con utilidades que **escriban directamente a disco** para que no pasen por el contexto.
5. **Nada de JSON crudo en salidas.** Informes y respuestas resumen y citan totales/ejemplos; no pegan volcados completos.
6. **Invalidación por cambios (CRÍTICO).** En cuanto pidas al usuario **modificar algo en CoinTracking** (editar/borrar/añadir operaciones, reimportar, corregir tipos), la caché queda **obsoleta**: márcala como inválida y **no la reutilices**. Antes de volver a dar cifras o informes, **confirma con el usuario que hizo el cambio** y **refresca** los datos (nueva consulta/volcado). Nunca mezcles hallazgos calculados con datos antiguos y datos nuevos.
7. **Verificación de remediaciones por lote, no una a una.** Al guiar al usuario para corregir varios hallazgos en la web de CoinTracking, **no llames al MCP después de cada corrección individual** para comprobarla. Guía el lote completo de correcciones aplicables primero (confirmando por chat, sin consultar la API entre medias); solo cuando el usuario indique que ha terminado la ronda, invalida la caché **una vez** y verifica **todos** los hallazgos corregidos con una consulta agregada. Si el usuario prefiere ir uno a uno, respétalo pero explica el coste extra de cuota (límite 60/hora).

## Consequences

- ✅ Menos tokens y menos llamadas a la API; más rápido y barato
- ✅ Compatible con el rigor: el cálculo determinista sobre datos volcados es más fiable y trazable (ADR-006, ADR-009)
- ⚠️ La caché contiene datos reales → **gitignored**; tratarla como sensible
- ⚠️ Requiere gestionar la frescura de la caché (marca de tiempo, invalidación al cambiar datos)

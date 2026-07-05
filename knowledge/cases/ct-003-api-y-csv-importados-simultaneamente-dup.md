---
id: "KB-C1-003"
title: "Caso CT-003: API y CSV importados simultáneamente (duplicado por doble fuente)"
level: "C"
domain: "cointracking"
source: "Análisis de casos reales auditados"
authority: "verified"
last_verified: "2026-07-05"
valid_from: "2024-01-01"
valid_until:
confidence: "high"
version: "1.0"
related_adr:
  - ADR-003
  - ADR-009
  - ADR-010
related_docs:
  - knowledge/patterns/INDEX.md
  - knowledge/cointracking/COST_BASIS_AND_VALIDATION.md
tags:
  - case
  - duplicados
  - verified
  - operativo
---




# CT-003: API y CSV importados simultáneamente (duplicado por doble fuente)

## Síntomas

- Doble número de operaciones para el mismo periodo
- Holdings superiores a los reales

## Causa Probable

**Hecho:** Existen transacciones con los mismos campos (fecha, importe, activo, exchange).

**Hipótesis:** Se importó la misma cuenta por API y por CSV sin desactivar una de las dos fuentes.

**Supuesto:** No se eliminó la importación previa antes de añadir la segunda.

## Evidencia Mínima

- Mismas fecha, importe, activo y exchange en ambas filas
- trade_id (si el exchange lo expone vía API) — con trade_id distinto, NO son duplicados (ver ADR-014)
- Revisión de Import Statistics para confirmar dos importaciones del mismo origen

## Pasos de Diagnóstico

1. Ordenar las transacciones del exchange por fecha.
1. Buscar pares con todos los campos idénticos.
1. {'Si hay trade_id disponible (API), comparar': 'distinto → no es duplicado; igual o ausente → aplicar la heurística de conteo de copias.'}
1. Revisar Import Statistics para confirmar el origen doble (API + CSV).

## Solución Recomendada

- Eliminar una de las dos importaciones completas (no filas sueltas) para evitar dejar mitades de operaciones enlazadas.
- Recalcular balances tras eliminar.
- Antes de borrar, seguir el protocolo de consentimiento de DECISIONS.md#ADR-014 (listar ejemplos concretos y confirmar con el usuario).

## Anti-patrón

Asumir que cualquier fila con todos los campos iguales es automáticamente un duplicado a eliminar.

## Por qué Falso Positivo

Sin trade_id, campos idénticos pueden corresponder a operaciones reales ejecutadas en el mismo segundo (ver CT-008 y CT-019); solo API+CSV de la misma cuenta en el mismo rango es evidencia suficiente por sí sola.

## Evaluación

- **Confianza:** verificado
- **Riesgo:** alto
- **Impacto fiscal:** Duplica beneficios/pérdidas y movimientos en el informe de ganancias.

## Señales Tempranas

- Holdings duplicados frente al balance real del exchange
- Import Statistics muestra dos importaciones solapadas del mismo exchange

## Validación Antes/Después

**Antes:**
- Operaciones repetidas

**Después:**
- Totales de balance coincidentes con el exchange real

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Nuevos detectores de duplicados de CoinTracking o cambios en qué exchanges exponen trade_id por API.
- **Fuente para revalidar:** knowledge/cointracking/CSV_FORMAT.md §9 (Duplicados exactos) y DECISIONS.md#ADR-014

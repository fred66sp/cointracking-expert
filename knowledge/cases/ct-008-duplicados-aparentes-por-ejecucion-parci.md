---
id: "KB-C1-008"
title: "Caso CT-008: Duplicados aparentes por ejecución parcial de una orden"
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

# CT-008: Duplicados aparentes por ejecución parcial de una orden

## Síntomas

- Varias operaciones muy similares (mismo precio, mismo par) en un intervalo corto

## Causa Probable

**Hecho:** Una orden se ejecutó en varias operaciones parciales (fills) con el mismo precio o precios muy próximos.

**Hipótesis:** Se interpretaron como duplicados por tener campos casi idénticos.

**Supuesto:** El exchange liquida órdenes grandes en varias ejecuciones (normal en libros de órdenes líquidos).

## Evidencia Mínima

- Order ID compartido entre las ejecuciones parciales (si el exchange lo expone)
- trade_id distinto para cada ejecución individual (ver DECISIONS.md#ADR-014)

## Pasos de Diagnóstico

1. Revisar si el exchange asocia un Order ID común a las operaciones sospechosas.
1. {'Comparar las cantidades parciales': 'si suman el tamaño de una orden razonable, es indicio de ejecución parcial legítima.'}
1. Aplicar la heurística de ADR-014 según el número de copias idénticas (2 = revisar; 3-10 = advertencia, no recomendar eliminar; >10 = probablemente legítimo, requiere confirmación en la API del exchange).

## Solución Recomendada

- No eliminar ninguna operación de una misma orden parcial.
- Si hay duda, verificar en la API del exchange que los trade_id son distintos antes de decidir.

## Anti-patrón

Eliminar automáticamente todas las operaciones con el mismo precio y activo por considerarlas duplicado.

## Por qué Falso Positivo

Las ejecuciones parciales de una misma orden son un comportamiento normal de los exchanges y no deben tratarse como duplicado.

## Evaluación

- **Confianza:** verificado
- **Riesgo:** medio
- **Impacto fiscal:** Eliminación indebida de operaciones reales, que reduce artificialmente el balance y puede generar saldo negativo (ver CT-019).

## Señales Tempranas

- Varias filas con el mismo Order ID pero trade_id distinto

## Validación Antes/Después

**Antes:**
- Sospecha de duplicado

**Después:**
- Orden íntegra conservada, sin eliminar ejecuciones legítimas

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en cómo los exchanges exponen Order ID/trade_id vía API.
- **Fuente para revalidar:** DECISIONS.md#ADR-014 y documentación de la API del exchange correspondiente

---
id: "KB-C1-009"
title: "Caso CT-009: Comisión (fee) omitida en la importación"
level: "C"
domain: "cointracking"
source: "Análisis de casos reales auditados"
authority: "verified"
last_verified: "2026-07-05"
valid_from: "2024-01-01"
valid_until: 2027-12-31
confidence: "medium"
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
  - casos_limite_espana
  - verified
  - operativo
---




# CT-009: Comisión (fee) omitida en la importación

## Síntomas

- Coste de adquisición inferior al esperado tras comparar con el extracto del exchange

## Causa Probable

**Hecho:** La fila importada no tiene valor en el campo de comisión.

**Hipótesis:** El CSV exportado por el exchange no incluye la columna de fee, o viene vacía para ese tipo de operación.

**Supuesto:** El exchange no reporta la comisión para ese tipo de operación concreto.

## Evidencia Mínima

- Extracto oficial del exchange mostrando la comisión cobrada
- Fila en CoinTracking con el campo de comisión vacío para la misma operación

## Pasos de Diagnóstico

1. Comparar el campo de comisión de la fila con el extracto oficial del exchange (ver knowledge/cointracking/CSV_FORMAT.md §5).

## Solución Recomendada

- Añadir la comisión manualmente solo si hay soporte documental (extracto o export oficial); no estimarla.

## Evaluación

- **Confianza:** probable
- **Riesgo:** medio
- **Impacto fiscal:** Coste de adquisición potencialmente incompleto, lo que infla la ganancia patrimonial calculada.

## Señales Tempranas

- Trades del mismo exchange con comisión vacía de forma sistemática

## Validación Antes/Después

**Antes:**
- Comisión vacía o ausente

**Después:**
- Comisión registrada y documentada con su fuente

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en el formato de exportación de comisiones de cada exchange.
- **Fuente para revalidar:** knowledge/cointracking/CSV_FORMAT.md §5 y guía de exportación del exchange concreto

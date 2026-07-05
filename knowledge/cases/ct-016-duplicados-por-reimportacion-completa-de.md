---
id: "KB-C1-016"
title: "Caso CT-016: Duplicados por reimportación completa del mismo periodo"
level: "C"
domain: "cointracking"
source: "Análisis de casos reales auditados"
authority: "verified"
last_verified: "2026-07-05"
valid_from: "2024-01-01"
valid_until: 2027-12-31
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




# CT-016: Duplicados por reimportación completa del mismo periodo

## Síntomas

- Bloque completo de operaciones duplicadas para un mismo rango de fechas

## Causa Probable

**Hecho:** El mismo lote de operaciones aparece dos veces.

**Hipótesis:** Se repitió una importación (API o CSV) sin eliminar la anterior.

**Supuesto:** Error operativo del usuario al reimportar tras un problema previo.

## Evidencia Mínima

- Import Statistics mostrando dos importaciones del mismo exchange y rango de fechas
- trade_id repetido en todas las filas del bloque (si está disponible)

## Pasos de Diagnóstico

1. Revisar el historial de importaciones (Import Statistics) para identificar el lote repetido.
1. Confirmar que corresponde exactamente al mismo rango y no a una superposición parcial legítima.

## Solución Recomendada

- Eliminar el lote de importación duplicado completo (no filas sueltas), aplicando el consentimiento explícito de DECISIONS.md#ADR-014.

## Evaluación

- **Confianza:** verificado
- **Riesgo:** alto
- **Impacto fiscal:** Duplicación de resultados (ganancias, pérdidas y volumen) en el informe fiscal.

## Señales Tempranas

- Duplicados masivos concentrados en el mismo rango de fechas de una importación concreta

## Validación Antes/Después

**Antes:**
- Lote duplicado

**Después:**
- Lote único, balances coincidentes con el exchange

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Mejoras del importador de CoinTracking que detecten reimportaciones automáticamente.
- **Fuente para revalidar:** knowledge/cointracking/CSV_FORMAT.md §9 y DECISIONS.md#ADR-014

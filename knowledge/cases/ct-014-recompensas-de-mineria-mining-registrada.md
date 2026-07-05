---
id: "KB-C1-014"
title: "Caso CT-014: Recompensas de minería (mining) registradas como depósito"
level: "C"
domain: "cointracking"
source: "Análisis de casos reales auditados"
authority: "verified"
last_verified: "2026-07-05"
valid_from: "2024-01-01"
valid_until:
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
  - rendimientos
  - verified
  - operativo
---




# CT-014: Recompensas de minería (mining) registradas como depósito

## Síntomas

- Ingresos por minería no visibles en el resumen de rendimientos

## Causa Probable

**Hecho:** Las entradas del pool de minería están clasificadas como "Deposit".

**Hipótesis:** Clasificación errónea por importación genérica.

**Supuesto:** El CSV del pool de minería no distingue el tipo de ingreso.

## Evidencia Mínima

- Historial del pool de minería con fecha e importe de cada pago

## Pasos de Diagnóstico

1. Comparar el historial del pool de minería con las entradas registradas en CoinTracking.

## Solución Recomendada

- Reclasificar como "Mining"/tipo de rendimiento correspondiente.

## Evaluación

- **Confianza:** probable
- **Riesgo:** medio
- **Impacto fiscal:** Clasificación de rendimiento inadecuada en el resumen fiscal.

## Señales Tempranas

- Entradas periódicas de importe variable coincidentes con pagos del pool

## Validación Antes/Después

**Antes:**
- Clasificado como Deposit

**Después:**
- Clasificado como Mining/ingreso correcto

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Nuevos tipos de transacción soportados por CoinTracking para minería.
- **Fuente para revalidar:** CoinTracking Help Center

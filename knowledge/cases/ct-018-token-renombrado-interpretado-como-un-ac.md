---
id: "KB-C1-018"
title: "Caso CT-018: Token renombrado interpretado como un activo distinto"
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
  - casos_limite_espana
  - verified
  - operativo
---




# CT-018: Token renombrado interpretado como un activo distinto

## Síntomas

- El balance de un proyecto aparece dividido entre dos tickers distintos

## Causa Probable

**Hecho:** Existen dos símbolos distintos para lo que es, económicamente, el mismo activo.

**Hipótesis:** El proyecto migró de ticker (rebrand/migración de red) y CoinTracking no fusionó ambos automáticamente.

**Supuesto:** El mapeo de activos de CoinTracking no reconoce todavía la migración.

## Evidencia Mínima

- Anuncio oficial del proyecto confirmando la migración/rebrand y el ratio de conversión

## Pasos de Diagnóstico

1. Verificar en fuentes oficiales del proyecto que se trata de una migración y no de dos activos distintos.
1. Revisar si CoinTracking tiene entrada de mapeo para el ticker nuevo.

## Solución Recomendada

- Actualizar/fusionar el mapeo de activos en CoinTracking siguiendo el ratio oficial de migración (ver también CT-010 de la base legacy sobre colisión de tickers en knowledge/cointracking/CSV_FORMAT.md §8).

## Evaluación

- **Confianza:** pendiente_verificar
- **Riesgo:** medio
- **Impacto fiscal:** Bases de coste separadas incorrectamente entre los dos tickers, distorsionando el FIFO de cada uno.

## Señales Tempranas

- Dos tickers para lo que el usuario reconoce como un único proyecto/activo

## Validación Antes/Después

**Antes:**
- Balance dividido en dos activos

**Después:**
- Balance unificado bajo el ticker vigente

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Nuevas migraciones de proyectos y actualizaciones del mapeo de tickers de CoinTracking.
- **Fuente para revalidar:** Anuncio oficial del proyecto y knowledge/cointracking/CSV_FORMAT.md §8 (colisión de tickers)

---
id: KB-C1-018
title: 'CT-018: Token renombrado interpretado como un activo distinto'
level: C
domain: cointracking
source: Caso auditado en proyecto real
authority: verified
last_verified: '2026-07-03'
valid_from: '2024-01-01'
valid_until: null
confidence: high
version: '1.0'
related_adr:
- ADR-014
- ADR-026
- ADR-004
tags:
- case-study
- casos_limite_espana
- cointracking
notes: 'Categoria: casos_limite_espana'
---

# CT-018: Token renombrado interpretado como un activo distinto

**Categoria:** casos_limite_espana | **Confianza:** pendiente_verificar | **Riesgo:** medio

## Sintomas

- El balance de un proyecto aparece dividido entre dos tickers distintos

## Solucion Recomendada

- Actualizar/fusionar el mapeo de activos en CoinTracking siguiendo el ratio oficial de migración (ver también CT-010 de la base legacy sobre colisión de tickers en knowledge/cointracking/CSV_FORMAT.md §8).

**Impacto fiscal:** Bases de coste separadas incorrectamente entre los dos tickers, distorsionando el FIFO de cada uno.

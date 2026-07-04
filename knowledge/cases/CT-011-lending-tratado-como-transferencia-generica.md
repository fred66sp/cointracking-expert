---
id: KB-C1-011
title: 'CT-011: Lending tratado como transferencia genérica'
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
- rendimientos
- cointracking
notes: 'Categoria: rendimientos'
---

# CT-011: Lending tratado como transferencia genérica

**Categoria:** rendimientos | **Confianza:** probable | **Riesgo:** medio

## Sintomas

- Ingresos por lending ausentes del resumen de rendimientos

## Solucion Recomendada

- Reclasificar como "Income"/tipo de rendimiento correspondiente (ver knowledge/cointracking/WEB_APP_GUIDE.md §5).

**Impacto fiscal:** Rendimientos potencialmente omitidos del resumen fiscal si permanecen como Transfer.

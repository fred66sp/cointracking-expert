---
id: KB-C1-008
title: 'CT-008: Duplicados aparentes por ejecución parcial de una orden'
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
- duplicados
- cointracking
notes: 'Categoria: duplicados'
---

# CT-008: Duplicados aparentes por ejecución parcial de una orden

**Categoria:** duplicados | **Confianza:** verificado | **Riesgo:** medio

## Sintomas

- Varias operaciones muy similares (mismo precio, mismo par) en un intervalo corto

## Solucion Recomendada

- No eliminar ninguna operación de una misma orden parcial.
- Si hay duda, verificar en la API del exchange que los trade_id son distintos antes de decidir.

**Impacto fiscal:** Eliminación indebida de operaciones reales, que reduce artificialmente el balance y puede generar saldo negativo (ver CT-019).

---
id: KB-C1-012
title: 'CT-012: Balance negativo por importación parcial vía API'
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
- saldos_imposibles_o_negativos
- cointracking
notes: 'Categoria: saldos_imposibles_o_negativos'
---

# CT-012: Balance negativo por importación parcial vía API

**Categoria:** saldos_imposibles_o_negativos | **Confianza:** verificado | **Riesgo:** alto

## Sintomas

- Balance negativo de un activo concreto

## Solucion Recomendada

- Completar el periodo faltante con un CSV exportado manualmente desde el exchange.

**Impacto fiscal:** Cálculo FIFO incompleto al faltar compras tempranas.

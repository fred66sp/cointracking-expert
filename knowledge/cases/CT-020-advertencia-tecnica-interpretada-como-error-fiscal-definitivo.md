---
id: KB-C1-020
title: 'CT-020: Advertencia técnica interpretada como error fiscal definitivo'
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

# CT-020: Advertencia técnica interpretada como error fiscal definitivo

**Categoria:** casos_limite_espana | **Confianza:** verificado | **Riesgo:** bajo

## Sintomas

- El usuario asume que cualquier warning del informe de ganancias implica una declaración incorrecta

## Solucion Recomendada

- Corregir únicamente cuando exista evidencia suficiente de un problema real; documentar los warnings revisados y descartados como tales.

**Impacto fiscal:** Variable según el origen real del warning; no asumir impacto sin verificar el caso concreto.

---
id: KB-C1-016
title: 'CT-016: Duplicados por reimportación completa del mismo periodo'
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

# CT-016: Duplicados por reimportación completa del mismo periodo

**Categoria:** duplicados | **Confianza:** verificado | **Riesgo:** alto

## Sintomas

- Bloque completo de operaciones duplicadas para un mismo rango de fechas

## Solucion Recomendada

- Eliminar el lote de importación duplicado completo (no filas sueltas), aplicando el consentimiento explícito de DECISIONS.md#ADR-014.

**Impacto fiscal:** Duplicación de resultados (ganancias, pérdidas y volumen) en el informe fiscal.

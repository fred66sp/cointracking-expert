---
id: KB-C1-017
title: 'CT-017: Coste cero por compra omitida de ejercicios anteriores'
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
- ventas_sin_base_de_coste
- cointracking
notes: 'Categoria: ventas_sin_base_de_coste'
---

# CT-017: Coste cero por compra omitida de ejercicios anteriores

**Categoria:** ventas_sin_base_de_coste | **Confianza:** verificado | **Riesgo:** critico

## Sintomas

- Ganancia excesiva en una venta concreta

## Solucion Recomendada

- Completar el histórico con la fuente anterior (CSV o API).

**Impacto fiscal:** Coste de adquisición inexistente, lo que sobreestima directamente la ganancia patrimonial declarada.

---
id: KB-C1-015
title: 'CT-015: Swap DeFi fragmentado en varias operaciones on-chain'
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
- permutas_complejas
- cointracking
notes: 'Categoria: permutas_complejas'
---

# CT-015: Swap DeFi fragmentado en varias operaciones on-chain

**Categoria:** permutas_complejas | **Confianza:** hipotesis | **Riesgo:** alto

## Sintomas

- Un único swap del usuario aparece como varias operaciones separadas

## Solucion Recomendada

- Reconstruir el flujo completo y registrar el swap como una única permuta (activo de entrada → activo de salida), no como varias operaciones independientes.

**Impacto fiscal:** Coste de adquisición mal asignado si se tratan los pasos intermedios como permutas independientes.

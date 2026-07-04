---
id: KB-C1-004
title: 'CT-004: Balance negativo por orden cronológico incorrecto (zona horaria)'
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

# CT-004: Balance negativo por orden cronológico incorrecto (zona horaria)

**Categoria:** saldos_imposibles_o_negativos | **Confianza:** probable | **Riesgo:** medio

## Sintomas

- Balance negativo temporal que se corrige más adelante en el tiempo

## Solucion Recomendada

- Corregir la fecha/hora de la operación mal ubicada usando el dato real del exchange.
- Reimportar si el desplazamiento afecta a muchas filas.

**Impacto fiscal:** Puede alterar el orden FIFO y por tanto qué lote de compra se consume en cada venta.

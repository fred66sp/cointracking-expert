---
id: KB-C1-003
title: 'CT-003: API y CSV importados simultáneamente (duplicado por doble fuente)'
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

# CT-003: API y CSV importados simultáneamente (duplicado por doble fuente)

**Categoria:** duplicados | **Confianza:** verificado | **Riesgo:** alto

## Sintomas

- Doble número de operaciones para el mismo periodo
- Holdings superiores a los reales

## Solucion Recomendada

- Eliminar una de las dos importaciones completas (no filas sueltas) para evitar dejar mitades de operaciones enlazadas.
- Recalcular balances tras eliminar.
- Antes de borrar, seguir el protocolo de consentimiento de DECISIONS.md#ADR-014 (listar ejemplos concretos y confirmar con el usuario).

**Impacto fiscal:** Duplica beneficios/pérdidas y movimientos en el informe de ganancias.

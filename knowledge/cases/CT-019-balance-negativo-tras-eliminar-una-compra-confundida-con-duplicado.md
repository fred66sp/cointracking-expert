---
id: KB-C1-019
title: 'CT-019: Balance negativo tras eliminar una compra confundida con duplicado'
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

# CT-019: Balance negativo tras eliminar una compra confundida con duplicado

**Categoria:** saldos_imposibles_o_negativos | **Confianza:** probable | **Riesgo:** alto

## Sintomas

- Balance negativo aparecido justo después de que el usuario eliminara una operación "duplicada"

## Solucion Recomendada

- Restaurar la operación eliminada (reimportar o registrar manualmente) si se confirma que era legítima.
- Aplicar en adelante el protocolo de DECISIONS.md#ADR-014 antes de eliminar cualquier duplicado sospechoso.

**Impacto fiscal:** FIFO alterado y pérdida real de balance del activo hasta restaurar desde backup.

---
id: KB-C1-005
title: 'CT-005: Recompensas de staking clasificadas como depósito genérico'
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

# CT-005: Recompensas de staking clasificadas como depósito genérico

**Categoria:** rendimientos | **Confianza:** probable | **Riesgo:** medio

## Sintomas

- No aparecen ingresos de staking en el resumen de rendimientos, aunque el exchange sí los pagó

## Solucion Recomendada

- Reclasificar manualmente como "Staking" o "Income" según corresponda (ver knowledge/cointracking/WEB_APP_GUIDE.md §5, edición masiva por tipo).

**Impacto fiscal:** Rendimiento potencialmente mal clasificado a efectos de IRPF; el tratamiento fiscal exacto de staking está marcado como pendiente en knowledge/taxation/spain/PENDIENTES.md — no asumir la calificación sin verificar ese documento.

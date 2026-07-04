---
id: KB-C1-007
title: 'CT-007: Transferencia interna confundida con venta'
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
- transferencias_huerfanas
- cointracking
notes: 'Categoria: transferencias_huerfanas'
---

# CT-007: Transferencia interna confundida con venta

**Categoria:** transferencias_huerfanas | **Confianza:** verificado | **Riesgo:** alto

## Sintomas

- Ganancia inesperada sin motivo aparente

## Solucion Recomendada

- Cambiar el tipo de la operación a "Transfer" (ver knowledge/cointracking/WEB_APP_GUIDE.md §2).

**Impacto fiscal:** Ganancia patrimonial ficticia si no se corrige antes de generar el informe fiscal.

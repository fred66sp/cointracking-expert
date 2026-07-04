---
id: KB-C1-001
title: 'CT-001: Transferencia entre exchanges importada solo en origen'
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

# CT-001: Transferencia entre exchanges importada solo en origen

**Categoria:** transferencias_huerfanas | **Confianza:** verificado | **Riesgo:** alto

## Sintomas

- Balance negativo en el exchange destino
- Advertencias de transacciones incompletas
- Holdings inferiores a los reales

## Solucion Recomendada

- Importar el historial faltante del exchange destino (API o CSV).
- Si no existe registro de origen, crear la transferencia manualmente enlazando ambos lados (ver knowledge/cointracking/WEB_APP_GUIDE.md §2).
- Validar que el balance del activo queda no negativo tras recalcular.

**Impacto fiscal:** Puede provocar ventas sin base de coste y balances erróneos que arrastran el cálculo FIFO.

---
id: KB-C1-006
title: 'CT-006: Binance Convert importado como venta y compra independientes'
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

# CT-006: Binance Convert importado como venta y compra independientes

**Categoria:** permutas_complejas | **Confianza:** probable | **Riesgo:** alto

## Sintomas

- Costes de adquisición incoherentes tras usar la función Convert de Binance

## Solucion Recomendada

- Si CoinTracking ya trata Convert como una permuta única (Trade), no requiere ajuste; verificarlo antes de tocar nada.
- Si aparece como dos filas independientes sin vínculo, ajustar manualmente para que la base de coste de la compra resultante sea el valor de mercado del activo entregado (ver knowledge/cointracking/COST_BASIS_AND_VALIDATION.md §2).

**Impacto fiscal:** Base de coste incorrecta en la permuta, que en España tributa como ganancia/pérdida patrimonial en el momento del cambio.

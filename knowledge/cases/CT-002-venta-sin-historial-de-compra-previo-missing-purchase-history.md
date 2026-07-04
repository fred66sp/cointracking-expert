---
id: KB-C1-002
title: 'CT-002: Venta sin historial de compra previo (Missing Purchase History)'
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

# CT-002: Venta sin historial de compra previo (Missing Purchase History)

**Categoria:** ventas_sin_base_de_coste | **Confianza:** verificado | **Riesgo:** critico

## Sintomas

- Advertencia "No hay una compra adecuada para esta venta"
- Ganancias extremadamente elevadas o coste por unidad irreal

## Solucion Recomendada

- Importar los años/exchanges anteriores que falten.
- Registrar manualmente la compra solo si existe evidencia documental (extracto, hash on-chain).
- Recalcular el informe de ganancias tras completar el historial.

**Impacto fiscal:** Cálculo incorrecto de ganancias patrimoniales (base de coste inexistente o cero); afecta directamente a la declaración.

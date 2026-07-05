---
id: KB-B1-XXX
title: "Untitled Document"
level: B
domain: cointracking
source: "Internal documentation"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: medium
version: 1.0

tags:
  - todo
  - needs-review

notes: "Metadatos agregados automáticamente. Verificar y actualizar conforme ADR-032."
---

---
id: "KB-C1-006"
title: "Caso CT-006: Binance Convert importado como venta y compra independientes"
level: "C"
domain: "cointracking"
source: "Análisis de casos reales auditados"
authority: "verified"
last_verified: "2026-07-05"
valid_from: "2024-01-01"
valid_until:
confidence: "medium"
version: "1.0"
related_adr:
  - ADR-003
  - ADR-009
  - ADR-010
related_docs:
  - knowledge/patterns/INDEX.md
  - knowledge/cointracking/COST_BASIS_AND_VALIDATION.md
tags:
  - case
  - permutas_complejas
  - verified
  - operativo
---

# CT-006: Binance Convert importado como venta y compra independientes

## Síntomas

- Costes de adquisición incoherentes tras usar la función Convert de Binance

## Causa Probable

**Hecho:** La conversión aparece como dos operaciones separadas (una venta y una compra) sin vínculo entre ellas.

**Hipótesis:** La importación es parcial o no reconoce el identificador de conversión de Binance.

**Supuesto:** El CSV/API de Binance Convert no expone un identificador único que CoinTracking pueda enlazar.

## Evidencia Mínima

- Historial de "Convert" de Binance con fecha, par de activos e importes
- Las dos filas correspondientes en CoinTracking (venta + compra) con la misma marca temporal

## Pasos de Diagnóstico

1. Comparar la operación en CoinTracking con el historial "Convert" de Binance.
1. Revisar si el importador tiene un modo específico para Convert (distinto del genérico de Trade).

## Solución Recomendada

- Si CoinTracking ya trata Convert como una permuta única (Trade), no requiere ajuste; verificarlo antes de tocar nada.
- Si aparece como dos filas independientes sin vínculo, ajustar manualmente para que la base de coste de la compra resultante sea el valor de mercado del activo entregado (ver knowledge/cointracking/COST_BASIS_AND_VALIDATION.md §2).

## Evaluación

- **Confianza:** probable
- **Riesgo:** alto
- **Impacto fiscal:** Base de coste incorrecta en la permuta, que en España tributa como ganancia/pérdida patrimonial en el momento del cambio.

## Señales Tempranas

- Diferencias de coste apreciables justo después de usar Convert

## Validación Antes/Después

**Antes:**
- Permuta con coste inconsistente

**Después:**
- Activos conciliados con coste de mercado en el momento del cambio

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en el formato de exportación de Binance Convert o en cómo lo interpreta el importador de CoinTracking.
- **Fuente para revalidar:** Binance (historial de Convert) y knowledge/cointracking/COST_BASIS_AND_VALIDATION.md §2

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
id: "KB-C1-002"
title: "Caso CT-002: Venta sin historial de compra previo (Missing Purchase History)"
level: "C"
domain: "cointracking"
source: "Análisis de casos reales auditados"
authority: "verified"
last_verified: "2026-07-05"
valid_from: "2024-01-01"
valid_until:
confidence: "high"
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
  - ventas_sin_base_de_coste
  - verified
  - operativo
---

# CT-002: Venta sin historial de compra previo (Missing Purchase History)

## Síntomas

- Advertencia "No hay una compra adecuada para esta venta"
- Ganancias extremadamente elevadas o coste por unidad irreal

## Causa Probable

**Hecho:** Existe una venta registrada para un activo.

**Hipótesis:** No existe la compra previa correspondiente en el pool de compras.

**Supuesto:** El historial de compras está incompleto (rango de fechas o exchange no importado).

## Evidencia Mínima

- Venta registrada con fecha, importe y activo
- Ausencia confirmada de compra previa en el pool tras revisar todo el historial importado
- Rango de fechas efectivamente cubierto por la importación

## Pasos de Diagnóstico

1. Revisar el detalle de la advertencia "Missing Purchase History" en el informe de ganancias.
1. Buscar compras anteriores del mismo activo en cualquier exchange/wallet del usuario.
1. Revisar si el rango temporal importado empieza después de la fecha real de adquisición.

## Solución Recomendada

- Importar los años/exchanges anteriores que falten.
- Registrar manualmente la compra solo si existe evidencia documental (extracto, hash on-chain).
- Recalcular el informe de ganancias tras completar el historial.

## Evaluación

- **Confianza:** verificado
- **Riesgo:** critico
- **Impacto fiscal:** Cálculo incorrecto de ganancias patrimoniales (base de coste inexistente o cero); afecta directamente a la declaración.

## Señales Tempranas

- Advertencia inmediata al generar el informe de ganancias tras importar ventas

## Validación Antes/Después

**Antes:**
- Advertencia "Missing Purchase History"

**Después:**
- Advertencia eliminada tras completar el historial

- Coste de adquisición > 0 para esa venta

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en el algoritmo de cálculo FIFO o en el texto/lógica de la advertencia.
- **Fuente para revalidar:** knowledge/cointracking/COST_BASIS_AND_VALIDATION.md §3.1 y §3.2

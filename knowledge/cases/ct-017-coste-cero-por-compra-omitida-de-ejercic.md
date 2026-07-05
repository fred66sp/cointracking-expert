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
id: "KB-C1-017"
title: "Caso CT-017: Coste cero por compra omitida de ejercicios anteriores"
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

# CT-017: Coste cero por compra omitida de ejercicios anteriores

## Síntomas

- Ganancia excesiva en una venta concreta

## Causa Probable

**Hecho:** No hay compra registrada para el activo vendido dentro del historial importado.

**Hipótesis:** El historial de un ejercicio anterior no se importó (p. ej. tras cambiar de exchange).

**Supuesto:** El usuario cambió de plataforma y no trajo el histórico completo de la anterior.

## Evidencia Mínima

- Extracto histórico del exchange/plataforma anterior mostrando la compra original

## Pasos de Diagnóstico

1. Buscar el origen real del activo (exchange o wallet anterior).
1. Revisar el rango de fechas cubierto por la importación actual frente a la fecha de la compra original.

## Solución Recomendada

- Completar el histórico con la fuente anterior (CSV o API).

## Evaluación

- **Confianza:** verificado
- **Riesgo:** critico
- **Impacto fiscal:** Coste de adquisición inexistente, lo que sobreestima directamente la ganancia patrimonial declarada.

## Señales Tempranas

- Advertencia "Missing Purchase History" asociada a esa venta

## Validación Antes/Después

**Antes:**
- Coste cero

**Después:**
- Coste real recuperado y FIFO correcto

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en el algoritmo de cálculo FIFO de CoinTracking.
- **Fuente para revalidar:** knowledge/cointracking/COST_BASIS_AND_VALIDATION.md §3.2

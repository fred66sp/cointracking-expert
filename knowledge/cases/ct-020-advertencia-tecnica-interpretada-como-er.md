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
id: "KB-C1-020"
title: "Caso CT-020: Advertencia técnica interpretada como error fiscal definitivo"
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
  - casos_limite_espana
  - verified
  - operativo
---

# CT-020: Advertencia técnica interpretada como error fiscal definitivo

## Síntomas

- El usuario asume que cualquier warning del informe de ganancias implica una declaración incorrecta

## Causa Probable

**Hecho:** Existe una advertencia técnica visible en el informe (p. ej. sobre FIAT extranjero o pools agotados).

**Hipótesis:** La advertencia es preventiva y no siempre implica un error real en los datos.

**Supuesto:** Falta revisión manual del caso concreto antes de concluir que hay un problema.

## Evidencia Mínima

- Texto exacto del warning
- Historial relacionado con la operación o activo señalado

## Pasos de Diagnóstico

1. Leer el detalle completo del warning (no solo el titular).
1. Contrastar con la evidencia real (balances, historial del exchange) antes de concluir que hay un error.
1. Verificar balances tras la revisión.

## Solución Recomendada

- Corregir únicamente cuando exista evidencia suficiente de un problema real; documentar los warnings revisados y descartados como tales.

## Anti-patrón

Asumir que todo warning implica automáticamente una declaración incorrecta.

## Por qué Falso Positivo

Muchos avisos de CoinTracking son preventivos (p. ej. FIAT negativo, ver knowledge/cointracking/COST_BASIS_AND_VALIDATION.md §4.1) y requieren validación manual antes de concluir que hay un error.

## Evaluación

- **Confianza:** verificado
- **Riesgo:** bajo
- **Impacto fiscal:** Variable según el origen real del warning; no asumir impacto sin verificar el caso concreto.

## Señales Tempranas

- Aparición de advertencias nuevas tras una importación o recálculo

## Validación Antes/Después

**Antes:**
- Warning sin revisar

**Después:**
- Diagnóstico documentado (confirmado como error real o descartado como preventivo)

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en las reglas de validación/advertencias de CoinTracking.
- **Fuente para revalidar:** knowledge/cointracking/COST_BASIS_AND_VALIDATION.md §3-4 y CoinTracking Help Center

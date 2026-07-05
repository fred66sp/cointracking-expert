---
id: "KB-C1-004"
title: "Caso CT-004: Balance negativo por orden cronológico incorrecto (zona horaria)"
level: "C"
domain: "cointracking"
source: "Análisis de casos reales auditados"
authority: "verified"
last_verified: "2026-07-05"
valid_from: "2024-01-01"
valid_until: 2027-12-31
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
  - saldos_imposibles_o_negativos
  - verified
  - operativo
---




# CT-004: Balance negativo por orden cronológico incorrecto (zona horaria)

## Síntomas

- Balance negativo temporal que se corrige más adelante en el tiempo

## Causa Probable

**Hecho:** Dos operaciones relacionadas (p. ej. venta y compra previa) aparecen en orden invertido en el libro mayor.

**Hipótesis:** El CSV se exportó con una zona horaria distinta a la configurada, desplazando timestamps cerca de la medianoche o de un cambio de DST.

**Supuesto:** El usuario no declaró o declaró mal la zona horaria de origen al importar (ver DECISIONS.md#ADR-005).

## Evidencia Mínima

- Par de operaciones consecutivas cuyo orden real (confirmado por el exchange) no coincide con el orden en CoinTracking
- Hora exacta de ambas operaciones en el exchange de origen

## Pasos de Diagnóstico

1. Revisar el timestamp exacto de las operaciones afectadas.
1. Comparar con el historial del exchange (que suele mostrar la hora en UTC o en la zona de la cuenta).
1. Verificar si la fecha cae cerca de un cambio de horario de verano/invierno (riesgo específico de Europe/Madrid, ver ADR-005).

## Solución Recomendada

- Corregir la fecha/hora de la operación mal ubicada usando el dato real del exchange.
- Reimportar si el desplazamiento afecta a muchas filas.

## Evaluación

- **Confianza:** probable
- **Riesgo:** medio
- **Impacto fiscal:** Puede alterar el orden FIFO y por tanto qué lote de compra se consume en cada venta.

## Señales Tempranas

- Balance negativo puntual que desaparece si se reordenan dos operaciones cercanas en el tiempo

## Validación Antes/Después

**Antes:**
- Operaciones fuera de orden real

**Después:**
- Cronología consistente con el exchange de origen

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en cómo CoinTracking exporta/importa timestamps, o corrección del riesgo residual de ADR-005.
- **Fuente para revalidar:** DECISIONS.md#ADR-005 y knowledge/cointracking/CSV_FORMAT.md §2

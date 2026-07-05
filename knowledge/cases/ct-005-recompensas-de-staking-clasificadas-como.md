---
id: "KB-C1-005"
title: "Caso CT-005: Recompensas de staking clasificadas como depósito genérico"
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
  - rendimientos
  - verified
  - operativo
---




# CT-005: Recompensas de staking clasificadas como depósito genérico

## Síntomas

- No aparecen ingresos de staking en el resumen de rendimientos, aunque el exchange sí los pagó

## Causa Probable

**Hecho:** Las entradas de staking están registradas con tipo "Deposit" en vez de "Staking"/"Income".

**Hipótesis:** El importador (API o CSV) no distingue el subtipo de la operación.

**Supuesto:** El exchange no etiqueta la recompensa de forma diferenciada en su export.

## Evidencia Mínima

- Historial del exchange mostrando el tipo original de la operación (p. ej. "Staking Reward")
- Fecha e importe de cada entrada periódica sospechosa

## Pasos de Diagnóstico

1. Revisar el tipo de transacción asignado en CoinTracking para esas entradas.
1. Comparar con el historial del exchange para confirmar que son recompensas de staking y no un depósito real.

## Solución Recomendada

- Reclasificar manualmente como "Staking" o "Income" según corresponda (ver knowledge/cointracking/WEB_APP_GUIDE.md §5, edición masiva por tipo).

## Evaluación

- **Confianza:** probable
- **Riesgo:** medio
- **Impacto fiscal:** Rendimiento potencialmente mal clasificado a efectos de IRPF; el tratamiento fiscal exacto de staking está marcado como pendiente en knowledge/taxation/spain/PENDIENTES.md — no asumir la calificación sin verificar ese documento.

## Señales Tempranas

- Entradas periódicas de importe pequeño y regular, sin contrapartida de compra

## Validación Antes/Después

**Antes:**
- Tipo "Deposit" incorrecto

**Después:**
- Clasificación coherente con el origen real de los fondos

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Evolución de los tipos de transacción soportados por CoinTracking, o cierre del pendiente fiscal de staking.
- **Fuente para revalidar:** CoinTracking Help Center y knowledge/taxation/spain/PENDIENTES.md

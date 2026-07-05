---
id: "KB-C1-019"
title: "Caso CT-019: Balance negativo tras eliminar una compra confundida con duplicado"
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




# CT-019: Balance negativo tras eliminar una compra confundida con duplicado

## Síntomas

- Balance negativo aparecido justo después de que el usuario eliminara una operación "duplicada"

## Causa Probable

**Hecho:** Se eliminó una compra que en realidad era una ejecución parcial legítima o una operación con trade_id distinto.

**Hipótesis:** Se aplicó la regla "duplicado = campos idénticos" sin comprobar el trade_id (ver CT-003, CT-008).

**Supuesto:** El usuario actuó sobre una detección de duplicados sin verificar contra la API del exchange.

## Evidencia Mínima

- Order ID / trade_id de la operación eliminada, si aún es recuperable
- Confirmación en la API del exchange de que existían varias operaciones legítimas con esos campos

## Pasos de Diagnóstico

1. Revisar el historial de cambios (REGISTRO-CAMBIOS.md si existe) para identificar qué se eliminó y cuándo.
1. Verificar en la API del exchange si la operación eliminada tenía un trade_id propio distinto de las que se conservaron.

## Solución Recomendada

- Restaurar la operación eliminada (reimportar o registrar manualmente) si se confirma que era legítima.
- Aplicar en adelante el protocolo de DECISIONS.md#ADR-014 antes de eliminar cualquier duplicado sospechoso.

## Anti-patrón

Eliminar cualquier operación con el mismo importe, activo y fecha por parecerse a otra.

## Por qué Falso Positivo

Puede tratarse de una ejecución parcial legítima (ver CT-008) o de un batching de operaciones reales en el mismo segundo, como el caso documentado en ADR-014.

## Evaluación

- **Confianza:** probable
- **Riesgo:** alto
- **Impacto fiscal:** FIFO alterado y pérdida real de balance del activo hasta restaurar desde backup.

## Señales Tempranas

- Balance negativo inmediatamente después de una limpieza de duplicados

## Validación Antes/Después

**Antes:**
- Balance negativo

**Después:**
- Balance normal tras restaurar la operación legítima

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en los importadores o en la disponibilidad de trade_id por exchange.
- **Fuente para revalidar:** DECISIONS.md#ADR-014 (incidente de origen documentado el 2026-07-03)

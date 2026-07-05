---
id: "KB-C1-012"
title: "Caso CT-012: Balance negativo por importación parcial vía API"
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
  - saldos_imposibles_o_negativos
  - verified
  - operativo
---




# CT-012: Balance negativo por importación parcial vía API

## Síntomas

- Balance negativo de un activo concreto

## Causa Probable

**Hecho:** El historial importado por API no cubre el periodo completo de operaciones del usuario.

**Hipótesis:** La API del exchange tiene una limitación de rango temporal o de número de resultados.

**Supuesto:** Existe una restricción temporal propia de la API de ese exchange.

## Evidencia Mínima

- Fecha de la primera operación devuelta por la API
- Confirmación (en el propio exchange) de operaciones anteriores a esa fecha

## Pasos de Diagnóstico

1. Comparar la fecha de la primera operación en CoinTracking con la fecha real de apertura de la cuenta en el exchange.

## Solución Recomendada

- Completar el periodo faltante con un CSV exportado manualmente desde el exchange.

## Evaluación

- **Confianza:** verificado
- **Riesgo:** alto
- **Impacto fiscal:** Cálculo FIFO incompleto al faltar compras tempranas.

## Señales Tempranas

- La primera compra registrada en CoinTracking es posterior a la primera compra real conocida por el usuario

## Validación Antes/Después

**Antes:**
- Balance negativo

**Después:**
- Balance correcto tras completar el historial

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en los límites de rango temporal de la API de cada exchange.
- **Fuente para revalidar:** knowledge/cointracking/MCP_API.md y documentación de la API del exchange

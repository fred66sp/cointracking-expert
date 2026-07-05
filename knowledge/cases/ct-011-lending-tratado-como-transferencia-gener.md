---
id: "KB-C1-011"
title: "Caso CT-011: Lending tratado como transferencia genérica"
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
  - rendimientos
  - verified
  - operativo
---




# CT-011: Lending tratado como transferencia genérica

## Síntomas

- Ingresos por lending ausentes del resumen de rendimientos

## Causa Probable

**Hecho:** Las entradas de lending están clasificadas como "Transfer" en vez de un tipo de rendimiento.

**Hipótesis:** Error de importación por un tipo genérico mal mapeado.

**Supuesto:** La plataforma de lending no distingue el pago de intereses de un movimiento normal en su export.

## Evidencia Mínima

- Historial de la plataforma de lending mostrando los pagos de interés
- Fecha e importe de cada entrada periódica

## Pasos de Diagnóstico

1. Revisar el tipo asignado en CoinTracking a esas entradas.
1. Comparar con el historial de la plataforma para confirmar que son intereses de lending.

## Solución Recomendada

- Reclasificar como "Income"/tipo de rendimiento correspondiente (ver knowledge/cointracking/WEB_APP_GUIDE.md §5).

## Evaluación

- **Confianza:** probable
- **Riesgo:** medio
- **Impacto fiscal:** Rendimientos potencialmente omitidos del resumen fiscal si permanecen como Transfer.

## Señales Tempranas

- Entradas periódicas de importe pequeño sin contrapartida de compra

## Validación Antes/Después

**Antes:**
- Clasificado como Transfer

**Después:**
- Clasificado como ingreso/rendimiento correcto

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Nuevos tipos de transacción soportados por CoinTracking para plataformas de lending.
- **Fuente para revalidar:** CoinTracking Help Center

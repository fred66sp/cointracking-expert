---
id: "KB-C1-015"
title: "Caso CT-015: Swap DeFi fragmentado en varias operaciones on-chain"
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
  - permutas_complejas
  - verified
  - operativo
---




# CT-015: Swap DeFi fragmentado en varias operaciones on-chain

## Síntomas

- Un único swap del usuario aparece como varias operaciones separadas

## Causa Probable

**Hecho:** El explorador on-chain muestra múltiples eventos (aprobación, paso por token puente, swap final) para una sola acción del usuario.

**Hipótesis:** El router del protocolo DeFi divide la operación internamente (p. ej. a través de varios pools o tokens intermedios).

**Supuesto:** El usuario solo percibió "un swap", pero on-chain generó varios movimientos técnicos.

## Evidencia Mínima

- Hash de la transacción(es) en el explorador on-chain
- Secuencia completa de eventos del mismo bloque/transacción

## Pasos de Diagnóstico

1. Seguir el hash de la transacción en el explorador on-chain para ver todos los eventos asociados.
1. Distinguir aprobaciones (sin efecto económico) de movimientos reales de valor.

## Solución Recomendada

- Reconstruir el flujo completo y registrar el swap como una única permuta (activo de entrada → activo de salida), no como varias operaciones independientes.

## Anti-patrón

Tratar cada movimiento on-chain visible como una operación fiscal independiente.

## Por qué Falso Positivo

Un único swap del usuario puede generar múltiples eventos técnicos (aprobación, tokens puente) que no son transmisiones económicas separadas.

## Evaluación

- **Confianza:** hipotesis
- **Riesgo:** alto
- **Impacto fiscal:** Coste de adquisición mal asignado si se tratan los pasos intermedios como permutas independientes.

## Señales Tempranas

- Muchos movimientos en el mismo bloque/transacción para una sola acción del usuario

## Validación Antes/Después

**Antes:**
- Operación fragmentada en varias filas

**Después:**
- Swap conciliado como una única permuta

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Nuevos protocolos DeFi con mecánicas de enrutamiento distintas.
- **Fuente para revalidar:** Documentación del protocolo DeFi concreto y explorador on-chain correspondiente

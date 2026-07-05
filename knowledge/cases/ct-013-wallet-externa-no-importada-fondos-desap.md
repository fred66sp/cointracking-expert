---
id: "KB-C1-013"
title: "Caso CT-013: Wallet externa no importada (fondos "desaparecidos")"
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
  - transferencias_huerfanas
  - verified
  - operativo
---




# CT-013: Wallet externa no importada (fondos "desaparecidos")

## Síntomas

- Fondos que salieron de un exchange y no aparecen en ningún otro lugar de CoinTracking

## Causa Probable

**Hecho:** Solo existe el registro de la retirada en el exchange de origen.

**Hipótesis:** La wallet de destino (autocustodia) no está importada en CoinTracking.

**Supuesto:** El usuario mueve el activo a una wallet propia (hardware wallet, wallet de software) que no ha añadido a la cuenta.

## Evidencia Mínima

- Tx Hash de la retirada
- Dirección de destino visible en el explorador on-chain

## Pasos de Diagnóstico

1. Buscar el depósito correspondiente entre las wallets/exchanges ya importados.
1. Si no aparece, confirmar con el Tx Hash en un explorador on-chain que la dirección de destino pertenece al usuario.

## Solución Recomendada

- Registrar la wallet externa en CoinTracking (manual o vía dirección pública) y añadir el depósito correspondiente.

## Anti-patrón

Asumir que toda retirada de un exchange implica una venta.

## Por qué Falso Positivo

Puede tratarse de un movimiento a autocustodia (wallet propia), que no es una transmisión a efectos fiscales.

## Evaluación

- **Confianza:** verificado
- **Riesgo:** alto
- **Impacto fiscal:** Venta ficticia si se interpreta la retirada como transmisión sin comprobar el destino.

## Señales Tempranas

- Balance del activo inferior al esperado tras una retirada, sin venta registrada

## Validación Antes/Después

**Antes:**
- Movimiento inconsistente (huérfano)

**Después:**
- Conciliado con la wallet de destino registrada

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Nuevos conectores de wallets/exploradores soportados por CoinTracking.
- **Fuente para revalidar:** knowledge/cointracking/WEB_APP_GUIDE.md §2

---
id: "KB-C1-010"
title: "Caso CT-010: Airdrop registrado como compra con coste artificial"
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
  - casos_limite_espana
  - verified
  - operativo
---

# CT-010: Airdrop registrado como compra con coste artificial

## Síntomas

- El activo recibido por airdrop aparece con un coste de adquisición mayor que cero sin que haya habido pago

## Causa Probable

**Hecho:** La entrada está clasificada con un tipo que implica compra (Trade/Deposit con contravalor).

**Hipótesis:** Error de clasificación manual, o el importador no distingue el origen "airdrop" de una compra normal.

**Supuesto:** El usuario desconocía el origen exacto del activo al momento de registrarlo.

## Evidencia Mínima

- Anuncio o comunicación oficial del proyecto confirmando el airdrop
- Fecha y cantidad recibida, contrastada con la wallet/exchange

## Pasos de Diagnóstico

1. Revisar el origen documentado del activo (anuncio del proyecto, explorador on-chain).
1. Confirmar que no hubo ninguna contraprestación pagada por el usuario.

## Solución Recomendada

- Reclasificar el tipo de operación acorde al origen documentado (p. ej. "Airdrop"/"Income" si CoinTracking lo soporta).
- El tratamiento fiscal exacto del airdrop en IRPF depende de si hubo contraprestación y de la calificación de la AEAT/DGT vigente; no asumir un criterio sin consultar knowledge/taxation/spain/ y, si no está cubierto, marcarlo `[VERIFICAR]` y consultar fuente oficial antes de declarar.

## Anti-patrón

Tratar toda entrada gratuita de un activo exactamente igual (todos los airdrops son lo mismo).

## Por qué Falso Positivo

El tratamiento fiscal depende de la naturaleza documentada del ingreso (p. ej. si requirió una acción del usuario, si hubo contraprestación); no puede generalizarse sin verificar caso a caso.

## Evaluación

- **Confianza:** pendiente_verificar
- **Riesgo:** medio
- **Impacto fiscal:** Clasificación fiscal potencialmente incorrecta; requiere verificación con fuente oficial antes de trasladarlo a un informe.

## Señales Tempranas

- Entrada de activo sin pago asociado ni contrapartida de venta

## Validación Antes/Después

**Antes:**
- Registrado como compra con coste artificial

**Después:**
- Clasificación revisada y respaldada por documentación del origen

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en el criterio interpretativo de la AEAT/DGT sobre airdrops.
- **Fuente para revalidar:** Agencia Tributaria (AEAT/DGT) y knowledge/taxation/spain/PENDIENTES.md

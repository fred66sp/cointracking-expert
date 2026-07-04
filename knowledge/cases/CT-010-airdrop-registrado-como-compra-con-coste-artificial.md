---
id: KB-C1-010
title: 'CT-010: Airdrop registrado como compra con coste artificial'
level: C
domain: cointracking
source: Caso auditado en proyecto real
authority: verified
last_verified: '2026-07-03'
valid_from: '2024-01-01'
valid_until: null
confidence: high
version: '1.0'
related_adr:
- ADR-014
- ADR-026
- ADR-004
tags:
- case-study
- casos_limite_espana
- cointracking
notes: 'Categoria: casos_limite_espana'
---

# CT-010: Airdrop registrado como compra con coste artificial

**Categoria:** casos_limite_espana | **Confianza:** pendiente_verificar | **Riesgo:** medio

## Sintomas

- El activo recibido por airdrop aparece con un coste de adquisición mayor que cero sin que haya habido pago

## Solucion Recomendada

- Reclasificar el tipo de operación acorde al origen documentado (p. ej. "Airdrop"/"Income" si CoinTracking lo soporta).
- El tratamiento fiscal exacto del airdrop en IRPF depende de si hubo contraprestación y de la calificación de la AEAT/DGT vigente; no asumir un criterio sin consultar knowledge/taxation/spain/ y, si no está cubierto, marcarlo `[VERIFICAR]` y consultar fuente oficial antes de declarar.

**Impacto fiscal:** Clasificación fiscal potencialmente incorrecta; requiere verificación con fuente oficial antes de trasladarlo a un informe.

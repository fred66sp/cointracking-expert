---
id: KB-B1-016
title: "Salida de Binance de la UE por MiCA (2026-07) — impacto en reconciliación"
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

# Salida de Binance de la UE por MiCA (2026-07) — impacto en reconciliación

**Tipo:** Contexto regulatorio (fuentes de prensa/noticias, no oficiales de CoinTracking ni de la AEAT)
**Última verificación:** 2026-07-03
**Vigencia:** situación en evolución activa a fecha de escritura — Binance planea volver a solicitar licencia MiCA; reverificar el estado antes de asumir que sigue vigente si esta fecha es antigua (ADR-008).
**Estado:** Fundamentado en fuentes de prensa (no oficiales), verificado en la sesión mediante búsqueda web

> ⚠️ Este documento **no es** doctrina fiscal ni conocimiento oficial de CoinTracking — es contexto operativo para saber qué buscar durante una auditoría cuando el usuario ha tenido que migrar de exchange por este motivo.

## 1. Qué ha pasado

El Reglamento MiCA (Markets in Crypto-Assets) de la UE exige a los proveedores de servicios cripto tener una licencia MiCA de al menos un Estado miembro antes del **1 de julio de 2026** para poder operar en los 27 países de la UE. **Binance no consiguió la licencia a tiempo**: retiró su solicitud en Grecia el 24 de junio de 2026 y, desde el 1 de julio de 2026, ha cortado a los usuarios de la UE las **nuevas órdenes spot, depósitos, altas de cuenta y productos Earn/staking**. Las **retiradas siguen abiertas** (los fondos existentes son accesibles). Binance planea volver a solicitar la licencia (posiblemente vía Francia).

Fuentes (prensa, verificadas en la sesión):
- [Binance tells EU users it will no longer provide services after failing to secure MiCA license (CoinDesk)](https://www.coindesk.com/policy/2026/06/26/binance-tells-eu-users-it-will-no-longer-provide-services-after-failing-to-secure-mica-license)
- [Binance to halt crypto services across EU countries after failing to secure MiCA approval (Euronews)](https://www.euronews.com/business/2026/06/25/binance-to-halt-crypto-services-across-eu-countries-after-failing-to-secure-mica-approval)
- [Binance is locked out of Europe on July 1. Here is what actually happened (crypto.news)](https://crypto.news/binance-eu-mica-license-lockout-july-2026-explained/)

## 2. Por qué importa para la auditoría/reconciliación (no es fiscalidad nueva)

**MiCA no crea una regla fiscal nueva** — la tributación de las operaciones cripto en España sigue siendo la de siempre (`../taxation/spain/CAPITAL_GAINS.md`). Lo que MiCA genera es un **evento operativo forzoso** en el exchange que puede dejar rastro en los datos de CoinTracking y hay que saber interpretar:

1. **Migración de activos entre exchanges (p. ej. Binance → Coinbase, caso real del usuario, 2026-07):** si el usuario retira todo o parte de su cartera de Binance por este motivo, es una **transferencia entre cuentas propias** — no tributa (`CAPITAL_GAINS.md` §1) — **siempre que CoinTracking la registre correctamente como par retirada/depósito emparejado**, no como venta. Auditar estas migraciones con el mismo protocolo de transferencias huérfanas ya establecido (`audit-cointracking/SKILL.md`, fase 3): Tx Hash si existe, o heurística de importe/ventana temporal.
2. **Posible conversión o liquidación forzosa antes de retirar:** si Binance obliga a convertir algún activo no conforme (p. ej. cierto stablecoin) antes de permitir la retirada, **eso sí sería una permuta cripto-cripto** con hecho imponible real (Art. 37.1.h LIRPF) — no confundirla con la transferencia en sí. Distinguir ambas cosas es crítico para no perder ni inventar un hecho imponible.
3. **Ya se observó un precedente del mismo patrón** (no MiCA-Binance-EU, pero mismo tipo de evento): la conversión forzosa USDT→USDC de Binance en Q1 2025 por otra normativa, ya detectada y verificada sin discrepancias en la auditoría de `agp2025` (ver `reports/output/agp2025/`).
4. **No asumir que la migración está completa o correcta solo porque CoinTracking la muestre coherente internamente** — aplica la misma regla de verificación en dos capas ya usada en el resto de la auditoría (`COST_BASIS_AND_VALIDATION.md` §4, "regla crítica"): cotejar contra el histórico real de Binance/Coinbase, no solo contra CoinTracking.

## 3. Qué hacer en la próxima auditoría de una cuenta afectada

- Buscar en el periodo de la migración (aprox. desde finales de junio de 2026) las retiradas de Binance y los depósitos correspondientes en el exchange de destino.
- Verificar que estén emparejadas (no huérfanas) y que los importes cuadren tras comisiones.
- Revisar si hay alguna operación de tipo `Trade`/conversión justo antes de la retirada que sugiera una liquidación forzosa de un activo concreto, y clasificarla como hecho imponible si corresponde.
- Si el usuario menciona que "tuvo que mover todo de Binance", no dar la migración por completa sin contrastar el histórico de Binance antes del cierre del acceso (las retiradas siguen abiertas, pero el acceso al histórico de operaciones puede degradarse con el tiempo si Binance restringe más la cuenta).

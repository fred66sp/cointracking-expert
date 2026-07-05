---
id: KB-B1-022
title: "Troubleshooting: Síntomas y Soluciones"
level: B
domain: cointracking
source: "Documentación interna"
authority: reference
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: medium
version: 1.0

tags:
  - navigation
  - reference

notes: "Documento de navegación/referencia de la base de conocimiento"
---

## Warnings literales de CoinTracking

| Warning / texto en CoinTracking | Significado breve | Gravedad | Ver |
|---|---|---|---|
| "There is no suitable purchase to this sale, all purchasing pools consumed" | Se vende un activo sin compra con base de coste registrada | Crítica (infla la ganancia) | `COST_BASIS_AND_VALIDATION.md` §3.1 · casos CT-002, CT-017 |
| Coste por unidad irreal / base de coste en 0 en el informe de ganancias | Falta la compra; CoinTracking asigna coste 0 | Crítica | `COST_BASIS_AND_VALIDATION.md` §3.2 · casos CT-002, CT-017 |
| Advertencia sobre FIAT extranjero (divisa distinta a la principal de la cuenta) | Solo la moneda principal de la cuenta tiene soporte completo | Baja-media | `COST_BASIS_AND_VALIDATION.md` §3.3 |
| Saldo FIAT negativo | Normalmente **no es un error**: es fiat gastado en cripto sin depósito FIAT previo registrado | Informativa (verificar antes de tratarla como error) | `COST_BASIS_AND_VALIDATION.md` §4.1 |
| "Reoccurring Duplicate Transactions" / operaciones repetidas | Patrón reconocido por CoinTracking; puede ser reimportación real o batching legítimo del exchange | Variable — **nunca eliminar sin verificar** | `COST_BASIS_AND_VALIDATION.md` §4.2 · `DECISIONS.md#ADR-014` · casos CT-003, CT-008, CT-016, CT-019 |
| Cualquier warning del informe de ganancias, sin más contexto | No todo warning implica declaración incorrecta; muchos son preventivos | Depende del caso concreto | caso CT-020 |

---


# Troubleshooting: de síntoma/warning a causa y solución

**Tipo:** Índice de enrutamiento (no añade conocimiento nuevo)
**Última verificación:** 2026-07-04
**Vigencia:** enruta a `COST_BASIS_AND_VALIDATION.md`, `CSV_FORMAT.md`, `DECISIONS.md` y `knowledge/patterns/cointracking_casos_v2.yaml`; revisar si esos documentos cambian (ADR-008).

Este documento **no repite** el conocimiento ya destilado — lo indexa por síntoma para encontrarlo rápido cuando CoinTracking muestra un warning concreto o el usuario describe un problema. Es el punto de entrada; el detalle completo (causa, evidencia mínima, pasos de diagnóstico, solución) vive en el documento o caso enlazado.

**Cómo usarlo:**
1. Localiza el síntoma/warning más parecido en la tabla.
2. Ve al documento/caso enlazado para la explicación completa y los pasos de diagnóstico.
3. Antes de **modificar o eliminar** cualquier dato en CoinTracking, aplica el consentimiento informado (ADR-009 §7) y, si es un posible duplicado, el protocolo de `trade_id` (ADR-014).
4. Si el síntoma no aparece aquí: no improvises (ADR-009 §2) — busca en el resto de `knowledge/`, y si tampoco está, en la fuente oficial (`reference/CATALOG.md`) o pregunta al usuario.



## Síntomas descritos por el usuario (sin warning literal)

| Síntoma | Causas más probables | Ver |
|---|---|---|
| Balance negativo de un activo | Transferencia solo importada en origen; orden cronológico invertido por zona horaria; importación API parcial; eliminación indebida de una operación confundida con duplicado | casos CT-001, CT-004, CT-012, CT-019 · `DECISIONS.md#ADR-005` |
| Fondos que "desaparecen" tras una retirada | Wallet externa (autocustodia) no importada — no asumir que es una venta | caso CT-013 |
| Ganancia inesperada sin venta real aparente | Transferencia interna clasificada por error como Trade/venta | caso CT-007 |
| Ingresos de staking/lending/minería que no aparecen en el resumen de rendimientos | Reclasificado por el importador como "Deposit"/"Transfer" genérico en vez del tipo correcto | casos CT-005, CT-011, CT-014 |
| Costes incoherentes tras usar Binance Convert / un swap DeFi | La conversión o swap se importó como operaciones independientes sin vincular | casos CT-006, CT-015 |
| Coste de adquisición menor de lo esperado frente al extracto del exchange | Comisión (fee) omitida en la importación | caso CT-009 |
| Balance de un proyecto dividido entre dos tickers | Token renombrado/migrado que CoinTracking no fusionó automáticamente | caso CT-018 |
| Activo aparece con coste de adquisición aunque fue un airdrop | Clasificado con un tipo que implica compra | caso CT-010 (tratamiento fiscal `pendiente_verificar`, ver `knowledge/taxation/spain/PENDIENTES.md`) |
| Doble número de operaciones para el mismo periodo/exchange | Misma cuenta importada dos veces (API + CSV, o reimportación completa) | casos CT-003, CT-016 |

---

## Antes de dar por buenos los balances (metodología oficial)

Orden recomendado por CoinTracking ("READ FIRST") para sanear una cuenta, ya integrado en la skill `audit-cointracking`:

1. Importar todos los exchanges/wallets (evitar seguimiento de un solo lado).
2. Comparar con los saldos reales — nunca dar por bueno un balance solo por ser internamente consistente.
3. Detectar y eliminar duplicados (con el protocolo de `trade_id`, ADR-014).
4. Resolver pequeñas discrepancias de saldo.
5. Corregir monedas faltantes / cálculos erróneos.
6. Revisar transacciones faltantes (depósitos sin retirada emparejada).
7. Considerar forks e ICOs.

Ver `COST_BASIS_AND_VALIDATION.md` §4.3 para el detalle.
